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
<style type="text/css">
body {
	font-family: "Open Sans", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", Helvetica, Arial, sans-serif;
}
.breadcrumb {
	font-size: xx-small;
}
.message {
	font-size: x-large;
}
.enrol-row {
	clear: both;
}
.enrol {
	display: grid;
}
.title {
	font-size: larger;
}
.logo {
	display: grid;
	float: left;
	margin-right: 10px;
}
</style>
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

	// @todo - make a function like now.ToDatetime() to show how long remaining until expiry, e.g. 'Expires in 3 days', or 'Expires in 5 hours and 10 minutes'

	var logo string
	if enrol.PublisherLogo != "" {
		logo = utils.Concat(`<img src="`, enrol.PublisherLogo, `" alt="`, enrol.Publisher, `">`)
	}

	var course string
	if enrol.EnrollStatus.Enabled() && now.BeforeEnd(enrol.EndDate) {
		// active
		course = utils.Hyper(utils.Concat("/", strconv.Itoa(enrol.EnrollID)), enrol.CourseTitle)
	} else {
		// expired, pending, etc.
		course = enrol.CourseTitle
	}

	var status string
	if enrol.Completed {
		status = "Completed"
	} else {
		status = "Incomplete"
	}

	return utils.Concat(`<div class="enrol-row"><div class="logo">`, logo, `</div><div class="enrol"><div class="title">`, course, `</div><div class="status">`, status, `</div><div class="enrol-start">`, enrol.StartDate.ToDate(), `</div><div class="enrol-end">`, enrol.EndDate.ToDate(), `</div></div></div>`)
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
