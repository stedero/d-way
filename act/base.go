package act

import (
	"net/http"
	"time"
)

const defaultTimeoutSeconds = 10

// NewHTTPClient creates a HTPP client with a default timeout.
func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: time.Second * defaultTimeoutSeconds}
}
