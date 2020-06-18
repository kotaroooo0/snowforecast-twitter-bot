package handler

import (
	"net/http"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

type JobHandler interface {
	HandleGetJobStatus(*gin.Context)
}

type JobHandlerImpl struct{}

func NewJobHandler() JobHandler {
	return JobHandlerImpl{}
}

func (jh JobHandlerImpl) HandleGetJobStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jobrunner.StatusJson())
}
