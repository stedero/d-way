package prc

import (
	"ibfd.org/d-way/act"
	"ibfd.org/d-way/doc"
	log "ibfd.org/d-way/log4u"
)

type actFunc func() (*act.StepResult, error)

// Exec executes all processing steps for a job.
func Exec(job *Job, src *doc.Source) (*JobResult, error) {
	var result *act.TimedResult
	var err error
	jobResult := NewJobResult(len(job.rule.Steps))
	var statResult *act.TimedResult
	for _, step := range job.rule.Steps {
		if !jobResult.Done() {
			statResult = nil
			switch step {
			case "CLEAN":
				result, err = exec(step, cleaner(job, src))
			case "GET":
				result, err = exec(step, getter(src))
			case "RESOLVE":
				result, err = exec(step, resolver(src))
			case "SDRM":
				result, err = exec(step, soda(job, src))
			case "STAT":
				statResult, err = exec(step, stat(job, src))
			case "XTOJ":
				result, err = exec(step, xtoj(job, src))
			default:
				log.Warnf("undefined step: %s", step)
			}
			if err != nil {
				return nil, err
			}
			if statResult != nil {
				jobResult.add(statResult)
			} else if result != nil {
				src = doc.ReaderSource(result.Reader())
				jobResult.add(result)
			}
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

func resolver(src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Resolve(src)
	}
}

func getter(src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Get(src)
	}
}

func cleaner(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Clean(src, job.cookies)
	}
}

func soda(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.SDRM(src, job.cookies)
	}
}

func stat(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.Stat(src, job.RequestModSince())
	}
}

func xtoj(job *Job, src *doc.Source) actFunc {
	return func() (*act.StepResult, error) {
		return act.XTOJ(src, job.cookies)
	}
}
