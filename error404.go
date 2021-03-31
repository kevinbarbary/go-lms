package main

import (
	api "github.com/kevinbarbary/go-lms/api"
	html "github.com/kevinbarbary/go-lms/html"
	utils "github.com/kevinbarbary/go-lms/utils"
	"net/http"
)

func error404(w http.ResponseWriter, r *http.Request, location, info string, crumbs []html.Crumb) {

	// get the auth token - first try cookies and if no cookie token found hit the /auth endpoint to get a fresh token
	token := api.GetToken(r)

	var message string

	result, _, newToken, user, err := api.CheckToken(token, r.UserAgent(), utils.GetSite(r), token)
	if err == nil {

		// save the new token to use in the next api call
		api.SaveToken(w, newToken, utils.GetDomain(r))

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

	breadcrumb := html.BreadcrumbTrail(crumbs)

	content := html.Error(message, info)

	html.Webpage(w, r, user, html.Page{html.ERROR, ""}, breadcrumb, content)
}
