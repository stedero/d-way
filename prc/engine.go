package prc

import (
	"fmt"
	"io"

	"ibfd.org/d-way/act"
)

// Exec executes all processes for a document
func Exec(job *Job, w io.Writer) {
	fmt.Fprintf(w, "Starting job for: %s\n", job.document)
	for i, process := range job.rule.Processes {
		fmt.Fprintf(w, "executing step %d %s\n", i+1, process)
		switch process {
		case "GET":
			action := act.NewActionGet()
			r, err := action.Get(job.document)
			if err != nil {
				fmt.Fprintf(w, "error %s %s: %v\n", process, job.document, err)
			} else {
				_, err := io.Copy(w, r)
				if err != nil {
					fmt.Fprintf(w, "error copying %s: %v\n", job.document, err)
				}
			}
		default:
		}
	}
	fmt.Fprintf(w, "Finished job\n")
}
