package snowforecast

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type ISnowforecastApiClient interface {
	GetSnowfallForecastBySkiResortSearchWord(string) (SnowfallForecast, error)
}

type SnowforecastApiClient struct{}

func NewSnowforecastApiClient() ISnowforecastApiClient {
	return &SnowforecastApiClient{}
}

// 以下の3パターンの予報が取得できる
// 1.本日の朝からの予報が見れる時
// 2.本日の昼からの予報が見れる時
// 3.本日の夜からの予報が見れる時
func (sc SnowforecastApiClient) GetSnowfallForecastBySkiResortSearchWord(skiResortSearchWord string) (SnowfallForecast, error) {
	doc, err := goquery.NewDocument("https://ja.snow-forecast.com/resorts/" + skiResortSearchWord + "/6day/top")
	if err != nil {
		return SnowfallForecast{}, err
	}

	snowfalls := make([]Snow, 0)
	forecastTableSnow := doc.Find("td.forecast-table-snow__cell")
	forecastTableSnow.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			if index == 0 {
				// 朝と昼の情報が取得できない時
				nightSnowfall := selectionToInt(s)
				snowfalls = append(snowfalls, NewSnow(0, 0, nightSnowfall))
			} else if index == 1 {
				// 朝の情報が取得できない時
				noonSnowfall := selectionToInt(forecastTableSnow.Eq(index - 1))
				nightSnowfall := selectionToInt(s)
				snowfalls = append(snowfalls, NewSnow(0, noonSnowfall, nightSnowfall))
			} else {
				// 朝昼晩の情報が取得できる時
				morningSnowfall := selectionToInt(forecastTableSnow.Eq(index - 2))
				noonSnowfall := selectionToInt(forecastTableSnow.Eq(index - 1))
				nightSnowfall := selectionToInt(s)
				snowfalls = append(snowfalls, NewSnow(morningSnowfall, noonSnowfall, nightSnowfall))
			}
		}
	})

	rainfalls := make([]Rain, 0)
	forecastTableRain := doc.Find("td.forecast-table-rain__cell")
	forecastTableRain.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			if index == 0 {
				// 朝と昼の情報が取得できない時
				nightRainfall := selectionToInt(s)
				rainfalls = append(rainfalls, NewRain(0, 0, nightRainfall))
			} else if index == 1 {
				// 朝の情報が取得できない時
				noonRainfall := selectionToInt(forecastTableRain.Eq(index - 1))
				nightRainfall := selectionToInt(s)
				rainfalls = append(rainfalls, NewRain(0, noonRainfall, nightRainfall))
			} else {
				// 朝昼晩の情報が取得できる時
				morningRainfall := selectionToInt(forecastTableRain.Eq(index - 2))
				noonRainfall := selectionToInt(forecastTableRain.Eq(index - 1))
				nightRainfall := selectionToInt(s)
				rainfalls = append(rainfalls, NewRain(morningRainfall, noonRainfall, nightRainfall))
			}
		}
	})

	return NewSnowfallForecast(snowfalls, rainfalls, skiResortSearchWord), nil
}

func selectionToInt(s *goquery.Selection) int {
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

type SnowfallForecast struct {
	Snows     []Snow
	Rains     []Rain
	SkiResort string
}

func NewSnowfallForecast(snows []Snow, rains []Rain, skiResort string) SnowfallForecast {
	return SnowfallForecast{
		Snows:     snows,
		Rains:     rains,
		SkiResort: skiResort,
	}
}

type Snow struct {
	Morning int
	Noon    int
	Night   int
}

func NewSnow(morning, noon, night int) Snow {
	return Snow{
		Morning: morning,
		Noon:    noon,
		Night:   night,
	}
}

type Rain struct {
	Morning int
	Noon    int
	Night   int
}

func NewRain(morning, noon, night int) Rain {
	return Rain{
		Morning: morning,
		Noon:    noon,
		Night:   night,
	}
}
