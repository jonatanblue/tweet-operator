package k8sclient

import (
	"context"

	v1 "github.com/jonatanblue/tweet-operator/pkg/apis/example.com/v1"

	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type tweetClient interface {
	Create(ctx context.Context, tweet *v1.Tweet, opts metav1.CreateOptions) (*v1.Tweet, error)
	UpdateStatus(ctx context.Context, tweet *v1.Tweet, opts metav1.UpdateOptions) (*v1.Tweet, error)
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
	_, err := c.tweetClient.UpdateStatus(
		context.TODO(),
		&v1.Tweet{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Status: v1.TweetStatus{
				ID:       tweet.Status.ID,
				Likes:    tweet.Status.Likes,
				Retweets: tweet.Status.Retweets,
				Replies:  tweet.Status.Replies,
			},
		},
		metav1.UpdateOptions{},
	)
	return err
}
