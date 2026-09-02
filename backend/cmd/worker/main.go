package main

import (
	"github.com/cryptosum/backend/internal/config"
	"github.com/cryptosum/backend/internal/db"
	"github.com/cryptosum/backend/internal/worker"
)

func main() {
	// 1. Load strict configuration
	config.LoadConfig()

	// 2. Initialize Database
	db.ConnectDB(config.AppConfig.DatabaseURL)

	worker.ConnectRedis()
	
	worker.StartTradingWorker()
}
