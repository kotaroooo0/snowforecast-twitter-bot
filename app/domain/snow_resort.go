package domain

import (
	"fmt"
	"net/url"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
)

type SnowResort struct {
	SearchWord string
	Label      string
}

type Tweet struct {
	ID             string
	UserScreenName string
	Text           string
}

type SnowResortService interface {
	ReplyForecast(SnowResort, Tweet) (SnowResort, error)
}

type SnowResortServiceImpl struct {
	// ドメイン層はどの層にも依存しない
	TwitterApiClient twitter.ITwitterApiClient
}

func (ss SnowResortServiceImpl) ReplyForecast(snowResort SnowResort, tweet Tweet) (SnowResort, error) {
	params := url.Values{}
	params.Set("in_reply_to_status_id", tweet.ID)
	_, err := ss.TwitterApiClient.PostTweet(fmt.Sprintf("@%s %s", tweet.UserScreenName, "content"), params)
	if err != nil {
		return SnowResort{}, nil
	}
	return SnowResort{}, nil
}

type SnowResortRepository interface {
	ListSnowResorts(string) ([]string, error)
	FindSnowResort(string) (SnowResort, error)
}
