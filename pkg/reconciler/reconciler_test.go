package reconciler

import (
	"testing"

	"github.com/pkg/errors"

	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ReconcileOne(t *testing.T) {
	tests := map[string]struct {
		reconciler TweetReconciler
		desired    *tweettypes.Tweet
		actual     *tweettypes.Tweet
		reconciled bool
		method     string
		calls      int
		err        error
	}{
		"tweet exists reconciled": {
			reconciler: TweetReconciler{},
			desired: &tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: "Hello World",
				},
			},
			actual: &tweettypes.Tweet{
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
			reconciled: true,
			method:     "",
			calls:      0,
			err:        nil,
		},
		"tweet does not exist reconciled": {
			reconciler: TweetReconciler{},
			desired:    &tweettypes.Tweet{},
			actual:     &tweettypes.Tweet{},
			reconciled: true,
			method:     "",
			calls:      0,
			err:        nil,
		},
		"tweet does not exist tweet created not reconciled": {
			reconciler: TweetReconciler{
				twitterClient: newTwitterClientMock(
					"PostTweet",
					&tweettypes.Tweet{
						Spec: tweettypes.TweetSpec{
							Text: "Hello World",
						},
					},
					nil,
					nil,
				),
			},
			desired: &tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: "Hello World",
				},
			},
			actual:     &tweettypes.Tweet{},
			reconciled: false,
			method:     "PostTweet",
			calls:      1,
			err:        nil,
		},
		"desired not found tweet deleted not reconciled": {
			reconciler: TweetReconciler{
				twitterClient: newTwitterClientMock(
					"DeleteTweet",
					&tweettypes.Tweet{
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
					nil,
					nil,
				),
			},
			desired: &tweettypes.Tweet{},
			actual: &tweettypes.Tweet{
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
			reconciled: false,
			method:     "DeleteTweet",
			calls:      1,
			err:        nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reconciled, err := test.reconciler.ReconcileOne(test.desired, test.actual)
			if err != test.err {
				t.Errorf("expected error %v, got %v", test.err, err)
			}
			if reconciled != test.reconciled {
				t.Errorf("expected reconciled %v, got %v", test.reconciled, reconciled)
			}
		})
		if test.calls > 0 {
			test.reconciler.twitterClient.(*twitterClientMock).AssertNumberOfCalls(t, test.method, test.calls)
		}
	}
}

func Test_getDesiredState(t *testing.T) {
	tests := map[string]struct {
		reconciler TweetReconciler
		tweetName  string
		desired    *tweettypes.Tweet
		err        error
	}{
		"no tweet object exists desired should be empty": {
			reconciler: TweetReconciler{
				k8sClient: NewK8sClientMockGetTweetNoError(
					"hello-world",
					&tweettypes.Tweet{},
				),
			},
			tweetName: "hello-world",
			desired:   &tweettypes.Tweet{},
			err:       nil,
		},
		"tweet object exists desired set": {
			reconciler: TweetReconciler{
				k8sClient: NewK8sClientMockGetTweetNoError(
					"hello-world",
					&tweettypes.Tweet{
						Spec: tweettypes.TweetSpec{
							Text: "Hello World",
						},
					}),
			},
			tweetName: "hello-world",
			desired: &tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: "Hello World",
				},
			},
			err: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := test.reconciler.getDesiredState(test.tweetName)
			assert.NoError(t, err)
			assert.Equal(t, test.desired, actual)
		})
	}
}

func Test_getActualState(t *testing.T) {
	tests := map[string]struct {
		reconciler TweetReconciler
		text       string
		expected   *tweettypes.Tweet
		calls      int
		err        error
	}{
		"tweet not found no error": {
			reconciler: TweetReconciler{
				twitterClient: newTwitterClientMock(
					"GetTweetsForUser",
					"bob",
					&tweettypes.Tweets{},
					nil,
				),
				twitterUserName: "bob",
			},
			text:     "Hello World",
			expected: &tweettypes.Tweet{},
			calls:    1,
			err:      nil,
		},
		"tweet found no error": {
			reconciler: TweetReconciler{
				twitterClient: newTwitterClientMock(
					"GetTweetsForUser",
					"bob",
					&tweettypes.Tweets{
						{
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
					},
					nil,
				),
				twitterUserName: "bob",
			},
			text: "Hello World",
			expected: &tweettypes.Tweet{
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
		"tweet not found error": {
			reconciler: TweetReconciler{
				twitterClient: newTwitterClientMock(
					"GetTweetsForUser",
					"bob",
					nil,
					errors.New("some error"),
				),
				twitterUserName: "bob",
			},
			text:     "Hello World",
			expected: nil,
			calls:    1,
			err:      errors.New("failed to get tweets: some error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := test.reconciler.getActualState(test.text)
			assertError(t, test.err, err)
			assert.Equal(t, test.expected, actual)
		})
		if test.calls > 0 {
			test.reconciler.twitterClient.(*twitterClientMock).AssertNumberOfCalls(t, "GetTweetsForUser", test.calls)
		}
	}

}

func assertError(t *testing.T, expected error, actual error) {
	if expected == nil {
		assert.Nil(t, actual)
	} else {
		assert.EqualError(t, expected, actual.Error())
	}
}

func newTwitterClientMock(methodName string, arg interface{}, ret interface{}, err error) *twitterClientMock {
	client := new(twitterClientMock)
	client.On(methodName, arg).Return(ret, err)
	return client
}

type twitterClientMock struct {
	mock.Mock
}

func (mock *twitterClientMock) GetTweetsForUser(userName string) (*tweettypes.Tweets, error) {
	args := mock.Called(userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tweettypes.Tweets), args.Error(1)
}

func (mock *twitterClientMock) PostTweet(tweet *tweettypes.Tweet) error {
	args := mock.Called(tweet)
	return args.Error(0)
}

func (mock *twitterClientMock) DeleteTweet(tweet *tweettypes.Tweet) error {
	args := mock.Called(tweet)
	return args.Error(0)
}

type k8sClientMock struct {
	mock.Mock
}

func (mock *k8sClientMock) GetTweet(name string) (*tweettypes.Tweet, error) {
	args := mock.Called(name)
	return args.Get(0).(*tweettypes.Tweet), args.Error(1)
}

func (mock *k8sClientMock) UpdateStatus(name string, tweet *tweettypes.Tweet) error {
	args := mock.Called(tweet)
	return args.Error(0)
}

func (mock *k8sClientMock) ListTweets() (*tweettypes.Tweets, error) {
	args := mock.Called()
	return args.Get(0).(*tweettypes.Tweets), args.Error(1)
}

func NewK8sClientMockGetTweetNoError(tweetName string, tweet *tweettypes.Tweet) *k8sClientMock {
	client := new(k8sClientMock)
	client.On("GetTweet", tweetName).Return(tweet, nil)
	return client
}
