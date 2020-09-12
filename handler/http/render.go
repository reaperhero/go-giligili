package http

import (
	"github.com/gin-gonic/gin"
	"go-giligili/model"
)

func (h *httpHandler) UserRegister(c *gin.Context) {
	var userRegister = model.UserRegister{}
	if err := c.ShouldBind(&userRegister); err != nil {
		c.JSON(200, model.ErrorResponse(err))
		return
	}
	// 表单验证
	if userRegister.PasswordConfirm != userRegister.Password {
		response := model.Response{
			Status: 40001,
			Msg:    "两次输入的密码不相同",
		}
		c.JSON(200, response)
	}
	var user model.User
	user, err := h.usecase.UserRegister(userRegister)
	if err != nil {
		c.JSON(200, err)
	}
	res := struct {
		ID        uint   `json:"id"`
		UserName  string `json:"user_name"`
		Nickname  string `json:"nickname"`
		Status    string `json:"status"`
		Avatar    string `json:"avatar"`
		CreatedAt int64  `json:"created_at"`
	}{
		ID:        user.ID,
		UserName:  user.UserName,
		Nickname:  user.Nickname,
		Status:    user.Status,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Unix(),
	}
	c.JSON(200, res)
}
