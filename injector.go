//+ wireinject
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

func initNewTwitterHandlerImpl(addr string) (handler.TwitterHandler, error) {
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
