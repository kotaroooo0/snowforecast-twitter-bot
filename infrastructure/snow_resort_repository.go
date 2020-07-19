package repository

import (
	"github.com/go-redis/redis/v7"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
)

func NewSnowResortRepositoryImpl(client *redis.Client) domain.SnowResortRepository {
	return &SnowResortRepositoryImpl{
		Client: client,
	}
}

// TODO: DomainModelを返すように修正
func (s SnowResortRepositoryImpl) ListSnowResorts(key string) ([]string, error) {
	result, err := s.Client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
