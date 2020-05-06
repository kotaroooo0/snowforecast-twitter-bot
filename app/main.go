package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
)

func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	EnvLoad()
	// season outしたためストップ
	// batch.Start()

	// userRepository := repository.NewUserPersistence()
	// userUseCase := usecase.NewUserUseCase()
	twitterHandler := handler.TwitterHandlerImpl{}
	jobHandler := handler.JobHandlerImpl{}

	r := gin.Default()
	r.GET("/twitter_webhook", twitterHandler.HandleTwitterGetCrcToken)
	r.POST("/twitter_webhook", twitterHandler.HandleTwitterPostWebhook)
	r.GET("/job_status", jobHandler.HandleGetJobStatus)

	r.Run(":3000")
}
