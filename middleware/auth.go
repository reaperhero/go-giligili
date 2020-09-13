package middleware

import (
	"github.com/gin-gonic/gin"
	"go-giligili/model"
)

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user_id"); user != nil {
			if _, ok := user.(*model.User); ok {
				c.Next()
				return
			}
		}

		c.JSON(200, model.Response{
			Status: 401,
			Msg:    "需要登录",
		})
		c.Abort()
	}
}
