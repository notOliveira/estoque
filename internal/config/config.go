package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	DBName   string
	Port     int
}

func LoadConfig() *Config {

	godotenv.Load()

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "storage"
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	return &Config{
		MongoURI: mongoUri,
		DBName:   dbName,
		Port:     port,
	}
}