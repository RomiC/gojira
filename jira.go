/*
Small package providing following operations via JIRA REST API

	Issues
		* GetIssue — Get iformation about issue
		* PostIssue — Add new issue
*/

package jira

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiUri string = "/rest/api/2"
	debug  bool   = false
)

type Jira struct {
	Url      string
	Login    string
	Password string

	Client *http.Client
}

// Jira constructor of jira-object
func NewJira(url, login, password string) *Jira {
	j := &Jira{Url: url, Login: login, Password: password}

	j.Client = &http.Client{}

	return j
}

// request is private method for making JIRA-requests
// and parsing respons
func (j *Jira) request(method, url, body string) ([]byte, error) {
	bodyReader := bytes.NewBuffer([]byte(body))

	req, err := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(j.Login, j.Password)

	resp, err := j.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return []byte(""), &JiraError{req.URL.String(), body, "", "Authorization failed!"}
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	return respBytes, nil
}

// parse is private method for parsing JIRA API responses
func (j *Jira) parse(resp []byte, obj JiraRespInterface) error {
	err := json.Unmarshal(resp, &obj)
	if err != nil {
		return err
	}

	err = obj.GetErrors()

	if err != nil {
		return err
	}

	return nil
}
