package prc

import (
	"io"
	"net/http"
	"time"

	"ibfd.org/d-way/act"
	"ibfd.org/d-way/rule"
)

// Job defines the rules and hold results for the steps to execute.
type Job struct {
	rule        *rule.Rule
	cookies     []*http.Cookie
	reqModSince *time.Time
}

// JobResult the result of executing a job.
type JobResult struct {
	steps            []*act.TimedResult
	last             *act.TimedResult
	total            *act.TimedResult
	lastModifiedDate string
	stopJob          bool
}

// NewJob creates a Job
func NewJob(r *rule.Rule, cookies []*http.Cookie, reqModSince *time.Time) *Job {
	return &Job{r, cookies, reqModSince}
}

// NewJobResult creates a job result.
func NewJobResult(stepCount int) *JobResult {
	return &JobResult{steps: make([]*act.TimedResult, 0, stepCount)}
}

// Start starts the job time.
func (jobResult *JobResult) Start() {
	jobResult.total = act.NewTimedResult("Total").Start()
}

// End stops the job timer.
func (jobResult *JobResult) End() {
	jobResult.total.End()
}

// Done determines whether the job is finished
// and any remaining steps can be skipped.
func (jobResult *JobResult) Done() bool {
	return jobResult.stopJob
}

func (jobResult *JobResult) add(timedResult *act.TimedResult) {
	jobResult.steps = append(jobResult.steps, timedResult)
	jobResult.last = timedResult
	if timedResult.Step == "STAT" {
		jobResult.lastModifiedDate = timedResult.Content()
		jobResult.stopJob = timedResult.StatusCode() == http.StatusNotModified
	}
}

// Content returns the content from the last step that was executed.
func (jobResult *JobResult) Content() string {
	return jobResult.last.Content()
}

// Reader returns the reader from the last step that was executed.
func (jobResult *JobResult) Reader() io.ReadCloser {
	return jobResult.last.Reader()
}

// ContentType returns the content type of the last step that was executed.
func (jobResult *JobResult) ContentType() string {
	contentType := ""
	response := jobResult.last.Response()
	if response != nil {
		contentType = response.Header["Content-Type"][0]
	}
	if contentType != "" {
		return contentType
	}
	return jobResult.last.MimeType()
}

// LastModifiedDate returns the modification date of the requested resource.
func (jobResult *JobResult) LastModifiedDate() string {
	return jobResult.lastModifiedDate
}

// StatusCode returns the HTTP status code.
func (jobResult *JobResult) StatusCode() int {
	return jobResult.last.StatusCode()
}

// RequestModSince returns the value of the request header If-modified-since.
func (job *Job) RequestModSince() *time.Time {
	return job.reqModSince
}

// Response get the response of the last step that was executed.
func (jobResult *JobResult) Response() *http.Response {
	return jobResult.last.Response()
}

// Steps returns the steps that where executed.
func (jobResult *JobResult) Steps() []*act.TimedResult {
	return jobResult.steps
}

// Total returns the job execution result.
func (jobResult *JobResult) Total() *act.TimedResult {
	return jobResult.total
}
