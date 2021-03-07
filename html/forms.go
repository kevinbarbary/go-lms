package html

import "../utils"

func SignIn(username, password, path string) string {
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
      	<div class="d-grid gap-2"><button class="btn btn-lg btn-outline-primary" type="submit">Sign in</button></div>
    </form>`)
}
