package http

import (
	"github.com/gin-gonic/gin"
	"go-giligili/middleware"
	"go-giligili/model"
	"go-giligili/model/usecase"
)

type httpHandler struct {
	usecase usecase.VidoeUsecase
}

func SetRouterApi(r *gin.Engine, u usecase.VidoeUsecase) {
	var api = httpHandler{usecase: u}
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, model.Response{
				Status: 0,
				Msg:    "Pong",
			})
		})
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

	}

	// 视频操作
	video := r.Group("/api/v1")
	{
		video.POST("video", api.CreateVideo)
		video.GET("video/:id", api.ShowVideo)
		video.GET("videos", api.ListVideo)
		video.PUT("video/:id", api.UpdateVideo)
		video.DELETE("video/:id", api.DeleteVideo)
		// 排行榜
		v1.GET("rank/daily", api.DailyRank)
	}

	authUserApi := r.Group("/api/v1/")
	authUserApi.Use(middleware.AuthRequired())
	{
		// User Routing
		authUserApi.GET("user/me", api.UserMe)
		authUserApi.DELETE("user/logout", api.UserLogout)
	}

}
