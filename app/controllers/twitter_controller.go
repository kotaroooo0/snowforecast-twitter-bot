package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kotaroooo0/snowforecast-twitter-bot/key"
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
	// TODO: crc_tokenが渡されなかったら返す?

	fmt.Print(req)

	res := responses.NewGetTwitterWebhookCrcCheckResponse()
	res.Token = twitter.CreateCRCToken(req.CrcToken)
	ctx.JSON(http.StatusOK, res)
}

func (c *TwitterController) PostWebhook() {
	c.engine.POST("/twitter_webhook", c.postWebhook)
}

func (c *TwitterController) postWebhook(ctx *gin.Context) {
	//Read the body of the tweet
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	//Initialize a webhok load obhject for json decoding
	var load WebhookLoad
	err := json.Unmarshal(body, &load)
	if err != nil {
		fmt.Println("An error occured: " + err.Error())
	}
	//Check if it was a tweet_create_event and tweet was in the payload and it was not tweeted by the bot
	if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr {
		return
	}
	params := url.Values{}
	params.Set("in_reply_to_status_id", load.TweetCreateEvent[0].IdStr)

	api := key.GetTwitterApi()
	t, err := api.PostTweet("@"+load.TweetCreateEvent[0].User.Handle+" Hello World", params)
	log.Println(t)

	if err != nil {
		fmt.Println("An error occured:")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Tweet sent successfully")
	}
}

//Struct to parse webhook load
type WebhookLoad struct {
	UserId           string  `json:"for_user_id"`
	TweetCreateEvent []Tweet `json:"tweet_create_events"`
}

//Struct to parse tweet
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

//Struct to parse user
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}
