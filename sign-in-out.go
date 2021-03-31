package main

import (
	api "github.com/kevinbarbary/go-lms/api"
	html "github.com/kevinbarbary/go-lms/html"
	utils "github.com/kevinbarbary/go-lms/utils"
	"net/http"
	"log"
)

func signIn(w http.ResponseWriter, r *http.Request, path string) {

	var user, content string

	name := r.FormValue("username")
	pass := r.FormValue("password")
	if name == "" && pass == "" {
		content = html.SignIn("", "", path)
		if r.Method == "POST" {
			content = utils.Concat(html.StyleMessage("Enter credentials and try again", "danger"), content)
		}
	} else {

		// get the auth token - first try cookies and if no cookie token found hit the /auth endpoint to get a fresh token
		token, u := api.Auth(utils.GetSite(r), name, pass, r.UserAgent(), false)

		api.SaveToken(w, token, utils.GetDomain(r))
log.Print("saved token")
log.Print(token)
log.Print(utils.GetDomain(r))
log.Print("redirecting")
		if u != "" {
			html.SetMessage(w, r, html.SIGNED_IN)
			http.Redirect(w, r, utils.Concat("/", r.FormValue("path")), 302)
		}

		user = u
		content = utils.Concat(html.StyleMessage("Sign in failed", "danger"), html.SignIn(name, pass, r.FormValue("path")))
	}

	breadcrumb := html.BreadcrumbTrail([]html.Crumb{{html.SIGN_IN, ""}})
	html.Webpage(w, r, user, html.Page{html.SIGN_IN, ""}, breadcrumb, content)
}

func signOut(w http.ResponseWriter, r *http.Request) {

	token := api.GetTokenIfSignedIn(r)
	if token != "" {

		// remove the user from the token by calling the auth endpoint without the LoginID and Password
		// i.e. get a new token without a user in (also invalidates the session in the current token so any tokens containing the same session can't be used again)
		newToken, user := api.Unauth(utils.GetSite(r), token, r.UserAgent())

		api.SaveToken(w, newToken, utils.GetDomain(r))
		if user == "" {
			html.SetMessage(w, r, html.SIGNED_OUT)
		} else {
			html.SetMessage(w, r, html.SIGN_OUT_FAIL)
		}
	}
	http.Redirect(w, r, "/", 302)
}
