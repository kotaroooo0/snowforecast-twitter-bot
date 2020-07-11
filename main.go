package main

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	//api := twitter.NewTwitterApiClient()
	jobrunner.Start()
	// jobrunner.Schedule("00 01 * * *", batch.TweetForecast{api, "Hakuba47", "TakasuSnowPark"})
	// jobrunner.Schedule("20 01 * * *", batch.TweetForecast{api, "MarunumaKogen", "TashiroKaguraMitsumata"})
}

func setupRouter() *gin.Engine {
	//twitterHandler,err := initNewTwitterHandlerImpl(os.Getenv("REDIS_HOST") + ":6379")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//jobHandler := handler.JobHandlerImpl{}

	r := gin.Default()
	//r.GET("/twitter_webhook", twitterHandler.HandleTwitterGetCrcToken)
	//r.POST("/twitter_webhook", twitterHandler.HandleTwitterPostWebhook)
	//r.GET("/job_status", jobHandler.HandleGetJobStatus)
	return r
}

func main() {
	envLoad()
	setupBatch()

	r := setupRouter()
	r.Run(":3000")
}
