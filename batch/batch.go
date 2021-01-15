package batch

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bamzi/jobrunner"
	"github.com/kotaroooo0/snowforecast-twitter-bot/apiclient/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/apiclient/twitter"
)

type Pair struct {
	First  string `yaml:"first"`
	Second string `yaml:"second"`
}

func NewPair(first, second string) Pair {
	return Pair{
		First:  first,
		Second: second,
	}
}

func TweetForecastRun(api twitter.IApiClient, pairs []Pair) error {
	jobrunner.Start()
	for i, p := range pairs {
		if p.First == "" || p.Second == "" {
			return fmt.Errorf("error: two elements are needed")
		}
		if err := jobrunner.Schedule(fmt.Sprintf("%02d 1 * * *", i*10%60), TweetForecast{api, p}); err != nil {
			return err
		}
	}
	return nil
}

type TweetForecast struct {
	Api         twitter.IApiClient
	SnowResorts Pair
}

func (t TweetForecast) Run() {
	apiClient := snowforecast.NewApiClient()
	tweetContentCreater := NewTweetContentCreater(apiClient)
	text, err := tweetContentCreater.TweetContent(t.SnowResorts)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := t.Api.PostTweet(text, nil); err != nil {
		log.Fatal(err)
	}
}

type TweetContentCreater struct {
	snowForecastApiClient snowforecast.IApiClient
}

func NewTweetContentCreater(snowForecastApiClient snowforecast.IApiClient) TweetContentCreater {
	return TweetContentCreater{
		snowForecastApiClient: snowForecastApiClient,
	}
}

func (c TweetContentCreater) TweetContent(pair Pair) (string, error) {
	firstData, err := c.snowForecastApiClient.GetForecastBySearchWord(pair.First)
	if err != nil {
		return "", err
	}
	secondData, err := c.snowForecastApiClient.GetForecastBySearchWord(pair.Second)
	if err != nil {
		return "", err
	}
	content := "今日 | 明日 | 明後日 (朝,昼,夜)\n"
	content += pair.First + "\n"
	content += areaLineString(firstData) + "\n"
	content += pair.Second + "\n"
	content += areaLineString(secondData) + "\n"
	return content, nil
}

func areaLineString(snowfallForecast *snowforecast.Forecast) string {
	content := strconv.Itoa(snowfallForecast.Snows[0].Morning) + addRainyChar(snowfallForecast.Rains[0].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Noon) + addRainyChar(snowfallForecast.Rains[0].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Night) + addRainyChar(snowfallForecast.Rains[0].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[1].Morning) + addRainyChar(snowfallForecast.Rains[1].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Noon) + addRainyChar(snowfallForecast.Rains[1].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Night) + addRainyChar(snowfallForecast.Rains[1].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[2].Morning) + addRainyChar(snowfallForecast.Rains[2].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Noon) + addRainyChar(snowfallForecast.Rains[2].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Night) + addRainyChar(snowfallForecast.Rains[2].Night) + "cm "
	return content
}

func addRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	}
	if rainfall > 0 {
		return "☂️"
	}
	return ""
}
