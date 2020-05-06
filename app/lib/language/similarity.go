package language

import (
	"strings"

	"github.com/kotaroooo0/snowforecast-twitter-bot/repository"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func GetSimilarSkiResort(target string) repository.SnowResort {
	client, err := repository.New("localhost:6379")
	if err != nil {
		panic(err)
	}
	r := repository.SnowResortRepositoryImpl{
		Client: client,
	}

	// 検索ワードとラベルで表記が大きく異なるスキー場も存在するので二つを合わせる
	lowercaseSnowResorts, err := r.ListSnowResorts("lowercase-snowresorts-searchword")
	snowResortLavels, err := r.ListSnowResorts("lowercase-snowresorts-label")
	sources := append(lowercaseSnowResorts, snowResortLavels...)

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
	distances := make([]int, len(sources))
	for i := 0; i < len(sources); i++ {
		distances[i] = levenshtein.DistanceForStrings([]rune(sources[i]), []rune(target), myOptions)
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

	sr, err := r.FindSnowResort(sources[minIdx])
	if err != nil {
		panic(err)
	}
	return sr
}
