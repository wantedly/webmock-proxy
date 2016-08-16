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
	db.LogMode(true)

	db.AutoMigrate(
		&Endpoint{},
		&Connection{},
		&Request{},
		&Response{},
		&ResponseBody{},
		&Header{},
	)
	return db, nil
}

func insertCache(e *Endpoint, db *gorm.DB) error {
	if err := db.Create(e).Error; err != nil {
		return err
	}
	return nil
}
