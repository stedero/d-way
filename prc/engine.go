package prc

import (
	"io"

	"ibfd.org/d-way/act"
)

// JobResult the result of executing a job.
type JobResult struct {
	Steps  []*act.StepResult
	Reader io.ReadCloser
}

// Exec executes all processing steps for a document
func Exec(job *Job) (*JobResult, error) {
	var result *act.StepResult
	var reader io.ReadCloser
	var err error
	jobResult := &JobResult{Steps: make([]*act.StepResult, 0, len(job.rule.Steps))}
	for _, step := range job.rule.Steps {
		switch step {
		case "GET":
			action := act.GetActionGet()
			result, err = action.Get(job.document)
			if err != nil {
				return nil, err
			}
			reader = result.Reader()
			jobResult.add(result)
		case "CLEAN":
			action := act.GetActionSan()
			result, err = action.Sanitize(reader, job.cookies)
			if err != nil {
				return nil, err
			}
			reader = result.Reader()
			jobResult.add(result)
		default:
		}
	}
	jobResult.Reader = reader
	return jobResult, err
}

func (jobResult *JobResult) add(stepResult *act.StepResult) {
	jobResult.Steps = append(jobResult.Steps, stepResult)
}
