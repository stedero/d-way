package prc

import (
	"io"
	"net/http"

	"ibfd.org/d-way/act"
	"ibfd.org/d-way/rule"
)

// Job defines the rules and hold results for the steps to execute.
type Job struct {
	rule    *rule.Rule
	cookies []*http.Cookie
}

// JobResult the result of executing a job.
type JobResult struct {
	steps  []*act.TimedResult
	reader io.ReadCloser
	last   *act.TimedResult
}

// NewJob creates a Job
func NewJob(r *rule.Rule, cookies []*http.Cookie) *Job {
	return &Job{r, cookies}
}

// NewJobResult creates a job result.
func NewJobResult(stepCount int) *JobResult {
	return &JobResult{steps: make([]*act.TimedResult, 0, stepCount)}
}

func (jobResult *JobResult) add(timedResult *act.TimedResult) {
	jobResult.steps = append(jobResult.steps, timedResult)
	jobResult.last = timedResult
	jobResult.reader = timedResult.Reader()
}

// Reader returns the reader from the last step that was executed.
func (jobResult *JobResult) Reader() io.ReadCloser {
	return jobResult.last.Reader()
}

// ContentType returns the content type of the last step that was executed.
func (jobResult *JobResult) ContentType() string {
	contentType := "";
	response := jobResult.last.Response()
	if response != nil {
		contentType = response.Header["Content-Type"][0]
	}
	if contentType != "" {
		return contentType
	} else {
		return jobResult.last.MimeType()
	}
}

// Response get the response of the last step that was executed.
func (jobResult *JobResult) Response() *http.Response {
	return jobResult.last.Response()
}

// Steps returns the steps that where executed.
func (jobResult *JobResult) Steps() []*act.TimedResult {
	return jobResult.steps
}
