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

func html(w http.ResponseWriter, r *http.Request, user, header, breadcrumb, content string) {

	var css string
	if header == SIGN_IN {
		css = `<link rel="stylesheet" type="text/css" href="/assets/css/sign-in.css">`
	}

	title := "LMS - "
	if header == "" {
		title = utils.Concat(title, "home")
	} else {
		title = utils.Concat(title, header)
	}

	// modal test...
	var menu string
	if user != "*" { // don't show the menu
		class := "btn btn-primary btn-sm"
		if user == "" {
			if header != SIGN_IN {
				menu = utils.HyperClass("/", SIGN_IN, class, "button")
			}
		} else if user == api.TokenUser() {
			menu = utils.HyperClass("/sign-out", "Sign out", class, "button")
		} else {
			menu = utils.HyperClass("/sign-out", utils.Concat("Sign out: ", user), class, "button")
		}
		if menu != "" {
			menu = utils.Concat(`<p class="menu">`, menu, `</p>`)
		}
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
		header,
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
	var enrolStr, expires, expiryClass string
	enrolStr = strconv.Itoa(enrol.EnrollID)
	if enrol.EnrollStatus.Enabled() && now.BeforeEnd(enrol.EndDate) {
		// active
		expires = utils.Concat("Expires in ", utils.FormatUntil(now.Until(enrol.EndDate)))
		hyper = true
		if enrol.Completed {
			completeStatus = `<span class="badge rounded-pill bg-success">Completed</span>`
		}
	} else {
		// expired, pending, etc.
		expires = "Expired"
		expiryClass = " expired"
	}

	row := utils.Concat(`<div class="border p-3 mb-3`, expiryClass, completeClass, `" id="enrol-id-`, enrolStr, `"><div class="logo">`, logo, `</div><div class="enrol"><div class="title">`, enrol.CourseTitle, `</div><div class="status">`, completeStatus, `</div><div class="enrol-start">Start Date: `, enrol.StartDate.ToDate(), `</div><div class="expires">`, expires, `</div></div></div>`)

	if hyper {
		row = utils.Hyper(utils.Concat("/", enrolStr), row)
	}

	return utils.Concat(`<div class="enrol-row">`, row, `</div>`)
}

func htmlEnrolEnd() string {
	return "</div>"
}

func progress(total, started, completed int) string {
	var start, complete = strconv.Itoa(100 / total * started), strconv.Itoa(100 / total * completed)
	return utils.Concat(`<div class="progress"><div class="progress-bar bg-success" role="progressbar" style="width: `, complete, `%;" aria-valuenow="`, complete, `" aria-valuemin="0" aria-valuemax="100">`, complete, `%</div><div class="progress-bar bg-warning" role="progressbar" style="width: `, start, `%;" aria-valuenow="`, start, `" aria-valuemin="0" aria-valuemax="100">`, start, `%</div></div>`)
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
				status = `<span class="badge rounded-pill bg-success">Completed</span>`
			} else {
				if tutorial.TimesAccessed > 0 {
					started++
					status = `<span class="badge rounded-pill bg-warning">Started</span>`
				} else {
					status = ""
				}
			}
			html = utils.Concat(html, `<div class="tutorial-row" id="tutorial-id-`, strconv.Itoa(tutorial.TutorialID), `">`, utils.Hyper(utils.Concat(tutorial.LaunchURL, "&returnHTTP=1&returnURL=", url.QueryEscape(utils.Concat("//", utils.Domain(), "/")), strconv.Itoa(data.EnrollID)), utils.Concat(`<div class="border p-2 mb-2"><div class="name">`, tutorial.TutorialTitle, `</div><div class="status">`, status, `</div></div>`)), `</div>`)
		}
	}

	// modal WIP...
	var modalContinue string
	if !lastAccessed.NotSet() {
		modalContinue = utils.Concat(`
<a href="#" class="btn btn-primary btn-sm" id="continue" data-bs-toggle="modal" data-bs-target="#exampleModal" data-url="`, utils.Concat(lastUrl, "&returnHTTP=1&forceExit=0&returnURL=", url.QueryEscape(utils.Concat("//", utils.Domain(), "/parent/")), strconv.Itoa(data.EnrollID)), `">
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
