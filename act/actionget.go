package act

import (
	"log"
	"os"

	"ibfd.org/d-way/config"
	"ibfd.org/d-way/doc"
)

// ActionGet defines the action that fetches document.
type ActionGet struct {
	basePath string
}

var publicationsBasePath string
var actionGet *ActionGet

func init() {
	config := cfg.GetMatcher()
	actionGet = &ActionGet{config.PublicationsBasePath}
}

// Get fetches a document.
func Get(document *doc.Document) (*StepResult, error) {
	return actionGet.get(document)
}

// Get fetches a document.
func (action *ActionGet) get(document *doc.Document) (*StepResult, error) {
	target := action.target(document.Path())
	log.Printf("Fetching: %s", target)
	reader, err := os.Open(target)
	return NewStepResult().SetReader(reader).SetMimeType(document.MimeType()), err
}

func (action *ActionGet) target(path string) string {
	return action.basePath + path
}
