package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robfig/cron"
	"go-giligili/handler/http"
	"go-giligili/middleware"
	"go-giligili/model"
	"go-giligili/model/repository"
	"go-giligili/model/usecase"
	"go-giligili/util"
	"reflect"
	"runtime"
	"time"
)

var (
	Dsn           = util.GetEnvWithDefault("MYSQL_URL", "root:123456@tcp(localhost:3306)/go-giligili?charset=utf8&parseTime=True&loc=Local")
	RedisAddr     = util.GetEnvWithDefault("REDIS_ADDR", "127.0.0.1:6379")
	RedisPassword = util.GetEnvWithDefault("REDIS_PASSWORD", "")
	SessionSecret = util.GetEnvWithDefault("SESSION_SECRET", "bKfDw9M2yMHV574I")
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

func GetRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	return client
}

func webService(usecase usecase.VidoeUsecase) {
	r := gin.Default()
	// 中间件, 顺序不能改
	r.Use(middleware.Session(SessionSecret))
	r.Use(middleware.Cors())
	http.SetRouterApi(r, usecase)
	r.Run(":8080")
}

// Run 运行
func Run(job func() error) {
	from := time.Now().UnixNano()
	err := job()
	to := time.Now().UnixNano()
	jobName := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	if err != nil {
		fmt.Printf("%s error: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	} else {
		fmt.Printf("%s success: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	}
}

// CronJob 定时任务
func CronJob(client *redis.Client) {
	Cron := cron.New()
	var RestartDailyRank = func() error {
		return client.Del("rank:daily").Err()
	}
	Cron.AddFunc("0 0 0 * * *", func() { Run(RestartDailyRank) })

	Cron.Start()

	fmt.Println("Cronjob start.....")
}

func main() {
	mysqlDb := getDB()
	redisDb := GetRedis()
	repository := repository.NewRepositoryImpl(mysqlDb, redisDb)

	usecase := usecase.NewVideoUsecaseImpl(repository)
	// 读取翻译文件
	if err := model.LoadLocales("conf/zh-cn.yaml"); err != nil {
		panic(err)
	}
	CronJob(redisDb)
	webService(usecase)
}
