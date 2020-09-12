package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type IApiClient interface {
	PostTweet(string, url.Values) (anaconda.Tweet, error)
}

type ApiClient *anaconda.TwitterApi

func NewApiClient(c *Config) IApiClient {
	anaconda.SetConsumerKey(c.ConsumerKey)
	anaconda.SetConsumerSecret(c.ConsumerSecret)
	return anaconda.NewTwitterApi(c.AccessTokenKey, c.AccessTokenSecret)
}

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenKey    string
	AccessTokenSecret string
}

func NewConfig(consumerKey string, consumerSecret string, accessTokenKey string, accessTokenSecret string) *Config {
	return &Config{
		ConsumerKey:       consumerKey,
		ConsumerSecret:    consumerSecret,
		AccessTokenKey:    accessTokenKey,
		AccessTokenSecret: accessTokenSecret,
	}
}
