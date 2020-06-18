package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

type TwitterHandler interface {
	HandleTwitterGetCrcToken(*gin.Context)
	HandleTwitterPostWebhook(*gin.Context)
}

type TwitterHandlerImpl struct {
	TwitterUsecase usecase.TwitterUsecase
}

func NewTwitterHandler(twitterUsecase usecase.TwitterUsecase) TwitterHandler {
	return TwitterHandlerImpl{
		TwitterUsecase: twitterUsecase,
	}
}

func (th TwitterHandlerImpl) HandleTwitterGetCrcToken(ctx *gin.Context) {
	req := th.TwitterUsecase.NewGetTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	res := th.TwitterUsecase.GetCrcTokenResponse(req)
	ctx.JSON(http.StatusOK, res)
}

func (th TwitterHandlerImpl) HandleTwitterPostWebhook(ctx *gin.Context) {
	req := th.TwitterUsecase.NewPostTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	res := th.TwitterUsecase.PostAutoReplyResponse(req)
	ctx.JSON(http.StatusOK, res)
}
