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

	// TODO: Redisにキャッシュしてある結果がある場合それを返す

	// リプライを取得
	replyText := req.TweetCreateEvents[0].Text
	// 漢字をひらがなに変換(ex:GALA湯沢 -> GALAゆざわ)
	replyText = kanjiToHiragana(replyText, yahoo.NewYahooApiClient()) // TODO: ClientもDIしたほうがいいかも
	// ひらがなをアルファベットに変換(ex:GALAゆざわ -> GALAyuzawa)
	replyText = toHebon(replyText)
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

type CharHebon struct {
	Char  string
	Hebon string
}

func charHebonByIndex(str string, index int) CharHebon {
	hebonMap := map[string]string{
		"あ": "a", "い": "i", "う": "u", "え": "e", "お": "o",
		"か": "ka", "き": "ki", "く": "ku", "け": "ke", "こ": "ko",
		"さ": "sa", "し": "shi", "す": "su", "せ": "se", "そ": "so",
		"た": "ta", "ち": "chi", "つ": "tsu", "て": "te", "と": "to",
		"な": "na", "に": "ni", "ぬ": "nu", "ね": "ne", "の": "no",
		"は": "ha", "ひ": "hi", "ふ": "fu", "へ": "he", "ほ": "ho",
		"ま": "ma", "み": "mi", "む": "mu", "め": "me", "も": "mo",
		"や": "ya", "ゆ": "yu", "よ": "yo",
		"ら": "ra", "り": "ri", "る": "ru", "れ": "re", "ろ": "ro",
		"わ": "wa", "ゐ": "i", "ゑ": "e", "を": "o",
		"ぁ": "a", "ぃ": "i", "ぅ": "u", "ぇ": "e", "ぉ": "o",
		"が": "ga", "ぎ": "gi", "ぐ": "gu", "げ": "ge", "ご": "go",
		"ざ": "za", "じ": "ji", "ず": "zu", "ぜ": "ze", "ぞ": "zo",
		"だ": "da", "ぢ": "ji", "づ": "zu", "で": "de", "ど": "do",
		"ば": "ba", "び": "bi", "ぶ": "bu", "べ": "be", "ぼ": "bo",
		"ぱ": "pa", "ぴ": "pi", "ぷ": "pu", "ぺ": "pe", "ぽ": "po",
		"きゃ": "kya", "きゅ": "kyu", "きょ": "kyo",
		"しゃ": "sha", "しゅ": "shu", "しょ": "sho",
		"ちゃ": "cha", "ちゅ": "chu", "ちょ": "cho", "ちぇ": "che",
		"にゃ": "nya", "にゅ": "nyu", "にょ": "nyo",
		"ひゃ": "hya", "ひゅ": "hyu", "ひょ": "hyo",
		"みゃ": "mya", "みゅ": "myu", "みょ": "myo",
		"りゃ": "rya", "りゅ": "ryu", "りょ": "ryo",
		"ぎゃ": "gya", "ぎゅ": "gyu", "ぎょ": "gyo",
		"じゃ": "ja", "じゅ": "ju", "じょ": "jo",
		"びゃ": "bya", "びゅ": "byu", "びょ": "byo",
		"ぴゃ": "pya", "ぴゅ": "pyu", "ぴょ": "pyo",
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
		"aa": true, "ee": true, "ii": false, // i は連続しても省略しない
		"oo": true, "ou": true, "uu": true,
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
				if strings.Index(nextCh.Hebon, "ch") == 0 {
					ch.Hebon = "t"
				} else {
					ch.Hebon = nextCh.Hebon[0:1]
				}
			}
		} else if ch.Char == "ん" {
			// B,M,P の前の "ん" は "M" とする。
			nextCh := charHebonByIndex(str, i+1)
			if nextCh.Hebon != "" && strings.Index("bmp", nextCh.Hebon[0:1]) != -1 {
				ch.Hebon = "m"
			} else {
				ch.Hebon = "n"
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
