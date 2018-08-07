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
	reader   io.ReadCloser
	start    time.Time
	end      time.Time
	Duration time.Duration
}

// NewHTTPClient creates a HTPP client with a default timeout.
func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: time.Second * defaultTimeoutSeconds}
}

// Reader gets the response reader.
func (stepResult *StepResult) Reader() io.ReadCloser {
	if stepResult.reader != nil {
		return stepResult.reader
	}
	return stepResult.response.Body
}

// NewstepResult creates a step result.
func NewStepResult(step string) *StepResult {
	return &StepResult{Step: step}
}

// SetReader sets the reader.
func (stepResult *StepResult) SetReader(reader io.ReadCloser) *StepResult {
	stepResult.reader = reader
	return stepResult
}

// SetResponse sets the HTTP response.
func (stepResult *StepResult) SetResponse(response *http.Response) *StepResult {
	stepResult.response = response
	return stepResult
}

// Start signals the start of the process step.
func (stepResult *StepResult) Start() *StepResult {
	stepResult.start = time.Now()
	return stepResult
}

// End signals the end of the process step.
func (stepResult *StepResult) End() *StepResult {
	stepResult.end = time.Now()
	stepResult.Duration = stepResult.end.Sub(stepResult.start)
	return stepResult
}
