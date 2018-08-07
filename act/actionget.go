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

// GetActionGet creates a document fetcher.
func GetActionGet() *ActionGet {
	return actionGet
}

// Get fetches a document.
func (action *ActionGet) Get(document *doc.Document) (*StepResult, error) {
	stepResult := NewStepResult("GET").Start()
	target := action.target(document.Path)
	log.Printf("Fetching: %s", target)
	reader, err := os.Open(target)
	return stepResult.SetReader(reader).End(), err
}

func (action *ActionGet) target(path string) string {
	return action.basePath + path
}
