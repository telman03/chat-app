package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
    Email    string `gorm:"uniqueIndex"`
    Password string
}