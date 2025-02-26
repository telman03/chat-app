package main

import (
    "chat-app/auth-service/config"
    "chat-app/auth-service/database"
    "chat-app/auth-service/handlers"
    "chat-app/auth-service/models"
    "github.com/gofiber/fiber/v2"
)

func main() {

    config.LoadEnv()

    database.Connect()
    database.DB.AutoMigrate(&models.User{})
    
    app := fiber.New()

    app.Post("/register", handlers.Register)
    app.Post("/login", handlers.Login)
    app.Post("/auth/verify", handlers.VerifyToken) // 🚀 New route

    app.Listen(":8080")
}