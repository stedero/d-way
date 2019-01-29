package act

import (
	"net/http"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionSan defines the action that sanitizes a document.
type ActionSan struct {
	url    string
	client *http.Client
}

var actionSan *ActionSan

func init() {
	config := cfg.GetMatcher()
	actionSan = &ActionSan{config.CleanURL, NewHTTPClient()}
}

// Clean calls the docsan service to clean a HTML document
func Clean(src *doc.Source, cookies []*http.Cookie) (*StepResult, error) {
	req, err := http.NewRequest("POST", actionSan.url, src.Reader())
	setUserAgent(req)
	req.Header.Set("Content-type", "text/html")
	copyCookies(req, cookies)
	resp, err := actionSan.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &ActionError{resp.StatusCode, "failed to sanitize document"}
	}
	return NewContentResult().SetResponse(resp), err
}
