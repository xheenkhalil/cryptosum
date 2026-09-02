package main

import (
	"os"

	"github.com/cryptosum/backend/internal/db"
	"github.com/cryptosum/backend/internal/worker"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=cryptosum password=password dbname=cryptosum port=5433 sslmode=disable"
	}
	db.ConnectDB(dsn)

	worker.ConnectRedis()
	
	worker.StartTradingWorker()
}
