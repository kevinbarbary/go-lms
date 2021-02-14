package main

import (
	"./api"
	"./utils"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func html(w http.ResponseWriter, r *http.Request, user, header, breadcrumb, content string) {

	title := "LMS - "
	if header == "" {
		title = utils.Concat(title, "home")
	} else {
		title = utils.Concat(title, header)
	}

	var menu string
	class := "btn btn-primary btn-sm"
	if user == "" {
		if header != "Sign In" {
			menu = utils.HyperClass("/", "Sign In", class, "button")
		}
	} else if user == api.TokenUser() {
		menu = utils.HyperClass("/sign-out", "Sign out", class, "button")
	} else {
		menu = utils.HyperClass("/sign-out", utils.Concat("Sign out: ", user), class, "button")
	}
	if menu != "" {
		menu = utils.Concat(`<p class="menu">`, menu, `</p>`)
	}

	var body string
	if message, kind := GetMessage(r); message != "" {
		body = utils.Concat(body, StyleMessage(message, kind))
		UnsetMessage(w)
	}
	body = utils.Concat(body, content)

	// @todo - add footer

	data := struct {
		Title      string
		Css        string
		Breadcrumb template.HTML
		Menu       template.HTML
		Header     string
		Content    template.HTML
	}{
		title,
		utils.Assets("CSS"),
		template.HTML(utils.Concat(`<span class="breadcrumb-trail breadcrumb-prefix">You are here:</span>`, breadcrumb)),
		template.HTML(menu),
		header,
		template.HTML(body),
	}

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type crumb struct {
	title, link string
}

func breadcrumbTrail(list []crumb) string {
	var trail = make([]string, len(list))
	for i, current := range list {
		if current.link == "" {
			trail[i] = utils.Concat(`<li class="breadcrumb-item active" aria-current="page">`, current.title, `</li>`)
		} else {
			trail[i] = utils.Concat(`<li class="breadcrumb-item">`, utils.Hyper(current.link, current.title), `</li>`)
		}
	}
	return utils.Concat(`<nav class="breadcrumb-trail" id="breadcrumb" aria-label="breadcrumb"><ol class="breadcrumb">`, strings.Join(trail, ""), `</ol></nav>`)
}

func htmlEnrolStart() string {
	return `<div class="enrolments">`
}

func htmlEnrolRow(enrol api.UserEnrol, now api.Timestamp) string {

	var logo string
	if enrol.PublisherLogo != "" {
		logo = utils.Concat(`<img src="`, enrol.PublisherLogo, `" alt="`, enrol.Publisher, `">`)
	}

	var completeStatus, completeClass string
	if enrol.Completed {
		completeStatus = "Completed"
		completeClass = " completed"
	} else {
		completeStatus = "Incomplete"
		completeClass = " incomplete"
	}

	hyper := false
	var enrolStr, expires, expiryClass string
	enrolStr = strconv.Itoa(enrol.EnrollID)
	if enrol.EnrollStatus.Enabled() && now.BeforeEnd(enrol.EndDate) {
		// active
		expires = utils.Concat("Expires in ", utils.FormatUntil(now.Until(enrol.EndDate)))
		hyper = true
		if enrol.Completed {
			completeStatus = `<span class="badge bg-success">Completed</span>`
		}
	} else {
		// expired, pending, etc.
		expires = "Expired"
		expiryClass = " expired"
	}

	row := utils.Concat(`<div class="border p-3 mb-3`, expiryClass, completeClass, `" id="enrol-id-`, enrolStr, `"><div class="logo">`, logo, `</div><div class="enrol"><div class="title">`, enrol.CourseTitle, `</div><div class="status">`, completeStatus, `</div><div class="enrol-start">Start Date: `, enrol.StartDate.ToDate(), `</div><div class="expires">`, expires, `</div></div></div>`)

	if hyper {
		row = utils.Hyper(utils.Concat("/", enrolStr), row)
	}

	return utils.Concat(`<div class="enrol-row">`, row, `</div>`)
}

func htmlEnrolEnd() string {
	return "</div>"
}

func tutorialsHTML(data api.UserEnrolment) string {
	var html, completeStr string
	for _, lesson := range data.Lessons {
		if lesson.Title != data.CourseTitle {
			html = utils.Concat(html, "<h2>", lesson.Title, "</h2>")
		}
		for _, tutorial := range lesson.Tutorials {
			if tutorial.Completed {
				completeStr = `<span class="badge bg-success">Completed</span>`
			} else {
				completeStr = ""
			}
			html = utils.Concat(html, `<div class="tutorial-row">`, utils.Hyper(utils.Concat(tutorial.LaunchURL, "&returnHTTP=1&returnURL=", url.QueryEscape(utils.Concat("//", utils.Domain(), "/")), strconv.Itoa(data.EnrollID)), utils.Concat(`<div class="border p-2 mb-2"><div class="name">`, tutorial.TutorialTitle, `</div><div class="status">`, completeStr, `</div></div>`)), `</div>`)
		}
	}
	return html
}
