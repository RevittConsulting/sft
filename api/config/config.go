package config

import (
	"os"
)

// DatabaseConfig structure
type DatabaseConfig struct {
	Dialect, Host, Port, Username, Password, DBName, TestDBName, SSLMode string
}

// Config structure
type Config struct {
	Database *DatabaseConfig
}

// GetDatabaseConfig function returns database configuration
func GetDatabaseConfig() *Config {
	return &Config{
		Database: &DatabaseConfig{
			Dialect:    GetEnv("DB_DIALECT", ""),
			Host:       GetEnv("DB_HOST", ""),
			Port:       GetEnv("DB_PORT", ""),
			Username:   GetEnv("DB_USERNAME", ""),
			Password:   GetEnv("DB_PASSWORD", ""),
			DBName:     GetEnv("DB_NAME", ""),
			TestDBName: GetEnv("TEST_DB_NAME", ""),
			SSLMode:    GetEnv("DB_SSL_MODE", ""),
		},
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
