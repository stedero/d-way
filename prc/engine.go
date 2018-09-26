package prc

import (
	"ibfd.org/d-way/act"
	"ibfd.org/d-way/doc"
)

type actFunc func() (*act.StepResult, error)

// Exec executes all processing steps for a job.
func Exec(job *Job, src *doc.Source) (*JobResult, error) {
	var result *act.TimedResult
	var err error
	jobResult := NewJobResult(len(job.rule.Steps))
	for _, step := range job.rule.Steps {
		switch step {
		case "RESOLVE":
			result, err = exec(step, resolver(job, src))
		case "GET":
			result, err = exec(step, getter(src))
		case "CLEAN":
			result, err = exec(step, cleaner(job, src))
		case "SDRM":
			result, err = exec(step, soda(job, src))
		default:
			result, err = nil, nil
		}
		if err != nil {
			return nil, err
		}
		if result != nil {
			src = doc.ReaderSource(result.Reader())
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

func resolver(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Resolve(src.String(), job.cookies)
	}
}

func getter(src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Get(src)
	}
}

func cleaner(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Clean(src.Reader(), job.cookies)
	}
}

func soda(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.SDRM(src, job.cookies)
	}
}
