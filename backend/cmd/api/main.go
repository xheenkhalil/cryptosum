package main

import (
	"log"

	"github.com/cryptosum/backend/internal/api"
	"github.com/cryptosum/backend/internal/config"
	"github.com/cryptosum/backend/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Load strict configuration
	config.LoadConfig()

	// 2. Initialize Database
	db.ConnectDB(config.AppConfig.DatabaseURL)

	app := fiber.New()
	app.Use(logger.New())
	
	// Secure CORS using config
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.AppConfig.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "api"})
	})

	api.SetupAuthRoutes(app)
	api.SetupMarketRoutes(app)
	api.SetupAIRoutes(app)
	api.SetupWebSockets(app)
	api.SetupTradeRoutes(app)

	log.Println("API starting on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("API failed to start: %v", err)
	}
}
