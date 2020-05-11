package main

import (
	"log"

	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"

	repository "github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/batch"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setupBatch() {
	api := twitter.GetTwitterApi()
	jobrunner.Start()
	jobrunner.Schedule("00 01 * * *", batch.TweetForecast{api, "Hakuba47", "TakasuSnowPark"})
	jobrunner.Schedule("20 01 * * *", batch.TweetForecast{api, "MarunumaKogen", "TashiroKaguraMitsumata"})
}

func setupRouter() *gin.Engine {
	twitterApiClient := twitter.NewTwitterApiClient()
	yahooApiClient := yahoo.NewYahooApiClient()
	snowforecastApiClient := snowforecast.NewSnowforecastApiClient()
	redisClient, err := repository.New("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	snowResortRepository := repository.SnowResortRepositoryImpl{Client: redisClient}
	snowResortService := domain.SnowResortServiceImpl{TwitterApiClient: twitterApiClient, SnowforecastApiClient: snowforecastApiClient}
	twitterUseCase := usecase.TwitterUseCaseImpl{SnowResortService: snowResortService, SnowResortRepository: snowResortRepository, YahooApiClient: yahooApiClient}
	twitterHandler := handler.TwitterHandlerImpl{TwitterUseCase: twitterUseCase}
	jobHandler := handler.JobHandlerImpl{}

	r := gin.Default()
	r.GET("/twitter_webhook", twitterHandler.HandleTwitterGetCrcToken)
	r.POST("/twitter_webhook", twitterHandler.HandleTwitterPostWebhook)
	r.GET("/job_status", jobHandler.HandleGetJobStatus)
	return r
}

func main() {
	envLoad()

	// season outしたためストップ
	// setupBatch()

	r := setupRouter()
	r.Run(":3000")
}
