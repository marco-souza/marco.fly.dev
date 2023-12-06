package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var cfg = config.Load()

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect with database")
	}

	// migrate
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Product{})

	return db
}
