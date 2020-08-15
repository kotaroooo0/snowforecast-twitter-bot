package searcher

import (
	"fmt"
	"testing"
)

// TODO
// 今はただの動作確認

func TestFindAll(t *testing.T) {
	dbConfig := &DBConfig{
		User:     "root",
		Password: "password",
		Addr:     "127.0.0.1",
		Port:     "3306",
		DB:       "snowforecast_twitter_bot",
	}
	dbClient, err := NewDBClient(dbConfig)
	if err != nil {
		t.Fatal(err)
	}
	if err = dbClient.Ping(); err != nil {
		t.Fatal(err)
	}
	r := NewSnowResortRepositoryImpl(dbClient)
	sr, err := r.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sr)
}
