package act

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"ibfd.org/d-way/cfg"
)

const defaultTimeoutSeconds = 10

// ResultType indicates whether the result contains content or an action such as redirect.
type ResultType int

const (
	// Content result type
	Content ResultType = iota
	// Action result type
	Action
)

// StepResult the results of executing one step.
type StepResult struct {
	resultType ResultType
	response   *http.Response
	content    string
	reader     io.ReadCloser
	mimeType   string
	statusCode int
}

// TimedResult the result of executing one step with timing
type TimedResult struct {
	Step     string
	start    time.Time
	end      time.Time
	Duration time.Duration
	StepResult
}

// ActionError defines errors that are created by actions.
type ActionError struct {
	StatusCode int
	Msg        string
}

func (e *ActionError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.Msg)
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

// StatusCode returns the HTTP status code.
func (stepResult *StepResult) StatusCode() int {
	if stepResult.statusCode > -1 {
		return stepResult.statusCode
	}
	return stepResult.response.StatusCode
}

// Content get the content created by a step.
func (stepResult *StepResult) Content() string {
	if stepResult.response != nil {
		result, _ := ioutil.ReadAll(stepResult.response.Body)
		return string(result)
	}
	return stepResult.content
}

// ResultType get the result type.
func (stepResult *StepResult) ResultType() ResultType {
	return stepResult.resultType
}

// NewContentResult creates a step result.
func NewContentResult() *StepResult {
	return &StepResult{statusCode: -1, resultType: Content}
}

// NewActionResult creates a step result.
func NewActionResult() *StepResult {
	return &StepResult{statusCode: -1, resultType: Action}
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

// SetContent set the step content.
func (stepResult *StepResult) SetContent(ct string) *StepResult {
	stepResult.content = ct
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

// SetStatusCode sets the status code.
func (stepResult *StepResult) SetStatusCode(sc int) *StepResult {
	stepResult.statusCode = sc
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
	timedResult.statusCode = stepResult.statusCode
	timedResult.content = stepResult.content
	timedResult.resultType = stepResult.resultType
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
