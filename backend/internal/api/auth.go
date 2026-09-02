package api

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/cryptosum/backend/internal/db"
	"github.com/cryptosum/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spruceid/siwe-go"
)

// A simple in-memory cache for nonces. In production, use Redis.
var nonces = make(map[string]time.Time)

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	nonce := hex.EncodeToString(b)
	nonces[nonce] = time.Now()
	return nonce
}

func SetupAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/api/auth")

	authGroup.Get("/nonce", func(c *fiber.Ctx) error {
		nonce := generateNonce()
		return c.JSON(fiber.Map{"nonce": nonce})
	})

	authGroup.Post("/verify", func(c *fiber.Ctx) error {
		var req struct {
			Message   string `json:"message"`
			Signature string `json:"signature"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		message, err := siwe.ParseMessage(req.Message)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid SIWE message"})
		}

		// Verify signature and nonce
		_, err = message.Verify(req.Signature, nil, nil, nil)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid signature"})
		}

		if _, exists := nonces[message.GetNonce()]; !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired nonce"})
		}
		delete(nonces, message.GetNonce())

		address := strings.ToLower(message.GetAddress().String())

		// Check if user exists, else create
		var user models.User
		res := db.DB.Where("wallet = ?", address).First(&user)
		if res.Error != nil {
			user = models.User{Wallet: address}
			if err := db.DB.Create(&user).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
			}
		}

		// Create JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":     user.ID,
			"wallet": user.Wallet,
			"exp":    time.Now().Add(time.Hour * 72).Unix(),
		})

		// Use a fixed secret for MVP
		t, err := token.SignedString([]byte("cryptosum-super-secret-key"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		return c.JSON(fiber.Map{
			"token": t,
			"user":  user,
		})
	})
}
