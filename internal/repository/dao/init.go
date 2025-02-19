package dao

import "gorm.io/gorm"

func CreateTale(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
