package main

import (
	"log"
	"os"

	"github.com/cryptosum/backend/internal/api"
	"github.com/cryptosum/backend/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=127.0.0.1 user=postgres password=password dbname=postgres port=15432 sslmode=disable"
	}
	db.ConnectDB(dsn)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
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
