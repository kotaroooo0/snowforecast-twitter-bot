package batch

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/kotaroooo0/snowforecast-twitter-bot/apiclient/snowforecast"
)

// TODO: 動作確認しかしておらずちゃんとテストを書く
func TestTweetContent(t *testing.T) {
	apiClient := snowforecast.NewApiClient()
	c := NewTweetContentCreater(apiClient)
	pair1 := NewPair("UrabandaiNekoma", "GetouKogen")
	pair2 := NewPair("Hakuba47", "MyokoSuginohara")
	pair3 := NewPair("TashiroKaguraMitsumata", "MarunumaKogen")
	pair4 := NewPair("TakasuSnowPark", "HachiKogen")

	pp.Println(c.TweetContent(pair1))
	pp.Println(c.TweetContent(pair2))
	pp.Println(c.TweetContent(pair3))
	pp.Println(c.TweetContent(pair4))
}
