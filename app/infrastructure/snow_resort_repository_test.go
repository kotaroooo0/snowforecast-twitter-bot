package repository

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
)

func NewMockRedis(t *testing.T) *redis.Client {
	t.Helper()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("unexpected error while createing test redis server '%#v'", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return client
}

func TestListSnowResorts(t *testing.T) {
	client := NewMockRedis(t)
	s := SnowResortRepositoryImpl{
		Client: client,
	}

	client.SAdd("lowercase-snowresorts-serchword", "TashiroKaguraMitsumata", "Akakura-Kumado", "Hakuba47")
	actual, err := s.ListSnowResorts("lowercase-snowresorts-serchword")
	if err != nil {
		t.Fatalf("unexpected error while ListSnowResorts '%#v'", err)
	}

	// sliceではcmp.Diffで順序が考慮されてしまうのでSetに変換して比較する
	expectedSet := make(map[string]struct{})
	for _, v := range []string{"TashiroKaguraMitsumata", "Akakura-Kumado", "Hakuba47"} {
		expectedSet[v] = struct{}{}
	}
	actualSet := make(map[string]struct{})
	for _, v := range actual {
		actualSet[v] = struct{}{}
	}
	if diff := cmp.Diff(actualSet, expectedSet); diff != "" {
		t.Errorf("Diff: (-got +want)\n%s", diff)
	}
}
