package app

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	Listen   string `env:"RUN_ADDRESS"`
	Database string `env:"DATABASE_URI"`
	Accrual  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func newConfig() config {
	var cfg config

	flag.StringVar(&cfg.Listen, "a", "", "host:port to listen on")
	flag.StringVar(&cfg.Database, "d", "", "database address")
	flag.StringVar(&cfg.Accrual, "r", "", "accrual system address")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
