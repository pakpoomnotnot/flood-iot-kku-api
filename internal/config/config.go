package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GinMode        string
	MongoURI       string
	MongoDatabase  string
	MongoCollection string
	CORSOrigins    []string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase:  getEnv("MONGODB_DATABASE", "flood_iot_kku"),
		MongoCollection: getEnv("MONGODB_COLLECTION", "sensor_data"),
		CORSOrigins:    strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
