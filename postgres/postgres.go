package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func Initialize() (*gorm.DB, error) {
	dbHost := getEnvWithFallback("DB_HOST", "localhost")
	dbPort := getEnvWithFallback("DB_PORT", "5432")
	dbName := getEnvWithFallback("DB_NAME", "postgres")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	return gorm.Open(postgres.Open(dsn))
}

func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
