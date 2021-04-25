package html

import (
	api "github.com/kevinbarbary/go-lms/api"
	utils "github.com/kevinbarbary/go-lms/utils"
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

	courseData, newToken, u, _, err := api.Courses(token, r.UserAgent(), utils.GetSite(r), index, tags)
	if err != nil {
		log.Fatal(err)
	}

	// @todo - better error handling

	var pageIndex string

	api.SaveToken(w, newToken, utils.GetDomain(r))
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
	filters := 0
	tagsFilter := `<nav class="collapse links" id="filter-nav" aria-label="Course Filter"><ul class="list-unstyled mb-0 py-3 pt-md-1">`
	for _, tagType := range courseData.Tags {
		expand := false
		items := ""
		alphanum := utils.AlphaNumeric(tagType.TagType)
		for _, tag := range tagType.Tags {
			_, selected := utils.FindInt(tags, tag.TagID)
			if selected {
				filters++
				expand = true
			}
			items = utils.Concat(items, "<li>", cbx(strconv.Itoa(tag.TagID), tag.Tag, "d-inline-flex rounded", selected), "</li>")
		}
		if expand {
			tagsFilter = utils.Concat(tagsFilter, `<li class="mb-1"><button type="button" class="btn d-inline-flex align-items-center rounded" data-bs-toggle="collapse" data-bs-target="#tag-`,
				alphanum, `-collapse" aria-expanded="true" aria-current="true">`, tagType.TagType, `</button><div class="collapse show" id="tag-`,
				alphanum, `-collapse"><div class="collapse show" id="forms-collapse"><ul class="list-unstyled fw-normal pb-1 small">`)
		} else {
			tagsFilter = utils.Concat(tagsFilter, `<li class="mb-1"><button type="button" class="btn d-inline-flex align-items-center rounded collapsed" data-bs-toggle="collapse" data-bs-target="#tag-`,
				alphanum, `-collapse" aria-expanded="false" aria-current="true">`, tagType.TagType, `</button><div class="collapse" id="tag-`,
				alphanum, `-collapse"><div class="collapse show" id="forms-collapse"><ul class="list-unstyled fw-normal pb-1 small">`)
		}
		tagsFilter = utils.Concat(tagsFilter, items, `</ul></div></li>`)
	}
	tagsFilter = utils.Concat(tagsFilter, `</ul></nav>`)

	badge := ""
	if filters > 0 {
		badge = utils.Concat(`s <span class="badge bg-primary">`, strconv.Itoa(filters), `</span>`)
	}
	tagsFilter = utils.Concat(`<button class="btn btn-outline-primary btn-sm mb-2" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasFilter" aria-controls="offcanvasFilter">Filter`, badge, `</button>
<div class="offcanvas offcanvas-start" tabindex="-1" id="offcanvasFilter" aria-labelledby="offcanvasFilterLabel">
	<form method="post" action="/courses/page-`, pageIndex, `">
		<div class="offcanvas-header">
			<h5 class="offcanvas-title" id="offcanvasFilterLabel"><button class="btn btn-outline-primary btn-sm" type="submit">Apply Filter</button></h5>
			<button type="button" class="btn-close text-reset" data-bs-dismiss="offcanvas" aria-label="Close"></button>
		</div>
		<div class="offcanvas-body pt-0">
			`, tagsFilter, `
		</div>
	</form>
</div>`)

	for _, course := range courseData.Courses {
		card = utils.Hyper("#", utils.Concat(`<div class="card mx-auto" style="width: 208px;">
  <div class="card-title"><div class="card-image mx-auto pt-3">
    <img src="`, course.Image, `" class="card-img-top" alt="`, course.CourseTitle, `">
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
