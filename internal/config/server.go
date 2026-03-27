package config

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:":8080"`
}
