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

func DomainAndPort() string {
	a, b := getTwo("Domain", "Port", "Domain")
log.Print("utils DomainAndPort...")
log.Print(Concat(a, ":", b))
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
	// get the subdomain if Multi-Site, otherwise empty string
	multi := get("MultiSite", "GetMultiSite")
	if multi == "" {
		return ""
	}
	sub := strings.Split(r.Host, ".")
	return sub[0]
}

func GetDomain(r *http.Request) string {
	h := r.Host
log.Print("utils GetDomain: ", h)
	return h
//log.Print("Domain: vlc.corelearn.net")
//	return "vlc.corelearn.net"
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
