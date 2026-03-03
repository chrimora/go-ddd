package config

import "fmt"

type DBConfig struct {
	DBHost     string `env:"POSTGRES_HOST"`
	DBName     string `env:"POSTGRES_DB"`
	DBUser     string `env:"POSTGRES_USER"`
	DBPassword string `env:"POSTGRES_PASSWORD"`

	// Optional read replica connection
	ReadDBHost     string `env:"READ_POSTGRES_HOST" envDefault:""`
	ReadDBName     string `env:"READ_POSTGRES_DB" envDefault:""`
	ReadDBUser     string `env:"READ_POSTGRES_USER" envDefault:""`
	ReadDBPassword string `env:"READ_POSTGRES_PASSWORD" envDefault:""`
}

func (cfg DBConfig) WriteConnString() string {
	return buildConnString(cfg.DBHost, cfg.DBName, cfg.DBUser, cfg.DBPassword)
}

func (cfg DBConfig) ReadConnString() string {
	return buildConnString(
		fallback(cfg.ReadDBHost, cfg.DBHost),
		fallback(cfg.ReadDBName, cfg.DBName),
		fallback(cfg.ReadDBUser, cfg.DBUser),
		fallback(cfg.ReadDBPassword, cfg.DBPassword),
	)
}

func buildConnString(host, name, user, password string) string {
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s",
		host,
		name,
		user,
		password,
	)
}

func fallback(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
