package repository

import "go-giligili/model"

type Repository interface {
	CheckUserNickName(nickName string) int
	CheckUserUserName(userName string) int
	CreateUser(user *model.User) error
}
