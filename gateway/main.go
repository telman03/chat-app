package main

import (
	"chat-app/gateway/middleware"
	"github.com/gofiber/fiber/v2"
	"fmt"
)

func main() {
	fmt.Println("🚀 API Gateway starting on port 5050...")

	app := fiber.New()



	// Protected routes using AuthMiddleware
	api := app.Group("/api", middleware.AuthMiddleware())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Gateway is running!")
	})
	api.Get("/protected", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"message": "Protected route accessed", "user_id": userID})
	})

	app.Listen(":5050") // API Gateway runs on port 5000
}