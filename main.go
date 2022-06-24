package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	tweetClient "github.com/jonatanblue/tweet-operator/pkg/client/clientset/versioned"

	"k8s.io/client-go/rest"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func Twitter(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(false),
	}
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func mustLookupEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("%s must be set", key)
	}
	return value
}

func inClusterConfigAvailable() bool {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	return len(host) > 0 && len(port) > 0
}

func getKubeConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error
	if inClusterConfigAvailable() {
		config, err = rest.InClusterConfig()
	} else {
		if value, ok := os.LookupEnv("KUBECONFIG"); ok {
			config, err = clientcmd.BuildConfigFromFlags("", value)
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
			if err != nil {
				return nil, err
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	kubeConfig, err := getKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	tweetClient := tweetClient.NewForConfigOrDie(kubeConfig)

	tweet, err := tweetClient.ExampleV1().Tweets("default").Get(context.Background(), "hello-world", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	creds := &Credentials{
		ConsumerKey:       mustLookupEnv("CONSUMER_KEY"),
		ConsumerSecret:    mustLookupEnv("CONSUMER_SECRET"),
		AccessToken:       mustLookupEnv("ACCESS_TOKEN"),
		AccessTokenSecret: mustLookupEnv("ACCESS_TOKEN_SECRET"),
	}

	log.Print("main: getting Twitter client...")
	twitterAPIClient, err := Twitter(creds)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("creds are good - but not tweeting just yet")
	os.Exit(0)

	log.Print("main: tweeting...")
	t, r, err := twitterAPIClient.Statuses.Update(tweet.Spec.Text, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("tweet: %+v\n", t)
	log.Printf("response: %+v\n", r)

}
