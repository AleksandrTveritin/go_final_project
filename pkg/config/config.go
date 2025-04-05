package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   int
	DBFile string
}

var AppConfig *Config

func LoadConfig() error {
	_ = godotenv.Load(".env")

	AppConfig = &Config{
		Port:   7540,           // Значение по умолчанию
		DBFile: "scheduler.db", // Значение по умолчанию
	}

	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Printf("Invalid port value: %s, using default %d", portStr, AppConfig.Port)
		} else {
			AppConfig.Port = port
		}
	}

	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		AppConfig.DBFile = dbFile
	}

	return nil
}
