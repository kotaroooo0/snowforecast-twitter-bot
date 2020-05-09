package usecase

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-redis/redis"
	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/pkg/errors"
)

func before() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type SnowResortServiceMock struct{}

func (sm SnowResortServiceMock) ReplyForecast(snowResort domain.SnowResort) (domain.SnowResort, error) {
	return snowResort, nil
}

type SnowResortRepositoryMock struct {
	Client *redis.Client
}

func testClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // 1のDBをテスト用とする
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

// TODO: DomainModelを返すように修正
func (s SnowResortRepositoryMock) ListSnowResorts(key string) ([]string, error) {
	result, err := s.Client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s SnowResortRepositoryMock) FindSnowResort(key string) (domain.SnowResort, error) {
	result, err := s.Client.HGetAll(key).Result()
	if err != nil {
		return domain.SnowResort{}, err
	}

	return domain.SnowResort{SearchWord: result["search_word"], Label: result["label"]}, nil
}

func createPostTwitterWebhookRequest(text string) PostTwitterWebhookRequest {
	req := PostTwitterWebhookRequest{UserID: "hoge"}
	req.TweetCreateEvents = append(req.TweetCreateEvents, TweetCreateEvent{})
	req.TweetCreateEvents[0].Text = text
	req.TweetCreateEvents[0].User.IDStr = "fuga"
	return req
}

func TestPostAutoReplyResponse(t *testing.T) {
	before()

	testClient, err := testClient()
	if err != nil {
		t.Error(err)
	}
	useCaseImpl := TwitterUseCaseImpl{
		SnowResortService:    SnowResortServiceMock{},
		SnowResortRepository: SnowResortRepositoryMock{Client: testClient},
	}

	cases := []struct {
		input  string
		output PostTwitterWebhookResponse
	}{
		{
			input:  "白馬47",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Hakuba 47"},
		},
		{
			input:  "hakuba",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Hakuba 47"},
		},
		{
			input:  "47",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Hakuba 47"},
		},
		{
			input:  "かぐら",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Kagura"},
		},
		{
			input:  "みつまた",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Kagura"},
		},
		{
			input:  "高鷲SP",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Takasu Snow Park"},
		},
		{
			input:  "GALA湯沢",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Gala Yuzawa"},
		},
		{
			input:  "今庄",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Imajo 365"},
		},
		{
			input:  "ニセコ",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Niseko Grand Hirafu"},
		},
		{
			input:  "石打丸山",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Ishiuchi Maruyama"},
		},
		{
			input:  "赤倉観光",
			output: PostTwitterWebhookResponse{SnowResortLabel: "Akakura Kanko"},
		},
	}

	for _, tt := range cases {
		req := createPostTwitterWebhookRequest(tt.input)
		if diff := cmp.Diff(useCaseImpl.PostAutoReplyResponse(req), tt.output); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}

func createGetMorphologicalAnalysisResponse(readings []string) yahoo.GetMorphologicalAnalysisResponse {
	res := yahoo.GetMorphologicalAnalysisResponse{}
	for _, r := range readings {
		res.MaResult.WordList.Words = append(res.MaResult.WordList.Words, yahoo.Word{Pos: "hoge", Reading: r, Surface: "fuga"})
	}
	return res
}

type ApiClientMock struct {
}

func (a *ApiClientMock) GetMorphologicalAnalysis(str string) (yahoo.GetMorphologicalAnalysisResponse, error) {
	switch str {
	case "白馬47":
		{
			return createGetMorphologicalAnalysisResponse([]string{"はくば", "47"}), nil
		}
	case "妙高杉ノ原":
		{
			return createGetMorphologicalAnalysisResponse([]string{"みょうこう", "すぎのはら"}), nil
		}
	case "高鷲スノーパーク":
		{
			return createGetMorphologicalAnalysisResponse([]string{"たかす", "すのー", "ぱーく"}), nil
		}
	case "GALA湯沢":
		{
			return createGetMorphologicalAnalysisResponse([]string{"GALA", "ゆざわ"}), nil
		}
	case "hakuba47":
		{
			return createGetMorphologicalAnalysisResponse([]string{"hakuba", "47"}), nil
		}
	case "hakuba 47":
		{
			return createGetMorphologicalAnalysisResponse([]string{"hakuba", " ", "47"}), nil
		}
	case "myokosuginohara":
		{
			return createGetMorphologicalAnalysisResponse([]string{"myokosuginohara"}), nil
		}
	case "hak妙47uba高":
		{
			return createGetMorphologicalAnalysisResponse([]string{"hak", "みょう", "47", "uba", "たか"}), nil
		}
	}
	return yahoo.GetMorphologicalAnalysisResponse{}, nil
}

func TestKanjiToHiragana(t *testing.T) {
	cases := []struct {
		kanji    string
		hiragana string
	}{
		{kanji: "白馬47", hiragana: "はくば47"},
		{kanji: "妙高杉ノ原", hiragana: "みょうこうすぎのはら"},
		{kanji: "高鷲スノーパーク", hiragana: "たかすすのーぱーく"},
		{kanji: "GALA湯沢", hiragana: "GALAゆざわ"},
		{kanji: "hakuba47", hiragana: "hakuba47"},
		{kanji: "hakuba 47", hiragana: "hakuba 47"},
		{kanji: "myokosuginohara", hiragana: "myokosuginohara"},
		{kanji: "hak妙47uba高", hiragana: "hakみょう47ubaたか"},
	}

	for _, tt := range cases {
		act := kanjiToHiragana(tt.kanji, &ApiClientMock{})
		if act != tt.hiragana {
			t.Error(fmt.Sprintf("%s is not %s", act, tt.kanji))
		}
	}
}

func TestCharHebonByIndex(t *testing.T) {
	cases := []struct {
		s     string
		index int
		ch    CharHebon
	}{
		{s: "はくば47", index: 2, ch: CharHebon{Char: "ば", Hebon: "ba"}},
		{s: "はくば47", index: 3, ch: CharHebon{Char: "4", Hebon: ""}},
		{s: "はくば47", index: 4, ch: CharHebon{Char: "7", Hebon: ""}},
		{s: "みょうこうすぎのはら", index: 0, ch: CharHebon{Char: "みょ", Hebon: "myo"}},
		{s: "みょうこうすぎのはら", index: 1, ch: CharHebon{Char: "ょ", Hebon: ""}},
		{s: "みょうこうすぎのはら", index: 2, ch: CharHebon{Char: "う", Hebon: "u"}},
		{s: "たかすすのーぱーく", index: 4, ch: CharHebon{Char: "の", Hebon: "no"}},
		{s: "たかすすのーぱーく", index: 5, ch: CharHebon{Char: "ー", Hebon: ""}},
		{s: "たかすすのーぱーく", index: 6, ch: CharHebon{Char: "ぱ", Hebon: "pa"}},
		{s: "GALAゆざわ", index: 0, ch: CharHebon{Char: "G", Hebon: ""}},
		{s: "GALAゆざわ", index: 1, ch: CharHebon{Char: "A", Hebon: ""}},
		{s: "GALAゆざわ", index: 4, ch: CharHebon{Char: "ゆ", Hebon: "yu"}},
	}

	for _, tt := range cases {
		act := charHebonByIndex(tt.s, tt.index)
		if act != tt.ch {
			t.Error(fmt.Sprintf("%s is not %s", act, tt.ch))
		}
	}
}

func TestToHebon(t *testing.T) {
	cases := []struct {
		hiragana string
		hebon    string
	}{
		{hiragana: "はくば47", hebon: "hakuba47"},
		{hiragana: "みょうこうすぎのはら", hebon: "myokosuginohara"},
		{hiragana: "たかすすのーぱーく", hebon: "takasusunopaku"},
		{hiragana: "GALAゆざわ", hebon: "GALAyuzawa"},
	}
	for _, tt := range cases {
		act := toHebon(tt.hiragana)
		if act != tt.hebon {
			t.Error(fmt.Sprintf("%s is not %s", act, tt.hebon))
		}
	}
}
