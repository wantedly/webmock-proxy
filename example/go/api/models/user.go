package models

type User struct {
	ID     uint     `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Name   string   `json:"name" form:"name"`
	Age    uint     `json:"age" form:"age"`
	Emails []string `json:"emails" form:"emails"`
}
