package controllers

import (
	"net/http"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

type JobController struct {
	engine *gin.Engine
}

func NewJobController(engine *gin.Engine) *JobController {
	return &JobController{engine}
}

func (c *JobController) GetJobStatus() {
	c.engine.GET("/job_status", c.getJobStatus)
}

func (c *JobController) getJobStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jobrunner.StatusJson())
}
