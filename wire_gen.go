// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/handler"
	"github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

// Injectors from injector.go:

func initNewTwitterHandlerImpl(tc *twitter.TwitterConfig, yc *yahoo.YahooConfig, rc *repository.RedisConfig) (handler.TwitterHandler, error) {
	client, err := repository.NewRedisClient(rc)
	if err != nil {
		return nil, err
	}
	snowResortRepository := repository.NewSnowResortRepositoryImpl(client)
	iYahooApiClient := yahoo.NewYahooApiClient(yc)
	iTwitterApiClient := twitter.NewTwitterApiClient(tc)
	iSnowforecastApiClient := snowforecast.NewSnowforecastApiClient()
	snowResortService := domain.NewSnowResortServiceImpl(snowResortRepository, iYahooApiClient, iTwitterApiClient, iSnowforecastApiClient)
	twitterUseCase := usecase.NewTwitterUseCaseImpl(snowResortService, iYahooApiClient)
	twitterHandler := handler.NewTwitterHandlerImpl(twitterUseCase)
	return twitterHandler, nil
}

func initNewJobHandlerImpl() handler.JobHandler {
	jobHandler := handler.NewJobHandlerImpl()
	return jobHandler
}
