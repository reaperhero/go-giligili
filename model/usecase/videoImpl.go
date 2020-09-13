package usecase

import (
	"go-giligili/model"
	"time"
)

func (u *VideoUsecaseImpl) CreateVideo(title string, info string, url string, avator string) *model.Response {
	var video = model.Video{
		Title:  title,
		Info:   info,
		URL:    url,
		Avatar: avator,
	}
	item, err := u.db.CreateVideo(video)
	if err != nil {
		return &model.Response{
			Status: 50001,
			Msg:    "视频保存失败",
			Error:  err.Error(),
		}
	}
	data := struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Info      string `json:"info"`
		URL       string `json:"url"`
		Avatar    string `json:"avatar"`
		View      uint64 `json:"view"`
		CreatedAt int64  `json:"created_at"`
	}{
		ID:        item.ID,
		Title:     item.Title,
		Info:      item.Info,
		URL:       item.URL,
		Avatar:    item.Avatar,
		View:      0,
		CreatedAt: item.CreatedAt.Unix(),
	}
	return &model.Response{
		Status: 10000,
		Data:   data,
		Msg:    "",
		Error:  "",
	}
}

func (u *VideoUsecaseImpl) ShowVideo(id string) *model.Response {
	video, count, err := u.db.AddVideoShowCountByID(id)
	if err != nil {
		return &model.Response{
			Status: 40005,
			Data:   nil,
			Msg:    "获取视频出错",
			Error:  err.Error(),
		}
	}
	var response = struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Info      string    `json:"info"`
		URL       string    `json:"url"`
		Avatar    string    `json:"avatar"`
		View      uint64    `json:"view"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID:        video.ID,
		Title:     video.Title,
		Info:      video.Info,
		URL:       video.URL,
		Avatar:    video.Avatar,
		View:      count,
		CreatedAt: video.CreatedAt,
	}
	return &model.Response{
		Status: 1000,
		Data:   response,
		Msg:    "",
		Error:  "",
	}
}

func (u *VideoUsecaseImpl) GetVideosList(limit int, start int) *model.Response {
	videos, total, err := u.db.ListVideos(limit, start)
	if err != nil {
		return &model.Response{
			Status: 50000,
			Data:   nil,
			Msg:    "获取视频出错",
			Error:  err.Error(),
		}
	}
	var data = struct {
		Items interface{} `json:"items"`
		Total int         `json:"total"`
	}{
		Items: videos,
		Total: total,
	}
	return &model.Response{
		Status: 1000,
		Data:   data,
		Msg:    "",
		Error:  "",
	}
}

func (u *VideoUsecaseImpl) UpdateVideo(id int, title string, info string, url string, avatar string) *model.Response {
	var video = model.Video{
		Title:  title,
		Info:   info,
		URL:    url,
		Avatar: avatar,
	}
	if video, err := u.db.UpdateVideo(id, video); err != nil {
		return &model.Response{
			Status: 50000,
			Data:   nil,
			Msg:    "更新失败",
			Error:  err.Error(),
		}
	} else {
		return &model.Response{
			Status: 1000,
			Data:   video,
			Msg:    "",
			Error:  "",
		}
	}
}

func (u *VideoUsecaseImpl) DeleteVideo(id int) *model.Response {
	video, err := u.db.DeleteVideoById(id)
	if err != nil {
		return &model.Response{
			Status: 50000,
			Data:   nil,
			Msg:    "删除失败",
			Error:  err.Error(),
		}
	}
	return &model.Response{
		Status: 1000,
		Data:   video,
		Msg:    "",
		Error:  "",
	}
}

func (u *VideoUsecaseImpl) GetRankVideo() *model.Response {
	videos, err := u.db.RedisVideoViewRank()
	if err != nil {
		return &model.Response{
			Status: 50000,
			Data:   nil,
			Msg:    "查询出错",
			Error:  err.Error(),
		}
	}
	return &model.Response{
		Status: 1000,
		Data:   videos,
		Msg:    "",
		Error:  "",
	}
}
