package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Mysql.USERNAME, Mysql.PASSWORD, Mysql.HOST, Mysql.PORT, Mysql.NAME)
	DBtemp, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	MYSQLDB = DBtemp
	log.Println("Connected to MySQL")
}

func InitRedis() {
	DBtemp := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Redis.HOST, Redis.PORT),
		Password: Redis.PASSWORD,
		DB:       Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := DBtemp.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect Redis" + err.Error())
	}

	REDISDB = DBtemp
	log.Println("Connected to Redis")
}

func InitTable() {
	if err := MYSQLDB.AutoMigrate(&model.User{}); err != nil {
		panic("failed to migrate User table" + err.Error())
	}
	if err := MYSQLDB.AutoMigrate(&model.Video{}); err != nil {
		panic("failed to migrate Video table" + err.Error())
	}
	if err := MYSQLDB.AutoMigrate(&model.Like{}); err != nil {
		panic("failed to migrate Like table" + err.Error())
	}
	if err := MYSQLDB.AutoMigrate(&model.Comment{}); err != nil {
		panic("failed to migrate Comment table" + err.Error())
	}
	if err := MYSQLDB.AutoMigrate(&model.SocialObject{}); err != nil {
		panic("failed to migrate SocialObject table" + err.Error())
	}
	if err := MYSQLDB.AutoMigrate(&model.Message{}); err != nil {
		panic("failed to migrate Message table" + err.Error())
	}
}
