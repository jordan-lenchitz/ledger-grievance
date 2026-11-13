package config

import (
	"os"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	Port             string
	LogLevel         string
	SLAGuaranteeMode bool
	StrictCompliance bool
}

func Load() *Config {
	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBUser:           getEnv("DB_USER", "ledger"),
		DBPassword:       getEnv("DB_PASSWORD", "ledgerpass"),
		DBName:           getEnv("DB_NAME", "grievance_ledger"),
		Port:             getEnv("PORT", "8000"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		SLAGuaranteeMode: getEnv("SLA_GUARANTEE_MODE", "false") == "true",
		StrictCompliance: getEnv("STRICT_COMPLIANCE", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
