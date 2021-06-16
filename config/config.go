package config

import "github.com/caarlos0/env/v6"

type Config struct {
	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
}

func New() (*Config, error) {
	conf := Config{}
	err := env.Parse(&conf)
	return &conf, err
}
