package utils

import (
	godotenv "github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(kind, id string) string {
	err := godotenv.Load()
	if err != nil {
		log.Print(Concat("Error loading .env file in get ", id, "()"))
	}
	return os.Getenv(kind)
}

func getTwo(one, two, id string) (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Print(Concat("Error loading .env file in get ", id, "()"))
	}
	return os.Getenv(one), os.Getenv(two)
}

func getThree(one, two, three, id string) (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Print(Concat("Error loading .env file in get ", id, "()"))
	}
	return os.Getenv(one), os.Getenv(two), os.Getenv(three)
}

func Assets(kind string) string {
	return get(kind, "Assets")
}

func Domain() string {
	a, b := getTwo("Domain", "Port", "Domain")
	return Concat(a, ":", b)
}

func Logo() (string, string) {
	return getTwo("SiteLogo", "SiteName", "Logo")
}

func Endpoint(path string) string {
	return Concat(get("API", "Endpoint"), path)
}

func Creds(site string) (string, string) {
	id, multi, key := getThree("SiteID", "MultiSite", "SiteKey", "Creds")
	if id == multi {
		return id, site
	}
	return id, key
}

func GetMultiSite(r *http.Request) string {
	site := get("MultiSite", "GetMultiSite")
	if site == "" {
		return ""
	}
	domain := strings.Split(r.Host, ".")
	return domain[0]
}

func GetSite(r *http.Request) string {
	// Multi-Site: get the SiteID from the subdomain, can be overridden with a SiteMapper entry in .env
	domain := strings.Split(r.Host, ".")
	site := domain[0]
	kind := Concat("SiteMapper-", site)
	if mapped := get(kind, "GetSite"); mapped != "" {
		return mapped
	}
	return site
}
