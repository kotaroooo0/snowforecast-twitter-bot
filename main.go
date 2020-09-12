package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/elasticsearch"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

func envLoad() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error: loading .env file")
	}
	return nil
}

// season outしたためストップ
func setupBatch() error {
	//api := twitter.NewApiClient(twitter.NewTwitterConfig(os.Getenv("CONSUMER_KEY"),os.Getenv("CONSUMER_SECRET"),os.Getenv("ACCESS_TOKEN_KEY"),os.Getenv("ACCESS_TOKEN_SECRET")))
	jobrunner.Start()
	// jobrunner.Schedule("00 01 * * *", batch.TweetForecast{api, "Hakuba47", "TakasuSnowPark"})
	// jobrunner.Schedule("20 01 * * *", batch.TweetForecast{api, "MarunumaKogen", "TashiroKaguraMitsumata"})
	return nil
}

func setupRouter() (*gin.Engine, error) {
	tc := twitter.NewTwitterConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))
	tac := twitter.NewApiClient(tc)
	sac := snowforecast.NewApiClient()
	ei, err := elasticsearch.NewSnowResortSearcherEsImpl()
	if err != nil {
		return nil, err
	}
	rs := domain.NewReplyServiceImpl(ei, tac, sac)
	ru := usecase.NewReplyUseCaseImpl(rs)
	rh := handler.NewReplyHandlerImpl(ru)
	jh := handler.NewJobHandlerImpl()
	r := gin.Default()
	r.GET("/twitter_webhook", rh.HandleTwitterGetCrcToken)
	r.POST("/twitter_webhook", rh.HandleTwitterPostWebhook)
	r.GET("/job_status", jh.HandleGetJobStatus)
	return r, nil
}

func run() error {
	if err := envLoad(); err != nil {
		return err
	}
	if err := setupBatch(); err != nil {
		return err
	}
	r, err := setupRouter()
	if err != nil {
		return err
	}
	return r.Run(":3000")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
