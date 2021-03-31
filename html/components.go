package html

import (
	utils "github.com/kevinbarbary/go-lms/utils"
	"strconv"
	"strings"
)

type Crumb struct {
	Title, Link string
}

func BreadcrumbTrail(list []Crumb) string {
	var trail = make([]string, len(list))
	for i, current := range list {
		if current.Link == "" {
			trail[i] = utils.Concat(`<li class="breadcrumb-item active" aria-current="page">`, current.Title, `</li>`)
		} else {
			trail[i] = utils.Concat(`<li class="breadcrumb-item">`, utils.Hyper(current.Link, current.Title), `</li>`)
		}
	}
	return utils.Concat(`<nav class="breadcrumb-trail" id="breadcrumb" aria-label="breadcrumb"><ol class="breadcrumb">`, strings.Join(trail, ""), `</ol></nav>`)
}

func cbx(id, text, class string, selected bool) string {
	var checked string
	if selected {
		checked = " checked"
	}
	return utils.Concat(`<div class="form-check `, class, `">
  <input class="form-check-input" type="checkbox"`, checked, ` name="tag`, id, `" id="tag`, id, `">
  <label class="form-check-label" for="tag`, id, `">`, text, `</label>
</div>`)
}

func Progress(total, started, completed int) string {
	var start, complete = 100 / total * started, 100 / total * completed
	none := 100 - start - complete
	startStr, completeStr, noneStr := strconv.Itoa(start), strconv.Itoa(complete), strconv.Itoa(none)
	return utils.Concat(`<div class="progress"><div class="progress-bar bg-success" role="progressbar" style="width: `,
		completeStr, `%;" aria-valuenow="`, completeStr, `" aria-valuemin="0" aria-valuemax="100">`, completeStr,
		`%</div><div class="progress-bar bg-warning" role="progressbar" style="width: `, startStr, `%;" aria-valuenow="`,
		startStr, `" aria-valuemin="0" aria-valuemax="100">`, startStr, `%</div><div class="progress-bar bg-secondary" role="progressbar" style="width: `,
		noneStr, `%;" aria-valuenow="`, noneStr, `" aria-valuemin="0" aria-valuemax="100">`, noneStr, `%</div></div>`)
}

func paginate(index, limit, total int, params string) string {

	// small: 7
	// _1_ ... _3_ ... _1_
	// x ... y-1 y y+1 ... z

	// medium: 11
	// _2_ ... _5_ ... _2_
	// x x+1 ... y-2 y-1 y y+1 y+2 ... z-1 z

	// large: 15
	// _3_ ... _7_ ... _3_
	// x x+1 x+2 ... y-3 y-2 y-1 y y+1 y+2 y+3 ... z-2 z-1 z

	var pages int
	if total < limit {
		pages = 1
	} else {
		pages = total / limit
	}

	s := makePageNav(pages, index, 1, params)  // max 7 links
	m := makePageNav(pages, index, 2, params)  // max 11 links
	l := makePageNav(pages, index, 3, params)  // max 15 links
	xl := makePageNav(pages, index, 4, params) // max 19 links

	return utils.Concat(s, m, l, xl)

}

func makePageNav(pages, index, size int, params string) string {

	if pages == 1 {
		return ""
	}

	var centre int
	var middle, end, pad string

	selected := make(map[int]string, pages+1)
	selected[index] = " disabled active"

	if index-size > size {
		if index+size+size-1 < pages {
			centre = index
		} else {
			centre = pages - (size * 2)
		}
	} else {
		centre = (size * 2) + 1
		if centre >= pages {
			centre = pages - 1
		}
	}

	// size -> max links...
	// 1 -> 7, 2 -> 11, 3 -> 15, etc.
	max := (size * 4) + 3
	if pages > max {
		pad = `
		<li class="page-item disabled"><a class="page-link" href="#" tabindex="-1" aria-disabled="true"><span>...</span></a></li>`
	}

	start := utils.Concat(`
<nav class="mt-3" aria-label="Page navigation">
	<ul class="pagination pagination-sm justify-content-center">
	<li class="page-item`, selected[1], `"><a class="page-link" href="/courses/page-1`, params, `" tabindex="1" aria-disabled="true">1</a></li>`)

	if pages > 1 {
		end = utils.Concat(`
		<li class="page-item`, selected[pages], `"><a class="page-link" href="/courses/page-`, strconv.Itoa(pages), params, `">`, strconv.Itoa(pages), `</a></li>
	</ul>
</nav>`)
		if pages > 2 {
			middle = utils.Concat(`
		<li class="page-item`, selected[centre], `"><a class="page-link" href="/courses/page-`, strconv.Itoa(centre), params, `">`, strconv.Itoa(centre), `</a></li>`)
		}
	}

	for i := 1; i <= size; i++ {

		if centre-i > 1 {
			middle = utils.Concat(`<li class="page-item`, selected[centre-i], `"><a class="page-link" href="/courses/page-`, strconv.Itoa(centre-i), params, `" tabindex="1" aria-disabled="true">`, strconv.Itoa(centre-i), `</a></li>`, middle)
		}
		if centre+i < pages {
			middle = utils.Concat(middle, `<li class="page-item`, selected[centre+i], `"><a class="page-link" href="/courses/page-`, strconv.Itoa(centre+i), params, `" tabindex="1" aria-disabled="true">`, strconv.Itoa(centre+i), `</a></li>`)
		}

		// might need to append to start
		if (i + 1) < (centre - i) {
			if (i < size) || (centre+i >= pages-i) || ((i == size) && (index >= centre) && (centre-i-1 <= i+1)) {
				if i+1 < centre-size {
					start = utils.Concat(start, `<li class="page-item`, selected[i+1], `"><a class="page-link" href="/courses/page-`, strconv.Itoa(i+1), params, `" tabindex="1" aria-disabled="true">`, strconv.Itoa(i+1), `</a></li>`)
				}
			} else {
				start = utils.Concat(start, pad)
			}
		} else {
			if (i == size) && (centre+i+1 < pages-i) {
				middle = utils.Concat(middle, pad)
			}
		}

		// might need prepend to end
		if (pages - i) > (centre + i) {
			if (i < size) || (centre-i <= i+1) || ((i == size) && (index <= centre) && (centre+i+1 >= pages-i)) {
				if pages-i > centre+size {
					end = utils.Concat(`<li class="page-item`, selected[pages-i], `"><a class="page-link" href="/courses/page-`, strconv.Itoa(pages-i), params, `" tabindex="1" aria-disabled="true">`, strconv.Itoa(pages-i), `</a></li>`, end)
				}
			} else {
				end = utils.Concat(pad, end)
			}
		} else {
			if (i == size) && (centre-i-1 > i) {
				middle = utils.Concat(pad, middle)
			}
		}

	}

	return utils.Concat(start, middle, end)
}
