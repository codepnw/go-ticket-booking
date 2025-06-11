package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppPort          string
	DBAddr           string
	JWTSecret        string
	JWTRefreshSecret string
}

func SetupConfig(envPath string) (*AppConfig, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("loading env failed: %v", err)
	}

	appPort, ok := os.LookupEnv("APP_PORT")
	if !ok {
		return nil, fmt.Errorf("APP_PORT not found")
	}

	dbAddr, ok := os.LookupEnv("DB_ADDR")
	if !ok {
		return nil, fmt.Errorf("DB_ADDR not found")

	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return nil, fmt.Errorf("JWT_SECRET not found")
	}

	jwtRefreshSecret, ok := os.LookupEnv("JWT_REFRESH_SECRET")
	if !ok {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET not found")
	}

	return &AppConfig{
		AppPort:          appPort,
		DBAddr:           dbAddr,
		JWTSecret:        jwtSecret,
		JWTRefreshSecret: jwtRefreshSecret,
	}, nil
}
