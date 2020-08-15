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
	Id        int    `db:"id"`
	Name      string `db:"name"`
	SearchKey string `db:"search_key"`
}

type SnowResortRepository interface {
	FindAll() ([]*SnowResort, error)
}

type Tweet struct {
	ID             string
	UserScreenName string
	Text           string
}

type ReplyService interface {
	ReplyForecast(*Tweet) (*SnowResort, error)
}

type ReplyServiceImpl struct {
	// ドメイン層は他の層にも依存しない
	SnowResortRepository  SnowResortRepository
	YahooApiClient        yahoo.IYahooApiClient
	TwitterApiClient      twitter.ITwitterApiClient
	SnowforecastApiClient snowforecast.ISnowforecastApiClient
}

func NewReplyServiceImpl(snowResortRepository SnowResortRepository, yahooApiClient yahoo.IYahooApiClient, twitterApiClient twitter.ITwitterApiClient, snowforecastApiClient snowforecast.ISnowforecastApiClient) ReplyService {
	return &ReplyServiceImpl{
		SnowResortRepository:  snowResortRepository,
		YahooApiClient:        yahooApiClient,
		TwitterApiClient:      twitterApiClient,
		SnowforecastApiClient: snowforecastApiClient,
	}
}

// 検索すること自体がビジネスロジックであるため、ドメイン層に含める
func (ss ReplyServiceImpl) ReplyForecast(tweet *Tweet) (*SnowResort, error) {
	// リプライを取得
	replyText := tweet.Text
	// @snowfall_botを消す
	replyText = strings.Replace(replyText, "@snowfall_bot ", "", -1)
	// スペースを消す
	replyText = strings.Replace(replyText, " ", "", -1)
	key := strings.Replace(replyText, "　", "", -1)

	// キャッシュがあればやりたくないゾーン
	// 漢字をひらがなに変換(ex:GALA湯沢 -> GALAゆざわ)
	replyText, err := toHiragana(key, ss.YahooApiClient)
	if err != nil {
		return &SnowResort{}, err
	}
	// ひらがなをアルファベットに変換(ex:GALAゆざわ -> GALAyuzawa)
	replyText = jaconv.ToHebon(replyText)

	srs, err := ss.SnowResortRepository.FindAll()
	if err != nil {
		return &SnowResort{}, err
	}
	sr := findSimilarSnowResort(replyText, srs)
	// キャッシュがあればやりたくないゾーン

	sf, err := ss.SnowforecastApiClient.GetSnowfallForecastBySkiResortSearchWord(sr.SearchKey)
	if err != nil {
		return &SnowResort{}, err
	}

	content, err := replyContent(sr.Name, sf)
	if err != nil {
		return &SnowResort{}, err
	}

	params := url.Values{}
	params.Set("in_reply_to_status_id", tweet.ID)
	_, err = ss.TwitterApiClient.PostTweet(fmt.Sprintf("@%s %s", tweet.UserScreenName, content), params)
	if err != nil {
		return &SnowResort{}, err
	}
	return sr, nil
}

func replyContent(name string, sf snowforecast.SnowfallForecast) (string, error) {
	// TODO: 仮の文章
	content := name + "\n"
	content += "今日 | 明日 | 明後日\n"
	content += "3日後 | 4日後 | 5日後\n"
	content += strconv.Itoa(sf.Snows[0].Morning) + addRainyChar(sf.Rains[0].Morning) + ", " + strconv.Itoa(sf.Snows[0].Noon) + addRainyChar(sf.Rains[0].Noon) + ", " + strconv.Itoa(sf.Snows[0].Night) + addRainyChar(sf.Rains[0].Night) + "cm | "
	content += strconv.Itoa(sf.Snows[1].Morning) + addRainyChar(sf.Rains[1].Morning) + ", " + strconv.Itoa(sf.Snows[1].Noon) + addRainyChar(sf.Rains[1].Noon) + ", " + strconv.Itoa(sf.Snows[1].Night) + addRainyChar(sf.Rains[1].Night) + "cm | "
	content += strconv.Itoa(sf.Snows[2].Morning) + addRainyChar(sf.Rains[2].Morning) + ", " + strconv.Itoa(sf.Snows[2].Noon) + addRainyChar(sf.Rains[2].Noon) + ", " + strconv.Itoa(sf.Snows[2].Night) + addRainyChar(sf.Rains[2].Night) + "cm\n"
	content += strconv.Itoa(sf.Snows[3].Morning) + addRainyChar(sf.Rains[3].Morning) + ", " + strconv.Itoa(sf.Snows[3].Noon) + addRainyChar(sf.Rains[3].Noon) + ", " + strconv.Itoa(sf.Snows[3].Night) + addRainyChar(sf.Rains[3].Night) + "cm |"
	content += strconv.Itoa(sf.Snows[4].Morning) + addRainyChar(sf.Rains[4].Morning) + ", " + strconv.Itoa(sf.Snows[4].Noon) + addRainyChar(sf.Rains[4].Noon) + ", " + strconv.Itoa(sf.Snows[4].Night) + addRainyChar(sf.Rains[4].Night) + "cm |"
	content += strconv.Itoa(sf.Snows[5].Morning) + addRainyChar(sf.Rains[5].Morning) + ", " + strconv.Itoa(sf.Snows[5].Noon) + addRainyChar(sf.Rains[5].Noon) + ", " + strconv.Itoa(sf.Snows[5].Night) + addRainyChar(sf.Rains[5].Night) + "cm"

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

// TODO: メソッドが大きすぎるので分割してもいかも
// Lower Filterしてからレーベンシュタイン距離により類似単語検索する
func findSimilarSnowResort(source string, targets []*SnowResort) *SnowResort {
	// sourceを小文字に直す(ex:GALAyuzawa -> galayuzawa)
	source = strings.ToLower(source)

	// targetsもNameとSearchKeysの両方を小文字へ直す
	targetNames := make([]string, len(targets))
	for i := 0; i < len(targetNames); i++ {
		targetNames[i] = strings.ToLower(targets[i].Name)
	}
	targetSearchKeys := make([]string, len(targets))
	for i := 0; i < len(targetSearchKeys); i++ {
		targetSearchKeys[i] = strings.ToLower(targets[i].SearchKey)
	}
	allTargets := append(targetNames, targetSearchKeys...)

	// レーベンシュタイン距離を計算する際の重みづけ
	// 削除の際の距離を小さくしている
	// TODO: 標準化や他の編集距離を考える必要もある、その際に評価をするためにラベル付おじさんにならないといけないかもしれない
	myOptions := levenshtein.Options{
		InsCost: 10,
		DelCost: 1,
		SubCost: 10,
		Matches: levenshtein.IdenticalRunes,
	}
	distances := make([]int, len(allTargets))
	for i := 0; i < len(allTargets); i++ {
		distances[i] = levenshtein.DistanceForStrings([]rune(source), []rune(allTargets[i]), myOptions)
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

	// NameとSearchKeyを足して二倍になってるので、余りをとる
	return targets[minIdx%len(targets)]
}

func toHiragana(str string, yahooApiClient yahoo.IYahooApiClient) (string, error) {
	res, err := yahooApiClient.GetMorphologicalAnalysis(str)
	if err != nil {
		return "", err
	}

	h := ""
	for _, w := range res.MaResult.WordList.Words {
		h += w.Reading
	}
	return h, nil
}
