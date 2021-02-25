package main

import (
	"./api"
	"./utils"
	"net/http"
	"strconv"
)

func courses(w http.ResponseWriter, r *http.Request, offset int) {

	var user, breadcrumb, card, content string

	token := api.GetToken(r)

	courseData, newToken, u, _, err := api.Courses(token, offset)
	if err == nil {
		api.SaveToken(w, newToken)
		user = u
	}

	// @todo - don't do this if err... (but do what instead? make an error page that's not 404?)

	breadcrumb = breadcrumbTrail([]crumb{{"Courses", ""}})

	tagsFilter := `<aside class="bd-sidebar"><nav class="collapse bd-links" id="bd-docs-nav" aria-label="Docs navigation"><ul class="list-unstyled mb-0 py-3 pt-md-1">`
	for _, tagType := range courseData.Tags {
		alphanum := utils.AlphaNumeric(tagType.TagType)
		tagsFilter = utils.Concat(tagsFilter, `<li class="mb-1"><button class="btn d-inline-flex align-items-center rounded" data-bs-toggle="collapse" data-bs-target="#tag-`,
			alphanum, `-collapse" aria-expanded="true" aria-current="true">`, tagType.TagType, `</button><div class="collapse show" id="tag-`,
			alphanum, `-collapse" style=""><div class="collapse show" id="forms-collapse"><ul class="list-unstyled fw-normal pb-1 small">`)

		for _, tag := range tagType.Tags {
			tagsFilter = utils.Concat(tagsFilter, "<li>", cbx(strconv.Itoa(tag.TagID), tag.Tag, "d-inline-flex rounded"), "</li>")
		}
		tagsFilter = utils.Concat(tagsFilter, `</ul></div></li>`)
	}
	tagsFilter = utils.Concat(tagsFilter, "</ul></nav></aside>")

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

	content = utils.Concat(tagsFilter, `<div id="cards" class="row row-cols-2 row-cols-md-3 row-cols-lg-4 row-cols-xl-5 row-cols-xxl-6 g-3 g-md-2 g-lg-2 g-xl-2 g-xxl-2">`, content, `</div>`)

	pagination := paginate(offset, courseData.Next.Limit, courseData.Total)

	content = utils.Concat(content, pagination)

	html(w, r, user, page{COURSES, COURSES}, breadcrumb, content)
}
