package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var cfg = config.Load()

var tables = []interface{}{
	&User{},
	&Order{},
	&Item{},
	&Product{},
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect with database")
	}

	// migrate
	db.AutoMigrate(tables...)

	return db
}

func Drop(db *gorm.DB) {
	migrator := db.Migrator()
	migrator.DropTable(tables...)
}
