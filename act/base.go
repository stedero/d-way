package act

import (
	"io"
	"net/http"
	"time"
)

const defaultTimeoutSeconds = 10

// StepResult the results of executing one step.
type StepResult struct {
	Step     string
	response *http.Response
}

// NewHTTPClient creates a HTPP client with a default timeout.
func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: time.Second * defaultTimeoutSeconds}
}

// Reader gets the response reader.
func (stepResult *StepResult) Reader() io.ReadCloser {
	return stepResult.response.Body
}
