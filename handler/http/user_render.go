package http

import (
	"github.com/gin-contrib/sessions"
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

func (h *httpHandler) UserLogin(c *gin.Context) {
	var requestJson = struct {
		UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
		Password string `form:"password" json:"password" binding:"required,min=8,max=40"`
	}{}
	if err := c.ShouldBind(&requestJson); err != nil {
		c.JSON(200, model.ErrorResponse(err))
	}
	user, response := h.usecase.UserLogin(requestJson.UserName, requestJson.Password)
	if response.Status == 1000 {
		// 设置Session
		s := sessions.Default(c)
		s.Clear()
		s.Set("user_id", user.ID)
		s.Save()
		var data = struct {
			ID        uint   `json:"id"`
			UserName  string `json:"user_name"`
			Nickname  string `json:"nickname"`
			Status    string `json:"status"`
			Avatar    string `json:"avatar"`
			CreatedAt int64  `json:"created_at"`
		}{ID: user.ID,
			UserName:  user.UserName,
			Nickname:  user.Nickname,
			Status:    user.Status,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Unix(),
		}
		response.Data = data
		c.JSON(200, response)
		return
	}
	c.JSON(200, response)
}

func (h *httpHandler) UserMe(c *gin.Context) {
	var data = struct {
		ID        uint   `json:"id"`
		UserName  string `json:"user_name"`
		Nickname  string `json:"nickname"`
		Status    string `json:"status"`
		Avatar    string `json:"avatar"`
		CreatedAt int64  `json:"created_at"`
	}{}
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			data.ID = u.ID
			data.UserName = u.UserName
			data.Nickname = u.Nickname
			data.Status = u.Status
			data.Avatar = u.Avatar
			data.CreatedAt = u.CreatedAt.Unix()
		}
		res := model.Response{
			Status: 1000,
			Data:   data,
			Msg:    "",
			Error:  "",
		}
		c.JSON(200, res)
		return
	}
	c.JSON(200, model.Response{
		Status: 0,
		Data:   nil,
		Msg:    "session为空",
		Error:  "",
	})
}

func (h *httpHandler) UserLogout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(200, model.Response{
		Status: 0,
		Msg:    "登出成功",
	})
}
