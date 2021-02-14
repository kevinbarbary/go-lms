package main

import (
	"./api"
	"./utils"
	"net/http"
)

func signIn(w http.ResponseWriter, r *http.Request, path string) {

	var user, content string

	name := r.FormValue("username")
	pass := r.FormValue("password")
	if name == "" && pass == "" {
		content = formSignIn("", "", path)
		if r.Method == "POST" {
			content = utils.Concat(StyleMessage("Enter credentials and try again", "danger"), content)
		}
	} else {

		// get the auth token - first try cookies and if no cookie token found hit the /auth endpoint to get a fresh token
		token, u := api.Auth(name, pass)

		api.SaveToken(w, token)
		if u != "" {
			SetMessage(w, SIGNED_IN)
			http.Redirect(w, r, utils.Concat("/", r.FormValue("path")), 302)
		}

		user = u
		content = utils.Concat(StyleMessage("Sign in failed", "danger"), formSignIn(name, pass, r.FormValue("path")))
	}

	breadcrumb := breadcrumbTrail([]crumb{{"Sign In", ""}})
	html(w, r, user, "Sign In", breadcrumb, content)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	if api.CheckSignedIn(r) {

		// remove the user from the token by calling the auth endpoint without the LoginID and Password
		token, user := api.Auth("", "")

		api.SaveToken(w, token)
		if user == "" {
			SetMessage(w, SIGNED_OUT)
		} else {
			SetMessage(w, SIGN_OUT_FAIL)
		}
	}
	http.Redirect(w, r, "/", 302)
}

func formSignIn(username, password, path string) string {
	var user, pass string
	if username != "" {
		user = utils.Concat(` value="`, username, `"`)
	}
	if password != "" {
		pass = utils.Concat(` value="`, password, `"`)
	}
	return utils.Concat(`    <form class="form-sign-in" method="post">
      	<label for="username" class="sr-only">Username</label>
      	<input type="text" id="username" name="username" class="form-control" placeholder="Username" autofocus`, user, `>
      	<label for="password" class="sr-only">Password</label>
      	<input type="password" id="password" name="password" class="form-control" placeholder="Password"`, pass, `>
		<input type="hidden" id="path" name="path" value="`, path, `">
      	<div class="d-grid gap-2"><button class="btn btn-lg btn-primary" type="submit">Sign in</button></div>
    </form>`)
}
