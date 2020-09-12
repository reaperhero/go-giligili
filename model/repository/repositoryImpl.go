package repository

import (
	"github.com/jinzhu/gorm"
	"go-giligili/model"
	"log"
)

type DbImpl struct {
	db *gorm.DB
}

func NewRepositoryImpl(db *gorm.DB) Repository {
	repo := &DbImpl{db: db}
	return repo
}

func (r *DbImpl) CheckUserNickName(nickName string) int {
	if r.db == nil {
		log.Println("ss")
	}
	var count int
	r.db.Model(&model.User{}).Where("nickname = ?", nickName).Count(&count)
	log.Println(count)
	return count
}

func (r *DbImpl) CheckUserUserName(userName string) int {
	var count int
	r.db.Model(&model.User{}).Where("user_name = ?", userName).Count(&count)
	return count
}

func (r *DbImpl) CreateUser(user *model.User) error {
	return r.db.Create(&user).Error
}

func (r *DbImpl) UserLogin(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("user_name = ?", username).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}
