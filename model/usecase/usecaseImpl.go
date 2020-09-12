package usecase

import (
	"go-giligili/model"
	"go-giligili/model/repository"
)

type VideoUsecaseImpl struct {
	db repository.Repository
}

func NewVideoUsecaseImpl(repository repository.Repository) VidoeUsecase {
	service := VideoUsecaseImpl{db: repository}
	return &service
}

func (u *VideoUsecaseImpl) UserRegister(userRegister model.UserRegister) (model.User, *model.Response) {
	user := model.User{
		Nickname: userRegister.Nickname,
		UserName: userRegister.UserName,
		Status:   model.Active,
	}
	if count := u.db.CheckUserNickName(userRegister.Nickname); count > 0 {
		return user, &model.Response{
			Status: 40001,
			Msg:    "昵称被占用",
		}
	}

	if count := u.db.CheckUserUserName(userRegister.UserName); count > 0 {
		return user, &model.Response{
			Status: 40001,
			Msg:    "用户名已经注册",
		}
	}

	// 加密密码
	if err := user.SetPassword(userRegister.Password); err != nil {
		return user, &model.Response{
			Status: 40002,
			Msg:    "密码加密失败",
		}
	}

	// 创建用户
	if err := u.db.CreateUser(&user); err != nil {
		return user, &model.Response{
			Status: 40002,
			Msg:    "注册失败",
		}
	}
	return user, nil
}
