package main

import (
	"log"

	"github.com/codepnw/go-ticket-booking/config"
	"github.com/codepnw/go-ticket-booking/internal/api"
)

const envPath = "dev.env"

func main() {
	cfg, err := config.SetupConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	api.StartServer(*cfg)
}