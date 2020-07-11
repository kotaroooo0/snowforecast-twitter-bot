package yahoo

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type IYahooApiClient interface {
	GetMorphologicalAnalysis(str string) (GetMorphologicalAnalysisResponse, error)
}

type YahooApiClient struct {
	YahooAppID string
}

func NewYahooApiClient(yahooConfig *YahooConfig) IYahooApiClient {
	return &YahooApiClient{
		YahooAppID: yahooConfig.YahooAppKey,
	}
}

type YahooConfig struct {
	YahooAppKey string
}

func NewYahooConfig(yahooApiKey string) *YahooConfig {
	return &YahooConfig{
		YahooAppKey: yahooApiKey,
	}
}

func (y *YahooApiClient) GetMorphologicalAnalysis(str string) (GetMorphologicalAnalysisResponse, error) {
	params := url.Values{}
	params.Set("appid", y.YahooAppID)
	params.Set("sentence", str)
	params.Set("result", "hiragana")

	r := GetMorphologicalAnalysisResponse{}
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

type GetMorphologicalAnalysisResponse struct {
	MaResult MaResult `xml:"ma_result,omitempty"`
}

type MaResult struct {
	FilteredCount int      `xml:"filtered_count,omitempty" `
	TotalCount    int      `xml:"total_count,omitempty"`
	WordList      WordList `xml:"word_list,omitempty"`
}

type WordList struct {
	Words []Word `xml:"word,omitempty"`
}

type Word struct {
	Pos     string `xml:"pos,omitempty"`
	Reading string `xml:"reading,omitempty"`
	Surface string `xml:"surface,omitempty"`
}
