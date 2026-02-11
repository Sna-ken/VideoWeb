package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MYSQLDB *gorm.DB
var REDISDB *redis.Client

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

var JWTConfig = struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  int64 // 秒为单位
	RefreshTokenExpiry int64
}{
	AccessTokenSecret:  "your-access-token-secret-key",
	RefreshTokenSecret: "your-refresh-token-secret-key",
	AccessTokenExpiry:  180,    // 30分钟
	RefreshTokenExpiry: 604800, // 7天
}
