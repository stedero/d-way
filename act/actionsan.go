package act

import (
	"io"
	"net/http"
)

// ActionSan defines the action that sanitizes a document.
type ActionSan struct {
	httpClient *http.Client
}

const docsanURL = "https://docsan.herokuapp.com"

// NewActionSan creates a document sanitizer.
func NewActionSan() *ActionSan {
	return &ActionSan{httpClient: NewHTTPClient()}
}

// Sanitize calls the sanitizer service to clean a HTML document
func (action *ActionSan) Sanitize(r io.Reader) (io.Reader, error) {
	response, err := action.httpClient.Post(docsanURL, "Content-type: text/html", r)
	if err != nil {
		return nil, err
	}
	return response.Body, err
}
