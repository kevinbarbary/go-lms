package html

import (
	api "github.com/kevinbarbary/go-lms/api"
	utils "github.com/kevinbarbary/go-lms/utils"
	"strconv"
)

const NOT_STARTED = `<span class="badge rounded-pill bg-secondary">Not Started</span>`
const STARTED = `<span class="badge rounded-pill bg-warning">Started</span>`
const COMPLETED = `<span class="badge rounded-pill bg-success">Completed</span>`

func Enrol(data []api.UserEnrol, now api.Timestamp) string {
	var content string
	for _, enrol := range data {
		content = utils.Concat(content, row(enrol, now))
	}
	if content == "" {
		content = "<p>You do not have any enrolments</p>"
	}
	return utils.Concat(start(), content, end())
}

func start() string {
	return `<div class="enrolments">`
}

func row(enrol api.UserEnrol, now api.Timestamp) string {

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

func end() string {
	return "</div>"
}
