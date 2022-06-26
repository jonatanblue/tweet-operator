package k8sclient

import (
	"context"

	v1 "github.com/jonatanblue/tweet-operator/pkg/apis/example.com/v1"

	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type tweetClient interface {
	Create(ctx context.Context, tweet *v1.Tweet, opts metav1.CreateOptions) (*v1.Tweet, error)
	Update(ctx context.Context, tweet *v1.Tweet, opts metav1.UpdateOptions) (*v1.Tweet, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Tweet, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.TweetList, error)
}

type K8sClient struct {
	tweetClient tweetClient
}

func NewK8sClient(tweetClient tweetClient) *K8sClient {
	return &K8sClient{
		tweetClient: tweetClient,
	}
}

func (c *K8sClient) GetTweet(name string) (*tweettypes.Tweet, error) {
	tweet, err := c.tweetClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &tweettypes.Tweet{
		Spec: tweettypes.TweetSpec{
			Name: tweet.Name,
			Text: tweet.Spec.Text,
		},
		Status: tweettypes.TweetStatus{
			ID:       tweet.Status.ID,
			Likes:    tweet.Status.Likes,
			Retweets: tweet.Status.Retweets,
			Replies:  tweet.Status.Replies,
		},
	}, nil
}

func (c *K8sClient) UpdateStatus(name string, tweet *tweettypes.Tweet) error {
	t, err := c.tweetClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	new := t.DeepCopy()
	new.Status = v1.TweetStatus{
		ID:       tweet.Status.ID,
		Likes:    tweet.Status.Likes,
		Retweets: tweet.Status.Retweets,
		Replies:  tweet.Status.Replies,
	}
	_, err = c.tweetClient.Update(
		context.TODO(),
		new,
		metav1.UpdateOptions{},
	)
	return err
}

func (c *K8sClient) ListTweets() (*tweettypes.Tweets, error) {
	list, err := c.tweetClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	tweets := tweettypes.Tweets{}
	for _, t := range list.Items {
		tweets = append(tweets, tweettypes.Tweet{
			Spec: tweettypes.TweetSpec{
				Name: t.Name,
				Text: t.Spec.Text,
			},
			Status: tweettypes.TweetStatus{
				ID:       t.Status.ID,
				Likes:    t.Status.Likes,
				Retweets: t.Status.Retweets,
				Replies:  t.Status.Replies,
			},
		})
	}
	return &tweets, nil
}
