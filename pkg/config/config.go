package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	MongoDbUri     string
	DbName         string
	CollectionName string
	Port           string
}

func init() {
	loadMongo()
}

func loadMongo() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func GetConfig() *Config {
	return &Config{
		MongoDbUri:     getEnv("MONGODB_URI", ""),
		DbName:         getEnv("DB_NAME", ""),
		CollectionName: getEnv("COLLECTION_NAME", ""),
		Port:           getEnv("PORT", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
