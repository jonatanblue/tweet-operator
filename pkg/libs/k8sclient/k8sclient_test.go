package k8sclient

import (
	"errors"
	"testing"

	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"context"

	v1 "github.com/jonatanblue/tweet-operator/pkg/apis/example.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_GetTweet(t *testing.T) {
	tests := map[string]struct {
		client *K8sClient
		name   string
		want   *tweettypes.Tweet
		err    error
	}{
		"tweet found no error": {
			client: NewK8sClient(
				newTweetClientMock(
					"Get",
					"hello-world",
					&v1.Tweet{
						ObjectMeta: metav1.ObjectMeta{
							Name: "hello-world",
						},
						Spec: v1.TweetSpec{
							Text: "Hello World",
						},
						Status: v1.TweetStatus{
							ID:       12345,
							Likes:    0,
							Retweets: 0,
							Replies:  0,
						},
					},
					nil,
				),
			),
			name: "hello-world",
			want: &tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: "Hello World",
				},
				Status: tweettypes.TweetStatus{
					ID:       12345,
					Likes:    0,
					Retweets: 0,
					Replies:  0,
				},
			},
			err: nil,
		},
		"tweet not found error": {
			client: NewK8sClient(
				newTweetClientMock(
					"Get",
					"hello-world",
					nil,
					errors.New("not found"),
				),
			),
			name: "hello-world",
			want: nil,
			err:  errors.New("not found"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tweet, err := test.client.GetTweet(test.name)
			if err != nil {
				assert.EqualError(t, err, test.err.Error())
			}
			assert.Equal(t, test.want, tweet)
		})
	}
}

func Test_UpdateStatus(t *testing.T) {
	tests := map[string]struct {
		client *K8sClient
		name   string
		in     *tweettypes.Tweet
		calls  int
		err    error
	}{
		"tweet status updated no error": {
			client: NewK8sClient(
				newTweetClientMockUpdateStatus(),
			),
			name: "hello-world",
			in: &tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: "Hello World",
				},
				Status: tweettypes.TweetStatus{
					ID:       12345,
					Likes:    1,
					Retweets: 2,
					Replies:  3,
				},
			},
			calls: 1,
			err:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.client.UpdateStatus(test.name, test.in)
			if err != nil {
				assert.EqualError(t, err, test.err.Error())
			}
			test.client.tweetClient.(*tweetClientMock).AssertNumberOfCalls(t, "Get", test.calls)
			test.client.tweetClient.(*tweetClientMock).AssertNumberOfCalls(t, "Update", test.calls)
		})
	}
}

func Test_ListTweets(t *testing.T) {
	tests := map[string]struct {
		client *K8sClient
		want   *tweettypes.Tweets
		err    error
	}{
		"found 1 tweet": {
			client: NewK8sClient(
				newTweetClientMock(
					"List",
					metav1.ListOptions{},
					&v1.TweetList{
						Items: []v1.Tweet{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "hello-world",
								},
								Spec: v1.TweetSpec{
									Text: "Hello World",
								},
								Status: v1.TweetStatus{
									ID:       12345,
									Likes:    0,
									Retweets: 0,
									Replies:  0,
								},
							},
						},
					},
					nil,
				),
			),
			want: &tweettypes.Tweets{
				{
					Spec: tweettypes.TweetSpec{
						Name: "hello-world",
						Text: "Hello World",
					},
					Status: tweettypes.TweetStatus{
						ID:       12345,
						Likes:    0,
						Retweets: 0,
						Replies:  0,
					},
				},
			},
		},
		"found 2 tweets": {
			client: NewK8sClient(
				newTweetClientMock(
					"List",
					metav1.ListOptions{},
					&v1.TweetList{
						Items: []v1.Tweet{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "hello-world",
								},
								Spec: v1.TweetSpec{
									Text: "Hello World",
								},
								Status: v1.TweetStatus{},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "hello-world-2",
								},
								Spec: v1.TweetSpec{
									Text: "Hello World 2",
								},
								Status: v1.TweetStatus{},
							},
						},
					},
					nil,
				),
			),
			want: &tweettypes.Tweets{
				{
					Spec: tweettypes.TweetSpec{
						Name: "hello-world",
						Text: "Hello World",
					},
					Status: tweettypes.TweetStatus{},
				},
				{
					Spec: tweettypes.TweetSpec{
						Name: "hello-world-2",
						Text: "Hello World 2",
					},
					Status: tweettypes.TweetStatus{},
				},
			},
		},
		"found no tweets": {
			client: NewK8sClient(
				newTweetClientMock(
					"List",
					metav1.ListOptions{},
					&v1.TweetList{
						Items: []v1.Tweet{},
					},
					nil,
				),
			),
			want: &tweettypes.Tweets{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tweets, err := test.client.ListTweets()
			if err != nil {
				assert.EqualError(t, err, test.err.Error())
			}
			assert.Equal(t, test.want, tweets)
		})
	}
}

func newTweetClientMock(methodName string, arg interface{}, ret interface{}, err error) *tweetClientMock {
	client := new(tweetClientMock)
	client.On(methodName, arg).Return(ret, err)
	return client
}

func newTweetClientMockUpdateStatus() *tweetClientMock {

	client := new(tweetClientMock)

	client.On(
		"Update",
		&v1.Tweet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "hello-world",
			},
			Status: v1.TweetStatus{
				ID:       12345,
				Likes:    1,
				Retweets: 2,
				Replies:  3,
			},
		},
		metav1.UpdateOptions{},
	).Return(
		nil,
		nil,
	).On(
		"Get",
		"hello-world",
	).Return(
		&v1.Tweet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "hello-world",
			},
			Status: v1.TweetStatus{
				ID:       12345,
				Likes:    1,
				Retweets: 2,
				Replies:  3,
			},
		},
		nil,
	)

	return client
}

type tweetClientMock struct {
	mock.Mock
}

func (mock *tweetClientMock) Create(ctx context.Context, tweet *v1.Tweet, opts metav1.CreateOptions) (*v1.Tweet, error) {
	args := mock.Called(tweet, opts)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*v1.Tweet), args.Error(1)
}

func (mock *tweetClientMock) Update(ctx context.Context, tweet *v1.Tweet, opts metav1.UpdateOptions) (*v1.Tweet, error) {
	args := mock.Called(tweet, opts)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*v1.Tweet), args.Error(1)
}

func (mock *tweetClientMock) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Tweet, error) {
	args := mock.Called(name)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*v1.Tweet), args.Error(1)
}

func (mock *tweetClientMock) List(ctx context.Context, opts metav1.ListOptions) (*v1.TweetList, error) {
	args := mock.Called(opts)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*v1.TweetList), args.Error(1)
}
