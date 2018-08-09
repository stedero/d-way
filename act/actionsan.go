package act

import (
	"io"
	"net/http"

	"ibfd.org/d-way/config"
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
func Clean(r io.ReadCloser, cookies []*http.Cookie) (*StepResult, error) {
	return actionSan.clean(r, cookies)
}

// Sanitize calls the docsan service to clean a HTML document
func (action *ActionSan) clean(r io.ReadCloser, cookies []*http.Cookie) (*StepResult, error) {
	defer r.Close()
	req, err := http.NewRequest("POST", action.url, r)
	req.Header.Set("Content-type", "text/html")
	copyCookies(req, cookies)
	resp, err := action.client.Do(req)
	if err != nil {
		return nil, err
	}
	return NewStepResult().SetResponse(resp), err
}
