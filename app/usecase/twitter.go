package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"golang.org/x/exp/utf8string"
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
	// リプライを取得
	replyText := req.TweetCreateEvents[0].Text
	// 漢字をひらがなに変換
	replyText = kanjiToHiragana(replyText, yahoo.NewYahooApiClient())
	// ひらがなをアルファベットに変換
	replyText = toHebon(replyText)

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
	for _, w := range res.MaResult.WordList {
		h += w.Reading
	}
	return h
}

type CharHebon struct {
	Char  string
	Hebon string
}

func charHebonByIndex(str string, index int) CharHebon {
	hebonMap := map[string]string{
		"あ": "A", "い": "I", "う": "U", "え": "E", "お": "O",
		"か": "KA", "き": "KI", "く": "KU", "け": "KE", "こ": "KO",
		"さ": "SA", "し": "SHI", "す": "SU", "せ": "SE", "そ": "SO",
		"た": "TA", "ち": "CHI", "つ": "TSU", "て": "TE", "と": "TO",
		"な": "NA", "に": "NI", "ぬ": "NU", "ね": "NE", "の": "NO",
		"は": "HA", "ひ": "HI", "ふ": "FU", "へ": "HE", "ほ": "HO",
		"ま": "MA", "み": "MI", "む": "MU", "め": "ME", "も": "MO",
		"や": "YA", "ゆ": "YU", "よ": "YO",
		"ら": "RA", "り": "RI", "る": "RU", "れ": "RE", "ろ": "RO",
		"わ": "WA", "ゐ": "I", "ゑ": "E", "を": "O",
		"ぁ": "A", "ぃ": "I", "ぅ": "U", "ぇ": "E", "ぉ": "O",
		"が": "GA", "ぎ": "GI", "ぐ": "GU", "げ": "GE", "ご": "GO",
		"ざ": "ZA", "じ": "JI", "ず": "ZU", "ぜ": "ZE", "ぞ": "ZO",
		"だ": "DA", "ぢ": "JI", "づ": "ZU", "で": "DE", "ど": "DO",
		"ば": "BA", "び": "BI", "ぶ": "BU", "べ": "BE", "ぼ": "BO",
		"ぱ": "PA", "ぴ": "PI", "ぷ": "PU", "ぺ": "PE", "ぽ": "PO",
		"きゃ": "KYA", "きゅ": "KYU", "きょ": "KYO",
		"しゃ": "SHA", "しゅ": "SHU", "しょ": "SHO",
		"ちゃ": "CHA", "ちゅ": "CHU", "ちょ": "CHO", "ちぇ": "CHE",
		"にゃ": "NYA", "にゅ": "NYU", "にょ": "NYO",
		"ひゃ": "HYA", "ひゅ": "HYU", "ひょ": "HYO",
		"みゃ": "MYA", "みゅ": "MYU", "みょ": "MYO",
		"りゃ": "RYA", "りゅ": "RYU", "りょ": "RYO",
		"ぎゃ": "GYA", "ぎゅ": "GYU", "ぎょ": "GYO",
		"じゃ": "JA", "じゅ": "JU", "じょ": "JO",
		"びゃ": "BYA", "びゅ": "BYU", "びょ": "BYO",
		"ぴゃ": "PYA", "ぴゅ": "PYU", "ぴょ": "PYO",
	}

	var hebon string
	var char string
	utfstr := utf8string.NewString(str)
	// 2文字ヒットするとき
	if index+1 < utf8.RuneCountInString(str) {
		char = utfstr.Slice(index, index+2)
		hebon = hebonMap[char]
	}
	// 2文字はヒットしないが1文字はヒットするとき
	if hebon == "" && index < utfstr.RuneCount() {
		char = utfstr.Slice(index, index+1)
		hebon = hebonMap[char]
	}
	return CharHebon{Char: char, Hebon: hebon}
}

func toHebon(str string) string {
	isOmitted := map[string]bool{
		"AA": true, "EE": true, "II": false, // I は連続しても省略しない
		"OO": true, "OU": true, "UU": true,
	}

	var hebon string
	var lastHebon string

	i := 0
	for {
		ch := charHebonByIndex(str, i)
		if ch.Char == "っ" {
			// "っち"
			nextCh := charHebonByIndex(str, i+1)
			if nextCh.Hebon != "" {
				if strings.Index(nextCh.Hebon, "CH") == 0 {
					ch.Hebon = "T"
				} else {
					ch.Hebon = nextCh.Hebon[0:1]
				}
			}
		} else if ch.Char == "ん" {
			// B,M,P の前の "ん" は "M" とする。
			nextCh := charHebonByIndex(str, i+1)
			if nextCh.Hebon != "" && strings.Index("BMP", nextCh.Hebon[0:1]) != -1 {
				ch.Hebon = "M"
			} else {
				ch.Hebon = "N"
			}
		} else if ch.Char == "ー" {
			// 長音は無視
			ch.Hebon = ""
		}

		if ch.Hebon != "" {
			// 変換できる文字の場合
			if lastHebon != "" {
				// 連続する母音の除去
				joinedHebon := lastHebon + ch.Hebon
				if len(joinedHebon) > 2 {
					joinedHebon = joinedHebon[len(joinedHebon)-2:]
				}
				if isOmitted[joinedHebon] {
					ch.Hebon = ""
				}
			}
			hebon += ch.Hebon
		} else {
			if ch.Char != "ー" {
				// 変換できない文字の場合
				hebon += ch.Char
			}
		}

		lastHebon = ch.Hebon
		i += utf8.RuneCountInString(ch.Char)
		if i >= utf8.RuneCountInString(str) {
			break
		}
	}

	return hebon
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
	UserID            string `json:"for_user_id" form:"for_user_id" binding:"required"`
	TweetCreateEvents []struct {
		TweetID    int64  `json:"id" form:"id" binding:"required"`
		TweetIDStr string `json:"id_str" form:"id_str" binding:"required"`
		User       struct {
			UserID     int64  `json:"id" form:"id" binding:"required"`
			IDStr      string `json:"id_str" form:"id_str" binding:"required"`
			ScreenName string `json:"screen_name" form:"screen_name" binding:"required"`
		} `json:"user" form:"user" binding:"required"`
		Text string `json:"text" form:"text" binding:"required"`
	} `json:"tweet_create_events" form:"tweet_create_events" binding:"required"`
}

type PostTwitterWebhookResponse struct {
	SnowResortLabel string `json:"snow_resort"`
}
