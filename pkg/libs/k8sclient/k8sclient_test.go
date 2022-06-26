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
				newTweetClientMockMultiArg(
					"UpdateStatus",
					[]interface{}{
						&v1.Tweet{
							ObjectMeta: metav1.ObjectMeta{
								Name: "hello-world",
							},
							Status: v1.TweetStatus{
								ID:       12345,
								Likes:    0,
								Retweets: 0,
								Replies:  0,
							},
						},
						metav1.UpdateOptions{},
					},
					nil,
					nil,
				),
			),
			name: "hello-world",
			in: &tweettypes.Tweet{
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
			test.client.tweetClient.(*tweetClientMock).AssertNumberOfCalls(t, "Update", test.calls)
		})
	}
}

func newTweetClientMock(methodName string, arg interface{}, ret interface{}, err error) *tweetClientMock {
	client := new(tweetClientMock)
	client.On(methodName, arg).Return(ret, err)
	return client
}

func newTweetClientMockMultiArg(methodName string, args []interface{}, ret interface{}, err error) *tweetClientMock {
	client := new(tweetClientMock)
	client.On(methodName, args...).Return(ret, err)
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
