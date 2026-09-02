package api

import (
	"github.com/gofiber/fiber/v2"
)

func SetupMarketRoutes(app *fiber.App) {
	marketGroup := app.Group("/api/market")

	marketGroup.Get("/trending", func(c *fiber.Ctx) error {
		// Stub data for trending tokens
		trending := []fiber.Map{
			{"symbol": "SOL", "name": "Solana", "price": 145.20, "change24h": 5.4},
			{"symbol": "ETH", "name": "Ethereum", "price": 3100.50, "change24h": 2.1},
			{"symbol": "PEPE", "name": "Pepe", "price": 0.000008, "change24h": -4.2},
		}
		return c.JSON(fiber.Map{"trending": trending})
	})

	marketGroup.Get("/token/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		
		// Stub response for token details
		return c.JSON(fiber.Map{
			"address": address,
			"symbol": "MOCK",
			"name": "Mock Token",
			"price": 1.05,
			"liquidity": 1500000,
			"volume24h": 450000,
			"marketCap": 10500000,
			"riskScore": 25, // 0-100, lower is better
			"holders": 12500,
		})
	})
}
