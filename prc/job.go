package prc

import (
	"net/http"

	"ibfd.org/d-way/doc"
	"ibfd.org/d-way/rule"
)

// Job describes the steps to be executed for a document
type Job struct {
	document *doc.Document
	rule     *rule.Rule
	cookies  []*http.Cookie
}

// NewJob creates a Job
func NewJob(d *doc.Document, r *rule.Rule, cookies []*http.Cookie) *Job {
	return &Job{d, r, cookies}
}
