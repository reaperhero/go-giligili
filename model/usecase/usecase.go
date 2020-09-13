package usecase

import "go-giligili/model"

type VidoeUsecase interface {
	UserRegister(userRegister model.UserRegister) (model.User, *model.Response)
	UserLogin(username string, password string) (model.User, *model.Response)
	CreateVideo(title string, info string, url string, avator string) *model.Response
	ShowVideo(id string) *model.Response
	GetVideosList(limit int, start int) *model.Response
	UpdateVideo(id int, title string, info string, url string, avatar string) *model.Response
	DeleteVideo(id int) *model.Response
	GetRankVideo() *model.Response
}
