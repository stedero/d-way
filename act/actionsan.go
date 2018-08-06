package act

import (
	"io"
	"net/http"
	"ibfd.org/d-way/config"
)

// ActionSan defines the action that sanitizes a document.
type ActionSan struct {
	url string
	client *http.Client
}

var actionSan *ActionSan

func init() {
	config := cfg.GetMatcher()
	actionSan = &ActionSan{config.CleanURL, NewHTTPClient()}
}

// GetActionSan returns the document sanitizer.
func GetActionSan() *ActionSan {
	return actionSan;
}

// Sanitize calls the sanitizer service to clean a HTML document
func (action *ActionSan) Sanitize(r io.ReadCloser) (io.ReadCloser, error) {
	defer r.Close()
	response, err := action.client.Post(action.url, "Content-type: text/html", r)
	if err != nil {
		return nil, err
	}
	return response.Body, err
}
