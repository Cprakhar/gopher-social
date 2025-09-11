package config

import (
	"time"

	"github.com/cprakhar/gopher-social/internal/env"
)

type Config struct {
	ApiURL  string
	Addr    string
	DB      DBConfig
	Env     string
	Version string
}

type DBConfig struct {
	Addr            string
	MaxConns        int32
	MinConns        int32
	MaxIdleTime     time.Duration
	MaxConnLifetime time.Duration
}

func Load() Config {
	cfg := Config{
		Addr: env.GetString("ADDR", ":8080"),
		ApiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		DB: DBConfig{
			Addr:            env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gopher-social?sslmode=disable"),
			MaxConns:        int32(env.GetInt("DB_MAX_CONNS", 30)),
			MinConns:        int32(env.GetInt("DB_MIN_CONNS", 10)),
			MaxIdleTime:     env.GetDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
			MaxConnLifetime: env.GetDuration("DB_MAX_CONN_LIFETIME", time.Hour),
		},
		Env:     env.GetString("ENV", "development"),
		Version: env.GetString("VERSION", "0.0.1"),
	}
	return cfg
}
