package api

import (
	"../utils"
	"encoding/json"
	"log"
)

type EnrolStatus string // @todo - rune

func (s EnrolStatus) Enabled() bool {
	// Status A = Active, G = Group, P = Pending, D = Disabled, etc.
	return s == "A" || s == "G"
}

type UserEnrol struct {
	EnrollID      int      `json:"EnrollID"`
	CourseID      int64    `json:"CourseID"`
	Type          string   `json:"Type"`
	CourseTitle   string   `json:"CourseTitle"`
	PublisherID   string   `json:"PublisherID"`
	Publisher     string   `json:"Publisher"`
	PublisherLogo string   `json:"PublisherLogo"`
	StartDate     JsonDate `json:"StartDate"` // date
	EndDate       JsonDate `json:"EndDate"`   // date
	TotalDuration int64    `json:"TotalDuration"`
	//LastAccessed	JsonDate  	`json:"LastAccessed"` 	// date // @todo - use this instead
	LastAccessed   string      `json:"LastAccessed"` // date
	Completed      bool        `json:"Completed"`
	EnrollStatus   EnrolStatus `json:"EnrollStatus"`
	CertificateURL string      `json:"CertificateURL"`
}

func UserEnrolments(token, useragent, site, loginId string) ([]UserEnrol, string, string, Timestamp, error) {

	response, err := Call("GET", utils.Endpoint(utils.Concat("/enrolments/", loginId)), token, useragent, site, nil, true)

	if err != nil {
		log.Print("UserEnrolments Error - invalid response from API call... ", err.Error())
		return nil, "", "", 0, err
	}

	data, e, help, now, token, user := extract(response)

	if e != "" {
		log.Print("UserEnrolments Error... ", e)
	}
	if help != "" {
		log.Print("UserEnrolments help... ", help)
	}

	if data == nil {
		log.Print("UserEnrolments... NO DATA")
		return nil, token, user, now, err
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		log.Print("UserEnrolments - Marshal fail... ", err.Error())
		return nil, token, user, now, err
	}

	var val []UserEnrol
	err = json.Unmarshal(byteData, &val)
	if err != nil {
		log.Print("UserEnrolments - Unmarshal fail... ", err.Error())
		return nil, token, user, now, err
	}

	return val, token, user, now, nil
}
