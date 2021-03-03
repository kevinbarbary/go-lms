package main

import (
	"./api"
	"./utils"
	"log"
	"net/http"
	"strconv"
)

const SIGN_IN = "Sign In"
const LEARN = "Enrolments"
const COURSES = "Courses"
const ERROR = "Error"
const PLAIN = "Plain"

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
		return

	case "/sign-in":
		signIn(w, r, "")
		return

	case "/courses":
		courses(w, r, 1)
		return

	default:
		path := route[1:]

		// paginated courses
		if len(path) > 13 && path[0:13] == "courses/page-" {
			index, e := strconv.Atoi(path[13:])
			if e == nil {
				courses(w, r, index)
				return
			}

			// @todo - add a more specific error message to the 404 page
			error404(w, r)
			return
		}

		// WIP: return from launch in modal
		if len(path) > 7 && path[0:7] == "parent/" {
			_, e := strconv.Atoi(path[7:])
			if e == nil {
				html(w, r, "", page{PLAIN, "Please wait..."}, "Loading...", utils.Concat(`<script type="text/javascript">window.parent.href="/`, path[7:], `";</script>`))
				return
			}

			// @todo - add a more specific error message to the 404 page
			error404(w, r)
			return
		}

		id, e := strconv.Atoi(path)
		if e == nil {
			if api.CheckSignedIn(r) {
				learn(w, r, id)
				return
			}
			signIn(w, r, path)
			return
		}

		error404(w, r)
		return
	}
}

func GetError(err error) string {
	e := err.Error()
	if e == "" {
		e = "An unknown error had occurred"
	}
	return utils.Concat(`<span class="error">`, e, `</span>`)
}
