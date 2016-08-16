package webmock

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Cache{})
	return db, nil
}

func insertCache(c *Cache, db *gorm.DB) error {
	if err := db.Create(c).Error; err != nil {
		return err
	}
	return nil
}
