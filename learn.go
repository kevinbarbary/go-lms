package main

import (
	"./api"
	"./html"
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
				breadcrumb = html.BreadcrumbTrail([]html.Crumb{{"Enrolments", ""}})
				title = "Enrolments"
				content = html.Enrol(enrol, now)

			} else {
				user = api.GetSignedInTokenFlag(token)
				breadcrumb = html.BreadcrumbTrail([]html.Crumb{{"Enrolments", "/"}, {html.ERROR, ""}})
				title = html.ERROR
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
					if enrol.NotValid() {

						// @todo - call /enrolment/history/{LoginID}/{EnrollID}

						error404(w, r, "enrolment", `Your <span class="text-secondary">enrolment</span> could not be found. Go <a href="/">back</a> and try again.`, []html.Crumb{{"Enrolments", "/"}, {"Enrolment Not Found", ""}})
						return

					} else {
						user = u
						breadcrumb = html.BreadcrumbTrail([]html.Crumb{{"Enrolments", "/"}, {enrol.CourseTitle, ""}})
						title = enrol.CourseTitle
						var continueModal, continueContent string
						var tutorials, started, completed int
						content, continueModal, tutorials, started, completed = html.Tutorials(enrol)
						if continueModal != "" {
							continueContent = utils.Concat(`<div class="mt-3">`, continueModal, `</div>`)
						}
						if tutorials > 1 {
							content = utils.Concat(`<div class="shadow-lg p-3 mb-3 bg-light rounded"><h5>Progress</h5>`, html.Progress(tutorials, started, completed), continueContent, `</div>`, content)
						}
					}
				}

			} else {
				var home string
				user = api.GetSignedInTokenFlag(token)
				if user == "" {
					home = html.SIGN_IN
				} else {
					home = "Enrolments"
				}
				breadcrumb = html.BreadcrumbTrail([]html.Crumb{{home, "/"}, {html.ERROR, ""}})
				title = html.ERROR
				content = utils.Concat("<p>", GetError(err), "</p>")
			}

		}

		var kind string
		if title == html.ERROR {
			kind = html.ERROR
		} else {
			kind = html.LEARN
		}

		html.Webpage(w, r, user, html.Page{kind, title}, breadcrumb, content)
	}
}
