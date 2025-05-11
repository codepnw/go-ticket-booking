package main

import (
	"log"
	"os"

	"github.com/codepnw/go-ticket-booking/internal/store"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("dev.env"); err != nil {
		log.Panicf("loading env failed: %v", err)
	}
}

func main() {
	dbAddr, ok := os.LookupEnv("DB_ADDR")
	if !ok {
		log.Fatal("DB_ADDR not found")
	}

	db, err := store.InitPostgresDB(dbAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("server started")
}