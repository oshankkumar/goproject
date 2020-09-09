package config

import "example.com/example/goproject/pkg/db"

func NewDBConfig() db.Config {
	return db.Config{
		AppName:         "",
		Driver:          "",
		User:            "",
		Password:        "",
		Host:            "",
		Port:            0,
		DBName:          "",
		MaxIdleConns:    0,
		MaxOpenConns:    0,
		ConnMaxLifeTime: 0,
	}
}
