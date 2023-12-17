package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DB_PG_STRING")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("Can't connect to DB")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Can't connect to DB")
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(90)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
