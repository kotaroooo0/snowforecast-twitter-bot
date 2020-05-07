package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
)

type TwitterUseCase interface {
	NewGetTwitterWebhookRequest() GetTwitterWebhookRequest
	NewPostTwitterWebhookRequest() PostTwitterWebhookRequest
	GetCrcTokenResponse(GetTwitterWebhookRequest) GetTwitterWebhookResponse
}

type TwitterUseCaseImpl struct {
	// SnowResortRepository repository.SnowResortRepository
}

func (tu TwitterUseCaseImpl) NewGetTwitterWebhookRequest() GetTwitterWebhookRequest {
	return GetTwitterWebhookRequest{}
}

// TwitterのWebhookの認証に用いる
// ref: https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks
func (tu TwitterUseCaseImpl) GetCrcTokenResponse(req GetTwitterWebhookRequest) GetTwitterWebhookResponse {
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(req.CrcToken))
	return GetTwitterWebhookResponse{
		Token: "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	}
}

func (tu TwitterUseCaseImpl) NewPostTwitterWebhookRequest() PostTwitterWebhookRequest {
	return PostTwitterWebhookRequest{}
}

type GetTwitterWebhookRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

type GetTwitterWebhookResponse struct {
	Token string `json:"response_token"`
}

type PostTwitterWebhookRequest struct {
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

type PostTwitterWebhookResponse struct {
}
