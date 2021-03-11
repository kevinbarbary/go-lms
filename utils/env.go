package utils

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
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

func Creds(site string) (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file in get Creds()")
	}
	id := os.Getenv("SiteID")
	if id == os.Getenv("MultiSite") {
		return id, site, nil
	}
	return id, os.Getenv("SiteKey"), nil
}

func GetSite(r *http.Request) string {
	domain := strings.Split(r.Host, ".")
	return domain[0]
}
