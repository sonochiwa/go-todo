package config

import (
	"os"
)

type Config struct {
	MongoDbUri     string
	DbName         string
	CollectionName string
	Port           string
}

func New() *Config {
	return &Config{
		MongoDbUri:     os.Getenv("MONGODB_URI"),
		DbName:         os.Getenv("DB_NAME"),
		CollectionName: os.Getenv("COLLECTION_NAME"),
		Port:           os.Getenv("PORT"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
