// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	repository "github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

func initNewTwitterHandlerImpl(tc *twitter.TwitterConfig, yc *yahoo.YahooConfig, rc *repository.RedisConfig) (handler.TwitterHandler, error) {
	wire.Build(
		yahoo.NewYahooApiClient,
		twitter.NewTwitterApiClient,
		snowforecast.NewSnowforecastApiClient,
		repository.NewRedisClient,
		repository.NewSnowResortRepositoryImpl,
		domain.NewSnowResortServiceImpl,
		usecase.NewTwitterUseCaseImpl,
		handler.NewTwitterHandlerImpl,
	)
	return &handler.TwitterHandlerImpl{}, nil
}

func initNewJobHandlerImpl() handler.JobHandler {
	wire.Build(
		handler.NewJobHandlerImpl,
	)
	return &handler.JobHandlerImpl{}
}
