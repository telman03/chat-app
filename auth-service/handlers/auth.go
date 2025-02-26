package handlers

import (
    "chat-app/auth-service/database"
    "chat-app/auth-service/models"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "time"

    "golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
    var data struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&data); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

    user := models.User{
        Username: data.Username,
        Email:    data.Email,
        Password: string(hashedPassword),
    }

    database.DB.Create(&user)

    return c.Status(201).JSON(fiber.Map{"message": "User registered"})
}

const SecretKey = "supersecretkey"

func Login(c *fiber.Ctx) error {
    var data struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&data); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    var user models.User
    database.DB.Where("email = ?", data.Email).First(&user)

    if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)) != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(SecretKey))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
    }

    return c.JSON(fiber.Map{"token": tokenString})
}