package main

import (
	"./api"
	"./utils"
	"net/http"
)

func error404(w http.ResponseWriter, r *http.Request) {

	// get the auth token - first try cookies and if no cookie token found hit the /auth endpoint to get a fresh token
	token := api.GetToken(r)

	var message string

	result, _, newToken, user, err := api.CheckToken(token, token)
	if err == nil {

		// save the new token to use in the next api call
		api.SaveToken(w, newToken)

		if result.LoginID == "" || result.LoginID == user {
			message = utils.Concat("Error - page not found")
		} else {
			message = utils.Concat("Session error - page not found")
		}
	} else {
		message = GetError(err)
		user = api.GetSignedInTokenFlag(token)
	}

	var home string
	if user == "" {
		home = "Sign In"
	} else {
		home = "Enrolments"
	}
	breadcrumb := breadcrumbTrail([]crumb{{home, "/"}, {"Page Not Found", ""}})
	html(w, r, user, "Page Not Found", breadcrumb, utils.Concat(`<p>`, message, `</p>`))
}
