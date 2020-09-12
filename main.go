package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-giligili/handler/http"
	"go-giligili/middleware"
	"go-giligili/model"
	"go-giligili/model/repository"
	"go-giligili/model/usecase"
	"go-giligili/util"
	"time"
)

var (
	Dsn            = util.GetEnvWithDefault("MYSQL_URL", "root:123456@tcp(localhost:3306)/go-giligili?charset=utf8&parseTime=True&loc=Local")
	SESSION_SECRET = util.GetEnvWithDefault("SESSION_SECRET", "bKfDw9M2yMHV574I")
)

func getDB() *gorm.DB {
	DB, err := gorm.Open("mysql", Dsn)
	DB.LogMode(true)
	// Error
	if err != nil {
		panic(err)
	}
	if gin.Mode() == "release" {
		DB.LogMode(false)
	}
	//设置连接池
	//空闲
	DB.DB().SetMaxIdleConns(20)
	//打开
	DB.DB().SetMaxOpenConns(100)
	//超时
	DB.DB().SetConnMaxLifetime(time.Second * 30)
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}).
		AutoMigrate(&model.Video{})
	return DB
}

func webService(usecase usecase.VidoeUsecase) {
	r := gin.Default()
	// 中间件, 顺序不能改
	r.Use(middleware.Session(SESSION_SECRET))
	r.Use(middleware.Cors())
	http.SetRouterApi(r, usecase)
	r.Run(":8080")
}

func main() {
	db := getDB()
	repository := repository.NewRepositoryImpl(db)

	usecase := usecase.NewVideoUsecaseImpl(repository)
	// 读取翻译文件
	if err := model.LoadLocales("conf/zh-cn.yaml"); err != nil {
		panic(err)
	}
	webService(usecase)
}
