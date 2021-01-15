package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/apiclient/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/apiclient/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/batch"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/elasticsearch"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
	"gopkg.in/yaml.v2"
)

// 環境変数をロード
func envLoad() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error: loading .env file")
	}
	return nil
}

// バッチで投稿する予報対象のスキー場をロード
func targetsLoad() ([]batch.Pair, error) {
	ps := []batch.Pair{}
	file, err := os.Open("batch.snow_resorts.yaml")
	if err != nil {
		return ps, fmt.Errorf("error: loading batch.snow_resorts.yaml file")
	}
	file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return ps, err
	}

	err = yaml.Unmarshal([]byte(body), &ps)
	return ps, err
}

func setupBatch(pairs []batch.Pair) error {
	if len(pairs) == 0 {
		// シーズン終わりなどバッチを動かさなくても良い時｀
		log.Println("Warning: No forecast targets ski resorts")
		return nil
	}
	api := twitter.NewApiClient(twitter.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET")))
	return batch.TweetForecastRun(api, pairs)
}

func setupRouter() (*gin.Engine, error) {
	tc := twitter.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))
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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := envLoad(); err != nil {
		return err
	}
	targets, err := targetsLoad()
	if err != nil {
		return err
	}
	if err := setupBatch(targets); err != nil {
		return err
	}
	r, err := setupRouter()
	if err != nil {
		return err
	}
	return r.Run(":3000")
}
