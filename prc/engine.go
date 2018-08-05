package prc

import (
	"fmt"
	"io"

	"ibfd.org/d-way/act"
)

// Exec executes all processes for a document
func Exec(job *Job, writer io.Writer) {
	fmt.Fprintf(writer, "Starting job for: %s\n", job.document)
	var reader io.Reader
	var err error
	for i, process := range job.rule.Processes {
		fmt.Fprintf(writer, "executing step %d %s\n", i+1, process)
		switch process {
		case "GET":
			action := act.NewActionGet()
			reader, err = action.Get(job.document)
			if err != nil {
				fmt.Fprintf(writer, "error %s %s: %v\n", process, job.document, err)
			}
		case "SAN":
			action := act.NewActionSan()
			reader, err = action.Sanitize(reader)
			if err != nil {
				fmt.Fprintf(writer, "error %s %s: %v\n", process, job.document, err)
			}
		default:
		}
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Fprintf(writer, "error copying %s: %v\n", job.document, err)
	}
	fmt.Fprintf(writer, "\nFinished job\n")
}
