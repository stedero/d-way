package act

import (
	"net/http"

	"ibfd.org/d-way/cfg"
	log "ibfd.org/d-way/log4u"
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
func Resolve(uid string, cookies []*http.Cookie) (*StepResult, error) {
	log.Debugf("Resolve: %s?%s\n", actionResolve.url, uid)
	req, err := http.NewRequest("GET", actionResolve.url, nil)
	setUserAgent(req)
	q := req.URL.Query()
	q.Add("uid", uid)
	req.URL.RawQuery = q.Encode()
	copyCookies(req, cookies)
	resp, err := actionResolve.client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Debugf("Resolve result status: %d", resp.StatusCode)
	return NewStepResult().SetResponse(resp), err
}
