package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	DBName         string
	CollectionName string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	cfg := Config{
		MongoURI:       os.Getenv("MONGO_URI"),
		DBName:         os.Getenv("MONGO_DB"),
		CollectionName: os.Getenv("MONGO_COLLECTION"),
	}

	if cfg.MongoURI == "" {
		return Config{}, fmt.Errorf("MONGO_URI is required")
	}

	if cfg.DBName == "" {
		return Config{}, fmt.Errorf("MONGO_DB is required")
	}
	
	if cfg.CollectionName == "" {
		return Config{}, fmt.Errorf("MONGO_COLLECTION is required")
	}

	return cfg, nil
}
