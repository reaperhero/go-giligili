package usecase

import "go-giligili/model"

type VidoeUsecase interface {
	UserRegister(userRegister model.UserRegister) (model.User, *model.Response)
}
