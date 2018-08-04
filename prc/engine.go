package prc

import (
	"fmt"
	"io"
)

// Exec executes all processes for a document
func Exec(job *Job, w io.Writer) {
	fmt.Fprintf(w, "Starting job for: %s\n", job.document)
	for i, process := range job.rule.Processes {
		fmt.Fprintf(w, "executing step %d %s\n", i+1, process)
	}
	fmt.Fprintf(w, "Finished job\n")
}
