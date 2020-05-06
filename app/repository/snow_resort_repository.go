package repository

import (
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
)

type SnowResortRepository interface {
	ListSnowResorts(offset, limit int) ([]string, error)
}

type SnowResortRepositoryImpl struct {
	Client *redis.Client
}

func New(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

func (s *SnowResortRepositoryImpl) ListSnowResorts(key string) ([]string, error) {
	result, err := s.Client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SnowResortRepositoryImpl) FindSnowResort(key string) (SnowResort, error) {
	result, err := s.Client.HGetAll(key).Result()
	if err != nil {
		return SnowResort{}, err
	}

	return SnowResort{SearchWord: result["search_word"], Label: result["label"]}, nil
}

type SnowResort struct {
	SearchWord string
	Label      string
}
