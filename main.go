package main

import (
	"log"
	"time"

	"github.com/kotaroooo0/snowforecast-twitter-bot/text"
)

func main() {
	scheduledTweet()

	// TODO: リプライ待ちの動作に変更
	time.Sleep(100 * time.Second)
}

func scheduledTweet() {

	// TODO: 開発中はAPIを叩かない
	// api := key.GetTwitterApi()

	go func() {

		// TODO: 頻度を変更
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				text := text.TweetContent()
				log.Println(text)

				// TODO: 開発中はAPIを叩かない
				// tweet, err := api.PostTweet(text, nil)
				// if err != nil {
				// 	panic(err)
				// }
				// fmt.Print(tweet.Text)
			}
		}
	}()
}
