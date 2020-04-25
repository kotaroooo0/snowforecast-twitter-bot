package language

import (
	"strings"
	"unicode/utf8"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"golang.org/x/exp/utf8string"
)

func KanjiToHiragana(str string, yahooApiClient yahoo.IYahooApiClient) string {
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

func ToHebon(str string) string {
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
