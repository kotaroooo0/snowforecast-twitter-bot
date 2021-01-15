package domain

// TODO:動作確認用に書いたが、モックを使ったしっかりしたテストを書くとありがたがられる
// func TestReplyContent(t *testing.T) {
// 	testClient := snowforecast.NewApiClient()

// 	cases := []struct {
// 		snowResort SnowResort
// 		output     string
// 	}{
// 		{
// 			snowResort: SnowResort{SearchWord: "Patagonia", Label: "Patagonia Heliski"},
// 			output:     "hoge",
// 		},
// 	}

// 	for _, tt := range cases {
// 		content, err := replyContent(tt.snowResort, testClient)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if diff := cmp.Diff(content, tt.output); diff != "" {
// 			t.Errorf("Diff: (-got +want)\n%s", diff)
// 		}
// 	}
// }

// func before() {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// }

// type ReplyServiceMock struct {
// 	SnowResortRepository SnowResortRepository
// 	YahooApiClient       yahoo.IYahooApiClient
// 	ApiClient            twitter.IApiClient
// }

// type SnowResortRepositoryMock struct {
// 	Client                   *redis.Client
// 	ListSnowResortsCallCount int
// 	FindSnowResortCallCount  int
// 	SetSnowResortCallCount   int
// }

// func testClient() (*redis.Client, error) {
// 	client := redis.NewClient(&redis.Options{
// 		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
// 		DB:   1, // 1のDBをテスト用とする
// 	})
// 	if err := client.Ping().Err(); err != nil {
// 		return nil, errors.Wrapf(err, "failed to ping redis server")
// 	}
// 	return client, nil
// }

// TODO: 以下の三つのメソッドはinfra層にほぼ同じ実装があるけどいいんだろうか
// infra層からimportしようとすると import cycle not allowed になる
// func (s *SnowResortRepositoryMock) ListSnowResorts(key string) ([]string, error) {
// 	s.ListSnowResortsCallCount++
// 	result, err := s.Client.SMembers(key).Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (s *SnowResortRepositoryMock) FindSnowResort(key string) (SnowResort, error) {
// 	s.FindSnowResortCallCount++
// 	result, err := s.Client.HGetAll(key).Result()
// 	if err != nil {
// 		return SnowResort{}, err
// 	}
// 	return SnowResort{SearchKey: result["search_word"], Name: result["label"]}, nil
// }

// func (s *SnowResortRepositoryMock) SetSnowResort(key string, snowResort SnowResort) error {
// 	s.SetSnowResortCallCount++
// 	err := s.Client.HMSet(key, map[string]interface{}{"search_word": snowResort.SearchKey, "label": snowResort.Name})
// 	return err.Err()
// }

// func TestGetSimilarSnowResortFromReply(t *testing.T) {
// 	before()

// 	testClient, err := testClient()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	snowResortRepositoryMock := SnowResortRepositoryMock{Client: testClient}

// 	snowResortServiceMock := SnowResortServiceImpl{
// 		SnowResortRepository: &snowResortRepositoryMock,
// 		YahooApiClient:       &ApiClientMock{},
// 	}

// 	t.Run("Get correct similar snow resort", func(t *testing.T) {
// 		cases := []struct {
// 			input  string
// 			output SnowResort
// 		}{
// 			{
// 				input:  "白馬47",
// 				output: SnowResort{Label: "Hakuba 47", SearchWord: "Hakuba47"},
// 			},
// 			{
// 				input:  "hakuba",
// 				output: SnowResort{Label: "Hakuba 47", SearchWord: "Hakuba47"},
// 			},
// 			{
// 				input:  "47",
// 				output: SnowResort{Label: "Hakuba 47", SearchWord: "Hakuba47"},
// 			},
// 			{
// 				input:  "@snowfall_bot    　かぐら",
// 				output: SnowResort{Label: "Kagura", SearchWord: "TashiroKaguraMitsumata"},
// 			},
// 			{
// 				input:  "@snowfall_bot 　みつ　また",
// 				output: SnowResort{Label: "Kagura", SearchWord: "TashiroKaguraMitsumata"},
// 			},
// 			{
// 				input:  "高鷲SP",
// 				output: SnowResort{Label: "Takasu Snow Park", SearchWord: "TakasuSnowPark"},
// 			},
// 			{
// 				input:  "@snowfall_bot GALA湯沢",
// 				output: SnowResort{Label: "Gala Yuzawa", SearchWord: "Gala-Yuzawa"},
// 			},
// 			{
// 				input:  "今庄",
// 				output: SnowResort{Label: "Imajo 365", SearchWord: "Imajo365"},
// 			},
// 			{
// 				input:  "ニセコ",
// 				output: SnowResort{Label: "Niseko Grand Hirafu", SearchWord: "Niseko"},
// 			},
// 			{
// 				input:  "石打丸山",
// 				output: SnowResort{Label: "Ishiuchi Maruyama", SearchWord: "IshiuchiMaruyama"},
// 			},
// 			{
// 				input:  "赤倉観光",
// 				output: SnowResort{Label: "Akakura Kanko", SearchWord: "Akakura-Shin-Akakura"},
// 			},
// 		}

// 		for _, tt := range cases {
// 			snowResort, err := snowResortServiceMock.GetSimilarSnowResortFromReply(tt.input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if diff := cmp.Diff(snowResort, tt.output); diff != "" {
// 				t.Errorf("Diff: (-got +want)\n%s", diff)
// 			}
// 		}
// 	})

// 	t.Run("Set amd Get cached data", func(t *testing.T) {
// 		key := "skijam"
// 		skijam := SnowResort{SearchWord: "SkiJamKatsuyama", Label: "Ski Jam Katsuyama"}
// 		// 1回目の検索
// 		snowResort, err := snowResortServiceMock.GetSimilarSnowResortFromReply(key)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if diff := cmp.Diff(snowResort, skijam); diff != "" {
// 			t.Errorf("Diff: (-got +want)\n%s", diff)
// 		}
// 		// キャッシュされているか確認
// 		cachedSnowResort, err := snowResortRepositoryMock.FindSnowResort(key)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if diff := cmp.Diff(cachedSnowResort, skijam); diff != "" {
// 			t.Errorf("Diff: (-got +want)\n%s", diff)
// 		}
// 		// 2回目の検索
// 		// キャッシュから返しているか
// 		listSnowResortsCallCount := snowResortRepositoryMock.ListSnowResortsCallCount
// 		findSnowResortCallCount := snowResortRepositoryMock.FindSnowResortCallCount
// 		setSnowResortCallCount := snowResortRepositoryMock.SetSnowResortCallCount

// 		snowResort, err = snowResortServiceMock.GetSimilarSnowResortFromReply(key)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if diff := cmp.Diff(snowResort, skijam); diff != "" {
// 			t.Errorf("Diff: (-got +want)\n%s", diff)
// 		}

// 		afterListSnowResortsCallCount := snowResortRepositoryMock.ListSnowResortsCallCount
// 		afterFindSnowResortCallCount := snowResortRepositoryMock.FindSnowResortCallCount
// 		afterSetSnowResortCallCount := snowResortRepositoryMock.SetSnowResortCallCount

// 		if afterListSnowResortsCallCount-listSnowResortsCallCount != 0 {
// 			t.Error("Not Cached")
// 		}
// 		if afterFindSnowResortCallCount-findSnowResortCallCount != 1 {
// 			t.Error("Not Cached")
// 		}
// 		if afterSetSnowResortCallCount-setSnowResortCallCount != 0 {
// 			t.Error("Not Cached")
// 		}
// 	})
// }

// func createGetMorphologicalAnalysisResponse(readings []string) yahoo.GetMorphologicalAnalysisResponse {
// 	res := yahoo.GetMorphologicalAnalysisResponse{}
// 	for _, r := range readings {
// 		res.MaResult.WordList.Words = append(res.MaResult.WordList.Words, yahoo.Word{Pos: "hoge", Reading: r, Surface: "fuga"})
// 	}
// 	return res
// }

// type ApiClientMock struct {
// }

// func (a *ApiClientMock) GetMorphologicalAnalysis(str string) (yahoo.GetMorphologicalAnalysisResponse, error) {
// 	switch str {
// 	case "白馬47":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"はくば", "47"}), nil
// 		}
// 	case "妙高杉ノ原":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"みょうこう", "すぎのはら"}), nil
// 		}
// 	case "高鷲スノーパーク":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"たかす", "すのー", "ぱーく"}), nil
// 		}
// 	case "高鷲SP":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"たかす", "SP"}), nil
// 		}
// 	case "今庄":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"いまじょう"}), nil
// 		}
// 	case "GALA湯沢":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"GALA", "ゆざわ"}), nil
// 		}
// 	case "石打丸山":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"いしうち", "まるやま"}), nil
// 		}
// 	case "赤倉観光":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"あかくら", "かんこう"}), nil
// 		}
// 	case "hakuba47":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"hakuba", "47"}), nil
// 		}
// 	case "hakuba 47":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"hakuba", " ", "47"}), nil
// 		}
// 	case "ニセコ":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"にせこ"}), nil
// 		}
// 	case "myokosuginohara":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"myokosuginohara"}), nil
// 		}
// 	case "hak妙47uba高":
// 		{
// 			return createGetMorphologicalAnalysisResponse([]string{"hak", "みょう", "47", "uba", "たか"}), nil
// 		}
// 	}
// 	return createGetMorphologicalAnalysisResponse([]string{str}), nil
// }

// func TestToHiragana(t *testing.T) {
// 	cases := []struct {
// 		kanji    string
// 		hiragana string
// 	}{
// 		{kanji: "白馬47", hiragana: "はくば47"},
// 		{kanji: "妙高杉ノ原", hiragana: "みょうこうすぎのはら"},
// 		{kanji: "高鷲スノーパーク", hiragana: "たかすすのーぱーく"},
// 		{kanji: "GALA湯沢", hiragana: "GALAゆざわ"},
// 		{kanji: "hakuba47", hiragana: "hakuba47"},
// 		{kanji: "hakuba 47", hiragana: "hakuba 47"},
// 		{kanji: "myokosuginohara", hiragana: "myokosuginohara"},
// 		{kanji: "hak妙47uba高", hiragana: "hakみょう47ubaたか"},
// 	}

// 	for _, tt := range cases {
// 		act, err := toHiragana(tt.kanji, &ApiClientMock{})
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if act != tt.hiragana {
// 			t.Error(fmt.Sprintf("%s is not %s", act, tt.kanji))
// 		}
// 	}
// }
