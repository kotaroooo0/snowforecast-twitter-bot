package text

import (
	"strconv"

	"github.com/kotaroooo0/snowforecast-twitter-bot/batch"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/snowforecast"
)

type TweetContentCreater struct {
	ApiClient snowforecast.IApiClient
}

func (c TweetContentCreater) TweetContent(pair batch.Pair) (string, error) {
	firstData, err := c.ApiClient.GetForecastBySearchWord(pair.First)
	if err != nil {
		return "", err
	}
	secondData, err := c.ApiClient.GetForecastBySearchWord(pair.Second)
	if err != nil {
		return "", err
	}
	content := "今日 | 明日 | 明後日 (朝,昼,夜)\n"
	content += pair.First + "\n"
	content += areaLineString(firstData) + "\n"
	content += pair.Second + "\n"
	content += areaLineString(secondData) + "\n"
	return content, nil
}

func areaLineString(snowfallForecast *snowforecast.Forecast) string {
	content := strconv.Itoa(snowfallForecast.Snows[0].Morning) + addRainyChar(snowfallForecast.Rains[0].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Noon) + addRainyChar(snowfallForecast.Rains[0].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Night) + addRainyChar(snowfallForecast.Rains[0].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[1].Morning) + addRainyChar(snowfallForecast.Rains[1].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Noon) + addRainyChar(snowfallForecast.Rains[1].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Night) + addRainyChar(snowfallForecast.Rains[1].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[2].Morning) + addRainyChar(snowfallForecast.Rains[2].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Noon) + addRainyChar(snowfallForecast.Rains[2].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Night) + addRainyChar(snowfallForecast.Rains[2].Night) + "cm "
	return content
}

func addRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	}
	if rainfall > 0 {
		return "☂️"
	}
	return ""
}
