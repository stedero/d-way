package act

import (
	"os"

	"ibfd.org/d-way/log"

	"ibfd.org/d-way/config"
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
func Get(document *doc.Document) (*StepResult, error) {
	target := actionGet.target(document.Path())
	log.Debugf("Fetching: %s", target)
	reader, err := os.Open(target)
	return NewStepResult().SetReader(reader).SetMimeType(document.MimeType()), err
}

func (action *ActionGet) target(path string) string {
	return action.basePath + path
}
