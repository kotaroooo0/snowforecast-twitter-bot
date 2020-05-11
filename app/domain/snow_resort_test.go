package domain

// TODO:動作確認用に書いたが、モックを使ったしっかりしたテストを書くとありがたがられる
// func TestReplyContent(t *testing.T) {
// 	testClient := snowforecast.NewSnowforecastApiClient()

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
