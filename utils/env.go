package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Domain() string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in getDomain()")
	}
	return os.Getenv("Domain")
}

func Endpoint(path string) string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in getEndpoint()")
	}
	return Concat(os.Getenv("API"), path)
}

func Creds() (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in getCreds()")
	}
	return os.Getenv("SiteID"), os.Getenv("SiteKey"), nil
}
