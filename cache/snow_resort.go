package cache

import (
	"os"
	"reflect"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/pkg/errors"
)

type SnowResortCache interface {
	Get(string) (*domain.SnowResort, error)
	Set(string, *domain.SnowResort) error
}

type SnowResortCacheImpl struct {
	Client *redis.Client
}

func NewSnowResortCacheImpl(c *redis.Client) SnowResortCache {
	return &SnowResortCacheImpl{
		Client: c,
	}
}

func NewRedisClient(redisConfig *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr + ":" + os.Getenv("REDIS_PORT"),
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

type RedisConfig struct {
	Addr string
}

func NewRedisConfig(addr string) *RedisConfig {
	return &RedisConfig{
		Addr: addr,
	}
}

func (s SnowResortCacheImpl) Get(key string) (*domain.SnowResort, error) {
	sr := &domain.SnowResort{}
	result, err := s.Client.HGetAll(key).Result()
	if err != nil {
		return sr, err
	}

	// map[string]string -> map[string]interface{}へ型変換
	m := make(map[string]interface{}, len(result))
	for k, v := range result {
		// TODO: このロジックはリフレクションに持っていけないか
		i, err := strconv.Atoi(v)
		if err != nil {
			m[k] = v
		} else {
			m[k] = i
		}
	}

	// 構造体へ変換
	toStruct(m, sr)
	return sr, nil
}

func (s SnowResortCacheImpl) Set(key string, snowResort *domain.SnowResort) error {
	m := toMap(snowResort)
	err := s.Client.HMSet(key, m)
	return err.Err()
}

func toMap(v interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		ftv := rt.Field(i)
		fv := rv.Field(i)
		if rv.CanSet() {
			m[ftv.Name] = fv.Interface()
		}
	}
	return m
}

func toStruct(m map[string]interface{}, s interface{}) {
	rv := reflect.ValueOf(s).Elem()
	for k, v := range m {
		rv.FieldByName(k).Set(reflect.ValueOf(v))
	}
}
