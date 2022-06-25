package twitterclient

import (
	"net/http"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetTweetsForUser(t *testing.T) {
	tests := map[string]struct {
		client *TwitterClient
		name   string
		want   *tweettypes.Tweets
		calls  int
		err    error
	}{
		"found 1 tweet": {
			client: NewTwitterClient(
				nil,
				newTimelineClientMock(
					"UserTimeline",
					[]interface{}{
						&twitter.UserTimelineParams{
							UserID:     0,
							ScreenName: "bob",
							Count:      200,
							SinceID:    0,
							MaxID:      0,
						},
					},
					[]twitter.Tweet{
						{
							ID:        12345,
							Text:      "Hello World",
							CreatedAt: "2020-01-01T00:00:00Z",
							User: &twitter.User{
								Name: "Bob",
							},
							FavoriteCount: 1,
							RetweetCount:  2,
							ReplyCount:    3,
						},
					},
					nil,
				),
			),
			name: "bob",
			want: &tweettypes.Tweets{
				{
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
			},
			calls: 1,
			err:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tweets, err := test.client.GetTweetsForUser(test.name)
			if err != nil {
				assert.EqualError(t, test.err, err.Error())
			}
			assert.Equal(t, test.want, tweets)
			test.client.timelineClient.(*timelineClientMock).AssertNumberOfCalls(t, "UserTimeline", test.calls)
		})
	}
}

func Test_PostTweet(t *testing.T) {
	tests := map[string]struct {
		client *TwitterClient
		in     tweettypes.Tweet
		calls  int
		err    error
	}{
		"post tweet success": {
			client: NewTwitterClient(
				newStatusClientMock(
					"Update",
					[]interface{}{
						"Hello World",
						&twitter.StatusUpdateParams{
							Status:            "Hello World",
							InReplyToStatusID: 0,
						},
					},
					nil,
					nil,
				),
				nil,
			),
			in: tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: "Hello World",
				},
			},
			calls: 1,
			err:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.client.PostTweet(&test.in)
			if err != nil {
				assert.EqualError(t, test.err, err.Error())
			}
			test.client.statusClient.(*statusClientMock).AssertNumberOfCalls(t, "Update", test.calls)
		})
	}
}

func newStatusClientMock(method string, args []interface{}, ret interface{}, err error) *statusClientMock {
	client := new(statusClientMock)
	client.On(method, args...).Return(ret, err)
	return client
}

type statusClientMock struct {
	mock.Mock
}

func (mock *statusClientMock) Update(status string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error) {
	args := mock.Called(status, params)
	if args.Get(0) == nil {
		return nil, nil, args.Error(1)
	}
	return args.Get(0).(*twitter.Tweet), nil, args.Error(1)
}

func newTimelineClientMock(method string, args []interface{}, ret interface{}, err error) *timelineClientMock {
	client := new(timelineClientMock)
	client.On(method, args...).Return(ret, err)
	return client
}

type timelineClientMock struct {
	mock.Mock
}

func (mock *timelineClientMock) UserTimeline(params *twitter.UserTimelineParams) ([]twitter.Tweet, *http.Response, error) {
	args := mock.Called(params)
	if args.Get(0) == nil {
		return nil, nil, args.Error(1)
	}
	return args.Get(0).([]twitter.Tweet), nil, args.Error(1)
}
