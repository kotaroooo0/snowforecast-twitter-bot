package main

import (
	"log"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/kotaroooo0/snowforecast-twitter-bot/key"
	"github.com/kotaroooo0/snowforecast-twitter-bot/text"
)

func main() {
	api := key.GetTwitterApi()

	jobrunner.Start()
	jobrunner.Schedule("00 01 * * *", TweetForecast{api, "Hakuba47", "MyokoSuginohara"})
	jobrunner.Schedule("20 01 * * *", TweetForecast{api, "MarunumaKogen", "IshiuchiMaruyama"})
	jobrunner.Schedule("40 01 * * *", TweetForecast{api, "TakasuSnowPark", "BiwakoValley"})
	// jobrunner.Schedule("06 01 * * *", TweetForecast{api, "Niseko", "SapporoKokusai"})

	jobrunner.Schedule("30 12 * * *", TweetForecast{api, "Hakuba47", "MyokoSuginohara"})
	jobrunner.Schedule("00 13 * * *", TweetForecast{api, "MarunumaKogen", "IshiuchiMaruyama"})
	jobrunner.Schedule("30 13 * * *", TweetForecast{api, "TakasuSnowPark", "BiwakoValley"})
	// jobrunner.Schedule("36 12 * * *", TweetForecast{api, "Niseko", "SapporoKokusai"})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/jobrunner/status", JobJSON)
	r.Run(":3000")
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

func JobJSON(c *gin.Context) {
	c.JSON(http.StatusOK, jobrunner.StatusJson())
}
