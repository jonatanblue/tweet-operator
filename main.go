package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	tweetClient "github.com/jonatanblue/tweet-operator/pkg/client/clientset/versioned"
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

func main() {

	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("Building config from flags, %s", err.Error())
	}

	tweetClient := tweetClient.NewForConfigOrDie(config)

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

	twitterAPIClient, err := Twitter(creds)
	if err != nil {
		log.Fatal(err)
	}

	t, r, err := twitterAPIClient.Statuses.Update(tweet.Spec.Text, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("tweet: %+v\n", t)
	log.Printf("response: %+v\n", r)

}
