package api

import (
	"../utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// @todo - Time.Unix ?
type Timestamp int64

type JsonDate time.Time

type Params map[string]string

func MergeParams(a, b Params) Params {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

//func (j JsonDate) MarshalJSON() ([]byte, error) {
//	return json.Marshal(j)
//}

func (j JsonDate) ToTime() time.Time {
	return time.Time(j)
}

func (j JsonDate) Format(s string) string {
	return j.ToTime().Format(s)
}

func (j JsonDate) ToDate() string {
	return j.Format("2006-01-02")
}

func Creds(site, key string) (Params, error) {
	return Params{"SiteID": site, "SiteKey": key}, nil
}

type response struct {
	Data      interface{} `json:"data"`
	Error     string      `json:"error"`
	Help      string      `json:"help"`
	Timestamp Timestamp   `json:"timestamp"`
	Token     string      `json:"token"`
	User      string      `json:"user"`
	Session   int64       `json:"session"`
}

func extract(data string) (interface{}, string, string, Timestamp, string, string) {

	jsonData := []byte(data)

	var resp response
	err := json.Unmarshal(jsonData, &resp)
	if err != nil {
		log.Print("JSON Extract Error... ", err.Error())
		return nil, "", "", 0, "", ""
	}

	return resp.Data, resp.Error, resp.Help, resp.Timestamp, resp.Token, resp.User
}

// @todo - reduce duplication...
func Call(method, endpoint, token string, payload Params) (string, error) {

	if method == "GET" {

		client := http.Client{
			Timeout: time.Second * 5,
		}

		request, err := http.NewRequest(method, endpoint, bytes.NewBuffer(nil))
		if err != nil {
			log.Print("API GET Request Error... ", err.Error())
			return "", err
		}

		request.Header.Add("Accept", "application/json")
		if token != "" {
			request.Header.Set("Authorization", utils.Concat("Bearer ", token))
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Print("API Call Error - invalid response from GET... ", err.Error())
			return "", err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("API GET Response - no body... ", err.Error())
			return "", err
		}

		return string(body), nil

	} else if method == "POST" {

		data, err := json.Marshal(payload)
		if err != nil {
			log.Print("API Payload Marshall Error... ", err.Error())
			return "", err
		}

		client := http.Client{
			Timeout: time.Second * 5,
		}

		request, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
		if err != nil {
			log.Print("API POST Request Error... ", err.Error())
			return "", err
		}

		request.Header.Add("Accept", "application/json")
		request.Header.Set("Content-type", "application/json")
		if token != "" {
			request.Header.Set("Authorization", utils.Concat("Bearer ", token))
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Print("API Call Error - invalid response from POST... ", err.Error())
			return "", err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("API POST Response - no body... ", err.Error())
			return "", err
		}

		return string(body), nil

	}

	return "", nil
}

func (t Timestamp) ToTime() time.Time {
	return time.Unix(int64(t)/10000, 0)
}

func (t Timestamp) ToUnix() int64 {
	return t.ToTime().Unix()
}

func (t Timestamp) ToDate() string {
	return t.ToTime().Format("2006-01-02")
}

func (t Timestamp) ToDatetime() string {
	return t.ToTime().Format("2006-01-02 15:04:05")
}

func (t Timestamp) BeforeEnd(d JsonDate) bool {
	// returns if t is before the end of d, i.e. t < (d + 1 day)

	// convert t to time.Time (t is already a date) -> t2
	// convert d to date and add one day then convert to time.Time -> d2
	// return t2 < d2

	d1 := d.ToTime()
	d2 := time.Date(d1.Year(), d1.Month(), d1.Day(), 0, 0, 0, 0, d1.Location()).AddDate(0, 0, 1)

	return t.ToTime().Before(d2)
}
