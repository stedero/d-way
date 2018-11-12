package act

import (
	"net/http"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionXTOJ defines the action that converts a XML document to JSON.
type ActionXTOJ struct {
	url    string
	client *http.Client
}

var actionXTOJ *ActionXTOJ

func init() {
	config := cfg.GetMatcher()
	actionSan = &ActionSan{config.CleanURL, NewHTTPClient()}
}

// XTOJ calls the xtoj service to to convert a XML document to JSON.
func XTOJ(src *doc.Source, cookies []*http.Cookie) (*StepResult, error) {
	return nil, &ActionError{405, "XTOJ action not implemented yet."}
}
