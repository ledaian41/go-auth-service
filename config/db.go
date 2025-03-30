package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func InitDatabase() *gorm.DB {
	var err error
	db, err := gorm.Open(sqlite.Open(Env.CachePath), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to database")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	return db
}
