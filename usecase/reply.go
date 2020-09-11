package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
)

type ReplyUseCase interface {
	NewGetTwitterWebhookRequest() GetTwitterWebhookRequest
	NewPostTwitterWebhookRequest() PostTwitterWebhookRequest
	GetCrcTokenResponse(GetTwitterWebhookRequest) GetTwitterWebhookResponse
	PostAutoReplyResponse(PostTwitterWebhookRequest) PostTwitterWebhookResponse
}

type ReplyUseCaseImpl struct {
	ReplyService domain.ReplyService
}

func NewReplyUseCaseImpl(rs domain.ReplyService) ReplyUseCase {
	return &ReplyUseCaseImpl{
		ReplyService: rs,
	}
}

func (tu ReplyUseCaseImpl) NewGetTwitterWebhookRequest() GetTwitterWebhookRequest {
	return GetTwitterWebhookRequest{}
}

// TwitterのWebhookの認証に用いる
// ref: https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks
func (tu ReplyUseCaseImpl) GetCrcTokenResponse(req GetTwitterWebhookRequest) GetTwitterWebhookResponse {
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(req.CrcToken))
	return GetTwitterWebhookResponse{
		Token: "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	}
}

func (tu ReplyUseCaseImpl) NewPostTwitterWebhookRequest() PostTwitterWebhookRequest {
	return PostTwitterWebhookRequest{}
}

func (tu ReplyUseCaseImpl) PostAutoReplyResponse(req PostTwitterWebhookRequest) PostTwitterWebhookResponse {
	// リプライがない、もしくはユーザが不正な場合は空を返す
	if len(req.TweetCreateEvents) < 1 || req.UserID == req.TweetCreateEvents[0].User.IDStr {
		return PostTwitterWebhookResponse{}
	}
	tweet := domain.Tweet{
		ID:             req.TweetCreateEvents[0].TweetIDStr,
		Text:           req.TweetCreateEvents[0].Text,
		UserScreenName: req.TweetCreateEvents[0].User.ScreenName,
	}

	// リプライから全世界のスキー場の中で最も適切なスキー場を求める
	sr, err := tu.ReplyService.ReplyForecast(&tweet)
	if err != nil {
		return PostTwitterWebhookResponse{}
	}
	return PostTwitterWebhookResponse{sr.Name}
}

type GetTwitterWebhookRequest struct {
	CrcToken string `json:"crc_token" form:"crc_token" binding:"required"`
}

type GetTwitterWebhookResponse struct {
	Token string `json:"response_token"`
}

type PostTwitterWebhookRequest struct {
	UserID            string             `json:"for_user_id" form:"for_user_id" binding:"required"`
	TweetCreateEvents []TweetCreateEvent `json:"tweet_create_events" form:"tweet_create_events" binding:"required"`
}

type TweetCreateEvent struct {
	TweetID    int64  `json:"id" form:"id" binding:"required"`
	TweetIDStr string `json:"id_str" form:"id_str" binding:"required"`
	Text       string `json:"text" form:"text" binding:"required"`
	User       struct {
		UserID     int64  `json:"id" form:"id" binding:"required"`
		IDStr      string `json:"id_str" form:"id_str" binding:"required"`
		ScreenName string `json:"screen_name" form:"screen_name" binding:"required"`
	} `json:"user" form:"user" binding:"required"`
}

type PostTwitterWebhookResponse struct {
	SnowResortLabel string `json:"snow_resort"`
}
