package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type key string

const KeyMetrics = key("metrics")
const KeyUUID = key("uuid")

type Config struct {
	Service  Service
	Postgres Postgres
	Metrics  Metrics
	Platform Platform
}

type Service struct {
	Port string `env:"COMMUNITY_SERVICE_PORT"`
}

type Postgres struct {
	User     string `env:"COMMUNITY_SERVICE_POSTGRES_USER"`
	Password string `env:"COMMUNITY_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"COMMUNITY_SERVICE_POSTGRES_DB"`
	Host     string `env:"COMMUNITY_SERVICE_POSTGRES_HOST"`
	Port     string `env:"COMMUNITY_SERVICE_POSTGRES_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
