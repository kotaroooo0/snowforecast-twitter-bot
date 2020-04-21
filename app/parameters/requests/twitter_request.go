package requests

// GET Request (CRC Check)
type GetTwitterWebhookRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

func NewGetTwitterWebhookRequest() GetTwitterWebhookRequest {
	return GetTwitterWebhookRequest{}
}

// POST Request (Event)
type PostTwitterWebHookRequest struct {
	UserID            string `json:"for_user_id" form:"for_user_id" binding:"required"`
	TweetCreateEvents []struct {
		TweetID    int64  `json:"id" form:"id" binding:"required"`
		TweetIDStr string `json:"id_str" form:"id_str" binding:"required"`
		User       struct {
			UserID     int64  `json:"id" form:"id" binding:"required"`
			IDStr      string `json:"id_str" form:"id_str" binding:"required"`
			ScreenName string `json:"screen_name" form:"screen_name" binding:"required"`
		} `json:"user" form:"user" binding:"required"`
		Text string
	} `json:"tweet_create_events" form:"tweet_create_events" binding:"required"`
}

func NewPostTwitterWebHookRequest() PostTwitterWebHookRequest {
	return PostTwitterWebHookRequest{}
}
