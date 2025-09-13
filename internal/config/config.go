package config

import (
	"time"

	"github.com/cprakhar/gopher-social/internal/env"
	"github.com/cprakhar/gopher-social/internal/ratelimiter"
)

type Config struct {
	ApiURL      string
	Addr        string
	DB          DBConfig
	Env         string
	Version     string
	Mail        MailConfig
	WebURL      string
	Auth        authConfig
	Redis       redisConfig
	RateLimiter ratelimiter.Config
}

type redisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type authConfig struct {
	Basic basicConfig
	Token tokenConfig
}

type basicConfig struct {
	Username string
	Password string
}

type tokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
	Aud    string
}

type MailConfig struct {
	Exp    time.Duration
	ApiKey string
	Sender string
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
		Addr:   env.GetString("ADDR", ":8080"),
		ApiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		WebURL: env.GetString("WEB_URL", "http://localhost:3000"),
		DB: DBConfig{
			Addr:            env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gopher-social?sslmode=disable"),
			MaxConns:        int32(env.GetInt("DB_MAX_CONNS", 30)),
			MinConns:        int32(env.GetInt("DB_MIN_CONNS", 10)),
			MaxIdleTime:     env.GetDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
			MaxConnLifetime: env.GetDuration("DB_MAX_CONN_LIFETIME", time.Hour),
		},
		Env:     env.GetString("ENV", "development"),
		Version: env.GetString("VERSION", "0.0.1"),
		Mail: MailConfig{
			Exp:    env.GetDuration("MAIL_EXP", 3*24*time.Hour),
			ApiKey: env.GetString("MAIL_API_KEY", ""),
			Sender: env.GetString("MAIL_SENDER", ""),
		},
		Auth: authConfig{
			Basic: basicConfig{
				Username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				Password: env.GetString("BASIC_AUTH_PASSWORD", "adminpassword"),
			},
			Token: tokenConfig{
				Secret: env.GetString("AUTH_TOKEN", ""),
				Exp:    env.GetDuration("AUTH_TOKEN_EXP", 3*24*time.Hour),
				Iss:    env.GetString("AUTH_TOKEN_ISS", "gopher-social"),
				Aud:    env.GetString("AUTH_TOKEN_AUD", "gopher-social"),
			},
		},
		Redis: redisConfig{
			Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			Password: env.GetString("REDIS_PASSWORD", ""),
			DB:       env.GetInt("REDIS_DB", 0),
			Enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		RateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATELIMITER_ENABLED", true),
		},
	}
	return cfg
}
