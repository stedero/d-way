package act

import (
	"io/ioutil"
	"log"
	"strings"

	"ibfd.org/d-way/doc"
)

// ActionGet defines the action that fetches document.
type ActionGet struct {
}

var host string
var publicationsBasePath string

func init() {
	host = "https://dev-online.ibfd.org"
	publicationsBasePath = "/collections"
}

// NewActionGet creates a document fetcher.
func NewActionGet() *ActionGet {
	return &ActionGet{}
}

// Get fetches a document.
func (action *ActionGet) Get(document *doc.Document) ([]byte, error) {
	target := host + publicationsBasePath + strings.TrimPrefix(document.Path, "/data")
	log.Printf("Fetching: %s", target)
	client := NewHTTPClient()
	response, err := client.Get(target)
	if err != nil {
		return nil, err
	}
	//	response.StatusCode
	return ioutil.ReadAll(response.Body)
}
