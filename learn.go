package main

import (
	"./api"
	"./utils"
	"net/http"
	"strconv"
)

func learn(w http.ResponseWriter, r *http.Request, enrollId int) {

	var enrolStr string
	if enrollId > 0 {
		enrolStr = strconv.Itoa(enrollId)
	}

	if r.Method == "POST" {
		signIn(w, r, enrolStr)
		return
	} else {

		var user, breadcrumb, title, content string

		token := api.GetToken(r)

		if !api.CheckTokenSignedIn(token) {
			signIn(w, r, enrolStr)
			return
		}

		if enrollId == 0 {

			enrol, newToken, u, now, err := api.UserEnrolments(token, api.TokenUser())
			if err == nil {

				api.SaveToken(w, newToken)

				if u == "" {
					signIn(w, r, "")
					return
				}

				user = u
				breadcrumb = breadcrumbTrail([]crumb{{"Enrolments", ""}})
				title = "Enrolments"
				content = enrolHTML(enrol, now)

			} else {
				user = api.GetSignedInTokenFlag(token)
				breadcrumb = breadcrumbTrail([]crumb{{"Enrolments", "/"}, {"Error", ""}})
				title = "Error"
				content = utils.Concat("<p>", GetError(err), "</p>")
			}

		} else {

			enrol, newToken, u, err := api.UserTutorials(token, api.TokenUser(), enrollId)
			if err == nil {

				api.SaveToken(w, newToken)

				if u == "" {
					signIn(w, r, strconv.Itoa(enrollId))
					return
				} else {
					user = u
					if enrol.NotValid() {

						// @todo - call /enrolment/history/{LoginID}/{EnrollID}

						breadcrumb = breadcrumbTrail([]crumb{{"Enrolments", "/"}, {"Enrolment not found", ""}})
						title = "Error"
						content = "<p>Enrolment not found</p>"
					} else {
						breadcrumb = breadcrumbTrail([]crumb{{"Enrolments", "/"}, {enrol.CourseTitle, ""}})
						title = enrol.CourseTitle
						content = tutorialsHTML(enrol)
					}
				}

			} else {
				var home string
				user = api.GetSignedInTokenFlag(token)
				if user == "" {
					home = SIGN_IN
				} else {
					home = "Enrolments"
				}
				breadcrumb = breadcrumbTrail([]crumb{{home, "/"}, {"Error", ""}})
				title = "Error"
				content = utils.Concat("<p>", GetError(err), "</p>")
			}

		}

		html(w, r, user, title, breadcrumb, content)
	}
}

func enrolHTML(data []api.UserEnrol, now api.Timestamp) string {
	var html string
	for _, enrol := range data {
		html = utils.Concat(html, htmlEnrolRow(enrol, now))
	}
	if html == "" {
		html = "<p>You do not have any enrolments</p>"
	}
	return utils.Concat(htmlEnrolStart(), html, htmlEnrolEnd())
}
