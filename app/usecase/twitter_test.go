package usecase

import (
	"fmt"
	"testing"

	"github.com/kotaroooo0/snowforecast-twitter-bot/parameters/responses"
)

func CreateDummyPostKanjiToHiraganaResponse(readings []string) responses.PostKanjiToHiraganaResponse {
	res := responses.NewPostKanjiToHiraganaResponse()
	for _, r := range readings {
		res.MaResult.WordList = append(res.MaResult.WordList, responses.Word{Reading: r})
	}
	return res
}

type ApiClientMock struct {
}

func (a *ApiClientMock) GetMorphologicalAnalysis(str string) (responses.PostKanjiToHiraganaResponse, error) {
	switch str {
	case "白馬47":
		{
			return CreateDummyPostKanjiToHiraganaResponse([]string{"はくば", "47"}), nil
		}
	case "妙高杉ノ原":
		{
			return CreateDummyPostKanjiToHiraganaResponse([]string{"みょうこう", "すぎのはら"}), nil
		}
	case "高鷲スノーパーク":
		{
			return CreateDummyPostKanjiToHiraganaResponse([]string{"たかす", "すのー", "ぱーく"}), nil
		}
	case "GALA湯沢":
		{
			return CreateDummyPostKanjiToHiraganaResponse([]string{"GALA", "ゆざわ"}), nil
		}
	}
	return responses.NewPostKanjiToHiraganaResponse(), nil
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
		{kanji: "GALA", hiragana: "GALA"},
		{kanji: "hakuba47", hiragana: "hakuba47"},
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
		{s: "はくば47", index: 2, ch: CharHebon{Char: "ば", Hebon: "BA"}},
		{s: "はくば47", index: 3, ch: CharHebon{Char: "4", Hebon: ""}},
		{s: "はくば47", index: 4, ch: CharHebon{Char: "7", Hebon: ""}},
		{s: "みょうこうすぎのはら", index: 0, ch: CharHebon{Char: "みょ", Hebon: "MYO"}},
		{s: "みょうこうすぎのはら", index: 1, ch: CharHebon{Char: "ょ", Hebon: ""}},
		{s: "みょうこうすぎのはら", index: 2, ch: CharHebon{Char: "う", Hebon: "U"}},
		{s: "たかすすのーぱーく", index: 4, ch: CharHebon{Char: "の", Hebon: "NO"}},
		{s: "たかすすのーぱーく", index: 5, ch: CharHebon{Char: "ー", Hebon: ""}},
		{s: "たかすすのーぱーく", index: 6, ch: CharHebon{Char: "ぱ", Hebon: "PA"}},
		{s: "GALAゆざわ", index: 0, ch: CharHebon{Char: "G", Hebon: ""}},
		{s: "GALAゆざわ", index: 1, ch: CharHebon{Char: "A", Hebon: ""}},
		{s: "GALAゆざわ", index: 4, ch: CharHebon{Char: "ゆ", Hebon: "YU"}},
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
		{hiragana: "はくば47", hebon: "HAKUBA47"},
		{hiragana: "みょうこうすぎのはら", hebon: "MYOKOSUGINOHARA"},
		{hiragana: "たかすすのーぱーく", hebon: "TAKASUSUNOPAKU"},
		{hiragana: "GALAゆざわ", hebon: "GALAYUZAWA"},
	}
	for _, tt := range cases {
		act := toHebon(tt.hiragana)
		if act != tt.hebon {
			t.Error(fmt.Sprintf("%s is not %s", act, tt.hebon))
		}
	}
}
