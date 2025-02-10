package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	DSN      string        `yaml:"dsn"`
	GRPC     GRPCConfig    `yaml:"grpc" env-required:"true"`
}

type DSNStruct struct {
	POSTGRES_USER     string `env:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD"`
	POSTGRES_DB       string `env:"POSTGRES_DB"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT"`
	POSTGRES_HOST     string `env:"POSTGRES_HOST"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	cfg := Config{}

	path := "./config/local.yaml"
	if _, err := os.Stat(path); err != nil {
		panic("Config file not found on path: " + path)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	var user, password, dbName, port, postgresHost string

	if os.Getenv("POSTGRES_USER") == "" {
		// читаем env если есть, для локального запуска
		envPath := ".env"
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			log.Fatalf("Файл .env не найден: %s", envPath)
		}
		var DSN DSNStruct
		err := cleanenv.ReadConfig(envPath, &DSN)
		if err != nil {
			panic("Failed to read config: " + err.Error())
		}
		user = DSN.POSTGRES_USER
		password = DSN.POSTGRES_PASSWORD
		dbName = DSN.POSTGRES_DB
		port = DSN.POSTGRES_PORT
		postgresHost = DSN.POSTGRES_HOST
	} else {
		user = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbName = os.Getenv("POSTGRES_DB")
		port = os.Getenv("POSTGRES_PORT")
		postgresHost = os.Getenv("POSTGRES_HOST")

	}

	// читаем переменные из окружения

	cfg.DSN = "postgres://" + user + ":" + password + "@" + postgresHost + ":" + port + "/" + dbName + "?sslmode=disable"

	return &cfg
}
