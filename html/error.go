package html

import (
	utils "github.com/kevinbarbary/go-lms/utils"
)

func Error(message, info string) string {
	if info != "" {
		info = utils.Concat(`<p>`, info, `</p>`)
	}
	return utils.Concat(`
	<div class="row text-center">
        <div class="col-lg-6 offset-lg-3 col-sm-6 offset-sm-3 col-12 p-3 error-main">
			<div class="row">
				<div class="col-lg-8 col-12 col-sm-10 offset-lg-2 offset-sm-1">
					<h1 class="m-0">404</h1>
					<h6>`, message, `</h6>
					`, info, `
				</div>
			</div>
        </div>
  	</div>`)
}
