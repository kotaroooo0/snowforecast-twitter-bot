// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kotaroooo0/snowforecast-twitter-bot/cache"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	repository "github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

func initNewTwitterHandlerImpl(tc *twitter.TwitterConfig, yc *yahoo.YahooConfig, rc *cache.RedisConfig) (handler.ReplyHandler, error) {
	wire.Build(
		twitter.NewTwitterApiClient,
		yahoo.NewYahooApiClient,
		snowforecast.NewSnowforecastApiClient,
		cache.NewRedisClient,
		repository.NewSnowResortRepositoryImpl,
		domain.NewReplyServiceImpl,
		usecase.NewReplyUseCaseImpl,
		handler.NewReplyHandlerImpl,
	)
	return &handler.ReplyHandlerImpl{}, nil
}

func initNewJobHandlerImpl() handler.JobHandler {
	wire.Build(
		handler.NewJobHandlerImpl,
	)
	return &handler.JobHandlerImpl{}
}
