package webmock

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&Endpoint{},
		&Connection{},
		&Request{},
		&Response{},
	)
	return db, nil
}

func insertEndpoint(e *Endpoint, db *gorm.DB) error {
	if err := db.Create(e).Error; err != nil {
		return err
	}
	return nil
}

func findEndpoint(db *gorm.DB, method, url string) *Endpoint {
	var endpoint Endpoint
	db.Preload("Connections").
		Preload("Connections.Request", "method = ?", method).
		Preload("Connections.Response").
		Where(Endpoint{URL: url}).
		Find(&endpoint)
	return &endpoint
}
