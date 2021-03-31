package html

import (
	api "github.com/kevinbarbary/go-lms/api"
	utils "github.com/kevinbarbary/go-lms/utils"
	"html/template"
	"net/http"
)

const SIGN_IN = "Sign In"
const SIGN_OUT = "Sign Out"
const LEARN = "Enrolments"
const COURSES = "Courses"
const ERROR = "Error"
const PLAIN = "Plain"

type Page struct {
	Kind, Header string
}

func Webpage(w http.ResponseWriter, r *http.Request, user string, page Page, breadcrumb, content string) {

	title := "LMS - "
	if page.Header == "" {
		title = utils.Concat(title, "home")
	} else if page.Kind == SIGN_IN {
		title = utils.Concat(title, SIGN_IN)
	} else {
		title = utils.Concat(title, page.Header)
	}

	var css, signInOut, menu, menuSpacing, siteItem, learnItem, learnOutline, learnDisabled, coursesOutline, coursesDisabled string

	multi := utils.GetMultiSite(r)

	switch page.Kind {
	case SIGN_IN:
		css = `<link rel="stylesheet" type="text/css" href="/assets/css/sign-in.css">`
		coursesOutline = "outline-"
		if multi != "" {
			siteItem = utils.Concat(`<a href="/" class="btn btn-primary btn-sm disabled">`, SIGN_IN, ` - `, multi, `</a>`)
		}
	case LEARN:
		coursesOutline = "outline-"
		if page.Header == LEARN {
			learnDisabled = " disabled"
		}
	case COURSES:
		learnOutline = "outline-"
		coursesDisabled = " disabled"
	case ERROR:
		learnOutline = "outline-"
		coursesOutline = "outline-"
	}

	if page.Kind != PLAIN {
		class := "btn btn-outline-primary btn-sm"
		if multi != "" {
			menuSpacing = " ms-3"
		}
		if user == "" {
			if page.Kind != SIGN_IN {
				if multi == "" {
					signInOut = utils.HyperClass("/", SIGN_IN, class)
					menuSpacing = " ms-3"
				} else {
					signInOut = utils.HyperClass("/", utils.Concat(SIGN_IN, " - ", multi), class)
				}
			}
		} else {
			if page.Kind != SIGN_IN && multi == "" {
				menuSpacing = " ms-3"
			}
			if user == api.TokenUser() {
				signInOut = utils.HyperClass("/sign-out", SIGN_OUT, class)
			} else {
				signInOut = utils.HyperClass("/sign-out", utils.Concat(SIGN_OUT, ": ", user), class)
			}
			learnItem = utils.Concat(`<a href="/" class="btn btn-`, learnOutline, `primary btn-sm`, learnDisabled, `">Learn</a>`)
		}

		menu = utils.Concat(siteItem, `<div class="btn-group`, menuSpacing, `" role="group" aria-label="Menu">`,
			learnItem, `<a href="/courses" class="btn btn-`, coursesOutline, `primary btn-sm`, coursesDisabled, `">Browse</a></div>`)

		menu = utils.Concat(`<div class="menu mb-3">`, signInOut, menu, `</div>`)
	}

	var body string
	if message, kind := GetMessage(r); message != "" {
		body = utils.Concat(body, StyleMessage(message, kind))
		UnsetMessage(w, r)
	}
	body = utils.Concat(body, content)

	var logo, name string
	if logo, name = utils.Logo(); logo != "" {
		if name == "" {
			logo = utils.Concat(`<img src="`, logo, `" alt="logo">`)
		} else {
			logo = utils.Concat(`<img src="`, logo, `" alt="`, name, `">`)
		}
	} else if name != "" {
		logo = utils.Concat(`<mark>`, name, `</mark>`)
	}

	// @todo - add footer

	data := struct {
		Title      string
		Css        template.HTML
		Version    string
		Breadcrumb template.HTML
		Logo       template.HTML
		Menu       template.HTML
		Header     string
		Content    template.HTML
	}{
		title,
		template.HTML(css),
		utils.Assets("CSS"),
		template.HTML(utils.Concat(`<span class="breadcrumb-trail breadcrumb-prefix">You are here:</span>`, breadcrumb)),
		template.HTML(logo),
		template.HTML(menu),
		page.Header,
		template.HTML(body),
	}

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
