package domain

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"golang.org/x/exp/utf8string"

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
	TwitterApiClient      twitter.ITwitterApiClient
	SnowforecastApiClient snowforecast.ISnowforecastApiClient
}

func (ss SnowResortServiceImpl) ReplyForecast(snowResort SnowResort, tweet Tweet) (SnowResort, error) {
	params := url.Values{}
	params.Set("in_reply_to_status_id", tweet.ID)
	content, err := replyContent(snowResort, ss.SnowforecastApiClient)
	if err != nil {
		return SnowResort{}, err
	}
	_, err = ss.TwitterApiClient.PostTweet(fmt.Sprintf("@%s %s", tweet.UserScreenName, content), params)
	if err != nil {
		return SnowResort{}, err
	}
	return SnowResort{}, nil
}

func replyContent(snowResort SnowResort, snowforecastApiClient snowforecast.ISnowforecastApiClient) (string, error) {
	snowfallForecast, err := snowforecastApiClient.GetSnowfallForecastBySkiResortSearchWord(snowResort.SearchWord)
	if err != nil {
		return "", err
	}

	// TODO: 仮の文章
	content := snowResort.Label + "\n"
	content += "今日 | 明日 | 明後日\n"
	content += "3日後 | 4日後 | 5日後\n"
	content += strconv.Itoa(snowfallForecast.Snows[0].Morning) + addRainyChar(snowfallForecast.Rains[0].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Noon) + addRainyChar(snowfallForecast.Rains[0].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Night) + addRainyChar(snowfallForecast.Rains[0].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[1].Morning) + addRainyChar(snowfallForecast.Rains[1].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Noon) + addRainyChar(snowfallForecast.Rains[1].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Night) + addRainyChar(snowfallForecast.Rains[1].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[2].Morning) + addRainyChar(snowfallForecast.Rains[2].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Noon) + addRainyChar(snowfallForecast.Rains[2].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Night) + addRainyChar(snowfallForecast.Rains[2].Night) + "cm\n"
	content += strconv.Itoa(snowfallForecast.Snows[3].Morning) + addRainyChar(snowfallForecast.Rains[3].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[3].Noon) + addRainyChar(snowfallForecast.Rains[3].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[3].Night) + addRainyChar(snowfallForecast.Rains[3].Night) + "cm |"
	content += strconv.Itoa(snowfallForecast.Snows[4].Morning) + addRainyChar(snowfallForecast.Rains[4].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[4].Noon) + addRainyChar(snowfallForecast.Rains[4].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[4].Night) + addRainyChar(snowfallForecast.Rains[4].Night) + "cm |"
	content += strconv.Itoa(snowfallForecast.Snows[5].Morning) + addRainyChar(snowfallForecast.Rains[5].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[5].Noon) + addRainyChar(snowfallForecast.Rains[5].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[5].Night) + addRainyChar(snowfallForecast.Rains[5].Night) + "cm"

	// 140字までに切り詰めて返す
	if len([]rune(content)) > 140 {
		return utf8string.NewString(content).Slice(0, 140), nil
	}
	return content, nil
}

func addRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	} else if rainfall > 0 {
		return "☂️"
	}
	return ""
}

type SnowResortRepository interface {
	ListSnowResorts(string) ([]string, error)
	FindSnowResort(string) (SnowResort, error)
}
