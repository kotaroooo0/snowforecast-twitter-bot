package domain

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"golang.org/x/exp/utf8string"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
)

type SnowResort struct {
	Name      string `db:"name" json:"name"`
	SearchKey string `db:"search_key" json:"search_key"`
}

type SnowResortSearcher interface {
	FindSimilarSnowResort(string) (*SnowResort, error)
}

type SnowResortRepository interface {
	FindAll() ([]*SnowResort, error)
}

type Tweet struct {
	ID             string
	UserScreenName string
	Text           string
}

type ReplyService interface {
	ReplyForecast(*Tweet) (*SnowResort, error)
}

type ReplyServiceImpl struct {
	// ドメイン層は他の層にも依存しない
	SnowResortSearcher    SnowResortSearcher
	YahooApiClient        yahoo.IYahooApiClient
	TwitterApiClient      twitter.ITwitterApiClient
	SnowforecastApiClient snowforecast.ISnowforecastApiClient
}

func NewReplyServiceImpl(snowResortSearcher SnowResortSearcher, yahooApiClient yahoo.IYahooApiClient, twitterApiClient twitter.ITwitterApiClient, snowforecastApiClient snowforecast.ISnowforecastApiClient) ReplyService {
	return &ReplyServiceImpl{
		SnowResortSearcher:    snowResortSearcher,
		YahooApiClient:        yahooApiClient,
		TwitterApiClient:      twitterApiClient,
		SnowforecastApiClient: snowforecastApiClient,
	}
}

// 検索すること自体がビジネスロジックであるため、ドメイン層に含める
func (r ReplyServiceImpl) ReplyForecast(tweet *Tweet) (*SnowResort, error) {
	// リプライを取得
	replyText := tweet.Text
	// @snowfall_botを消す
	replyText = strings.Replace(replyText, "@snowfall_bot ", "", -1)

	sr, err := r.SnowResortSearcher.FindSimilarSnowResort(replyText)
	if err != nil {
		return &SnowResort{}, err
	}

	sf, err := r.SnowforecastApiClient.GetSnowfallForecastBySkiResortSearchWord(sr.SearchKey)
	if err != nil {
		return &SnowResort{}, err
	}

	content, err := replyContent(sr.Name, sf)
	if err != nil {
		return &SnowResort{}, err
	}

	params := url.Values{}
	params.Set("in_reply_to_status_id", tweet.ID)
	_, err = r.TwitterApiClient.PostTweet(fmt.Sprintf("@%s %s", tweet.UserScreenName, content), params)
	if err != nil {
		return &SnowResort{}, err
	}
	return sr, nil
}

func replyContent(name string, sf snowforecast.SnowfallForecast) (string, error) {
	// TODO: 仮の文章
	content := name + "\n"
	content += "今日 | 明日 | 明後日\n"
	content += "3日後 | 4日後 | 5日後\n"
	content += strconv.Itoa(sf.Snows[0].Morning) + addRainyChar(sf.Rains[0].Morning) + ", " + strconv.Itoa(sf.Snows[0].Noon) + addRainyChar(sf.Rains[0].Noon) + ", " + strconv.Itoa(sf.Snows[0].Night) + addRainyChar(sf.Rains[0].Night) + "cm | "
	content += strconv.Itoa(sf.Snows[1].Morning) + addRainyChar(sf.Rains[1].Morning) + ", " + strconv.Itoa(sf.Snows[1].Noon) + addRainyChar(sf.Rains[1].Noon) + ", " + strconv.Itoa(sf.Snows[1].Night) + addRainyChar(sf.Rains[1].Night) + "cm | "
	content += strconv.Itoa(sf.Snows[2].Morning) + addRainyChar(sf.Rains[2].Morning) + ", " + strconv.Itoa(sf.Snows[2].Noon) + addRainyChar(sf.Rains[2].Noon) + ", " + strconv.Itoa(sf.Snows[2].Night) + addRainyChar(sf.Rains[2].Night) + "cm\n"
	content += strconv.Itoa(sf.Snows[3].Morning) + addRainyChar(sf.Rains[3].Morning) + ", " + strconv.Itoa(sf.Snows[3].Noon) + addRainyChar(sf.Rains[3].Noon) + ", " + strconv.Itoa(sf.Snows[3].Night) + addRainyChar(sf.Rains[3].Night) + "cm |"
	content += strconv.Itoa(sf.Snows[4].Morning) + addRainyChar(sf.Rains[4].Morning) + ", " + strconv.Itoa(sf.Snows[4].Noon) + addRainyChar(sf.Rains[4].Noon) + ", " + strconv.Itoa(sf.Snows[4].Night) + addRainyChar(sf.Rains[4].Night) + "cm |"
	content += strconv.Itoa(sf.Snows[5].Morning) + addRainyChar(sf.Rains[5].Morning) + ", " + strconv.Itoa(sf.Snows[5].Noon) + addRainyChar(sf.Rains[5].Noon) + ", " + strconv.Itoa(sf.Snows[5].Night) + addRainyChar(sf.Rains[5].Night) + "cm"

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
