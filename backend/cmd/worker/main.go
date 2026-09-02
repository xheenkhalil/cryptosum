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
		dsn = "host=127.0.0.1 user=postgres password=password dbname=postgres port=15432 sslmode=disable"
	}
	db.ConnectDB(dsn)

	worker.ConnectRedis()
	
	worker.StartTradingWorker()
}
