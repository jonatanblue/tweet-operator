package twitterclient

import (
	"errors"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	tweettypes "github.com/jonatanblue/tweet-operator/pkg/types"
)

type StatusClient interface {
	Update(status string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error)
}

type TimelineClient interface {
	UserTimeline(params *twitter.UserTimelineParams) ([]twitter.Tweet, *http.Response, error)
}

type TwitterClient struct {
	statusClient   StatusClient
	timelineClient TimelineClient
}

func NewTwitterClient(
	statusClient StatusClient,
	timelineClient TimelineClient,
) *TwitterClient {
	return &TwitterClient{
		statusClient:   statusClient,
		timelineClient: timelineClient,
	}
}

func (c *TwitterClient) GetTweetsForUser(userName string) (*tweettypes.Tweets, error) {
	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      200,
	}
	tweets, _, err := c.timelineClient.UserTimeline(params)
	if err != nil {
		return nil, err
	}
	if len(tweets) == 0 {
		return nil, err
	}
	var result tweettypes.Tweets
	for _, tweet := range tweets {
		result = append(
			result,
			tweettypes.Tweet{
				Spec: tweettypes.TweetSpec{
					Text: tweet.Text,
				},
				Status: tweettypes.TweetStatus{
					ID:       tweet.ID,
					Likes:    int64(tweet.FavoriteCount),
					Retweets: int64(tweet.RetweetCount),
					Replies:  int64(tweet.ReplyCount),
				},
			},
		)
	}
	return &result, nil
}

func (c *TwitterClient) PostTweet(tweet *tweettypes.Tweet) error {
	_, _, err := c.statusClient.Update(
		tweet.Spec.Text,
		&twitter.StatusUpdateParams{
			Status: tweet.Spec.Text,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *TwitterClient) DeleteTweet(tweet *tweettypes.Tweet) error {
	return errors.New("todo")
}

func NewTwitterAPIClient(creds *Credentials) (*twitter.Client, error) {
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

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}
