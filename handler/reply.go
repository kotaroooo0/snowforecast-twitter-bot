package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotaroooo0/snowforecast-twitter-bot/usecase"
)

type ReplyHandler interface {
	HandleTwitterGetCrcToken(*gin.Context)
	HandleTwitterPostWebhook(*gin.Context)
}

type ReplyHandlerImpl struct {
	ReplyUseCase usecase.ReplyUseCase
}

func NewReplyHandlerImpl(replyUseCase usecase.ReplyUseCase) ReplyHandler {
	return &ReplyHandlerImpl{
		ReplyUseCase: replyUseCase,
	}
}

func (th ReplyHandlerImpl) HandleTwitterGetCrcToken(ctx *gin.Context) {
	req := th.ReplyUseCase.NewGetTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	res := th.ReplyUseCase.GetCrcTokenResponse(req)
	ctx.JSON(http.StatusOK, res)
}

func (th ReplyHandlerImpl) HandleTwitterPostWebhook(ctx *gin.Context) {
	req := th.ReplyUseCase.NewPostTwitterWebhookRequest()
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	res := th.ReplyUseCase.PostAutoReplyResponse(req)
	ctx.JSON(http.StatusOK, res)
}
