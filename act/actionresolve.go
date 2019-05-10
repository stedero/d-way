package act

import (
	"fmt"
	"net/http"
	"strings"

	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
)

// ActionResolve defines the action that resolves a UID to a document path.
type ActionResolve struct {
	url    string
	client *http.Client
}

var actionResolve *ActionResolve

func init() {
	config := cfg.GetMatcher()
	actionResolve = &ActionResolve{config.ResolveURL, NewHTTPClient()}
}

// Resolve calls the Linkresolver service to resolve a UID to a document path.
func Resolve(src *doc.Source) (*StepResult, error) {
	req, err := http.NewRequest("GET", actionResolve.url, nil)
	uid := uid(src)
	setUserAgent(req)
	q := req.URL.Query()
	q.Add("uid", uid)
	req.URL.RawQuery = q.Encode()
	resp, err := actionResolve.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &ActionError{resp.StatusCode, fmt.Sprintf("failed to resolve UID: %s", uid)}
	}
	return NewActionResult().SetResponse(resp).SetStatusCode(http.StatusFound), err
}

func uid(src *doc.Source) string {
	parts := strings.Split(src.String(), "/")
	return strings.Split(parts[len(parts)-1], "#")[0]
}
