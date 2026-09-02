package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var ctx = context.Background()

type TradeRequest struct {
	UserID  uint    `json:"user_id"`
	Symbol  string  `json:"symbol"`
	Amount  float64 `json:"amount"`
	IsBuy   bool    `json:"is_buy"`
	MaxSlip float64 `json:"max_slippage"`
}

func ConnectRedis() {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Redis connected successfully")
}

func StartTradingWorker() {
	pubsub := Rdb.Subscribe(ctx, "trade_requests")
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Println("Trading Worker is listening for requests...")

	for msg := range ch {
		var req TradeRequest
		if err := json.Unmarshal([]byte(msg.Payload), &req); err != nil {
			log.Printf("Invalid trade request: %v\n", err)
			continue
		}

		executeTrade(req)
	}
}

func executeTrade(req TradeRequest) {
	log.Printf("Processing Trade: User %d | %s | Amount: %.2f | Buy: %t", req.UserID, req.Symbol, req.Amount, req.IsBuy)

	// Risk Gate Stub
	if req.Amount > 10000 {
		log.Printf("Risk Gate Triggered: Trade too large for User %d\n", req.UserID)
		publishResult(req.UserID, req.Symbol, false, "Amount exceeds limit")
		return
	}

	// Simulation Stub
	log.Println("Simulating transaction...")
	
	// Execution Stub
	log.Println("Executing transaction on chain...")
	
	// Success
	publishResult(req.UserID, req.Symbol, true, "Trade executed successfully")
}

func publishResult(userID uint, symbol string, success bool, message string) {
	result := map[string]interface{}{
		"user_id": userID,
		"symbol":  symbol,
		"success": success,
		"message": message,
	}
	
	payload, _ := json.Marshal(result)
	Rdb.Publish(ctx, fmt.Sprintf("trade_results:%d", userID), string(payload))
}
