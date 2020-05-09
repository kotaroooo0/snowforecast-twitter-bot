package twitter

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

type ITwitterApiClient interface {
	PostTweet(string, url.Values) (anaconda.Tweet, error)
}

type TwitterApiClient *anaconda.TwitterApi

func NewTwitterApiClient() ITwitterApiClient {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))
	return api
}

// TODO: バッチ処理の部分もインターフェースで置き換えたら消す
func GetTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))
	return api
}
