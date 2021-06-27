package config

import "github.com/caarlos0/env/v6"

type Config struct {
	DBPort        int    `env:"DB_PORT,required"`
	DBHost        string `env:"DB_HOST,required"`
	DBName        string `env:"DB_NAME,required"`
	DBUser        string `env:"DB_USER,required"`
	DBPassword    string `env:"DB_PASSWORD,required"`
	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func New() (*Config, error) {
	conf := Config{}
	err := env.Parse(&conf)
	return &conf, err
}
