package text

import (
	"strconv"

	"github.com/kotaroooo0/snowforecast-twitter-bot/batch"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/scriping"
)

func TweetContent(pair batch.Pair) (string, error) {
	firstData, err := scriping.GetSnowfallForecastBySkiResort(pair.First)
	if err != nil {
		return "", err
	}
	secondData, err := scriping.GetSnowfallForecastBySkiResort(pair.Second)
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

func areaLineString(snowfallForecast *scriping.SnowfallForecast) string {
	content := strconv.Itoa(snowfallForecast.Snowfalls[0].MorningSnowfall) + addRainyChar(snowfallForecast.Rainfalls[0].MorningRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[0].NoonSnowfall) + addRainyChar(snowfallForecast.Rainfalls[0].NoonRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[0].NightSnowfall) + addRainyChar(snowfallForecast.Rainfalls[0].NightRainfall) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snowfalls[1].MorningSnowfall) + addRainyChar(snowfallForecast.Rainfalls[1].MorningRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[1].NoonSnowfall) + addRainyChar(snowfallForecast.Rainfalls[1].NoonRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[1].NightSnowfall) + addRainyChar(snowfallForecast.Rainfalls[1].NightRainfall) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snowfalls[2].MorningSnowfall) + addRainyChar(snowfallForecast.Rainfalls[2].MorningRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[2].NoonSnowfall) + addRainyChar(snowfallForecast.Rainfalls[2].NoonRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[2].NightSnowfall) + addRainyChar(snowfallForecast.Rainfalls[2].NightRainfall) + "cm "
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
