package act

import (
	"fmt"
	"net/http"
	"strings"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionXTOJ defines the action that converts a XML document to JSON.
type ActionXTOJ struct {
	url    string
	client *http.Client
}

var actionXTOJ *ActionXTOJ

func init() {
	config := cfg.GetMatcher()
	actionXTOJ = &ActionXTOJ{config.XtojURL, NewHTTPClient()}
}

// XTOJ calls the xtoj service to to convert a XML document to JSON.
func XTOJ(document *doc.Source, cookies []*http.Cookie) (*StepResult, error) {
	target := actionXTOJ.target(document)
	req, err := http.NewRequest("GET", target, nil)
	setUserAgent(req)
	req.Header.Set("Accept", "application/json")
	copyCookies(req, cookies)
	resp, err := actionXTOJ.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &ActionError{resp.StatusCode, fmt.Sprintf("failed to get JSON for %s", document.Path())}
	}
	return NewStepResult().SetResponse(resp), err
}

func (action *ActionXTOJ) target(document *doc.Source) string {
	return action.url + replaceExtensionIfNeeded(document)
}

func replaceExtensionIfNeeded(document *doc.Source) string {
	if document.Extension() == ".json" {
		return strings.Replace(document.Path(), ".json", ".xml", -1)
	}
	return document.Path()
}
