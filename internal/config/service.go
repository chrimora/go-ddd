package config

import "fmt"

type Env string

const (
	LocalEnvEnum Env = "local"
	TestEnvEnum  Env = "test"
	DevEnvEnum   Env = "dev"
	ProdEnvEnum  Env = "prod"
)

type ServiceConfig struct {
	Name string `env:"SERVICE_NAME" envDefault:"service"`
	Env  Env    `env:"SERVICE_ENV" envDefault:"local"`
}

func (s ServiceConfig) IsLocal() bool { return s.Env == LocalEnvEnum }
func (s ServiceConfig) IsTest() bool  { return s.Env == TestEnvEnum }
func (s ServiceConfig) IsDev() bool   { return s.Env == DevEnvEnum }
func (s ServiceConfig) IsProd() bool  { return s.Env == ProdEnvEnum }

func (e *Env) UnmarshalText(text []byte) error {
	switch string(text) {
	case string(LocalEnvEnum), string(TestEnvEnum), string(DevEnvEnum), string(ProdEnvEnum):
		*e = Env(text)
		return nil
	default:
		return fmt.Errorf("invalid SERVICE_ENV: %s", text)
	}
}
