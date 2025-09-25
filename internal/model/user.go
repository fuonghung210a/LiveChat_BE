package model

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
