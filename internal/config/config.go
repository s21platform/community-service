package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Metrics  Metrics
	Platform Platform
	School   School
	Logger   Logger
	Cache    Cache
	Kafka    Kafka
}

type Cache struct {
	Host string `env:"COMMUNITY_SERVICE_REDIS_HOST"`
	Port string `env:"COMMUNITY_SERVICE_REDIS_PORT"`
}

type Logger struct {
	Host string `env:"LOGGER_SERVICE_HOST"`
	Port string `env:"LOGGER_SERVICE_PORT"`
}

type Service struct {
	Port string `env:"COMMUNITY_SERVICE_PORT"`
	Name string `env:"COMMUNITY_SERVICE_NAME"`
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

type School struct {
	Host string `env:"SCHOOL_SERVICE_HOST"`
	Port string `env:"SCHOOL_SERVICE_PORT"`
}

type Kafka struct {
	Host  string `env:"KAFKA_HOST"`
	Port  string `env:"KAFKA_PORT"`
	LevelChangeTopic string `env:"PARTICIPANT_LEVEL_CHANGED"`
	ExpLevelChanged string `env:"PARTICIPANT_EXP_LEVEL_CHANGED"`
	StatusChanged string `env:"PARTICIPANT_EXP_LEVEL_CHANGED"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("failed to read env variables: %s", err)
	}
	return cfg
}
