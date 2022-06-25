package types

type Tweet struct {
	Spec   TweetSpec
	Status TweetStatus
}

type TweetSpec struct {
	Text string
}

type TweetStatus struct {
	ID       int64
	Likes    int64
	Retweets int64
	Replies  int64
}

type Tweets []Tweet
