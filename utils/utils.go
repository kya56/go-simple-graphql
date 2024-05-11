package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetValue(key string) string {
	var err error = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading .env file\n")
	}
	return os.Getenv(key)
}
