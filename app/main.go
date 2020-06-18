package main

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	repository "github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
	"log"
)

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// season outしたためストップ
func setupBatch() {
	// api := twitter.GetTwitterApi()
	jobrunner.Start()
	// jobrunner.Schedule("00 01 * * *", batch.TweetForecast{api, "Hakuba47", "TakasuSnowPark"})
	// jobrunner.Schedule("20 01 * * *", batch.TweetForecast{api, "MarunumaKogen", "TashiroKaguraMitsumata"})
}

func setupRouter() *gin.Engine {
	twitterApiClient := twitter.NewTwitterApiClient()
	yahooApiClient := yahoo.NewYahooApiClient()
	snowforecastApiClient := snowforecast.NewISnowforecastApiClient()
	redisClient, err := repository.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}
	snowResortRepository := repository.NewSnowResortRepository(redisClient)
	snowResortService := domain.NewSnowResortService(snowResortRepository, yahooApiClient, twitterApiClient, snowforecastApiClient)
	twitterUsecase := usecase.NewTwitterUsecase(snowResortService)
	twitterHandler := handler.NewTwitterHandler(twitterUsecase)
	jobHandler := handler.NewJobHandler()

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
