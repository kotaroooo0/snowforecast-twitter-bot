package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"

	"github.com/kotaroooo0/gojaconv/jaconv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type TwitterUseCase interface {
	NewGetTwitterWebhookRequest() GetTwitterWebhookRequest
	NewPostTwitterWebhookRequest() PostTwitterWebhookRequest
	GetCrcTokenResponse(GetTwitterWebhookRequest) GetTwitterWebhookResponse
	PostAutoReplyResponse(PostTwitterWebhookRequest) PostTwitterWebhookResponse
}

type TwitterUseCaseImpl struct {
	SnowResortService    domain.SnowResortService
	SnowResortRepository domain.SnowResortRepository
}

func (tu TwitterUseCaseImpl) NewGetTwitterWebhookRequest() GetTwitterWebhookRequest {
	return GetTwitterWebhookRequest{}
}

// TwitterのWebhookの認証に用いる
// ref: https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks
func (tu TwitterUseCaseImpl) GetCrcTokenResponse(req GetTwitterWebhookRequest) GetTwitterWebhookResponse {
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(req.CrcToken))
	return GetTwitterWebhookResponse{
		Token: "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	}
}

func (tu TwitterUseCaseImpl) NewPostTwitterWebhookRequest() PostTwitterWebhookRequest {
	return PostTwitterWebhookRequest{}
}

func (tu TwitterUseCaseImpl) PostAutoReplyResponse(req PostTwitterWebhookRequest) PostTwitterWebhookResponse {
	// リプライがない、もしくはユーザが不正な場合は空を返す
	if len(req.TweetCreateEvents) < 1 || req.UserID == req.TweetCreateEvents[0].User.IDStr {
		return PostTwitterWebhookResponse{}
	}

	// TODO: Redisにキャッシュしてある結果がある場合それを返す

	// リプライを取得
	replyText := req.TweetCreateEvents[0].Text
	// 漢字をひらがなに変換(ex:GALA湯沢 -> GALAゆざわ)
	replyText = kanjiToHiragana(replyText, yahoo.NewYahooApiClient()) // TODO: ClientもDIしたほうがいいかも
	// ひらがなをアルファベットに変換(ex:GALAゆざわ -> GALAyuzawa)
	replyText = jaconv.ToHebon(replyText)
	// 残った大文字を小文字に直す(ex:GALAyuzawa -> galayuzawa)
	replyText = strings.ToLower(replyText)

	// リプライから全世界のスキー場の中で最も適切なスキー場を求める
	lowercaseSnowResorts, err := tu.SnowResortRepository.ListSnowResorts("lowercase-snowresorts-searchword")
	if err != nil {
		return PostTwitterWebhookResponse{}
	}
	snowResortLabels, err := tu.SnowResortRepository.ListSnowResorts("lowercase-snowresorts-label")
	if err != nil {
		return PostTwitterWebhookResponse{}
	}
	sources := append(lowercaseSnowResorts, snowResortLabels...)
	similarSkiResortString := getSimilarSkiResort(replyText, sources)
	similarSkiResort, err := tu.SnowResortRepository.FindSnowResort(similarSkiResortString)
	if err != nil {
		return PostTwitterWebhookResponse{}
	}

	skiResort, err := tu.SnowResortService.ReplyForecast(similarSkiResort)
	if err != nil {
		return PostTwitterWebhookResponse{}
	}
	return PostTwitterWebhookResponse{skiResort.Label}
}

func kanjiToHiragana(str string, yahooApiClient yahoo.IYahooApiClient) string {
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

type GetTwitterWebhookRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

type GetTwitterWebhookResponse struct {
	Token string `json:"response_token"`
}

type PostTwitterWebhookRequest struct {
	UserID            string             `json:"for_user_id" form:"for_user_id" binding:"required"`
	TweetCreateEvents []TweetCreateEvent `json:"tweet_create_events" form:"tweet_create_events" binding:"required"`
}

type TweetCreateEvent struct {
	TweetID    int64  `json:"id" form:"id" binding:"required"`
	TweetIDStr string `json:"id_str" form:"id_str" binding:"required"`
	User       struct {
		UserID     int64  `json:"id" form:"id" binding:"required"`
		IDStr      string `json:"id_str" form:"id_str" binding:"required"`
		ScreenName string `json:"screen_name" form:"screen_name" binding:"required"`
	} `json:"user" form:"user" binding:"required"`
	Text string `json:"text" form:"text" binding:"required"`
}

type PostTwitterWebhookResponse struct {
	SnowResortLabel string `json:"snow_resort"`
}
