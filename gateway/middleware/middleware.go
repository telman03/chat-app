package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

const AuthServiceURL = "http://localhost:8080/auth/verify"

// AuthMiddleware verifies JWT tokens by calling auth-service
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			fmt.Println("❌ No Authorization header found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		// Remove 'Bearer ' prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		fmt.Println("🔍 Sending token to auth-service:", token)

		// Prepare request body
		body, _ := json.Marshal(map[string]string{"token": token})
		req, _ := http.NewRequest("POST", AuthServiceURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Send request to auth-service
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Failed to reach auth-service:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to contact auth-service"})
		}
		defer resp.Body.Close()

		// Read response from auth-service
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Println("✅ Auth-service response:", string(respBody))

		var authResponse map[string]interface{}
		json.Unmarshal(respBody, &authResponse)

		// Check if token is valid
		if valid, ok := authResponse["valid"].(bool); !ok || !valid {
			fmt.Println("❌ Invalid or expired token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Attach user_id to request context
		c.Locals("user_id", authResponse["user_id"])
		fmt.Println("✅ Token is valid! User ID:", authResponse["user_id"])

		return c.Next()
	}
}