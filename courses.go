package main

import (
	"./api"
	"./utils"
	"net/http"
)

func courses(w http.ResponseWriter, r *http.Request, offset int) {

	var user, breadcrumb, card, content string

	token := api.GetToken(r)

	courseData, newToken, u, _, err := api.Courses(token, offset)
	if err == nil {
		api.SaveToken(w, newToken)
		user = u
	}

	breadcrumb = breadcrumbTrail([]crumb{{"Courses", ""}})

	for _, course := range courseData.Courses {
		card = utils.Hyper("#", utils.Concat(`<div class="card mx-auto" style="width: 208px;">
  <div class="card-title"><div class="card-image mx-auto" style="width: 172px;">
    <img src="`, course.Image, `" class="card-img-top pt-3" alt="`, course.CourseTitle, `" style="width: 172px; height:82px">
  </div></div>
  <div class="card-body pt-0">
    <p class="card-text">`, course.CourseTitle, `</p>
  </div>
</div>`))
		content = utils.Concat(content, `<div class="course col">`, card, `</div>`)
	}

	content = utils.Concat(`<div id="cards" class="row row-cols-2 row-cols-md-3 row-cols-lg-4 row-cols-xl-5 row-cols-xxl-6 g-3 g-md-2 g-lg-2 g-xl-2 g-xxl-2">`, content, `</div>`)

	pagination := `<nav class="mt-3" aria-label="Page navigation">
  <ul class="pagination pagination-sm justify-content-center">
    <li class="page-item disabled"><a class="page-link" href="#" tabindex="-1" aria-disabled="true"><span>First</span></a></li>
    <li class="page-item disabled"><a class="page-link" href="#" tabindex="-1" aria-disabled="true"><span>Previous</span></a></li>
    <li class="page-item disabled"><a class="page-link" href="#" tabindex="-1" aria-disabled="true"><span>...</span></a></li>
    <li class="page-item"><a class="page-link" href="#">11</a></li>
    <li class="page-item"><a class="page-link" href="#">12</a></li>
    <li class="page-item disabled active"><a class="page-link" href="#">13</a></li>
    <li class="page-item"><a class="page-link" href="#">14</a></li>
    <li class="page-item"><a class="page-link" href="#">15</a></li>
    <li class="page-item disabled"><a class="page-link" href="#" tabindex="-1" aria-disabled="true"><span>...</span></a></li>
    <li class="page-item"><a class="page-link" href="#">Next</a></li>
    <li class="page-item"><a class="page-link" href="#">Last (50)</a></li>
  </ul>
</nav>`

	content = utils.Concat(content, pagination)

	html(w, r, user, page{COURSES, COURSES}, breadcrumb, content)
}
