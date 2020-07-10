//+ wireinject
package main

import (
	"github.com/go-redis/redis"
	"github.com/google/wire"
	repository "github.com/kotaroooo0/snowforecast-twitter-bot/infrastructure"
)

func initSnowResortRepositoryImpl(client *redis.Client) *repository.SnowResortRepositoryImpl {
	wire.Build(
		repository.NewSnowResortRepositoryImpl,
		repository.NewRedisClient,
	)
	return nil
}