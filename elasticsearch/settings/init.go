package elasticsearch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
)

func main() {
	// elasticsearchのデータ初期化用のテキストファイル
	file, err := os.Create("snow_resorts.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// スキー場の地域一覧
	regions := []string{
		"exotic", "summer", "22-9", "1", "2", "6", "7", "8", "9", "10", "11", "12", "13", "15", "16", "17", "18", "19", "20", "22", "25", "28", "29",
		"30", "31", "32", "33", "34", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
		"50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "60", "61", "62", "63", "64", "67", "71", "72", "73", "76", "78", "79",
		"80", "81", "84", "86", "87", "88", "90", "91", "92", "93", "96", "97", "98", "99",
		"100", "101", "102", "103", "104", "109", "110", "113", "116", "117", "119", "120", "123", "125", "126", "130", "131", "132", "133", "134", "135", "136", "137", "139",
		"142", "143", "144", "145", "146", "147", "148", "149", "151", "152", "153", "155", "156", "157", "158", "159", "160", "161", "162", "163", "164", "165", "166", "167", "168", "169",
		"171", "172", "173", "175", "177", "178", "179", "180", "181", "182", "184", "197", "199",
		"200", "202", "203", "204", "205", "206", "207", "208", "209", "210", "211", "212", "213", "214", "215", "216", "217", "218", "219", "220", "221", "225",
	}

	// リクエスト先のサーバに負荷がかかりすぎないようにへ並行処理数を30までにする
	m := new(sync.Mutex)
	ch := make(chan int, 30)
	wg := sync.WaitGroup{}
	var id int
	for _, r := range regions {
		ch <- 1
		wg.Add(1)
		go func(r string) {
			// 地域からスキー場を取得するリクエスト
			res, err := http.Get("https://ja.snow-forecast.com/resorts/list_by_feature/" + r + "?v=2")
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			snowResorts := parseStringToSnowResorts(string(body))
			for i := 0; i < len(snowResorts); i++ {
				m.Lock()
				// {"index": {"_index":"snow_resorts_v1", "_type":"_doc", "_id":"1"}}
				// {"title": "title0", "name":"name0", "age":10, "created":"2019-08-01"}
				// TODO: インデックス名のハードコーディング🙅‍♂️
				file.WriteString(fmt.Sprintf("{\"index\":{\"_index\":\"%s\", \"_type\":\"_doc\", \"_id\":\"%d\"}}\n", "snow_resorts_v1", id))
				file.WriteString(fmt.Sprintf("{\"name\":\"%s\",\"search_key\":\"%s\"}\n", snowResorts[i].Name, snowResorts[i].SearchKey))
				id++
				m.Unlock()
			}
			<-ch
			wg.Done()
		}(r)
	}
	wg.Wait()
}

// stringの配列から[]Snowresortを生成する
// example
// input: "[[\"\",[[\"HiddenValley2\",\"Hidden Valley Ski\"],[\"Snow-Creek\",\"Snow Creek\"]]]]"
// output: []SnowResort{
// 	SnowResort{"HiddenValley2", "Hidden Valley Ski"},
// 	SnowResort{"Snow-Creek", "Snow Creek"},
// }
// TODO: ASTつくってスマートにparseするライブラリつくりたい
func parseStringToSnowResorts(str string) []domain.SnowResort {
	isTarget := false
	word := ""
	stringSlice := []string{}
	for i := 0; i < len([]rune(str)); i++ {
		if string(str[i]) == "\"" {
			if isTarget && word != "" {
				stringSlice = append(stringSlice, word)
			}
			isTarget = !isTarget
			word = ""
		}
		if isTarget && string(str[i]) != "\"" {
			word += string(str[i])
		}
	}
	snowResorts := []domain.SnowResort{}
	for i := 1; i < len(stringSlice); i += 2 {
		snowResorts = append(snowResorts, domain.SnowResort{
			Name:      stringSlice[i-1],
			SearchKey: stringSlice[i],
		})
	}
	return snowResorts
}
