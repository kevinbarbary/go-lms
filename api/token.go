package api

import (
	"../utils"
	"encoding/json"
	"log"
	"net/http"
)

type TokenInfo struct {
	Expires          string
	Now              string
	SecondsRemaining int64
	URL              string
	SiteID           string
	LoginID          string
}

func CheckToken(token, check string) (TokenInfo, Timestamp, string, string, error) {

	response, err := Call("POST", utils.Endpoint("/auth/check"), token, Params{"Token": check})
	if err != nil {
		log.Print("CheckToken Error - invalid response from API call... ", err.Error())
		return TokenInfo{}, 0, "", "", err
	}

	data, e, help, t, newToken, user := extract(response)

	if e != "" {
		log.Print("CheckToken Error... ", e)
	}
	if help != "" {
		log.Print("CheckToken help... ", help)
	}

	timestamp := t / 10000
	if data == nil {
		log.Print("CheckToken... NO DATA")
		return TokenInfo{}, timestamp, newToken, user, err
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		log.Print("CheckToken - Marshal fail... ", err.Error())
		return TokenInfo{}, timestamp, newToken, user, err
	}

	var result TokenInfo
	err = json.Unmarshal(byteData, &result)
	if err != nil {
		log.Print("CheckToken - Unmarshal fail... ", err.Error())
		return TokenInfo{}, timestamp, newToken, user, err
	}

	return result, timestamp, newToken, user, err
}

func saveTokenCookie(w http.ResponseWriter, token string) {
	utils.SaveCookie(w, "token", token)
}

func SaveToken(w http.ResponseWriter, token string) {
	if token != "" {
		saveTokenCookie(w, token)
	}
}

func GetToken(r *http.Request) string {
	token, err := utils.GetCookieValue(r, "token")
	if err != nil || token == "" {
		token, _ = Auth("", "")
	}
	return token
}

func TokenUser() string {
	// the Course-Source API allows using ! in place of a LoginID to use the user in the token,
	// it also prefixes the auth token with a ! if the token contains a user
	return "!"
}

func GetSignedInTokenFlag(token string) string { // @todo - change the return type to rune ?
	if token != "" && token[:1] == TokenUser() {
		// if the token has an ! prefix then it contains a user
		return TokenUser()
	}
	return ""
}

func CheckTokenSignedIn(token string) bool {
	return GetSignedInTokenFlag(token) == TokenUser()
}

func CheckSignedIn(r *http.Request) bool {
	token, err := utils.GetCookieValue(r, "token")
	return err == nil && CheckTokenSignedIn(token)
}
