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
		switch process {
		case "GET":
			// action := act.NewActionGet()
			// result, err := action.Get(job.document)
			// if err != nil {
			// 	fmt.Fprintf(w, "error fetching %s: %v\n", job.document, err)
			// } else {
			// 	fmt.Fprintf(w, "\tresult length: %d\n", len(result))
			// 	fmt.Fprintf(w, "\tresult: %s\n", string(result))
			// }
		default:
		}
	}
	fmt.Fprintf(w, "Finished job\n")
}
