package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type jira struct {
	Url      string
	Login    string
	Password string

	Session string

	Client *http.Client

	Debug bool
}

var auth_uri string = "/rest/auth/latest/session"
var api_uri string = "/rest/api/latest"

// Jira constructor
func Jira(url, login, password string) (*jira, error) {
	j := &jira{Url: url, Login: login, Password: password}
	j.Debug = false

	cj, _ := cookiejar.New(new(cookiejar.Options))

	j.Client = &http.Client{
		Jar: cj,
	}

	err := j.auth()

	if err != nil {
		return nil, err
	}

	return j, nil
}

type jiraResp struct {
	Session       map[string]interface{}
	ErrorMessages []string
	Errors        json.RawMessage
	Key           string
	Fields        map[string]interface{}
}

func (j *jira) request(req *http.Request) (*jiraResp, error) {
	resp, err := j.Client.Do(req)

	if err != nil {
		return nil, err
	}

	respStr, err := ioutil.ReadAll(resp.Body)

	if j.Debug {
		fmt.Printf(">> %s\n", req.URL)
		fmt.Printf("<< %s\n\n", string(respStr))
	}

	respObj := &jiraResp{}
	err = json.Unmarshal(respStr, respObj)

	if err != nil {
		return nil, err
	}

	if respObj.ErrorMessages != nil {
		return nil, &JiraError{respObj.ErrorMessages[0]}
	}

	return respObj, nil
}

func (j *jira) auth() error {
	jsonStr := []byte(`{"username":"` + j.Login + `","password":"` + j.Password + `"}`)
	req, _ := http.NewRequest("POST", j.Url+auth_uri, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	respObj, err := j.request(req)

	if err != nil {
		return err
	}

	j.Session = respObj.Session["value"].(string)

	return nil
}

func (j *jira) GetIssue(issue string) (string, error) {
	if j.Session == "" {
		return "", &JiraError{Message: "Auth required!"}
	}

	req, _ := http.NewRequest("GET", j.Url+api_uri+"/issue/"+issue+"?fields=summary", nil)
	req.AddCookie(&http.Cookie{Name: "JSESSIONID", Value: j.Session})
	req.Header.Set("Content-Type", "application/json")

	respObj, err := j.request(req)

	if err != nil {
		return "", err
	}

	var ret bytes.Buffer
	fmt.Fprintf(&ret, "%s — %s", respObj.Key, respObj.Fields["summary"])

	return ret.String(), nil
}

type JiraError struct {
	Message string
}

func (je *JiraError) Error() string {
	return je.Message
}
