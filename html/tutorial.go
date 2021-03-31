package html

import (
	api "github.com/kevinbarbary/go-lms/api"
	utils "github.com/kevinbarbary/go-lms/utils"
	"net/url"
	"strconv"
)

func Tutorials(data api.UserEnrolment) (string, string, int, int, int) {
	var tutorials, started, completed int
	var content, status string
	var lastAccessed api.JsonDateTime

	var lastUrl string
	for _, lesson := range data.Lessons {
		if lesson.Title != data.CourseTitle {
			content = utils.Concat(content, "<h2>", lesson.Title, "</h2>")
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
			content = utils.Concat(content, `<div class="tutorial-row" id="tutorial-id-`, strconv.Itoa(tutorial.TutorialID), `">`,
				utils.Hyper(utils.Concat(tutorial.LaunchURL, "&returnHTTP=1&returnURL=", url.QueryEscape(utils.Concat("//", utils.DomainAndPort(), "/")),
					strconv.Itoa(data.EnrollID)), utils.Concat(`<div class="border p-2 mb-2"><div class="name">`, tutorial.TutorialTitle,
					`</div><div class="status">`, status, `</div></div>`)), `</div>`)
		}
	}

	// WIP: launch in modal...
	var modalContinue string
	if !lastAccessed.NotSet() {
		modalContinue = utils.Concat(`
<a href="#" class="btn btn-outline-primary btn-sm" id="continue" data-bs-toggle="modal" data-bs-target="#exampleModal" data-url="`,
			utils.Concat(lastUrl, "&returnHTTP=1&noForceRedirect=0&returnURL=", url.QueryEscape(utils.Concat("//",
				utils.DomainAndPort(), "/parent/")), strconv.Itoa(data.EnrollID)), `">
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

	return content, modalContinue, tutorials, started, completed
}
