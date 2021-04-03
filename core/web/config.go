package web

import (
	"fmt"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Dialect  string
}

func (db DatabaseConfig) ConnectionInfo() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name,
	)
}

type SentryConfig struct {
	Dns string
}

var appConf *AppConfig

type AppConfig struct {
	IsCorrectGenerated bool
	Port               int
	Env                string
	Pepper             string
	HMACKey            string
	Database           DatabaseConfig
	Sentry             SentryConfig
}

func (c AppConfig) IsProd() bool {
	return c.Env == "prod"
}

func NewConfig() *AppConfig {
	InitEnv()

	sentryConf := SentryConfig{
		Dns: getEnv("SENTRY_WORKER_DNS", ""),
	}

	dbConf := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		Name:     getEnv("DB_NAME", "camp"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		Dialect:  getEnv("DB_DIALECT", "postgres"),
	}

	appConf = &AppConfig{
		IsCorrectGenerated: true,
		Port:               getEnvAsInt("WEB_PORT", 3000),
		Env:                getEnv("WEB_ENV", "dev"),
		Pepper:             getEnv("PEPPER", "secret-random-string"),
		HMACKey:            getEnv("HMACKey", "secret-hmac-key"),
		Database:           dbConf,
		Sentry:             sentryConf,
	}
	return appConf
}

func LoadConfig() *AppConfig {
	if &appConf != nil || !appConf.IsCorrectGenerated {
		return NewConfig()
	}
	return appConf
}
