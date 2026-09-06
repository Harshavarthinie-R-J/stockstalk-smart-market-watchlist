package config

import "os"

type Config struct {
	Port           string
	Environment    string
	DatabaseURL    string
	JWTSecret      string
	MarketDataFile string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	return Config{
		Port:           port,
		Environment:    environment,
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		MarketDataFile: "mock/market_data.json",
	}
}
