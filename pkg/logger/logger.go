package logger

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
)

type Format string

const (
	JSON Format = "json"
	Text Format = "text"
)

type Config struct {
	Format Format
	Level  string
}

func Init(conf Config) {
	if strings.ToLower(string(conf.Format)) == string(JSON) {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	lvl, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		log.Panic("logrus.ParseLevel: invalid log level", err)
	}

	logrus.SetLevel(lvl)
}
