package prc

import (
	"io"
	"net/http"

	"ibfd.org/d-way/act"
	"ibfd.org/d-way/doc"
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
			result, err = exec(step, getter(job.document))
		case "CLEAN":
			result, err = exec(step, cleaner(reader, job.cookies))
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

func getter(doc *doc.Document) actFunc {
	return func() (*act.StepResult, error) {
		return act.Get(doc)
	}
}

func cleaner(r io.ReadCloser, cookies []*http.Cookie) actFunc {
	return func() (*act.StepResult, error) {
		return act.Clean(r, cookies)
	}
}

func soda(job *Job) actFunc {
	return func() (*act.StepResult, error) {
		return act.SDRM(job.document, job.cookies)
	}
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
