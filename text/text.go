package text

import (
	"fmt"
	"strconv"

	"github.com/kotaroooo0/snowforecast-twitter-bot/scriping"
)

func TweetContent(skiResort1, skiResort2 string) string {
	toLabel := map[string]string{
		"Hakuba47":         "白馬",
		"MyokoSuginohara":  "妙高",
		"IshiuchiMaruyama": "湯沢",
		"TakasuSnowPark":   "高鷲",
	}

	data1 := scriping.GetSnowfallForecastBySkiResort(skiResort1)
	data2 := scriping.GetSnowfallForecastBySkiResort(skiResort2)

	content := "今日 | 明日 | 明後日 (昼,夜)\n"
	content += toLabel[skiResort1] + "\n"
	content += AreaLineString(data1) + "\n"
	content += toLabel[skiResort2] + "\n"
	content += AreaLineString(data2) + "\n"
	return content
}

func AreaLineString(snowfallForecast *scriping.SnowfallForecast) string {
	fmt.Println(snowfallForecast.Snowfalls[0].DaySnowfall)
	fmt.Println(strconv.Itoa(snowfallForecast.Snowfalls[0].DaySnowfall))
	content := strconv.Itoa(snowfallForecast.Snowfalls[0].DaySnowfall) + "cm" + AddRainyChar(snowfallForecast.Rainfalls[0].DayRainfall) + " " + strconv.Itoa(snowfallForecast.Snowfalls[0].NightSnowfall) + "cm" + AddRainyChar(snowfallForecast.Rainfalls[0].NightRainfall) + " " + "|"
	content += strconv.Itoa(snowfallForecast.Snowfalls[1].DaySnowfall) + "cm" + AddRainyChar(snowfallForecast.Rainfalls[1].DayRainfall) + " " + strconv.Itoa(snowfallForecast.Snowfalls[1].NightSnowfall) + "cm" + AddRainyChar(snowfallForecast.Rainfalls[1].NightRainfall) + " " + "|"
	content += strconv.Itoa(snowfallForecast.Snowfalls[2].DaySnowfall) + "cm" + AddRainyChar(snowfallForecast.Rainfalls[2].DayRainfall) + " " + strconv.Itoa(snowfallForecast.Snowfalls[2].NightSnowfall) + "cm" + AddRainyChar(snowfallForecast.Rainfalls[2].NightRainfall) + " "
	return content
}

func AddRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	} else if rainfall > 0 {
		return "☂️"
	}
	return ""
}
