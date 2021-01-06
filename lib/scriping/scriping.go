package scriping

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type SnowfallForecast struct {
	Snowfalls []Snowfall
	Rainfalls []Rainfall
	SkiResort string
}

func NewSnowfallForecast(snowfalls []Snowfall, rainfalls []Rainfall, skiResort string) *SnowfallForecast {
	return &SnowfallForecast{
		Snowfalls: snowfalls,
		Rainfalls: rainfalls,
		SkiResort: skiResort,
	}
}

type Snowfall struct {
	MorningSnowfall int
	NoonSnowfall    int
	NightSnowfall   int
}

func NewSnowfall(morningSnowfall, noonSnowfall, nightSnowfall int) Snowfall {
	return Snowfall{
		MorningSnowfall: morningSnowfall,
		NoonSnowfall:    noonSnowfall,
		NightSnowfall:   nightSnowfall,
	}
}

type Rainfall struct {
	MorningRainfall int
	NoonRainfall    int
	NightRainfall   int
}

func NewRainfall(morningRainfall, noonRainfall, nightRainfall int) Rainfall {
	return Rainfall{
		MorningRainfall: morningRainfall,
		NoonRainfall:    noonRainfall,
		NightRainfall:   nightRainfall,
	}
}

// 以下の3パターンの予報が取得できる
// 1.本日の朝からの予報が見れる時
// 2.本日の昼からの予報が見れる時
// 3.本日の夜からの予報が見れる時
func GetSnowfallForecastBySkiResort(snowResort string) (*SnowfallForecast, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://ja.snow-forecast.com/resorts/%s/6day/top", snowResort))
	if err != nil {
		return nil, err
	}

	var convertErrFlag bool
	snowfalls := make([]Snowfall, 0)
	forecastTableSnow := doc.Find("td.forecast-table-snow__cell")
	forecastTableSnow.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			morning, noon := -1, -1
			// 夜の情報は必ず取得できる
			night := convertInt(s.Text(), &convertErrFlag)
			// 昼の情報も取得できる時
			if index > 0 {
				noon = convertInt(forecastTableSnow.Eq(index-1).Text(), &convertErrFlag)
			}
			// 朝の情報も取得できる時
			if index > 1 {
				morning = convertInt(forecastTableSnow.Eq(index-2).Text(), &convertErrFlag)
			}
			snowfalls = append(snowfalls, NewSnowfall(morning, noon, night))
		}
	})
	if convertErrFlag {
		return nil, fmt.Errorf("error: convert error occur")
	}

	rainfalls := make([]Rainfall, 0)
	forecastTableRain := doc.Find("td.forecast-table-rain__cell")
	forecastTableRain.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			morning, noon := -1, -1
			// 夜の情報は必ず取得できる
			night := convertInt(s.Text(), &convertErrFlag)
			// 昼の情報も取得できる時
			if index > 0 {
				noon = convertInt(forecastTableRain.Eq(index-1).Text(), &convertErrFlag)
			}
			// 朝の情報も取得できる時
			if index > 1 {
				morning = convertInt(forecastTableRain.Eq(index-2).Text(), &convertErrFlag)
			}
			rainfalls = append(rainfalls, NewRainfall(morning, noon, night))
		}
	})
	if convertErrFlag {
		return nil, fmt.Errorf("error: convert error occur")
	}

	return NewSnowfallForecast(snowfalls, rainfalls, snowResort), nil
}

// "-"の場合は特別に0を返すキャスト
func convertInt(s string, convertErrFlag *bool) int {
	if s == "-" {
		return 0
	}
	ret, err := strconv.Atoi(s)
	if err != nil {
		*convertErrFlag = true
	}
	return ret
}
