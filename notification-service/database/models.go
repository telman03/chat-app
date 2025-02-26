package database

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID    string `gorm:"index"`
	Message   string
	Timestamp string
}