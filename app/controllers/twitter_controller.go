package controllers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kotaroooo0/snowforecast-twitter-bot/lib/twitter"
	"github.com/kotaroooo0/snowforecast-twitter-bot/parameters/requests"
	"github.com/kotaroooo0/snowforecast-twitter-bot/parameters/responses"
)

type TwitterController struct {
	engine *gin.Engine
}

func NewTwitterController(engine *gin.Engine) *TwitterController {
	return &TwitterController{engine}
}

func (c *TwitterController) GetCrcToken() {
	c.engine.GET("/twitter_webhook", c.getCrcToken)
}

func (c *TwitterController) getCrcToken(ctx *gin.Context) {
	req := requests.NewGetTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	res := responses.NewGetTwitterWebhookCrcCheckResponse()
	res.Token = twitter.CreateCRCToken(req.CrcToken)
	ctx.JSON(http.StatusOK, res)
}

func (c *TwitterController) PostWebhook() {
	c.engine.POST("/twitter_webhook", c.postWebhook)
}

func (c *TwitterController) postWebhook(ctx *gin.Context) {
	req := requests.NewPostTwitterWebHookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	if len(req.TweetCreateEvents) < 1 || req.UserID == req.TweetCreateEvents[0].User.IDStr {
		return
	}

	// 自動でリプライを返す
	api := twitter.GetTwitterApi()
	params := url.Values{}
	params.Set("in_reply_to_status_id", req.TweetCreateEvents[0].TweetIDStr)
	_, err := api.PostTweet("@"+req.TweetCreateEvents[0].User.ScreenName+" Hello World", params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	} else {
		ctx.Status(200)
	}
}
