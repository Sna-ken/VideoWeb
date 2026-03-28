package config

import "os"

var Mysql = struct {
	USERNAME string
	PASSWORD string
	HOST     string
	PORT     string
	NAME     string
}{
	USERNAME: getEnv("DB_USER", "root"),
	PASSWORD: getEnv("DB_PASSWORD", "root"),
	HOST:     getEnv("DB_HOST", "localhost"),
	PORT:     getEnv("DB_PORT", "3306"),
	NAME:     getEnv("DB_NAME", "videodb"),
}

var Redis = struct {
	HOST     string
	PORT     string
	PASSWORD string
	DB       int
}{
	HOST:     "redis",
	PORT:     "6379",
	PASSWORD: "",
	DB:       0,
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
