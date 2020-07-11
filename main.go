package main

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	repository "github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"log"
	"os"
)

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// season outしたためストップ
func setupBatch() {
	//api := twitter.NewTwitterApiClient(twitter.NewTwitterConfig(os.Getenv("CONSUMER_KEY"),os.Getenv("CONSUMER_SECRET"),os.Getenv("ACCESS_TOKEN_KEY"),os.Getenv("ACCESS_TOKEN_SECRET")))
	jobrunner.Start()
	// jobrunner.Schedule("00 01 * * *", batch.TweetForecast{api, "Hakuba47", "TakasuSnowPark"})
	// jobrunner.Schedule("20 01 * * *", batch.TweetForecast{api, "MarunumaKogen", "TashiroKaguraMitsumata"})
}

func setupRouter() *gin.Engine {
	twitterHandler,err := initNewTwitterHandlerImpl(
		twitter.NewTwitterConfig(os.Getenv("CONSUMER_KEY"),os.Getenv("CONSUMER_SECRET"),os.Getenv("ACCESS_TOKEN_KEY"),os.Getenv("ACCESS_TOKEN_SECRET")),
		yahoo.NewYahooConfig(os.Getenv("YAHOO_APP_ID")),
		repository.NewRedisConfig(os.Getenv("REDIS_HOST")),
	)
	if err != nil {
		log.Fatal(err)
	}
	jobHandler := initNewJobHandlerImpl()

	r := gin.Default()
	r.GET("/twitter_webhook", twitterHandler.HandleTwitterGetCrcToken)
	r.POST("/twitter_webhook", twitterHandler.HandleTwitterPostWebhook)
	r.GET("/job_status", jobHandler.HandleGetJobStatus)
	return r
}

func main() {
	envLoad()
	setupBatch()

	r := setupRouter()
	r.Run(":3000")
}
