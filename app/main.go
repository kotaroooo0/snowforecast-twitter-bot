package main

import (
	"log"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/controllers"
	"github.com/kotaroooo0/snowforecast-twitter-bot/text"
)

func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	EnvLoad()
	// api := key.GetTwitterApi()
	jobrunner.Start()
	// jobrunner.Schedule("00 01 * * *", TweetForecast{api, "Hakuba47", "TakasuSnowPark"})
	// jobrunner.Schedule("20 01 * * *", TweetForecast{api, "MarunumaKogen", "TashiroKaguraMitsumata"})
	r := gin.Default()
	r.GET("/jobrunner/status", JobJSON)

	twitter := controllers.NewTwitterController(r)
	twitter.GetCrcToken()
	twitter.PostWebhook()

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
