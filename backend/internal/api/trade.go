package api

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func SetupTradeRoutes(app *fiber.App) {
	tradeGroup := app.Group("/api/trade")

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	tradeGroup.Post("/execute", func(c *fiber.Ctx) error {
		var req struct {
			UserID  uint    `json:"user_id"`
			Symbol  string  `json:"symbol"`
			Amount  float64 `json:"amount"`
			IsBuy   bool    `json:"is_buy"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		// Push to Redis Pub/Sub so the worker can pick it up
		payload, _ := json.Marshal(req)
		err := rdb.Publish(context.Background(), "trade_requests", string(payload)).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to enqueue trade"})
		}

		// Also listen for results in the WsHub to broadcast to frontend
		// (In a real app, this should be a dedicated listener process, but for MVP we can do it here or in ws.go)

		return c.JSON(fiber.Map{"status": "queued", "message": "Trade request submitted"})
	})
}
