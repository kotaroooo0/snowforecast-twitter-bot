package main

import (
	"log"
	"os"

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
	"go.uber.org/dig"
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
	r := gin.Default()
	c := dig.New()

	c.Provide(NewTwitterApiClient)
	c.Provide(NewYahooApiClient)
	c.Provide(NewSnowForecastApiClient)

	c.Invoke(func(twitterHandler handler.TwitterHandlerImpl) {
		r.GET("/twitter_webhook", twitterHandler.HandleTwitterGetCrcToken)
		r.POST("/twitter_webhook", twitterHandler.HandleTwitterPostWebhook)
	})
	c.Invoke(func(jobHandler handler.JobHandlerImpl) {
		r.GET("/job_status", jobHandler.HandleGetJobStatus)
	})

	twitterApiClient := twitter.NewTwitterApiClient()
	yahooApiClient := yahoo.NewYahooApiClient()
	snowforecastApiClient := snowforecast.NewSnowforecastApiClient()
	redisClient, err := repository.New(os.Getenv("REDIS_HOST") + ":6379")
	if err != nil {
		log.Fatal(err)
	}
	snowResortRepository := repository.SnowResortRepositoryImpl{Client: redisClient}
	snowResortService := domain.SnowResortServiceImpl{SnowResortRepository: snowResortRepository, TwitterApiClient: twitterApiClient, SnowforecastApiClient: snowforecastApiClient, YahooApiClient: yahooApiClient}
	twitterUsecase := usecase.TwitterUsecaseImpl{SnowResortService: snowResortService}
	twitterHandler := handler.TwitterHandlerImpl{TwitterUsecase: twitterUsecase}
	jobHandler := handler.JobHandlerImpl{}

	return r
}

func main() {
	envLoad()
	setupBatch()

	r := setupRouter()
	r.Run(":3000")
}
