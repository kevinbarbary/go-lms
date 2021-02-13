package main

import (
	"./api"
	"./utils"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func html(w http.ResponseWriter, r *http.Request, user, title, breadcrumb, content string) {
	var html = "<!DOCTYPE html>"
	html = utils.Concat(html, "<html>")
	html = utils.Concat(html, head(w, title))
	html = utils.Concat(html, body(w, r, title, breadcrumb, user, content))
	html = utils.Concat(html, "</html>")
	io.WriteString(w, html)
}

func head(w http.ResponseWriter, title string) string {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page := "LMS - "
	if title == "" {
		page = utils.Concat(page, "home")
	} else {
		page = utils.Concat(page, title)
	}
	return utils.Concat(`<head>
<title>`, page, `</title>
<link rel="stylesheet" type="text/css" href="assets/css/custom.css?v=`, utils.Assets("CSS"), `">
</head>
`)
}

func body(w http.ResponseWriter, r *http.Request, title, breadcrumb, user, content string) string {
	var html = "<body>"
	html = utils.Concat(html, `<p class="breadcrumb">You are here: `, breadcrumb, "</p>")
	if user == "" {
		html = utils.Concat(html, "<p>", utils.Hyper("/", "Sign In"), "</p>")
	} else if user == api.TokenUser() {
		html = utils.Concat(html, "<p>", utils.Hyper("/sign-out", "Sign out"), "</p>")
	} else {
		html = utils.Concat(html, "<p>", utils.Hyper("/sign-out", utils.Concat("Sign out: ", user)), "</p>")
	}
	html = utils.Concat(html, "<h1>", title, "</h1>")
	if message := GetMessage(r); message != "" {
		html = utils.Concat(html, `<p class="message">`, message, `</p>`)
		UnsetMessage(w)
	}
	html = utils.Concat(html, content)

	// @todo - add footer here

	html = utils.Concat(html, "</body>")
	return html
}

type crumb struct {
	title, link string
}

func breadcrumbTrail(list []crumb) string {
	var trail = make([]string, len(list))
	for i, current := range list {
		if current.link == "" {
			trail[i] = current.title
		} else {
			trail[i] = utils.Hyper(current.link, current.title)
		}
	}
	return utils.Concat(`<span class="crumb">`, strings.Join(trail, `</span> &rarr; <span class="crumb">`), `</span>`)
}

func htmlEnrolStart() string {
	return `<div class="enrolments">`
}

func htmlEnrolRow(enrol api.UserEnrol, now api.Timestamp) string {

	var logo string
	if enrol.PublisherLogo != "" {
		logo = utils.Concat(`<img src="`, enrol.PublisherLogo, `" alt="`, enrol.Publisher, `">`)
	}

	var enrolStr, course, expires string
	enrolStr = strconv.Itoa(enrol.EnrollID)
	if enrol.EnrollStatus.Enabled() && now.BeforeEnd(enrol.EndDate) {
		// active
		course = utils.Hyper(utils.Concat("/", enrolStr), enrol.CourseTitle)
		//expires = fmt.Sprintf("Expires in %v", now.Until(enrol.EndDate))
		expires = utils.Concat("Expires in ", utils.FormatUntil(now.Until(enrol.EndDate)))
	} else {
		// expired, pending, etc.
		course = enrol.CourseTitle
		expires = "Expired"
	}

	var status string
	if enrol.Completed {
		status = "Completed"
	} else {
		status = "Incomplete"
	}

	return utils.Concat(`<div class="enrol-row" id="enrol-id-`, enrolStr, `"><div class="logo">`, logo, `</div><div class="enrol"><div class="title">`, course, `</div><div class="status">`, status, `</div><div class="expires">`, expires, `</div><div class="enrol-start">`, enrol.StartDate.ToDate(), `</div><div class="enrol-end">`, enrol.EndDate.ToDate(), `</div></div></div>`)
}

func htmlEnrolEnd() string {
	return "</div>"
}

func tutorialsHTML(data api.UserEnrolment) string {
	var html string
	for _, lesson := range data.Lessons {
		html = utils.Concat(html, "<h2>", lesson.Title, "</h2>")
		for _, tutorial := range lesson.Tutorials {
			html = utils.Concat(html, `<div class="tutorial-row">`, utils.Hyper(utils.Concat(tutorial.LaunchURL, "&returnHTTP=1&returnURL=", url.QueryEscape(utils.Concat("//", utils.Domain(), "/")), strconv.Itoa(data.EnrollID)), tutorial.TutorialTitle), `</div>`)
		}
	}
	return html
}
