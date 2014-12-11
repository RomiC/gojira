package jira

import (
	"errors"
	"fmt"
)

// JiraRespInterface desribes interface for JiraResponse
type JiraRespInterface interface {
	GetErrors() error
}

// JiraResp base structure to store jira response
type JiraResp struct {
	ErrorMessages []string
	Errors        map[string]interface{}
}

// GetErrors return error in response
func (jr *JiraResp) GetErrors() error {
	errMsg := ""

	if jr.ErrorMessages != nil {
		for _, v := range jr.ErrorMessages {
			errMsg += (v + ",")
		}
	}

	if jr.Errors != nil {
		for i, v := range jr.Errors {
			errMsg += fmt.Sprintf("(%s) %s,", i, v)
		}
	}

	if len(errMsg) > 0 {
		return errors.New(errMsg[:len(errMsg)-1])
	} else {
		return nil
	}
}
