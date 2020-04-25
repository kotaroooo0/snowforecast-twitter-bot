package yahoo

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/kotaroooo0/snowforecast-twitter-bot/parameters/responses"
)

type IYahooApiClient interface {
	GetMorphologicalAnalysis(str string) (responses.PostKanjiToHiraganaResponse, error)
}

type YahooApiClient struct {
	YahooAppID string
}

func (y *YahooApiClient) GetMorphologicalAnalysis(str string) (responses.PostKanjiToHiraganaResponse, error) {
	params := url.Values{}
	params.Set("appid", y.YahooAppID)
	params.Set("sentence", str)
	params.Set("result", "hiragana")

	r := responses.NewPostKanjiToHiraganaResponse()
	res, err := http.Get("https://jlp.yahooapis.jp/MAService/V1/parse?" + params.Encode())
	if err != nil {
		return r, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	err = xml.Unmarshal(body, &r)
	return r, err
}

func NewYahooApiClient() IYahooApiClient {
	return &YahooApiClient{
		YahooAppID: os.Getenv("YAHOO_APP_ID"),
	}
}
