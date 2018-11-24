package act

import (
	"net/http"
	"os"
	"time"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionStat defines the action that fetches a document.
type ActionStat struct {
	basePath string
}

var actionStat *ActionStat

func init() {
	config := cfg.GetMatcher()
	actionStat = &ActionStat{config.PublicationsBasePath}
}

// Stat checks a file date.
func Stat(document *doc.Source, reqModSince *time.Time) (*StepResult, error) {
	path := document.Path()
	if path == "" {
		return nil, &ActionError{http.StatusBadRequest, "no file specified"}
	}
	target := actionStat.target(document.Path())
	file, err := os.Open(target)
	defer file.Close()
	if err != nil {
		return nil, &ActionError{http.StatusNotFound, err.Error()}
	}
	finfo, err := file.Stat()
	modTimeStr := finfo.ModTime().Format(cfg.LastModifiedDateFormat)
	modTime, _ := time.Parse(cfg.LastModifiedDateFormat, modTimeStr)
	statusCode := statusCode(modTime, reqModSince)
	return NewStepResult().SetContent(modTimeStr).SetStatusCode(statusCode), err
}

func (action *ActionStat) target(path string) string {
	return action.basePath + path
}

func statusCode(fileMod time.Time, reqMod *time.Time) int {
	if reqMod == nil || fileMod.After(*reqMod) {
		return http.StatusOK
	}
	return http.StatusNotModified
}
