package main

import (
	"./api"
	"./utils"
	"log"
	"net/http"
	"strconv"
)

const SIGN_IN = "Sign In"

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/sign-out", signOut)
	http.HandleFunc("/", route)
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
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

	case "/courses":
		courses(w, r, 0)

	default:
		path := route[1:]

		// WIP: return from launch in modal
		if len(path) > 6 {
			if path[0:6] == "parent" {
				_, e := strconv.Atoi(path[7:])
				if e == nil {
					html(w, r, "*", "Please wait...", "Loading...", utils.Concat(`<script type="text/javascript">window.parent.href="/`, path[7:], `";</script>`))
					return
				}

				error404(w, r)
				return
			}
		}

		id, e := strconv.Atoi(path)
		if e == nil {
			if api.CheckSignedIn(r) {
				learn(w, r, id)
				// modal test...
				return
			}

			// modal test...
			signIn(w, r, path)
			return
		}

		error404(w, r)
	}
}

func GetError(err error) string {
	e := err.Error()
	if e == "" {
		e = "An unknown error had occurred"
	}
	return utils.Concat(`<span class="error">`, e, `</span>`)
}
