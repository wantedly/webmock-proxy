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

func insertEndpoint(endpoint *Endpoint, db *gorm.DB) error {
	if err := db.Create(endpoint).Error; err != nil {
		return err
	}
	return nil
}

func updateEndpoint(ce, endpoint *Endpoint, db *gorm.DB) {
	db.Model(ce).Update(endpoint)
}

func deleteConnection(conn *Connection, db *gorm.DB) {
	db.Delete(conn)
}

func findEndpoint(method, url string, db *gorm.DB) *Endpoint {
	var endpoint Endpoint
	db.Preload("Connections").
		Preload("Connections.Request", "method = ?", method).
		Preload("Connections.Response").
		Where(Endpoint{URL: url}).
		Find(&endpoint)
	return &endpoint
}

func readEndpoint(url string, db *gorm.DB) *Endpoint {
	var endpoint Endpoint
	db.Preload("Connections").
		Preload("Connections.Request").
		Preload("Connections.Response").
		Where(Endpoint{URL: url}).
		Find(&endpoint)
	return &endpoint
}
