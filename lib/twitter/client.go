package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

type ITwitterApiClient interface {
	PostTweet(string, url.Values) (anaconda.Tweet, error)
}

type TwitterApiClient *anaconda.TwitterApi

func NewTwitterApiClient(twitterConfig *TwitterConfig) ITwitterApiClient {
	anaconda.SetConsumerKey(twitterConfig.ConsumerKey)
	anaconda.SetConsumerSecret(twitterConfig.ConsumerSecret)
	return anaconda.NewTwitterApi(twitterConfig.AccessTokenKey, twitterConfig.AccessTokenSecret)
}

type TwitterConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenKey    string
	AccessTokenSecret string
}

func NewTwitterConfig(consumerKey string, consumerSecret string, accessTokenKey string, accessTokenSecret string) *TwitterConfig {
	return &TwitterConfig{
		ConsumerKey:       consumerKey,
		ConsumerSecret:    consumerSecret,
		AccessTokenKey:    accessTokenKey,
		AccessTokenSecret: accessTokenSecret,
	}
}
