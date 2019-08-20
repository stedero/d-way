package act

import (
	"net/http"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionForbid defines the action that forbids access to a document.
type ActionForbid struct {
	basePath string
}

var actionForbid *ActionForbid

func init() {
	config := cfg.GetMatcher()
	actionForbid = &ActionForbid{config.PublicationsBasePath}
}

// Forbid denies access to a document
func Forbid(document *doc.Source) (*StepResult, error) {
	target := actionForbid.target(document.Path())
	return nil, &ActionError{http.StatusForbidden, target + " is forbidden"}
}

func (action *ActionForbid) target(path string) string {
	return action.basePath + path
}
