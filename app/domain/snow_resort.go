package domain

type SnowResort struct {
	SearchWord string
	Label      string
}

type SnowResortService interface {
	ReplyForecast(SnowResort) (SnowResort, error)
}

type SnowResortServiceImpl struct {
	// ドメイン層はどこにも依存しない
}

func (ss SnowResortServiceImpl) ReplyForecast(snowResort SnowResort) (SnowResort, error) {
	// 自動でリプライを返す
	// api := twitter.GetTwitterApi()
	// params := url.Values{}
	// params.Set("in_reply_to_status_id", req.TweetCreateEvents[0].TweetIDStr)
	// _, err := api.PostTweet("@"+req.TweetCreateEvents[0].User.ScreenName+" Hello World", params)
	// if err != nil {
	//      ctx.JSON(http.StatusBadRequest, err)
	// } else {
	//      ctx.Status(200)
	// }
	return SnowResort{}, nil
}

type SnowResortRepository interface {
	ListSnowResorts(string) ([]string, error)
	FindSnowResort(string) (SnowResort, error)
}
