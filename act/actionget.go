package act

import (
	"net/http"
	"os"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionGet defines the action that fetches a document.
type ActionGet struct {
	basePath string
}

var actionGet *ActionGet

func init() {
	config := cfg.GetMatcher()
	actionGet = &ActionGet{config.PublicationsBasePath}
}

// Get fetches a document.
func Get(document *doc.Source) (*StepResult, error) {
	path := document.Path()
	if path == "" {
		return nil, &ActionError{http.StatusBadRequest, "no file specified"}
	}
	target := actionGet.target(document.Path())
	reader, err := os.Open(target)
	if err != nil {
		return nil, &ActionError{http.StatusNotFound, err.Error()}
	}
	return NewStepResult().SetReader(reader).SetMimeType(document.MimeType()).SetStatusCode(http.StatusOK), err
}

func (action *ActionGet) target(path string) string {
	return action.basePath + path
}
