package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                    string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	DBSSLMode               string
	SuperTokensConnectionURI string
	SuperTokensAPIKey       string
}

func Load() *Config {
	return &Config{
		Port:                    getEnv("PORT", "8080"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "postgres"),
		DBPassword:              getEnv("DB_PASSWORD", ""),
		DBName:                  getEnv("DB_NAME", "userdb"),
		DBSSLMode:               getEnv("DB_SSLMODE", "disable"),
		SuperTokensConnectionURI: getEnv("SUPERTOKENS_CONNECTION_URI", "http://localhost:3567"),
		SuperTokensAPIKey:       getEnv("SUPERTOKENS_API_KEY", ""),
	}
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
