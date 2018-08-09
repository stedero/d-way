package prc

import (
	"io"

	"ibfd.org/d-way/act"
)

type actFunc func() (*act.StepResult, error)

// Exec executes all processing steps for a document
func Exec(job *Job) (*JobResult, error) {
	var result *act.TimedResult
	var reader io.ReadCloser
	var err error
	jobResult := NewJobResult(len(job.rule.Steps))
	for _, step := range job.rule.Steps {
		switch step {
		case "GET":
			result, err = exec(step, getter(job))
		case "CLEAN":
			result, err = exec(step, cleaner(job, reader))
		case "SDRM":
			result, err = exec(step, soda(job))
		default:
			result, err = nil, nil
		}
		if err != nil {
			return nil, err
		}
		if result != nil {
			reader = result.Reader()
			jobResult.add(result)
		}
	}
	return jobResult, err
}

func exec(step string, action actFunc) (*act.TimedResult, error) {
	timedResult := act.NewTimedResult(step).Start()
	stepResult, err := action()
	if err != nil {
		return nil, err
	}
	timedResult.End()
	timedResult.SetStepResult(stepResult)
	return timedResult, err
}

func getter(job *Job) actFunc {
	return func() (*act.StepResult, error) {
		return act.Get(job.document)
	}
}

func cleaner(job *Job, r io.ReadCloser) actFunc {
	return func() (*act.StepResult, error) {
		return act.Clean(r, job.cookies)
	}
}

func soda(job *Job) actFunc {
	return func() (*act.StepResult, error) {
		return act.SDRM(job.document, job.cookies)
	}
}
