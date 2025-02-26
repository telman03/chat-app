package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = db
    fmt.Println("Database connection successful!")
}

type Message struct {
    ID         uint   `gorm:"primaryKey"`
    SenderID   string `gorm:"index"`
    ReceiverID string `gorm:"index"`
    Content    string
    Timestamp  string
}