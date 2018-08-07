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

// GetActionSan returns the document sanitizer.
func GetActionSan() *ActionSan {
	return actionSan
}

// Sanitize calls the sanitizer service to clean a HTML document
func (action *ActionSan) Sanitize(r io.ReadCloser, cookies []*http.Cookie) (*StepResult, error) {
	defer r.Close()
	stepResult := NewStepResult("CLEAN").Start()
	req, err := http.NewRequest("POST", action.url, r)
	req.Header.Set("Content-type", "text/html")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := action.client.Do(req)
	if err != nil {
		return nil, err
	}
	return stepResult.SetResponse(resp).End(), err
}
