package search

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSimilarSnowResortFromReply(t *testing.T) {
	before()

	testClient, err := testClient()
	if err != nil {
		t.Error(err)
	}

	snowResortRepositoryMock := SnowResortRepositoryMock{Client: testClient}

	snowResortServiceMock := SnowResortServiceImpl{
		SnowResortRepository: &snowResortRepositoryMock,
		YahooApiClient:       &ApiClientMock{},
	}

	t.Run("Get correct similar snow resort", func(t *testing.T) {
		cases := []struct {
			input  string
			output SnowResort
		}{
			{
				input:  "白馬47",
				output: SnowResort{Label: "Hakuba 47", SearchWord: "Hakuba47"},
			},
			{
				input:  "hakuba",
				output: SnowResort{Label: "Hakuba 47", SearchWord: "Hakuba47"},
			},
			{
				input:  "47",
				output: SnowResort{Label: "Hakuba 47", SearchWord: "Hakuba47"},
			},
			{
				input:  "@snowfall_bot    　かぐら",
				output: SnowResort{Label: "Kagura", SearchWord: "TashiroKaguraMitsumata"},
			},
			{
				input:  "@snowfall_bot 　みつ　また",
				output: SnowResort{Label: "Kagura", SearchWord: "TashiroKaguraMitsumata"},
			},
			{
				input:  "高鷲SP",
				output: SnowResort{Label: "Takasu Snow Park", SearchWord: "TakasuSnowPark"},
			},
			{
				input:  "@snowfall_bot GALA湯沢",
				output: SnowResort{Label: "Gala Yuzawa", SearchWord: "Gala-Yuzawa"},
			},
			{
				input:  "今庄",
				output: SnowResort{Label: "Imajo 365", SearchWord: "Imajo365"},
			},
			{
				input:  "ニセコ",
				output: SnowResort{Label: "Niseko Grand Hirafu", SearchWord: "Niseko"},
			},
			{
				input:  "石打丸山",
				output: SnowResort{Label: "Ishiuchi Maruyama", SearchWord: "IshiuchiMaruyama"},
			},
			{
				input:  "赤倉観光",
				output: SnowResort{Label: "Akakura Kanko", SearchWord: "Akakura-Shin-Akakura"},
			},
		}

		for _, tt := range cases {
			snowResort, err := snowResortServiceMock.GetSimilarSnowResortFromReply(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(snowResort, tt.output); diff != "" {
				t.Errorf("Diff: (-got +want)\n%s", diff)
			}
		}
	})

	t.Run("Set amd Get cached data", func(t *testing.T) {
		key := "skijam"
		skijam := SnowResort{SearchWord: "SkiJamKatsuyama", Label: "Ski Jam Katsuyama"}
		// 1回目の検索
		snowResort, err := snowResortServiceMock.GetSimilarSnowResortFromReply(key)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(snowResort, skijam); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
		// キャッシュされているか確認
		cachedSnowResort, err := snowResortRepositoryMock.FindSnowResort(key)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(cachedSnowResort, skijam); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
		// 2回目の検索
		// キャッシュから返しているか
		listSnowResortsCallCount := snowResortRepositoryMock.ListSnowResortsCallCount
		findSnowResortCallCount := snowResortRepositoryMock.FindSnowResortCallCount
		setSnowResortCallCount := snowResortRepositoryMock.SetSnowResortCallCount

		snowResort, err = snowResortServiceMock.GetSimilarSnowResortFromReply(key)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(snowResort, skijam); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}

		afterListSnowResortsCallCount := snowResortRepositoryMock.ListSnowResortsCallCount
		afterFindSnowResortCallCount := snowResortRepositoryMock.FindSnowResortCallCount
		afterSetSnowResortCallCount := snowResortRepositoryMock.SetSnowResortCallCount

		if afterListSnowResortsCallCount-listSnowResortsCallCount != 0 {
			t.Error("Not Cached")
		}
		if afterFindSnowResortCallCount-findSnowResortCallCount != 1 {
			t.Error("Not Cached")
		}
		if afterSetSnowResortCallCount-setSnowResortCallCount != 0 {
			t.Error("Not Cached")
		}
	})
}
