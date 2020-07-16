package database

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func NewDB(driver, dataSource string, connMaxLifeTime, maxIdleConns, maxOpenConns int, debug bool) *gorm.DB {
	db, err := gorm.Open(driver, dataSource)

	if err != nil {
		log.Panic(err)
	}

	if err = db.DB().Ping(); err != nil {
		log.Panic(err)
	}

	if connMaxLifeTime != 0 && maxIdleConns != 0 && maxOpenConns != 0 {
		db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeTime) * time.Second)
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
	}

	db.LogMode(debug)

	return db
}