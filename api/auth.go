package api

import (
	"../utils"
	"log"
)

func Auth(username, password string) (string, string) {

	site, key, err := utils.Creds()
	if err != nil {
		log.Print("Auth Error - invalid response from API call... ", err.Error())
		return "", ""
	}

	payload, err := Creds(site, key)
	if err != nil {
		log.Print("Auth Error - failed to build auth payload for api call... ", err.Error())
		return "", ""
	}

	if username != "" || password != "" {
		payload = MergeParams(payload, Params{"LoginID": username, "Password": password})
	}

	response, err := Call("POST", utils.Endpoint("/auth"), "", payload)
	if err != nil {
		log.Print("Auth Error - API call failed... ", err.Error())
		return "", ""
	}

	_, e, help, timestamp, token, user := extract(response)
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

	return token, user
}
