package yahoo

import (
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
)

func before() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// 実際にAPIを叩くテスト
func TestGetMorphologicalAnalysis(t *testing.T) {
	before()
	cases := []struct {
		input  string
		output GetMorphologicalAnalysisResponse
	}{
		{
			input: "白馬47",
			output: GetMorphologicalAnalysisResponse{
				MaResult{
					FilteredCount: 2,
					TotalCount:    2,
					WordList: WordList{
						[]Word{
							{Pos: "名詞", Reading: "はくば", Surface: "白馬"},
							{Pos: "名詞", Reading: "47", Surface: "47"},
						},
					},
				},
			},
		},
		{
			input: "hakuba 47",
			output: GetMorphologicalAnalysisResponse{
				MaResult{
					FilteredCount: 3,
					TotalCount:    3,
					WordList: WordList{
						[]Word{
							{Pos: "名詞", Reading: "hakuba", Surface: "hakuba"},
							{Pos: "特殊", Reading: " ", Surface: " "},
							{Pos: "名詞", Reading: "47", Surface: "47"},
						},
					},
				},
			},
		},
		{
			input: "Myokosuginohara",
			output: GetMorphologicalAnalysisResponse{
				MaResult{
					FilteredCount: 1,
					TotalCount:    1,
					WordList: WordList{
						[]Word{
							{Pos: "名詞", Reading: "Myokosuginohara", Surface: "Myokosuginohara"},
						},
					},
				},
			},
		},
		{
			input: "hak妙47uba高",
			output: GetMorphologicalAnalysisResponse{
				MaResult{
					FilteredCount: 5,
					TotalCount:    5,
					WordList: WordList{
						[]Word{
							{Pos: "名詞", Reading: "hak", Surface: "hak"},
							{Pos: "名詞", Reading: "みょう", Surface: "妙"},
							{Pos: "名詞", Reading: "47", Surface: "47"},
							{Pos: "名詞", Reading: "uba", Surface: "uba"},
							{Pos: "名詞", Reading: "たか", Surface: "高"},
						},
					},
				},
			},
		},
	}

	client := NewYahooApiClient(NewYahooConfig(os.Getenv("YAHOO_APP_ID")))
	for _, tt := range cases {
		act, err := client.GetMorphologicalAnalysis(tt.input)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(act, tt.output); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}
