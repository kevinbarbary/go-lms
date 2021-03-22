package utils

import "net/http"

func createCookie(name, value string) *http.Cookie {
	var age int
	if value == "" {
		age = -1
	} else {
		age = 99999
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   "localhost" // @todo - make this work on all domains ("*" ?)
		Path:     "/",
		MaxAge:   age,
		HttpOnly: true,
	}
}

func SaveCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, createCookie(name, value))
}

func GetCookieValue(r *http.Request, name string) (string, error) {
	var cookie, err = r.Cookie(name)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", err
	}
	return cookie.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, createCookie(name, ""))
}
