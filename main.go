package main

import (
	api "github.com/kevinbarbary/go-lms/api"
	html "github.com/kevinbarbary/go-lms/html"
	utils "github.com/kevinbarbary/go-lms/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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
		html.Courses(w, r, 1, []int{})
		return

	default:
		path := route[1:]

		// paginated courses
		if len(path) > 13 && path[0:13] == "courses/page-" {

			var pageParam string
			var tags []int
			i := strings.Index(path[13:], "/")
			if i > 0 {
				// page number and tags in URL
				pageParam = path[13 : 13+i]
				tagsParam := path[13+i+1:]
				if len(tagsParam) > 4 && tagsParam[:4] == "tag-" {
					var err bool
					tags, err = utils.Atoi(strings.Split(tagsParam[4:], "-"))
					if err {
						log.Print("Error: invalid CourseTags list querystring parameter '", tagsParam[4:], "'")
					}
				}
			} else {
				// just page number in URL
				pageParam = path[13:]
			}

			index, e := strconv.Atoi(pageParam)
			if e == nil {
				html.Courses(w, r, index, tags)
				return
			}

			error404(w, r, "courses", `The <span class="text-secondary">courses</span> page could not be found. Go <a href="javascript:history.back()">back</a> and try again.`, []html.Crumb{{"Courses", "/courses"}, {"Page Not Found", ""}})
			return
		}

		// WIP: return from launch in modal
		if len(path) > 7 && path[0:7] == "parent/" {
			_, e := strconv.Atoi(path[7:])
			if e == nil {
				html.Webpage(w, r, "", html.Page{html.PLAIN, "Please wait..."}, "Loading...", utils.Concat(`<script type="text/javascript">window.parent.href="/`, path[7:], `";</script>`))
				return
			}

			error404(w, r, "enrolment", `Your <span class="text-secondary">enrolment</span> could not be found. Go <a href="/">home</a> and try again.`, []html.Crumb{{"Enrolments", "/"}, {"Page Not Found", ""}})
			return
		}

		// enrolment
		id, e := strconv.Atoi(path)
		if e == nil {
			if api.CheckSignedIn(r) {
				learn(w, r, id)
				return
			}
			signIn(w, r, path)
			return
		}

		error404(w, r, "", `Your page could not be found. Go <a href="/">home</a> and try again.`, []html.Crumb{{"Home", "/"}, {"Page Not Found", ""}})
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
