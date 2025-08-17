package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload" // Loads .env
)

type DBConfig struct {
	DBHost     string `env:"POSTGRES_HOST"`
	DBName     string `env:"POSTGRES_DB"`
	DBUser     string `env:"POSTGRES_USER"`
	DBPassword string `env:"POSTGRES_PASSWORD"`
}

func NewDBConfig() *DBConfig {
	cfg := &DBConfig{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	// Check required values are present
	if cfg.DBHost == "" || cfg.DBName == "" || cfg.DBUser == "" || cfg.DBPassword == "" {
		panic(fmt.Errorf("Missing env."))
	}
	return cfg
}
