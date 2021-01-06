package batch

import (
	"fmt"
	"log"

	"github.com/bamzi/jobrunner"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/text"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
)

type Pair struct {
	First  string `yaml:"first"`
	Second string `yaml:"second"`
}

func TweetForecastRun(api twitter.IApiClient, pairs []Pair) error {
	jobrunner.Start()
	for i, p := range pairs {
		if p.First == "" || p.Second == "" {
			return fmt.Errorf("error: two elements are needed")
		}
		if err := jobrunner.Schedule(fmt.Sprintf("00 %02d * * *", i), TweetForecast{api, p}); err != nil {
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
	text, err := text.TweetContent(t.SnowResorts)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := t.Api.PostTweet(text, nil); err != nil {
		log.Fatal(err)
	}
}
