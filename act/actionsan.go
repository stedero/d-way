package act

import (
	"io"
	"net/http"

	"ibfd.org/d-way/cfg"
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
func Clean(r io.Reader, cookies []*http.Cookie) (*StepResult, error) {
	req, err := http.NewRequest("POST", actionSan.url, r)
	setUserAgent(req)
	req.Header.Set("Content-type", "text/html")
	copyCookies(req, cookies)
	resp, err := actionSan.client.Do(req)
	if err != nil {
		return nil, err
	}
	return NewStepResult().SetResponse(resp), err
}
