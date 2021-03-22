package api

import (
	utils "github.com/kevinbarbary/go-lms/utils"
	"encoding/json"
	"log"
	"strconv"
)

type userTutorial struct {
	TutorialID    int          `json:"TutorialID"`
	TutorialTitle string       `json:"TutorialTitle"`
	LaunchURL     string       `json:"LaunchURL"`
	HighestScore  int          `json:"HighestScore"`
	Completed     bool         `json:"Completed"`
	CompletedDate string       `json:"CompletedDate"` // date / datetime ?
	FirstAccessed string       `json:"FirstAccessed"` // datetime
	LastAccessed  JsonDateTime `json:"LastAccessed"`
	TimesAccessed int          `json:"TimesAccessed"`
	Height        int          `json:"Height"`
	Width         int          `json:"Width"`
}

type userLesson struct {
	Title     string         `json:"Title"`
	Tutorials []userTutorial `json:"Tutorials"`
}

type UserEnrolment struct {
	EnrollID       int          `json:"EnrollID"`
	CourseID       int          `json:"CourseID"`
	CourseTitle    string       `json:"CourseTitle"`
	CourseOverview string       `json:"CourseOverview"`
	StartDate      string       `json:"StartDate"`     // date
	EndDate        string       `json:"EndDate"`       // date
	CompletedDate  string       `json:"CompletedDate"` // date / datetime ?
	TimesAccessed  int          `json:"TimesAccessed"`
	FirstAccessed  string       `json:"FirstAccessed"` // datetime
	LastAccessed   string       `json:"LastAccessed"`  // datetime
	Completed      bool         `json:"Completed"`
	Lessons        []userLesson `json:"Lessons"`
	CertificateURL string       `json:"CertificateURL"`
	CourseType     string       `json:"CourseType"`
	TotalDuration  int          `json:"TotalDuration"`
}

func (e UserEnrolment) NotValid() bool {
	return e.EnrollID == 0
}

func UserTutorials(token, useragent, site, loginId string, enrollId int) (UserEnrolment, string, string, error) {

	response, err := Call("GET", utils.Endpoint(utils.Concat("/enrolment/", loginId, "/", strconv.Itoa(enrollId))), token, useragent, site, nil, true)
	if err != nil {
		log.Print("UserTutorials Error - invalid response from API call... ", err.Error())
		return UserEnrolment{}, "", "", err
	}

	data, e, help, _, token, user := extract(response)

	if e != "" {
		log.Print("UserTutorials Error... ", e)
	}
	if help != "" {
		log.Print("UserTutorials help... ", help)
	}

	if data == nil {
		log.Print("UserTutorials... NO DATA")
		return UserEnrolment{}, token, user, nil
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		log.Print("UserTutorials - Marshal fail... ", err.Error())
		return UserEnrolment{}, token, user, err
	}

	var val UserEnrolment
	err = json.Unmarshal(byteData, &val)
	if err != nil {
		log.Print("UserTutorials - Unmarshal fail... ", err.Error())
		return UserEnrolment{}, token, user, err
	}

	return val, token, user, nil
}
