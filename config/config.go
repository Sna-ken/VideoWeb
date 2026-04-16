package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var MYSQLDB *gorm.DB
var REDISDB *redis.Client

type MysqlConfig struct {
	USERNAME string `mapstructure:"username"`
	PASSWORD string `mapstructure:"password"`
	HOST     string `mapstructure:"host"`
	PORT     string `mapstructure:"port"`
	NAME     string `mapstructure:"name"`
}

type RedisConfig struct {
	HOST     string `mapstructure:"host"`
	PORT     string `mapstructure:"port"`
	PASSWORD string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
	AccessTokenExpiry  int64  `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry int64  `mapstructure:"refresh_token_expiry"`
}

var Mysql *MysqlConfig
var Redis *RedisConfig
var JWT *JWTConfig

func Init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config fail: %v", err)
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed: %s", e.Name)

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("read config fail: %v", err)
			return
		}

		log.Printf("config reloaded successfully")
	})

	if err := viper.Sub("mysql").Unmarshal(&Mysql); err != nil {
		log.Fatalf("unmarshal mysql config fail: %v", err)
	}
	if err := viper.Sub("redis").Unmarshal(&Redis); err != nil {
		log.Fatalf("unmarshal redis config fail: %v", err)
	}
	if err := viper.Sub("jwt").Unmarshal(&JWT); err != nil {
		log.Fatalf("unmarshal jwt config fail: %v", err)
	}

}
