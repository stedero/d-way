package act

import (
	"net/http"
	"fmt"
	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionSDRM defines the action that adds Social DRM statement to a document
type ActionSDRM struct {
	url    string
	client *http.Client
}

var actionSDRM *ActionSDRM

func init() {
	config := cfg.GetMatcher()
	actionSDRM = &ActionSDRM{config.SdrmURL, NewHTTPClient()}
}

// SDRM calls the Soda service to add Social DRM to a document.
func SDRM(document *doc.Source, cookies []*http.Cookie) (*StepResult, error) {
	target := actionSDRM.target(document.Path())
	req, err := http.NewRequest("GET", target, nil)
	setUserAgent(req)
	req.Header.Set("Accept", "application/pdf")
	copyCookies(req, cookies)
	resp, err := actionSDRM.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &ActionError{resp.StatusCode, fmt.Sprintf("failed to add social DRM to %s", document.Path())}
	}
	return NewStepResult().SetResponse(resp), err
}

func (action *ActionSDRM) target(path string) string {
	return action.url + path
}
