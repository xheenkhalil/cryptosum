package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cryptosum/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func SetupAIRoutes(app *fiber.App) {
	api := app.Group("/api/ai")

	api.Post("/analyze", func(c *fiber.Ctx) error {
		var payload struct {
			Symbol    string  `json:"symbol"`
			Price     float64 `json:"price"`
			Volume24h float64 `json:"volume24h"`
			Liquidity float64 `json:"liquidity"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		ctx := context.Background()
		apiKey := config.AppConfig.GeminiAPIKey
		if apiKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI not configured"})
		}

		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to init AI client"})
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-2.5-flash") // Defaulting to an available gemini model
		model.SetTemperature(0.2)
		model.ResponseMIMEType = "application/json"

		prompt := fmt.Sprintf(`Analyze this token data and return a JSON response with:
- direction ("bullish", "bearish", "neutral")
- confidence (0-100)
- risk ("low", "medium", "high")
- reasoning (short string)

Token: %s
Price: $%f
Volume 24h: $%f
Liquidity: $%f
`, payload.Symbol, payload.Price, payload.Volume24h, payload.Liquidity)

		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI generation failed"})
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Empty AI response"})
		}

		part := resp.Candidates[0].Content.Parts[0]
		textPart, ok := part.(genai.Text)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid AI response format"})
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(textPart), &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse AI JSON"})
		}

		return c.JSON(result)
	})
}
