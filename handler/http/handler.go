package http

import (
	"github.com/gin-gonic/gin"
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
}
