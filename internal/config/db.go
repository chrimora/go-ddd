package config

type DBConfig struct {
	DBHost     string `env:"POSTGRES_HOST"`
	DBName     string `env:"POSTGRES_DB"`
	DBUser     string `env:"POSTGRES_USER"`
	DBPassword string `env:"POSTGRES_PASSWORD"`
}
