package main

import (
	"./api"
	"./utils"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const NOT_STARTED = `<span class="badge rounded-pill bg-secondary">Not Started</span>`
const STARTED = `<span class="badge rounded-pill bg-warning">Started</span>`
const COMPLETED = `<span class="badge rounded-pill bg-success">Completed</span>`

type page struct {
	kind, header string
}

func html(w http.ResponseWriter, r *http.Request, user string, page page, breadcrumb, content string) {

	var css string
	if page.kind == SIGN_IN {
		css = `<link rel="stylesheet" type="text/css" href="/assets/css/sign-in.css">`
	}

	title := "LMS - "
	if page.header == "" {
		title = utils.Concat(title, "home")
	} else if page.kind == SIGN_IN {
		title = utils.Concat(title, SIGN_IN)
	} else {
		title = utils.Concat(title, page.header)
	}

	var signInOut, menu, menuSpacing, learnItem, learn, courses string

	switch page.kind {
	case SIGN_IN:
		courses = "outline-"
	case LEARN:
		courses = "outline-"
	case COURSES:
		learn = "outline-"
	}

	if page.kind != PLAIN {
		class := "btn btn-outline-primary btn-sm"
		if user == "" {
			if page.kind != SIGN_IN {
				signInOut = utils.HyperClass("/", SIGN_IN, class)
				menuSpacing = " ms-3"
			}
		} else {
			if page.kind != SIGN_IN {
				menuSpacing = " ms-3"
			}
			if user == api.TokenUser() {
				signInOut = utils.HyperClass("/sign-out", "Sign out", class)
			} else {
				signInOut = utils.HyperClass("/sign-out", utils.Concat("Sign out: ", user), class)
			}
			learnItem = utils.Concat(`<a href="/" class="btn btn-`, learn, `primary btn-sm">Learn</a>`)
		}

		menu = utils.Concat(`<div class="btn-group`, menuSpacing, `" role="group" aria-label="Menu">`,
			learnItem, `<a href="/courses" class="btn btn-`, courses, `primary btn-sm">Browse</a></div>`)

		menu = utils.Concat(`<div class="menu mb-3">`, signInOut, menu, `</div>`)
	}

	var body string
	if message, kind := GetMessage(r); message != "" {
		body = utils.Concat(body, StyleMessage(message, kind))
		UnsetMessage(w)
	}
	body = utils.Concat(body, content)

	var logo, name string
	if logo, name = utils.Logo(); logo != "" {
		if name == "" {
			logo = utils.Concat(`<img src="`, logo, `" alt="logo">`)
		} else {
			logo = utils.Concat(`<img src="`, logo, `" alt="`, name, `">`)
		}
	} else if name != "" {
		logo = utils.Concat(`<mark>`, name, `</mark>`)
	}

	// @todo - add footer

	data := struct {
		Title      string
		Css        template.HTML
		Version    string
		Breadcrumb template.HTML
		Logo       template.HTML
		Menu       template.HTML
		Header     string
		Content    template.HTML
	}{
		title,
		template.HTML(css),
		utils.Assets("CSS"),
		template.HTML(utils.Concat(`<span class="breadcrumb-trail breadcrumb-prefix">You are here:</span>`, breadcrumb)),
		template.HTML(logo),
		template.HTML(menu),
		page.header,
		template.HTML(body),
	}

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type crumb struct {
	title, link string
}

func breadcrumbTrail(list []crumb) string {
	var trail = make([]string, len(list))
	for i, current := range list {
		if current.link == "" {
			trail[i] = utils.Concat(`<li class="breadcrumb-item active" aria-current="page">`, current.title, `</li>`)
		} else {
			trail[i] = utils.Concat(`<li class="breadcrumb-item">`, utils.Hyper(current.link, current.title), `</li>`)
		}
	}
	return utils.Concat(`<nav class="breadcrumb-trail" id="breadcrumb" aria-label="breadcrumb"><ol class="breadcrumb">`, strings.Join(trail, ""), `</ol></nav>`)
}

func htmlEnrolStart() string {
	return `<div class="enrolments">`
}

func htmlEnrolRow(enrol api.UserEnrol, now api.Timestamp) string {

	var logo string
	if enrol.PublisherLogo != "" {
		logo = utils.Concat(`<img src="`, enrol.PublisherLogo, `" alt="`, enrol.Publisher, `">`)
	}

	var completeStatus, completeClass string
	if enrol.Completed {
		completeStatus = "Completed"
		completeClass = " completed"
	} else {
		completeStatus = "Incomplete"
		completeClass = " incomplete"
	}

	hyper := false
	var enrolStr, expires, statusClass, expiryClass string
	enrolStr = strconv.Itoa(enrol.EnrollID)
	if enrol.EnrollStatus.Enabled() && now.BeforeEnd(enrol.EndDate) {
		// active
		expires = utils.Concat("Expires in ", utils.FormatUntil(now.Until(enrol.EndDate)))
		hyper = true
		if enrol.Completed {
			completeStatus = COMPLETED
		} else {
			if enrol.TotalDuration > 0 {
				completeStatus = STARTED
			} else {
				completeStatus = NOT_STARTED
			}
		}
		statusClass = " my-1"
	} else {
		// expired, pending, etc.
		expires = "Expired"
		expiryClass = " expired"
		statusClass = ""
	}

	row := utils.Concat(`<div class="border p-3 mb-3`, expiryClass, completeClass, `" id="enrol-id-`, enrolStr,
		`"><div class="logo">`, logo, `</div><div class="enrol"><div class="title">`, enrol.CourseTitle,
		`</div><div class="status`, statusClass, `">`, completeStatus, `</div><div class="enrol-start">Start Date: `,
		enrol.StartDate.ToDate(), `</div><div class="expires">`, expires, `</div></div></div>`)

	if hyper {
		row = utils.Hyper(utils.Concat("/", enrolStr), row)
	}

	return utils.Concat(`<div class="enrol-row">`, row, `</div>`)
}

func htmlEnrolEnd() string {
	return "</div>"
}

func progress(total, started, completed int) string {
	var start, complete = 100 / total * started, 100 / total * completed
	none := 100 - start - complete
	startStr, completeStr, noneStr := strconv.Itoa(start), strconv.Itoa(complete), strconv.Itoa(none)
	return utils.Concat(`<div class="progress"><div class="progress-bar bg-success" role="progressbar" style="width: `,
		completeStr, `%;" aria-valuenow="`, completeStr, `" aria-valuemin="0" aria-valuemax="100">`, completeStr,
		`%</div><div class="progress-bar bg-warning" role="progressbar" style="width: `, startStr, `%;" aria-valuenow="`,
		startStr, `" aria-valuemin="0" aria-valuemax="100">`, startStr, `%</div><div class="progress-bar bg-secondary" role="progressbar" style="width: `,
		noneStr, `%;" aria-valuenow="`, noneStr, `" aria-valuemin="0" aria-valuemax="100">`, noneStr, `%</div></div>`)
}

func tutorialsHTML(data api.UserEnrolment) (string, string, int, int, int) {
	var tutorials, started, completed int
	var html, status string
	var lastAccessed api.JsonDateTime

	var lastUrl string
	for _, lesson := range data.Lessons {
		if lesson.Title != data.CourseTitle {
			html = utils.Concat(html, "<h2>", lesson.Title, "</h2>")
		}
		for _, tutorial := range lesson.Tutorials {
			if tutorials == 0 {
				lastAccessed = tutorial.LastAccessed
				lastUrl = tutorial.LaunchURL
			} else {
				if tutorial.LastAccessed.After(lastAccessed) {
					lastAccessed = tutorial.LastAccessed
					lastUrl = tutorial.LaunchURL
				}
			}
			tutorials++
			if tutorial.Completed {
				completed++
				status = COMPLETED
			} else {
				if tutorial.TimesAccessed > 0 {
					started++
					status = STARTED
				} else {
					status = NOT_STARTED
				}
			}
			html = utils.Concat(html, `<div class="tutorial-row" id="tutorial-id-`, strconv.Itoa(tutorial.TutorialID), `">`,
				utils.Hyper(utils.Concat(tutorial.LaunchURL, "&returnHTTP=1&returnURL=", url.QueryEscape(utils.Concat("//", utils.Domain(), "/")),
					strconv.Itoa(data.EnrollID)), utils.Concat(`<div class="border p-2 mb-2"><div class="name">`, tutorial.TutorialTitle,
					`</div><div class="status">`, status, `</div></div>`)), `</div>`)
		}
	}

	// WIP: launch in modal...
	var modalContinue string
	if !lastAccessed.NotSet() {
		modalContinue = utils.Concat(`
<a href="#" class="btn btn-outline-primary btn-sm" id="continue" data-bs-toggle="modal" data-bs-target="#exampleModal" data-url="`,
			utils.Concat(lastUrl, "&returnHTTP=1&forceExit=0&returnURL=", url.QueryEscape(utils.Concat("//",
				utils.Domain(), "/parent/")), strconv.Itoa(data.EnrollID)), `">
Continue
</a>
<div class="modal fade" id="exampleModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
 <div class="modal-dialog modal-fullscreen">
   <div class="modal-content">
     <div class="modal-body">
       <iframe frameborder="0" style="overflow:hidden;height:100%;width:100%" height="100%" width="100%" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
     </div>
   </div>
 </div>
</div>
<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.1/jquery.min.js"></script>
<script type="text/javascript">
$("#continue").click(function () {
  var theModal = $(this).data("bs-target");
  $(theModal + ' iframe').attr('src', $(this).attr("data-url"));
});
</script>`)
	}

	return html, modalContinue, tutorials, started, completed
}
