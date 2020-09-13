package http

import (
	"github.com/gin-gonic/gin"
	"go-giligili/model"
	"strconv"
)

func (h *httpHandler) CreateVideo(c *gin.Context) {
	requestJson := struct {
		Title  string `form:"title" json:"title" binding:"required,min=2,max=100"`
		Info   string `form:"info" json:"info" binding:"max=3000"`
		URL    string `form:"url" json:"url"`
		Avatar string `form:"avatar" json:"avatar"`
	}{}
	if err := c.ShouldBind(&requestJson); err == nil {
		response := h.usecase.CreateVideo(requestJson.Title, requestJson.Info, requestJson.URL, requestJson.Avatar)
		c.JSON(200, response)
	} else {
		c.JSON(200, model.ErrorResponse(err))
	}
}

func (h *httpHandler) ShowVideo(c *gin.Context) {
	videoId := c.Param("id")
	response := h.usecase.ShowVideo(videoId)
	c.JSON(200, response)
}

func (h *httpHandler) ListVideo(c *gin.Context) {
	requestJson := struct {
		Limit int `form:"limit"`
		Start int `form:"start"`
	}{}
	if err := c.ShouldBind(&requestJson); err != nil {
		c.JSON(200, model.ErrorResponse(err))
		return
	}
	response := h.usecase.GetVideosList(requestJson.Limit, requestJson.Start)
	c.JSON(200, response)
}

func (h *httpHandler) UpdateVideo(c *gin.Context) {
	var requestJson = struct {
		Title  string `form:"title" json:"title" binding:"required,min=2,max=30"`
		Info   string `form:"info" json:"info" binding:"max=300"`
		Url    string `form:"url" json:"url" binding:"max=300"`
		Avatar string `form:"avatar" json:"avatar" binding:"max=300"`
	}{}
	videoId, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBind(&requestJson); err != nil || videoId == 0 {
		c.JSON(200, model.ErrorResponse(err))
		return
	}
	response := h.usecase.UpdateVideo(videoId, requestJson.Title, requestJson.Info, requestJson.Url, requestJson.Avatar)
	c.JSON(200, response)
}

func (h *httpHandler) DeleteVideo(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.Param("id"))
	response := h.usecase.DeleteVideo(videoId)
	c.JSON(200, response)
}

func (h *httpHandler) DailyRank(c *gin.Context) {
	response := h.usecase.GetRankVideo()
	c.JSON(200, response)
}
