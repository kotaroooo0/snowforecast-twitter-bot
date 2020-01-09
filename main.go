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
	jobrunner.Schedule("03 01 * * *", TweetForecast{api, "IshiuchiMaruyama", "TakasuSnowPark"})

	jobrunner.Schedule("00 12 * * *", TweetForecast{api, "Hakuba47", "MyokoSuginohara"})
	jobrunner.Schedule("03 12 * * *", TweetForecast{api, "IshiuchiMaruyama", "TakasuSnowPark"})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/jobrunner/status", JobJSON)
	r.Run(":8080")
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
