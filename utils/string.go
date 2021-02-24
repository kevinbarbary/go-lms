package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

func Concat(list ...string) string {
	var sb strings.Builder
	for _, it := range list {
		sb.WriteString(it)
	}
	return sb.String()
}

func Hyper(path, page string) string {
	return Concat(`<a href="`, path, `">`, page, `</a>`)
}

func HyperClass(path, page, class string) string {
	return Concat(`<a class="`, class, `" href="`, path, `">`, page, `</a>`)
}

func FormatUntil(t time.Duration) string {
	t = t.Round(time.Minute)
	h := t / time.Hour
	if h < 24 {
		// less than one day
		t -= h * time.Hour
		m := t / time.Minute
		if h < 1 {
			// less than 1 hour
			if m == 1 {
				return "1 minute"
			}
			return fmt.Sprintf("%2d minutes", m)
		} else {
			// between 1 hour and one day
			if h == 1 {
				if m == 1 {
					return fmt.Sprintf("1 hour, 1 minute", h, m)
				}
				return fmt.Sprintf("1 hour, %2d minutes", h, m)
			}
			if m == 1 {
				return fmt.Sprintf("%2d hours, 1 minute", h, m)
			}
			return fmt.Sprintf("%2d hours, %2d minutes", h, m)
		}
	}
	d := h / 24
	if d < 7 {
		// between 1 day and 1 week
		if d == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%2d days", d)
	}
	w := d / 7
	if w < 5 {
		// between 1 week and 1 month
		if w == 1 {
			return "1 week"
		}
		return fmt.Sprintf("%2d weeks", w)
	}
	m := d / 30
	if m < 12 {
		// between 1 month and 1 year
		if m == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%2d months", m)
	}
	// over 1 year
	y := d / 365
	if y < 2 {
		return "1 year"
	}
	return fmt.Sprintf("%2d years", y)
}

func AlphaNumeric(in string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Print("Error converting to alphanumeric: ", in)
		return in
	}
	return reg.ReplaceAllString(in, "")
}
