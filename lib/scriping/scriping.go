package scriping

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type SnowfallForecast struct {
	Snowfalls []Snowfall
	Rainfalls []Rainfall
	SkiResort string
}

func NewSnowfallForecast(snowfalls []Snowfall, rainfalls []Rainfall, skiResort string) SnowfallForecast {
	return SnowfallForecast{
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
func GetSnowfallForecastBySkiResort(skiResort string) (SnowfallForecast, error) {
	doc, err := goquery.NewDocument("https://ja.snow-forecast.com/resorts/" + skiResort + "/6day/top")
	if err != nil {
		panic(err)
	}

	snowfalls := make([]Snowfall, 0)
	forecastTableSnow := doc.Find("td.forecast-table-snow__cell")
	forecastTableSnow.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			if index == 0 {
				// 朝と昼の情報が取得できない時
				nightSnowfall := SelectionToInt(s)
				snowfalls = append(snowfalls, NewSnowfall(0, 0, nightSnowfall))
			} else if index == 1 {
				// 朝の情報が取得できない時
				noonSnowfall := SelectionToInt(forecastTableSnow.Eq(index - 1))
				nightSnowfall := SelectionToInt(s)
				snowfalls = append(snowfalls, NewSnowfall(0, noonSnowfall, nightSnowfall))
			} else {
				// 朝昼晩の情報が取得できる時
				morningSnowfall := SelectionToInt(forecastTableSnow.Eq(index - 2))
				noonSnowfall := SelectionToInt(forecastTableSnow.Eq(index - 1))
				nightSnowfall := SelectionToInt(s)
				snowfalls = append(snowfalls, NewSnowfall(morningSnowfall, noonSnowfall, nightSnowfall))
			}
		}
	})

	rainfalls := make([]Rainfall, 0)
	forecastTableRain := doc.Find("td.forecast-table-rain__cell")
	forecastTableRain.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			if index == 0 {
				// 朝と昼の情報が取得できない時
				nightRainfall := SelectionToInt(s)
				rainfalls = append(rainfalls, NewRainfall(0, 0, nightRainfall))
			} else if index == 1 {
				// 朝の情報が取得できない時
				noonRainfall := SelectionToInt(forecastTableRain.Eq(index - 1))
				nightRainfall := SelectionToInt(s)
				rainfalls = append(rainfalls, NewRainfall(0, noonRainfall, nightRainfall))
			} else {
				// 朝昼晩の情報が取得できる時
				morningRainfall := SelectionToInt(forecastTableRain.Eq(index - 2))
				noonRainfall := SelectionToInt(forecastTableRain.Eq(index - 1))
				nightRainfall := SelectionToInt(s)
				rainfalls = append(rainfalls, NewRainfall(morningRainfall, noonRainfall, nightRainfall))
			}
		}
	})

	return NewSnowfallForecast(snowfalls, rainfalls, skiResort)
}

func SelectionToInt(s *goquery.Selection) int {
	fall := s.Text()
	if fall == "-" {
		fall = "0"
	}
	fallInt, err := strconv.Atoi(fall)
	if err != nil {
		panic(err)
	}
	return fallInt
}
