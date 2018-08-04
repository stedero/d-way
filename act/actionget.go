package act

import (
	"io"
	"log"
	"os"

	"ibfd.org/d-way/doc"
)

// ActionGet defines the action that fetches document.
type ActionGet struct {
}

var publicationsBasePath string

func init() {
	publicationsBasePath = "/Users/steef/Desktop/d-way"
}

// NewActionGet creates a document fetcher.
func NewActionGet() *ActionGet {
	return &ActionGet{}
}

// Get fetches a document.
func (action *ActionGet) Get(document *doc.Document) (io.Reader, error) {
	target := publicationsBasePath + document.Path
	log.Printf("Fetching: %s", target)
	return os.Open(target)
}
