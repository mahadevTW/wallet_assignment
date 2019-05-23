package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	port, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 0, 64)
	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     os.Getenv("DB_HOST"),
			Port:     int(port),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
			Charset:  "utf8",
		},
	}
}
