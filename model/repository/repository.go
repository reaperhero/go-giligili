package repository

import "go-giligili/model"

type Repository interface {
	CheckUserNickName(nickName string) int
	CheckUserUserName(userName string) int
	CreateUser(user *model.User) error
	UserLogin(username string) (*model.User, error)
	CreateVideo(video model.Video) (*model.Video, error)
	AddVideoShowCountByID(id string) (video *model.Video, count uint64, err error)
	ListVideos(limit int, start int) (videos []model.Video, total int, err error)
	UpdateVideo(id int, modifyVideo model.Video) (video *model.Video, err error)
	DeleteVideoById(id int) (video *model.Video, err error)

	RedisVideoViewRank() ([]model.Video, error)
}
