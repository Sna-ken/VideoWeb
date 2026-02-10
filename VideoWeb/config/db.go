package config

var Mysql = struct {
	USERNAME string
	PASSWORD string
	HOST     string
	PORT     string
	NAME     string
}{
	USERNAME: "root",
	PASSWORD: "Aa133944",
	HOST:     "127.0.0.1",
	PORT:     "3306",
	NAME:     "videoweb",
}

var Redis = struct {
	HOST     string
	PORT     string
	PASSWORD string
	DB       int
}{
	HOST:     "127.0.0.1",
	PORT:     "6379",
	PASSWORD: "Aa133944",
	DB:       0,
}
