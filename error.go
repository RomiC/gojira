package jira

import (
	"fmt"
)

// JiraError represt error during JIRA API reqeust
type JiraError struct {
	RequestURL  string
	RequestData string
	Response    string
	Message     string
}

func (je *JiraError) Error() string {
	format := ""

	if debug {
		format = `JIRA error: %s
Debug:
	>>> %s, '%s'
	<<< %s
`
	} else {
		format = "JIRA error: %s\n"
	}
	return fmt.Sprintf(format, je.Message, je.RequestURL, je.RequestData, je.Response)
}
