package act

import (
	"io"
	"net/http"
	"time"
	"ibfd.org/d-way/cfg"
)

const defaultTimeoutSeconds = 10

// StepResult the results of executing one step.
type StepResult struct {
	response *http.Response
	reader   io.ReadCloser
	mimeType string
}

// TimedResult the result of executing one step with timing
type TimedResult struct {
	Step     string
	start    time.Time
	end      time.Time
	Duration time.Duration
	StepResult
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

// Response gets the HTTP response.
func (stepResult *StepResult) Response() *http.Response {
	return stepResult.response
}

// MimeType returns the mime type of this result.
func (stepResult *StepResult) MimeType() string {
	return stepResult.mimeType
}

// NewStepResult creates a step result.
func NewStepResult() *StepResult {
	return &StepResult{}
}

// NewTimedResult creates a timed result.
func NewTimedResult(step string) *TimedResult {
	return &TimedResult{Step: step}
}

// SetReader sets the reader.
func (stepResult *StepResult) SetReader(reader io.ReadCloser) *StepResult {
	stepResult.reader = reader
	return stepResult
}

// SetResponse sets the HTTP response.
func (stepResult *StepResult) SetResponse(response *http.Response) *StepResult {
	stepResult.response = response
	stepResult.mimeType = response.Header["Content-Type"][0]
	return stepResult
}

// SetMimeType sets the mime type.
func (stepResult *StepResult) SetMimeType(mimeType string) *StepResult {
	stepResult.mimeType = mimeType
	return stepResult
}

// Start signals the start of the process step.
func (timedResult *TimedResult) Start() *TimedResult {
	timedResult.start = time.Now()
	return timedResult
}

// End signals the end of the process step.
func (timedResult *TimedResult) End() *TimedResult {
	timedResult.end = time.Now()
	timedResult.Duration = timedResult.end.Sub(timedResult.start)
	return timedResult
}

// SetStepResult added the step result.
func (timedResult *TimedResult) SetStepResult(stepResult *StepResult) *TimedResult {
	timedResult.response = stepResult.response
	timedResult.reader = stepResult.reader
	timedResult.mimeType = stepResult.mimeType
	return timedResult
}

func setUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", cfg.GetUserAgent())
}

func copyCookies(request *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
}
