package reconciler

import (
	"log"

	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"
	"github.com/pkg/errors"
)

type K8sClient interface {
	GetTweet(name string) (*tweettypes.Tweet, error)
	UpdateStatus(name string, tweet *tweettypes.Tweet) error
	ListTweets() (*tweettypes.Tweets, error)
}

type TwitterClient interface {
	GetTweetsForUser(userName string) (result *tweettypes.Tweets, err error)
	PostTweet(tweet *tweettypes.Tweet) error
	DeleteTweet(tweet *tweettypes.Tweet) error
}

type TweetReconciler struct {
	k8sClient       K8sClient
	twitterClient   TwitterClient
	twitterUserName string
}

func NewTweetReconciler(
	k8sClient K8sClient,
	twitterClient TwitterClient,
	twitterUserName string,
) *TweetReconciler {
	return &TweetReconciler{
		k8sClient:       k8sClient,
		twitterClient:   twitterClient,
		twitterUserName: twitterUserName,
	}
}

func (reconciler *TweetReconciler) Reconcile() (bool, error) {
	desiredTweets, err := reconciler.k8sClient.ListTweets()
	if err != nil {
		return false, errors.Wrapf(err, "failed to get tweet list from k8s")
	}
	log.Printf("Got tweets from k8s, %+v", desiredTweets)

	for _, t := range *desiredTweets {
		log.Printf("Reconciling tweet, %+v", t)
		desired, err := reconciler.getDesiredState(t.Spec.Name)
		if err != nil {
			return false, errors.Wrapf(err, "failed to get desired state for %s", t.Spec.Name)
		}
		actual, err := reconciler.getActualState(desired.Spec.Text)
		if err != nil {
			return false, errors.Wrapf(err, "failed to get actual state for %s", t.Spec.Name)
		}
		reconciled, err := reconciler.ReconcileOne(desired, actual)
		if err != nil {
			return false, errors.Wrapf(err, "failed to reconcile %s", t.Spec.Name)
		}

		// Update custom resource with latest status
		err = reconciler.k8sClient.UpdateStatus(t.Spec.Name, actual)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update status for %s", t.Spec.Name)
		}

		if !reconciled {
			return false, nil
		}
	}
	return true, nil
}

func (reconciler *TweetReconciler) ReconcileOne(desired, actual *tweettypes.Tweet) (reconciled bool, err error) {
	if desired.Spec.Text == "" {
		if actual.Spec.Text != "" {
			log.Printf("Deleting tweet with ID, %v", actual.Status.ID)
			err := reconciler.twitterClient.DeleteTweet(actual)
			if err != nil {
				return false, errors.Wrap(err, "failed to delete tweet")
			}
			return false, nil
		}
	} else {
		if actual.Spec.Text == "" {
			err := reconciler.twitterClient.PostTweet(desired)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}
	return true, nil
}

func (reconciler *TweetReconciler) getDesiredState(name string) (*tweettypes.Tweet, error) {
	desired, err := reconciler.k8sClient.GetTweet(name)
	if err != nil {
		return nil, err
	}
	return desired, nil
}

func (reconciler *TweetReconciler) getActualState(text string) (*tweettypes.Tweet, error) {
	log.Printf("Getting tweets for user %s...", reconciler.twitterUserName)
	tweets, err := reconciler.twitterClient.GetTweetsForUser(reconciler.twitterUserName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tweets")
	}

	for _, tweet := range *tweets {
		if tweet.Spec.Text == text {
			return &tweet, nil
		}
	}

	return &tweettypes.Tweet{}, nil
}
