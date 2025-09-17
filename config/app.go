package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT, DB_URL, DB_NAME, SECRETKEY string
)

func InitEnv() error {
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("load env: %w", err)
		}
	} else {
		_ = godotenv.Load()
	}

	if v := os.Getenv("PORT"); v != "" {
		PORT = v
	} else {
		PORT = "8080"
	}
	if v := os.Getenv("DB_URL"); v == "" {
		DB_URL = "mongodb://localhost:27017/"
	} else {
		DB_URL = v
	}
	if v := os.Getenv("DB_NAME"); v == "" {
		DB_NAME = "cctv_db"
	} else {
		DB_NAME = v
	}
	if v := os.Getenv("SECRETKEY"); v == "" {
		SECRETKEY = ""
	} else {
		SECRETKEY = v
	}
	return nil
}

func Init() {
	if err := InitEnv(); err != nil {
		log.Fatal(err)
	}
}
