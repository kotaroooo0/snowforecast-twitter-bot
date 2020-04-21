package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/controllers"
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
	r := gin.Default()
	job := controllers.NewJobController(r)
	job.GetJobStatus()
	twitter := controllers.NewTwitterController(r)
	twitter.GetCrcToken()
	twitter.PostWebhook()
	r.Run(":3000")
}
