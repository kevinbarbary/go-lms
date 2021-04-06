package html

import utils "github.com/kevinbarbary/go-lms/utils"

func SignIn(username, password, path string) string {
	var user, pass string
	if username != "" {
		user = utils.Concat(` value="`, username, `"`)
	}
	if password != "" {
		pass = utils.Concat(` value="`, password, `"`)
	}
	return utils.Concat(`
	<form class="form form-sign-in" method="post">
      	<label for="username" class="sr-only">Username</label>
      	<input type="text" id="username" name="username" autocomplete="username" class="form-control" placeholder="Username" autofocus`, user, `>
      	<label for="password" class="sr-only">Password</label>
      	<input type="password" id="password" name="password" autocomplete="current-password" class="form-control" placeholder="Password"`, pass, `>
		<input type="hidden" id="path" name="path" value="`, path, `">
      	<div class="d-grid gap-2">
			<button class="btn btn-lg btn-outline-primary" id="sign-in-btn" onclick="signIn()" type="submit">
				<span id="sign-in-none" role="status" aria-hidden="true"></span>
				<span id="sign-in-text">`, SIGN_IN, `</span>
			</button>
		</div>
    </form>`)
}
