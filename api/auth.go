package api

import (
	"../utils"
	"log"
)

func Auth(site, username, password, useragent string, retry bool) (string, string) {
	return authenticate(useragent, site, username, password, "", retry)
}

func Unauth(site, token, useragent string) (string, string) {
	// i.e. sign out - remove user from token
	return authenticate(useragent, site, "", "", token, false)
}

func authenticate(useragent, site, username, password, token string, retry bool) (string, string) {

	site, key := utils.Creds(site)
	payload, err := Creds(site, key)
	if err != nil {
		log.Print("Auth Error - failed to build auth payload for api call... ", err.Error())
		return "", ""
	}

	if username != "" || password != "" {
		payload = MergeParams(payload, Params{"LoginID": username, "Password": password})
	}

	response, err := Call("POST", utils.Endpoint("/auth"), token, useragent, site, payload, retry)
	if err != nil {
		log.Print("Auth Error - API call failed... ", err.Error())
		return "", ""
	}

	_, e, help, timestamp, newToken, user := extract(response)
	if timestamp == 0 {
		log.Print("Auth Error - invalid timestamp... ", timestamp)
		return "", ""
	}
	if e != "" {
		log.Print("Auth Fail... ", e)
		if help != "" {
			log.Print("Auth help... ", help)
		}
		return "", ""
	}

	return newToken, user
}
