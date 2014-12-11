package jira

import (
	"fmt"
)

// JiraPostIssueResp struct describes response
// from Jira post issue method
type JiraIssue struct {
	JiraResp
	Id     string
	Key    string
	Self   string
	Fields map[string]interface{}
}

// PostIssue calls /rest/api/…/issue method via POST
// to add one more issue
func (j *Jira) PostIssue(projectId int, summary, description string) (*JiraIssue, error) {
	body := `
{
	"fields": {
        "project": {
            "id": "` + fmt.Sprintf("%d", projectId) + `"
        },
        "summary": "` + summary + `",
        "description": "` + description + `",
        "issuetype": {
        	"id": "3"
        }
    }
}
`

	respBytes, err := j.request("POST", j.Url+apiUri+"/issue/", body)
	if err != nil {
		return nil, err
	}

	respObj := &JiraIssue{}

	err = j.parse(respBytes, respObj)
	if err != nil {
		return nil, &JiraError{j.Url + apiUri + "/issue/", body, string(respBytes), err.Error()}
	}

	return respObj, nil
}

// GetIssue makes GET-request to /rest/api/…/issue
// to get info about the issue
func (j *Jira) GetIssue(issue string, fields ...string) (string, error) {
	respBytes, err := j.request("GET", j.Url+apiUri+"/issue/"+issue+"", "")

	if err != nil {
		return "", err
	}

	fmt.Println(string(respBytes))

	// respObj := &JiraRespGetIssue{}

	return "", nil
}
