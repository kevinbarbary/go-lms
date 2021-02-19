package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Assets(kind string) string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in get Assets()")
	}
	return os.Getenv(kind)
}

func Domain() string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in get Domain()")
	}
	return Concat(os.Getenv("Domain"), ":", os.Getenv("Port"))
}

func Logo() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in get Logo()")
	}
	return os.Getenv("SiteLogo"), os.Getenv("SiteName")
}

func Endpoint(path string) string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in get Endpoint()")
	}
	return Concat(os.Getenv("API"), path)
}

func Creds() (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in get Creds()")
	}
	return os.Getenv("SiteID"), os.Getenv("SiteKey"), nil
}
