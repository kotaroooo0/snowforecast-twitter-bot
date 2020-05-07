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
	TwitterUseCase usecase.TwitterUseCase
}

func (th TwitterHandlerImpl) HandleTwitterGetCrcToken(ctx *gin.Context) {
	req := th.TwitterUseCase.NewGetTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	res := th.TwitterUseCase.GetCrcTokenResponse(req)
	ctx.JSON(http.StatusOK, res)
}

func (th TwitterHandlerImpl) HandleTwitterPostWebhook(ctx *gin.Context) {
	req := th.TwitterUseCase.NewPostTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	res := th.TwitterUseCase.PostAutoReplyResponse(req)
	ctx.JSON(http.StatusOK, res)
}
