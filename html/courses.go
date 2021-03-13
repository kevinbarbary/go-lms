package html

import (
	"../api"
	"../utils"
	"log"
	"net/http"
	"strconv"
)

func Courses(w http.ResponseWriter, r *http.Request, index int, tags []int) {

	var user, breadcrumb, card, content string

	token := api.GetToken(r)

	if index < 1 {
		index = 1
	}

	// get selected tags
	var tagsParam string
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		var formTags []int
		for i := range r.Form {
			if len(i) > 3 && i[:3] == "tag" {
				t, err := strconv.Atoi(i[3:])
				if err == nil {
					formTags = append(formTags, t)
					tagsParam = utils.Concat(tagsParam, "-", i[3:])
				} else {
					log.Print("Error: ignoring invalid CourseTag '", i[3:], "' in Courses filter")
				}
			}
		}

		if len(formTags) > 0 {
			tags = formTags
		}

	} else if len(tags) > 0 {
		for i := range tags {
			tagsParam = utils.Concat(tagsParam, "-", strconv.Itoa(i))
		}
	}
	if tagsParam != "" {
		tagsParam = utils.Concat("/tag", tagsParam)
	}

	courseData, newToken, u, _, err := api.Courses(token, utils.GetSite(r), index, tags)
	if err != nil {
		log.Fatal(err)
	}

	// @todo - better error handling

	var pageIndex string

	api.SaveToken(w, newToken)
	user = u

	// build breadcrumb
	if index == 1 {
		pageIndex = "1"
		breadcrumb = BreadcrumbTrail([]Crumb{{"Courses", ""}})
	} else {
		pageIndex = strconv.Itoa(index)
		breadcrumb = BreadcrumbTrail([]Crumb{{"Courses", "/courses"}, {utils.Concat("Page ", pageIndex), ""}})
	}

	// build filter
	tagsFilter := `<aside class="bd-sidebar"><nav class="collapse bd-links" id="bd-docs-nav" aria-label="Docs navigation"><ul class="list-unstyled mb-0 py-3 pt-md-1">`
	for _, tagType := range courseData.Tags {

		alphanum := utils.AlphaNumeric(tagType.TagType)
		tagsFilter = utils.Concat(tagsFilter, `<li class="mb-1"><button class="btn d-inline-flex align-items-center rounded" data-bs-toggle="collapse" data-bs-target="#tag-`,
			alphanum, `-collapse" aria-expanded="true" aria-current="true">`, tagType.TagType, `</button><div class="collapse show" id="tag-`,
			alphanum, `-collapse"><div class="collapse show" id="forms-collapse"><ul class="list-unstyled fw-normal pb-1 small">`)

		for _, tag := range tagType.Tags {
			_, selected := utils.FindInt(tags, tag.TagID)
			tagsFilter = utils.Concat(tagsFilter, "<li>", cbx(strconv.Itoa(tag.TagID), tag.Tag, "d-inline-flex rounded", selected), "</li>")
		}

		tagsFilter = utils.Concat(tagsFilter, `</ul></div></li>`)
	}
	tagsFilter = utils.Concat(tagsFilter, "</ul></nav></aside>")

	tagsFilter = utils.Concat(`<form method="post" action="/courses/page-`, pageIndex, `"><button class="btn btn-outline-primary btn-sm" type="submit">Filter</button>`, tagsFilter, `</form>`)
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

	pagination := paginate(index, courseData.Next.Limit, courseData.Total, tagsParam)

	content = utils.Concat(content, pagination)

	Webpage(w, r, user, Page{COURSES, COURSES}, breadcrumb, content)
}