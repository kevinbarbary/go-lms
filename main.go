package main

import (
	"./api"
	"./utils"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", route)
	http.ListenAndServe(":8888", nil)
}

func route(w http.ResponseWriter, r *http.Request) {
	switch route := r.URL.Path; route {
	case "/":
		if api.CheckSignedIn(r) {
			learn(w, r, 0)
		} else {
			signIn(w, r, "")
		}
	case "/sign-in":
		signIn(w, r, "")
	case "/sign-out":
		signOut(w, r)
	default:
		path := route[1:]
		id, e := strconv.Atoi(path)
		if e == nil {
			if api.CheckSignedIn(r) {
				learn(w, r, id)
			} else {
				signIn(w, r, path)
			}
		} else {
			error404(w, r)
		}
	}
}

func GetError(err error) string {
	e := err.Error()
	if e == "" {
		e = "An unknown error had occurred"
	}
	return utils.Concat(`<span class="error">`, e, `</span>`)
}
