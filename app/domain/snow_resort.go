package domain

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kotaroooo0/gojaconv/jaconv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"golang.org/x/exp/utf8string"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
)

type SnowResort struct {
	SearchWord string
	Label      string
}

type Tweet struct {
	ID             string
	UserScreenName string
	Text           string
}

type SnowResortService interface {
	ReplyForecast(SnowResort, Tweet) (SnowResort, error)
	GetSimilarSnowResortFromReply(string) (SnowResort, error)
}

type SnowResortServiceImpl struct {
	// ドメイン層はどの層にも依存しない
	SnowResortRepository  SnowResortRepository
	YahooApiClient        yahoo.IYahooApiClient
	TwitterApiClient      twitter.ITwitterApiClient
	SnowforecastApiClient snowforecast.ISnowforecastApiClient
}

func (ss SnowResortServiceImpl) ReplyForecast(snowResort SnowResort, tweet Tweet) (SnowResort, error) {
	params := url.Values{}
	params.Set("in_reply_to_status_id", tweet.ID)
	content, err := replyContent(snowResort, ss.SnowforecastApiClient)
	if err != nil {
		return SnowResort{}, err
	}
	_, err = ss.TwitterApiClient.PostTweet(fmt.Sprintf("@%s %s", tweet.UserScreenName, content), params)
	if err != nil {
		return SnowResort{}, err
	}
	return SnowResort{}, nil
}

// TODO: メソッドが大きすぎるので分割してもいかも
func (ss SnowResortServiceImpl) GetSimilarSnowResortFromReply(reply string) (SnowResort, error) {
	// @snowfall_botを消す
	replyText := strings.Replace(reply, "@snowfall_bot ", "", -1)
	// スペースを消す
	replyText = strings.Replace(replyText, " ", "", -1)
	key := strings.Replace(replyText, "　", "", -1)

	// Redisにキャッシュしてある場合それを返す
	cachedSnowResort, err := ss.SnowResortRepository.FindSnowResort(key)
	if err != nil {
		return SnowResort{}, err
	}
	if (cachedSnowResort != SnowResort{}) {
		return cachedSnowResort, nil
	}

	// 漢字をひらがなに変換(ex:GALA湯沢 -> GALAゆざわ)
	replyText = toHiragana(key, ss.YahooApiClient)
	// ひらがなをアルファベットに変換(ex:GALAゆざわ -> GALAyuzawa)
	replyText = jaconv.ToHebon(replyText)
	// 残った大文字を小文字に直す(ex:GALAyuzawa -> galayuzawa)
	replyText = strings.ToLower(replyText)

	lowercaseSnowResorts, err := ss.SnowResortRepository.ListSnowResorts("lowercase-snowresorts-searchword")
	if err != nil {
		return SnowResort{}, err
	}
	snowResortLabels, err := ss.SnowResortRepository.ListSnowResorts("lowercase-snowresorts-label")
	if err != nil {
		return SnowResort{}, err
	}
	sources := append(lowercaseSnowResorts, snowResortLabels...)
	fmt.Println(sources)
	fmt.Println(len(sources))
	similarSkiResortString := getSimilarSkiResort(replyText, sources)

	targetSnowResort, err := ss.SnowResortRepository.FindSnowResort(similarSkiResortString)
	if err != nil {
		return SnowResort{}, err
	}

	// Redisにキャッシュする
	err = ss.SnowResortRepository.SetSnowResort(key, targetSnowResort)
	if err != nil {
		return SnowResort{}, err
	}

	return targetSnowResort, nil
}

func replyContent(snowResort SnowResort, snowforecastApiClient snowforecast.ISnowforecastApiClient) (string, error) {
	snowfallForecast, err := snowforecastApiClient.GetSnowfallForecastBySkiResortSearchWord(snowResort.SearchWord)
	if err != nil {
		return "", err
	}

	// TODO: 仮の文章
	content := snowResort.Label + "\n"
	content += "今日 | 明日 | 明後日\n"
	content += "3日後 | 4日後 | 5日後\n"
	content += strconv.Itoa(snowfallForecast.Snows[0].Morning) + addRainyChar(snowfallForecast.Rains[0].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Noon) + addRainyChar(snowfallForecast.Rains[0].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Night) + addRainyChar(snowfallForecast.Rains[0].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[1].Morning) + addRainyChar(snowfallForecast.Rains[1].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Noon) + addRainyChar(snowfallForecast.Rains[1].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Night) + addRainyChar(snowfallForecast.Rains[1].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[2].Morning) + addRainyChar(snowfallForecast.Rains[2].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Noon) + addRainyChar(snowfallForecast.Rains[2].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Night) + addRainyChar(snowfallForecast.Rains[2].Night) + "cm\n"
	content += strconv.Itoa(snowfallForecast.Snows[3].Morning) + addRainyChar(snowfallForecast.Rains[3].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[3].Noon) + addRainyChar(snowfallForecast.Rains[3].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[3].Night) + addRainyChar(snowfallForecast.Rains[3].Night) + "cm |"
	content += strconv.Itoa(snowfallForecast.Snows[4].Morning) + addRainyChar(snowfallForecast.Rains[4].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[4].Noon) + addRainyChar(snowfallForecast.Rains[4].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[4].Night) + addRainyChar(snowfallForecast.Rains[4].Night) + "cm |"
	content += strconv.Itoa(snowfallForecast.Snows[5].Morning) + addRainyChar(snowfallForecast.Rains[5].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[5].Noon) + addRainyChar(snowfallForecast.Rains[5].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[5].Night) + addRainyChar(snowfallForecast.Rains[5].Night) + "cm"

	// 140字までに切り詰めて返す
	if len([]rune(content)) > 140 {
		return utf8string.NewString(content).Slice(0, 140), nil
	}
	return content, nil
}

func addRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	} else if rainfall > 0 {
		return "☂️"
	}
	return ""
}

func toHiragana(str string, yahooApiClient yahoo.IYahooApiClient) string {
	res, err := yahooApiClient.GetMorphologicalAnalysis(str)
	if err != nil {
		panic(err)
	}

	h := ""
	for _, w := range res.MaResult.WordList.Words {
		h += w.Reading
	}
	return h
}

func getSimilarSkiResort(target string, skiresorts []string) string {
	// targetは小文字にしておく
	target = strings.ToLower(target)

	// レーベンシュタイン距離を計算する際の重みづけ
	// 削除の際の距離を小さくしている
	// TODO: 標準化や他の編集距離を考える必要もある、その際に評価をするためにラベル付おじさんにならないといけないかもしれない
	myOptions := levenshtein.Options{
		InsCost: 10,
		DelCost: 1,
		SubCost: 10,
		Matches: levenshtein.IdenticalRunes,
	}
	distances := make([]int, len(skiresorts))
	for i := 0; i < len(skiresorts); i++ {
		distances[i] = levenshtein.DistanceForStrings([]rune(skiresorts[i]), []rune(target), myOptions)
	}

	// 距離が最小のもののインデックスを取得する
	minIdx := 0
	minDistance := 1000000 // 十分大きな数
	for i := 0; i < len(distances); i++ {
		if distances[i] <= minDistance {
			minDistance = distances[i]
			minIdx = i
		}
	}

	return skiresorts[minIdx]
}

type SnowResortRepository interface {
	ListSnowResorts(string) ([]string, error)
	FindSnowResort(string) (SnowResort, error)
	SetSnowResort(string, SnowResort) error
}
