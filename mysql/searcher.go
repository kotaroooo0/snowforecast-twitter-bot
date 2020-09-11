package searcher

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kotaroooo0/gojaconv/jaconv"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/yahoo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type SnowResortSearcherMySQLImpl struct {
	SnowResortRepository domain.SnowResortRepository
	YahooApiClient       yahoo.IYahooApiClient
}

func (s SnowResortSearcherMysqlImpl) FindSimilarSnowResort(source string) (*domain.SnowResort, error) {
	// スペースを消す
	source = strings.Replace(source, " ", "", -1)
	key := strings.Replace(source, "　", "", -1)

	// 漢字をひらがなに変換(ex:GALA湯沢 -> GALAゆざわ)
	source, err := toHiragana(key, s.YahooApiClient)
	if err != nil {
		return &domain.SnowResort{}, err
	}
	// ひらがなをアルファベットに変換(ex:GALAゆざわ -> GALAyuzawa)
	source = jaconv.ToHebon(source)

	srs, err := s.SnowResortRepository.FindAll()
	if err != nil {
		return &domain.SnowResort{}, err
	}
	sr := findSimilarSnowResort(source, srs)
	return sr, nil
}

// TODO: メソッドが大きすぎるので分割してもいかも
// Lower Filterしてからレーベンシュタイン距離により類似単語検索する
func findSimilarSnowResort(source string, targets []*domain.SnowResort) *domain.SnowResort {
	// sourceを小文字に直す(ex:GALAyuzawa -> galayuzawa)
	source = strings.ToLower(source)

	// targetsもNameとSearchKeysの両方を小文字へ直す
	targetNames := make([]string, len(targets))
	for i := 0; i < len(targetNames); i++ {
		targetNames[i] = strings.ToLower(targets[i].Name)
	}
	targetSearchKeys := make([]string, len(targets))
	for i := 0; i < len(targetSearchKeys); i++ {
		targetSearchKeys[i] = strings.ToLower(targets[i].SearchKey)
	}
	allTargets := append(targetNames, targetSearchKeys...)

	// レーベンシュタイン距離を計算する際の重みづけ
	// 削除の際の距離を小さくしている
	// TODO: 標準化や他の編集距離を考える必要もある、その際に評価をするためにラベル付おじさんにならないといけないかもしれない
	myOptions := levenshtein.Options{
		InsCost: 10,
		DelCost: 1,
		SubCost: 10,
		Matches: levenshtein.IdenticalRunes,
	}
	distances := make([]int, len(allTargets))
	for i := 0; i < len(allTargets); i++ {
		distances[i] = levenshtein.DistanceForStrings([]rune(source), []rune(allTargets[i]), myOptions)
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

	// NameとSearchKeyを足して二倍になってるので、余りをとる
	return targets[minIdx%len(targets)]
}

func toHiragana(str string, yahooApiClient yahoo.IYahooApiClient) (string, error) {
	res, err := yahooApiClient.GetMorphologicalAnalysis(str)
	if err != nil {
		return "", err
	}

	h := ""
	for _, w := range res.MaResult.WordList.Words {
		h += w.Reading
	}
	return h, nil
}

type SnowResortRepositoryImpl struct {
	DB *sqlx.DB
}

func NewSnowResortRepositoryImpl(db *sqlx.DB) domain.SnowResortRepository {
	return &SnowResortRepositoryImpl{
		DB: db,
	}
}

func NewDBClient(dbConfig *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Password, dbConfig.Addr, dbConfig.Port, dbConfig.DB),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDBConfig(user, password, addr, port, db string) *DBConfig {
	return &DBConfig{
		User:     user,
		Password: password,
		Addr:     addr,
		Port:     port,
		DB:       db,
	}
}

type DBConfig struct {
	User     string
	Password string
	Addr     string
	Port     string
	DB       string
}

func (s SnowResortRepositoryImpl) FindAll() ([]*domain.SnowResort, error) {
	rows, err := s.DB.Queryx("select * from snow_resorts")
	if err != nil {
		return []*domain.SnowResort{}, err
	}

	var snowResorts []*domain.SnowResort
	for rows.Next() {
		var snowResort domain.SnowResort
		err := rows.StructScan(&snowResort)
		if err != nil {
			return []*domain.SnowResort{}, err
		}
		snowResorts = append(snowResorts, &snowResort)
	}
	return snowResorts, nil
}
