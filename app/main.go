package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/key"
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

	fmt.Println("faslgkhaslgkuhkal;sgh")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/jobrunner/status", JobJSON)
	r.GET("/webhook/twitter", GetWebhookTwitter)
	r.POST("/webhook/twitter", PostWebhookTwitter)

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

// type CRCResponse struct {
// 	ResponseToken string `json:"response_token"`
// }

func GetWebhookTwitter(c *gin.Context) {
	fmt.Println(c.Query("crc_token"))
	c.JSON(http.StatusOK, gin.H{"response_token": key.CreateCRCToken(c.Query("crc_token"))})

	// responseToken := CRCResponse{ResponseToken: key.CreateCRCToken(c.Request.FormValue("crc_token"))}
	// c.JSON(200, gin.H{
	// 	"response_token": key.CreateCRCToken(c.Request.FormValue("crc_token")),
	// })
}

func PostWebhookTwitter(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Status(http.StatusOK)
}
