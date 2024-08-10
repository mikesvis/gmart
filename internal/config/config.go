package config

import (
	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccuralSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() *Config {
	var config Config

	parseFlags(&config)
	env.Parse(&config)

	return &config
}

func parseFlags(c *Config) {
	flag.StringVarP(&c.RunAddress, "address", "a", "localhost:8080", "хост и порт запуска сервиса")
	flag.StringVarP(&c.DatabaseURI, "database_uri", "d", "", "адрес подключения к базе данных")
	flag.StringVarP(&c.AccuralSystemAddress, "accrural_address", "r", "", "адрес системы расчёта начислений")
	flag.Parse()
}
