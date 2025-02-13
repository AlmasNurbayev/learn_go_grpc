package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	DSN      string        `yaml:"DSN"`
	GRPC     GRPCConfig    `yaml:"grpc" env-required:"true"`
}

type dsnStruct struct {
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
	configPath := fetchConfigPath()
	cfg := ReadConfigByPath(configPath)

	var dsnStruct dsnStruct

	// если в конфиге нет DSN, то читаем env для получения DSN
	if cfg.DSN == "" {
		envPath := ".env"
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			log.Fatalf(".env not found : %s", envPath)
		}
		err := cleanenv.ReadConfig(envPath, dsnStruct)
		if err != nil {
			panic("Failed to read config from env: " + err.Error())
		}
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		db := os.Getenv("POSTGRES_DB")
		port := os.Getenv("POSTGRES_PORT")
		cfg.DSN = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + db + "?sslmode=disable"
	}
	return &cfg
}

func ReadConfigByPath(path string) Config {
	cfg := Config{}

	if _, err := os.Stat(path); err != nil {
		panic("Config file not found on path: " + path)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}
	cfg.DSN = os.Getenv("DSN")

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
