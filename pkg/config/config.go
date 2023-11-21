package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Server struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type MongoDB struct {
	MongoDbUri     string
	DbName         string
	CollectionName string
}

type Config struct {
	Server  Server
	MongoDB MongoDB
}

func init() {
	loadConfig()
}

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func GetConfig() *Config {
	return &Config{
		Server{
			Host:         getEnv("HOST", "0.0.0.0"),
			Port:         getEnv("PORT", "9000"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 60),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 60),
			IdleTimeout:  getEnvAsInt("IDLE_TIMEOUT", 60),
		},
		MongoDB{
			MongoDbUri:     getEnv("MONGODB_URI", ""),
			DbName:         getEnv("DB_NAME", ""),
			CollectionName: getEnv("COLLECTION_NAME", ""),
		},
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
