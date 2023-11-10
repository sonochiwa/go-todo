package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MongoDbUri     string
	DbName         string
	CollectionName string
	Port           string
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	return &Config{
		MongoDbUri:     getEnv("MONGODB_URI", ""),
		DbName:         getEnv("DB_NAME", ""),
		CollectionName: getEnv("COLLECTION_NAME", ""),
		Port:           getEnv("PORT", ""),
	}
}
