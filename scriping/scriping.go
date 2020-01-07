package scriping

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type SnowfallForecast struct {
	Snowfalls []*Snowfall
	Rainfalls []*Rainfall
	SkiResort string
}

func NewSnowfallForecast(snowfalls []*Snowfall, rainfalls []*Rainfall, skiResort string) *SnowfallForecast {
	return &SnowfallForecast{
		Snowfalls: snowfalls,
		Rainfalls: rainfalls,
		SkiResort: skiResort,
	}
}

type Snowfall struct {
	DaySnowfall   int
	NightSnowfall int
}

func NewSnowfall(daySnowfall, nightSnowfall int) *Snowfall {
	return &Snowfall{
		DaySnowfall:   daySnowfall,
		NightSnowfall: nightSnowfall,
	}
}

type Rainfall struct {
	DayRainfall   int
	NightRainfall int
}

func NewRainfall(dayRainfall, nightRainfall int) *Rainfall {
	return &Rainfall{
		DayRainfall:   dayRainfall,
		NightRainfall: nightRainfall,
	}
}

// 以下の3パターンの予報が取得できる
// 1.本日の朝からの予報が見れる時
// 2.本日の昼からの予報が見れる時
// 3.本日の夜からの予報が見れる時
func GetSnowfallForecastBySkiResort(skiResort string) *SnowfallForecast {
	doc, err := goquery.NewDocument("https://ja.snow-forecast.com/resorts/" + skiResort + "/6day/mid")
	if err != nil {
		panic(err)
	}

	snowfalls := make([]*Snowfall, 0)
	forecastTableSnow := doc.Find("td.forecast-table-snow__cell")
	daySnowfall := 0
	nightSnowfall := 0
	forecastTableSnow.Each(func(index int, s *goquery.Selection) {
		snowfall := s.Text()
		if snowfall == "-" {
			snowfall = "0"
		}
		snowfallInt, err := strconv.Atoi(snowfall)
		if err != nil {
			panic(err)
		}
		if s.HasClass("day-end") {
			nightSnowfall = snowfallInt
			snowfalls = append(snowfalls, NewSnowfall(daySnowfall, nightSnowfall))
			daySnowfall = 0
			nightSnowfall = 0
		} else {
			daySnowfall += snowfallInt
		}
	})

	rainfalls := make([]*Rainfall, 0)
	forecastTableRain := doc.Find("td.forecast-table-rain__cell")
	dayRainfall := 0
	nightRainfall := 0
	forecastTableRain.Each(func(index int, s *goquery.Selection) {
		rainfall := s.Text()
		if rainfall == "-" {
			rainfall = "0"
		}
		rainfallInt, err := strconv.Atoi(rainfall)
		if err != nil {
			panic(err)
		}
		if s.HasClass("day-end") {
			nightRainfall = rainfallInt
			rainfalls = append(rainfalls, NewRainfall(dayRainfall, nightRainfall))
			dayRainfall = 0
			nightRainfall = 0
		} else {
			dayRainfall += rainfallInt
		}
	})

	return NewSnowfallForecast(snowfalls, rainfalls, skiResort)
}
