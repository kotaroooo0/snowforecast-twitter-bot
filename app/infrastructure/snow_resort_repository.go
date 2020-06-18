package repository

import (
	"github.com/go-redis/redis/v7"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/pkg/errors"
	"os"
)

type SnowResortRepositoryImpl struct {
	Client *redis.Client
}

func NewSnowResortRepository(client *redis.Client) domain.SnowResortRepository {
	return SnowResortRepositoryImpl{
		Client: client,
	}
}

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":6379",
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

// TODO: DomainModelを返すように修正
func (s SnowResortRepositoryImpl) ListSnowResorts(key string) ([]string, error) {
	result, err := s.Client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s SnowResortRepositoryImpl) FindSnowResort(key string) (domain.SnowResort, error) {
	result, err := s.Client.HGetAll(key).Result()
	if err != nil {
		return domain.SnowResort{}, err
	}

	return domain.SnowResort{SearchWord: result["search_word"], Label: result["label"]}, nil
}

func (s SnowResortRepositoryImpl) SetSnowResort(key string, snowResort domain.SnowResort) error {
	err := s.Client.HMSet(key, "search_word", snowResort.SearchWord, "label", snowResort.Label)
	return err.Err()
}
