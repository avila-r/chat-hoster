package db

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
	once     sync.Once
)

func GetConnection() *gorm.DB {
	once.Do(
		func() {
			dsn := os.Getenv("DATABASE_DSN")

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

			if err != nil {
				log.Fatal("failed to connect database: ", err)
			}

			AutoMigrate(db)
			database = db
		})

	return database
}

func AutoMigrate(db *gorm.DB) {
	_ = db.AutoMigrate()
}
