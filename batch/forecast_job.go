package batch

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bamzi/jobrunner"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/text"
	"github.com/robfig/cron/v3"
)

func TweetForecastJobSchedule(spec string, job cron.Job) error {
	return jobrunner.Schedule(spec, job)
}

type TweetForecast struct {
	Api        *anaconda.TwitterApi
	SkiResort1 string
	SkiResort2 string
}

func (t TweetForecast) Run() {
	text := text.TweetContent(t.SkiResort1, t.SkiResort2)
	tweet, err := t.Api.PostTweet(text, nil)
	if err != nil {
		panic(err)
	}
	log.Println(tweet)
}
