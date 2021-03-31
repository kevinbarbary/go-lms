package utils

import (
	"net/http"
	"log"
)

func createCookie(name, value, domain string) *http.Cookie {
	var age int
	if value == "" {
		age = -1
	} else {
		age = 99999
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Path:     "/",
		MaxAge:   age,
//		HttpOnly: true,
	}
}

func SaveCookie(w http.ResponseWriter, name, value, domain string) {
	http.SetCookie(w, createCookie(name, value, domain))
}

func GetCookieValue(r *http.Request, name string) (string, error) {
log.Print("GetCookieValue...")
log.Print(r.Host)
log.Print("...GetCookieValue...")
	var cookie, err = r.Cookie(name)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", err
	}
log.Print(cookie.Value)
log.Print("...GetCookieValue")
	return cookie.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name, domain string) {
	http.SetCookie(w, createCookie(name, "", domain))
}
