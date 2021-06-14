package config

import "github.com/caarlos0/env/v6"

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
}

func New() (*Config, error) {
	conf := Config{}
	err := env.Parse(&conf)
	return &conf, err
}
