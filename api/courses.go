package api

import (
	utils "github.com/kevinbarbary/go-lms/utils"
	"encoding/json"
	"log"
	"strconv"
)

type CourseTag struct {
	TagID int    `json:"TagID"`
	Tag   string `json:"Tag"`
}

type CourseTagType struct {
	TagType string      `json:"TagType"`
	Tags    []CourseTag `json:"Tags"`
}

type Publisher struct {
	PublisherID   string `json:"PublisherID"`
	PublisherName string `json:"PublisherName"`
}

type Course struct {
	CourseID            int         `json:"CourseID"`
	CourseTitle         string      `json:"CourseTitle"`
	DurationID          int         `json:"DurationID"`
	Duration            int         `json:"Duration"`
	DurationDescription string      `json:"DurationDescription"`
	PublisherName       string      `json:"PublisherName"`
	PublisherID         string      `json:"PublisherID"`
	Image               string      `json:"Image"`
	Price               float32     `json:"Price"`
	Currency            string      `json:"Currency"`
	TrainingTime        string      `json:"TrainingTime"`
	CourseTags          []CourseTag `json:"CourseTags"`
}

type Pagination struct {
	Offset        int      `json:"CourseID"`
	Limit         int      `json:"Limit"`
	SortBy        string   `json:"SortBy"`
	SortAscending bool     `json:"SortAscending"`
	PublisherIDs  []string `json:"PublisherIDs"`
	CourseTags    []int    `json:"CourseTags"`
	Keywords      []string `json:"Keywords"`
}

type CoursesData struct {
	Courses    []Course        `json:"Courses"`
	Tags       []CourseTagType `json:"CourseTags"`
	Publishers []Publisher     `json:"Publishers"`
	Total      int             `json:"Total"`
	Previous   Pagination      `json:"Previous"`
	Next       Pagination      `json:"Next"`
}

func Courses(token, useragent, site string, index int, tags []int) (CoursesData, string, string, Timestamp, error) {

	// hard-coded for 24 courses per page

	response, err := Call("POST", utils.Endpoint("/courses"), token, useragent, site, Params{"Offset": strconv.Itoa((index - 1) * 24), "Limit": "24", "CourseTags": tags}, true)

	if err != nil {
		log.Print("Courses Error - invalid response from API call... ", err.Error())
		return CoursesData{}, "", "", 0, err
	}

	data, e, help, now, token, user := extract(response)

	if e != "" {
		log.Print("Courses Error... ", e)
	}
	if help != "" {
		log.Print("Courses help... ", help)
	}

	if data == nil {
		log.Print("Courses... NO DATA")
		return CoursesData{}, token, user, now, err
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		log.Print("Courses - Marshal fail... ", err.Error())
		return CoursesData{}, token, user, now, err
	}

	var val CoursesData
	err = json.Unmarshal(byteData, &val)
	if err != nil {
		log.Print("Courses - Unmarshal fail... ", err.Error())
		return CoursesData{}, token, user, now, err
	}

	return val, token, user, now, nil
}
