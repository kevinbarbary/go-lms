package main

import (
	"./api"
	"./utils"
	"net/http"
)

func error404(w http.ResponseWriter, r *http.Request, location, info string, crumbs []crumb) {

	// get the auth token - first try cookies and if no cookie token found hit the /auth endpoint to get a fresh token
	token := api.GetToken(r)

	var message string

	result, _, newToken, user, err := api.CheckToken(token, token)
	if err == nil {

		// save the new token to use in the next api call
		api.SaveToken(w, newToken)

		if result.LoginID == "" || result.LoginID == user {
			var pad string
			if location != "" {
				pad = " "
			}
			message = utils.Concat(`Error - `, location, pad, `page not found`)
		} else {
			message = utils.Concat("Session error - page not found")
		}
	} else {
		message = GetError(err)
		user = api.GetSignedInTokenFlag(token)
	}
	if message == "" {
		message = "Error - page not found"
	}

	breadcrumb := breadcrumbTrail(crumbs)

	content := errorHTML(message, info)

	html(w, r, user, page{ERROR, ""}, breadcrumb, content)
}
