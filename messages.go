package main

import (
	"./utils"
	"net/http"
	"strconv"
)

const SIGNED_OUT = 1
const SIGN_OUT_FAIL = 2
const SIGNED_IN = 3

var messages = map[int]string{
	1: "You are now signed-out",
	2: "Sign-out failed!",
	3: "You are now signed-in",
}

var kinds = map[int]string{
	1: "success",
	2: "danger",
	3: "success",
}

func GetMessage(r *http.Request) (string, string) {
	cookie, err := utils.GetCookieValue(r, "msg")
	if err != nil {
		return "", ""
	}
	i, err := strconv.Atoi(cookie)
	if err != nil {
		return "", ""
	}
	if max := len(messages); i > max {
		return "", ""
	}
	return messages[i], kinds[i]
}

func StyleMessage(message, kind string) string {
	return utils.Concat(`<div class="alert alert-`, kind, ` alert-dismissible fade show" role="alert">`, message, `<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button></div>`)
}

func SetMessage(w http.ResponseWriter, i int) {
	utils.SaveCookie(w, "msg", strconv.Itoa(i))
}

func UnsetMessage(w http.ResponseWriter) {
	utils.DeleteCookie(w, "msg")
}
