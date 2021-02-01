package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
)

type ReplyUseCase interface {
	GetCrcTokenResponse(GetTwitterWebhookRequest) (GetTwitterWebhookResponse, error)
	PostAutoReplyResponse(PostTwitterWebhookRequest) (PostTwitterWebhookResponse, error)
}

type ReplyUseCaseImpl struct {
	ReplyService domain.ReplyService
}

func NewReplyUseCaseImpl(rs domain.ReplyService) ReplyUseCase {
	return &ReplyUseCaseImpl{
		ReplyService: rs,
	}
}

func NewGetTwitterWebhookRequest() GetTwitterWebhookRequest {
	return GetTwitterWebhookRequest{}
}

// TwitterのWebhookの認証に用いる
// ref: https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks
func (tu ReplyUseCaseImpl) GetCrcTokenResponse(req GetTwitterWebhookRequest) (GetTwitterWebhookResponse, error) {
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	_, err := mac.Write([]byte(req.CrcToken))
	return GetTwitterWebhookResponse{
		Token: "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	}, err
}

func NewPostTwitterWebhookRequest() PostTwitterWebhookRequest {
	return PostTwitterWebhookRequest{}
}

func (tu ReplyUseCaseImpl) PostAutoReplyResponse(req PostTwitterWebhookRequest) (PostTwitterWebhookResponse, error) {
	// リプライがない、もしくはユーザが不正な場合は空を返す
	pp.Println(len(req.TweetCreateEvents))
	pp.Println(req.UserID)
	pp.Println(req.TweetCreateEvents[0].User.IDStr)
	if len(req.TweetCreateEvents) < 1 || req.UserID == req.TweetCreateEvents[0].User.IDStr {
		return PostTwitterWebhookResponse{}, fmt.Errorf("error: not found reply or invalid user")
	}
	tweet := domain.Tweet{
		ID:             req.TweetCreateEvents[0].TweetIDStr,
		Text:           req.TweetCreateEvents[0].Text,
		UserScreenName: req.TweetCreateEvents[0].User.ScreenName,
	}
	// リプライから全世界のスキー場の中で最も適切なスキー場を求める
	sr, err := tu.ReplyService.ReplyForecast(&tweet)
	return PostTwitterWebhookResponse{sr.Name}, err
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
