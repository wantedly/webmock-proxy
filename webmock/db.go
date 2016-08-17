package webmock

import (
	"net/http"

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

func selectCache(db *gorm.DB, r *http.Request, file *File) Endpoint {
	var endpoint Endpoint
	db.Preload("Connections").
		Preload("Connections.Request", "method = ?", r.Method).
		Preload("Connections.Request.Header").
		Preload("Connections.Response").
		Preload("Connections.Response.Header").
		Where(Endpoint{URL: file.URL}).
		Find(&endpoint)
	return endpoint
}
