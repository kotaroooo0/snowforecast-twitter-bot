package main

import (
	"fmt"

	"github.com/kotaroooo0/snowforecast-twitter-bot/key"
	"github.com/kotaroooo0/snowforecast-twitter-bot/text"
)

func main() {
	api := key.GetTwitterApi()
	text := text.TweetContent()

	tweet, err := api.PostTweet(text, nil)
	if err != nil {
		panic(err)
	}
	fmt.Print(tweet.Text)

	fmt.Println("Finished")
}
