package prc

import (
	"ibfd.org/d-way/doc"
	"ibfd.org/d-way/rule"
)

// Job describes the steps to be executed for a document
type Job struct {
	document *doc.Document
	rule     *rule.Rule
}

// NewJob creates a Job
func NewJob(d *doc.Document, r *rule.Rule) *Job {
	return &Job{d, r}
}
