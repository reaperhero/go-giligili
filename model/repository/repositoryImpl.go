package repository

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go-giligili/model"
	"log"
	"strconv"
	"strings"
)

type DbImpl struct {
	mysqlDb *gorm.DB
	redisDb *redis.Client
}

func NewRepositoryImpl(mysqlDb *gorm.DB, redisDb *redis.Client) Repository {
	repo := DbImpl{
		mysqlDb: mysqlDb,
		redisDb: redisDb,
	}
	return &repo
}

func (r *DbImpl) CheckUserNickName(nickName string) int {
	var count int
	r.mysqlDb.Model(&model.User{}).Where("nickname = ?", nickName).Count(&count)
	log.Println(count)
	return count
}

func (r *DbImpl) CheckUserUserName(userName string) int {
	var count int
	r.mysqlDb.Model(&model.User{}).Where("user_name = ?", userName).Count(&count)
	return count
}

func (r *DbImpl) CreateUser(user *model.User) error {
	return r.mysqlDb.Create(&user).Error
}

func (r *DbImpl) UserLogin(username string) (*model.User, error) {
	var user model.User
	if err := r.mysqlDb.Where("user_name = ?", username).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (r *DbImpl) CreateVideo(video model.Video) (*model.Video, error) {
	err := r.mysqlDb.Create(&video).Error
	return &video, err
}

func (r *DbImpl) AddVideoShowCountByID(id string) (video *model.Video, count uint64, err error) {
	var v model.Video
	if err := r.mysqlDb.First(&v, id).Error; err != nil {
		return nil, 0, err
	}
	// 增加点击数
	redisVideoId := fmt.Sprintf("view:video:%s", strconv.Itoa(int(v.ID)))
	r.redisDb.Incr(redisVideoId)
	// 增加排行点击数
	r.redisDb.ZIncrBy("rank:daily", 1, id)
	// 获取点击数
	countStr, _ := r.redisDb.Get(redisVideoId).Result()
	count, _ = strconv.ParseUint(countStr, 10, 64)
	return &v, count, nil
}

// List 视频列表
func (r *DbImpl) ListVideos(limit int, start int) (videos []model.Video, total int, err error) {
	videos = []model.Video{}
	total = 0

	if limit == 0 {
		limit = 6
	}

	if err := r.mysqlDb.Model(model.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.mysqlDb.Limit(limit).Offset(start).Find(&videos).Error; err != nil {
		return nil, 0, err
	}
	return videos, total, nil
}

// 更新视频
func (r *DbImpl) UpdateVideo(id int, modifyVideo model.Video) (video *model.Video, err error) {
	var queryVideo model.Video
	if err := r.mysqlDb.First(&queryVideo, id).Error; err != nil {
		return nil, err
	}
	queryVideo.Info = modifyVideo.Info
	queryVideo.Title = modifyVideo.Title
	queryVideo.URL = modifyVideo.URL
	queryVideo.Avatar = modifyVideo.Avatar
	if err := r.mysqlDb.Save(&queryVideo).Error; err != nil {
		return nil, err
	}
	return &queryVideo, nil
}

func (r *DbImpl) DeleteVideoById(id int) (video *model.Video, err error) {
	var queryVideo model.Video
	if err := r.mysqlDb.First(&queryVideo, id).Error; err != nil {
		return nil, err
	}
	r.mysqlDb.Delete(&queryVideo)
	return &queryVideo, nil
}

func (r *DbImpl) RedisVideoViewRank() ([]model.Video, error) {

	var videos []model.Video

	// 从redis读取点击前十的视频
	vids, _ := r.redisDb.ZRevRange("rank:daily", 0, 9).Result()

	if len(vids) > 1 {
		order := fmt.Sprintf("FIELD(id, %s)", strings.Join(vids, ","))
		err := r.mysqlDb.Where("id in (?)", vids).Order(order).Find(&videos).Error
		if err != nil {
			return nil, err
		}
	}
	return videos, nil
}
