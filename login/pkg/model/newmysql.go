package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MysqlConfig struct {
	DataSource  string
	IdleConns   int
	OpenConns   int
	IdleTimeout int64
	LifeTimeout int64
}

func NewMysql(config *MysqlConfig) *gorm.DB {
	if config == nil {
		panic("config cannot be nil")
	}

	db, err := gorm.Open("mysql", config.DataSource)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(config.IdleConns)
	db.DB().SetMaxOpenConns(config.OpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(config.IdleTimeout) * time.Second)
	db.DB().SetConnMaxIdleTime(time.Duration(config.LifeTimeout) * time.Second)

	return db
}
