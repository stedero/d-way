package act

import (
	"fmt"
	"net/http"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionOTCC defines the action that calls the OTCC service to get a document
type ActionOTCC struct {
	url    string
	client *http.Client
}

var actionOTCC *ActionOTCC

func init() {
	config := cfg.GetMatcher()
	actionOTCC = &ActionOTCC{config.OtccURL, NewHTTPClient()}
}

// OTCC calls the OTCC service to fetch a document.
func OTCC(document *doc.Source, cookies []*http.Cookie) (*StepResult, error) {
	target := actionOTCC.target(document.Path())
	req, err := http.NewRequest("GET", target, nil)
	setUserAgent(req)
	req.Header.Set("Accept", "application/pdf")
	copyCookies(req, cookies)
	resp, err := actionOTCC.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &ActionError{resp.StatusCode, fmt.Sprintf("failed to call OTCC for %s", document.Path())}
	}
	return NewContentResult().SetResponse(resp), err
}

func (action *ActionOTCC) target(path string) string {
	return action.url + path
}
