package prc

import (
	"io"

	"ibfd.org/d-way/act"
)

// JobResult the result of executing a job.
type JobResult struct {
	Steps  []act.StepResult
	reader io.ReadCloser
}

// Exec executes all processing steps for a document
func Exec(job *Job) (io.ReadCloser, error) {
	var result *act.StepResult
	var reader io.ReadCloser
	var err error
	for _, step := range job.rule.Steps {
		switch step {
		case "GET":
			action := act.GetActionGet()
			reader, err = action.Get(job.document)
			if err != nil {
				return nil, err
			}
		case "CLEAN":
			action := act.GetActionSan()
			result, err = action.Sanitize(reader, job.cookies)
			if err != nil {
				return nil, err
			}
			reader = result.Reader()
		default:
		}
	}
	return reader, err
}
