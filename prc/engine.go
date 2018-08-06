package prc

import (
	"fmt"
	"io"

	"ibfd.org/d-way/act"
)

// Exec executes all processing steps for a document
func Exec(job *Job, writer io.Writer) {
	fmt.Fprintf(writer, "Starting job for: %s\n", job.document)
	var reader io.ReadCloser
	var err error
	for i, step := range job.rule.Steps {
		fmt.Fprintf(writer, "executing step %d %s\n", i+1, step)
		switch step {
		case "GET":
			action := act.NewActionGet()
			reader, err = action.Get(job.document)
			if err != nil {
				fmt.Fprintf(writer, "error %s %s: %v\n", step, job.document, err)
			}
		case "CLEAN":
			action := act.NewActionSan()
			reader, err = action.Sanitize(reader)
			if err != nil {
				fmt.Fprintf(writer, "error %s %s: %v\n", step, job.document, err)
			}
		default:
		}
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Fprintf(writer, "error copying %s: %v\n", job.document, err)
	}
	reader.Close()
	fmt.Fprintf(writer, "\nFinished job\n")
}
