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

func (s *SnowResortRepositoryImpl) ListSnowResorts() ([]string, error) {
	key := "snowresorts-serchword"
	v := s.Client.SMembers(key)
	if v.Err() != nil {
		return nil, v.Err()
	}

	return v.Val(), nil
}
