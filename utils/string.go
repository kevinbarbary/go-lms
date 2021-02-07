package utils

import "strings"

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
