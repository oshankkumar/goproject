package db

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	AppName         string
	Driver          string
	User            string
	Password        string
	Host            string
	Port            int
	DBName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime time.Duration
}

func New(conf Config) (*sqlx.DB, error) {
	dsn := url.URL{
		User:   url.UserPassword(conf.User, conf.Password),
		Scheme: conf.Driver,
		Host:   fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:   conf.DBName,
		RawQuery: (&url.Values{
			"sslmode":          []string{"disable"},
			"application_name": []string{conf.AppName},
		}).Encode(),
	}

	db, err := sqlx.Open(conf.Driver, dsn.String())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	return db, nil
}
