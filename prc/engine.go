package prc

import (
	"io"

	"ibfd.org/d-way/act"
)

// Exec executes all processing steps for a document
func Exec(job *Job) (io.ReadCloser, error) {
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
			reader, err = action.Sanitize(reader)
			if err != nil {
				return nil, err
			}
		default:
		}
	}
	return reader, err
}
