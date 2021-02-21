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

	html(w, r, user, "Courses", breadcrumb, utils.Concat(`<div id="cards" class="row row-cols-2 row-cols-md-3 row-cols-lg-4 row-cols-xl-5 row-cols-xxl-6 g-3 g-md-2 g-lg-2 g-xl-2 g-xxl-2">`, content, `</div>`))
}
