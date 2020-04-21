package text

import (
	"strconv"

	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/scriping"
)

func TweetContent(skiResort1, skiResort2 string) string {
	toLabel := map[string]string{
		"Niseko":                 "ニセコ",
		"SapporoKokusai":         "札幌国際",
		"Hakuba47":               "白馬47",
		"MyokoSuginohara":        "赤倉",
		"TashiroKaguraMitsumata": "かぐら",
		"IshiuchiMaruyama":       "石打丸山",
		"MarunumaKogen":          "丸沼高原",
		"TakasuSnowPark":         "高鷲",
		"BiwakoValley":           "琵琶湖バレイ",
	}

	data1 := scriping.GetSnowfallForecastBySkiResort(skiResort1)
	data2 := scriping.GetSnowfallForecastBySkiResort(skiResort2)

	content := "今日 | 明日 | 明後日 (朝,昼,夜)\n"
	content += toLabel[skiResort1] + "\n"
	content += AreaLineString(data1) + "\n"
	content += toLabel[skiResort2] + "\n"
	content += AreaLineString(data2) + "\n"
	return content
}

func AreaLineString(snowfallForecast *scriping.SnowfallForecast) string {
	content := strconv.Itoa(snowfallForecast.Snowfalls[0].MorningSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[0].MorningRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[0].NoonSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[0].NoonRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[0].NightSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[0].NightRainfall) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snowfalls[1].MorningSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[1].MorningRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[1].NoonSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[1].NoonRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[1].NightSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[1].NightRainfall) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snowfalls[2].MorningSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[2].MorningRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[2].NoonSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[2].NoonRainfall) + ", " + strconv.Itoa(snowfallForecast.Snowfalls[2].NightSnowfall) + AddRainyChar(snowfallForecast.Rainfalls[2].NightRainfall) + "cm "
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
