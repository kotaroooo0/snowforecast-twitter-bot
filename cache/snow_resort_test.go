package cache

import (
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/pkg/errors"
)

func before() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func NewSnowResortCacheTestImpl() SnowResortCache {
	before()
	c, err := NewRedisTestClient()
	if err != nil {
		panic(err)
	}
	return NewSnowResortCacheImpl(c)
}

func NewRedisTestClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   1, // 1のDBをテスト用とする
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

func TestGet(t *testing.T) {
	src := NewSnowResortCacheTestImpl()
	cases := []struct {
		k    string
		want *domain.SnowResort
	}{
		{
			k:    "白馬",
			want: &domain.SnowResort{},
		},
		{
			k:    "かぐら",
			want: &domain.SnowResort{},
		},
	}

	for _, tt := range cases {
		act, err := src.Get(tt.k)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(act, tt.want); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}

func TestSet(t *testing.T) {
	src := NewSnowResortCacheTestImpl()
	cases := []struct {
		k string
		s *domain.SnowResort
	}{
		{
			k: "白馬",
			s: &domain.SnowResort{
				Id:        0,
				Name:      "Hakuba47",
				SearchKey: "hakuba47",
				Elevation: 1500,
				Region:    "japan-nagano",
			},
		},
		{
			k: "かぐら",
			s: &domain.SnowResort{
				Id:        12345,
				Name:      "Kagura",
				SearchKey: "kagura",
				Elevation: 1600,
				Region:    "japan-nigata",
			},
		},
	}

	for _, tt := range cases {
		err := src.Set(tt.k, tt.s)
		if err != nil {
			t.Error(err)
		}
		act, err := src.Get(tt.k)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(act, tt.s); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}

type User struct {
	Id    int
	Name  string
	Adult bool
}

func TestToMap(t *testing.T) {
	cases := []struct {
		s    *User
		want map[string]interface{}
	}{
		{
			s: &User{
				Id:    0,
				Name:  "kotaroooo0",
				Adult: true,
			},
			want: map[string]interface{}{"Id": 0, "Name": "kotaroooo0", "Adult": true},
		},
		{
			s: &User{
				Id:    12345,
				Name:  "adachikun",
				Adult: false,
			},
			want: map[string]interface{}{"Id": 12345, "Name": "adachikun", "Adult": false},
		},
	}

	for _, tt := range cases {
		act := toMap(tt.s)
		if diff := cmp.Diff(act, tt.want); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}

func TestToStruct(t *testing.T) {
	cases := []struct {
		m    map[string]interface{}
		want *User
	}{
		{
			m: map[string]interface{}{"Id": 0, "Name": "kotaroooo0", "Adult": true},
			want: &User{
				Id:    0,
				Name:  "kotaroooo0",
				Adult: true,
			},
		},
		{
			m: map[string]interface{}{"Id": 12345, "Name": "adachikun", "Adult": false},
			want: &User{
				Id:    12345,
				Name:  "adachikun",
				Adult: false,
			},
		},
	}

	for _, tt := range cases {
		u := &User{}
		toStruct(tt.m, u)
		if diff := cmp.Diff(u, tt.want); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}
