package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/jonatanblue/tweet-operator/pkg/libs/k8sclient"
	"github.com/jonatanblue/tweet-operator/pkg/libs/twitterclient"

	"github.com/jonatanblue/tweet-operator/pkg/reconciler"

	"k8s.io/client-go/rest"

	tweetclient "github.com/jonatanblue/tweet-operator/pkg/client/clientset/versioned"
)

type runMode string

const (
	runModeLoop    = runMode("loop")
	runModeRunOnce = runMode("run-once")
)

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
	runMode := runModeLoop
	// Lookup optional run mode env var
	if os.Getenv("RUN_MODE") == "run-once" {
		runMode = runModeRunOnce
	}

	// Kubernetes client
	kubeConfig, err := getKubeConfig()
	if err != nil {
		log.Fatal(err)
	}
	tweetClientSet := tweetclient.NewForConfigOrDie(kubeConfig)
	tweetClient := tweetClientSet.ExampleV1().Tweets("default")
	k8sClient := k8sclient.NewK8sClient(tweetClient)

	// Twitter client
	creds := twitterclient.Credentials{
		ConsumerKey:       mustLookupEnv("CONSUMER_KEY"),
		ConsumerSecret:    mustLookupEnv("CONSUMER_SECRET"),
		AccessToken:       mustLookupEnv("ACCESS_TOKEN"),
		AccessTokenSecret: mustLookupEnv("ACCESS_TOKEN_SECRET"),
	}
	apiClient, err := twitterclient.NewTwitterAPIClient(&creds)
	if err != nil {
		log.Fatal(err)
	}
	twitterClient := twitterclient.NewTwitterClient(
		apiClient.Statuses,
		apiClient.Timelines,
	)

	// Reconciler
	reconciler := reconciler.NewTweetReconciler(
		k8sClient,
		twitterClient,
		mustLookupEnv("TWITTER_USERNAME"),
	)

	log.Print("Starting reconciliation loop...")
	for true {
		reconciled, err := reconciler.Reconcile()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("main: reconciled=%v", reconciled)

		if runMode == runModeRunOnce {
			break
		}

		<-time.After(10 * time.Second)
	}
}
