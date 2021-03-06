package rule

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Rule defines a process rule.
type Rule struct {
	Name     string `json:"name"`
	Regex    string `json:"regex"`
	Regexc   *regexp.Regexp
	MimeType string   `json:"mimeType"`
	Steps    []string `json:"steps"`
}

// LogConfig defines the logging configuration.
type LogConfig struct {
	Filename string `json:"filename"`
	Level    string `json:"level"`
}

// Matcher defines the rules to match paths.
type Matcher struct {
	Comment              string     `json:"comment"`
	PublicationsBasePath string     `json:"publications_base_path"`
	CleanURL             string     `json:"clean_url"`
	OtccURL              string     `json:"otcc_url"`
	ResolveURL           string     `json:"resolve_url"`
	SdrmURL              string     `json:"sdrm_url"`
	XtojURL              string     `json:"xtoj_url"`
	CacheMaxAgeSeconds   int        `json:"cache_max_age_seconds"`
	Logging              *LogConfig `json:"logging"`
	Rules                []*Rule    `json:"rules"`
}

// NewMatcher creates a Matcher.
func NewMatcher(data []byte) (*Matcher, error) {
	var matcher Matcher
	err := json.Unmarshal(data, &matcher)
	if err != nil {
		return nil, err
	}
	matcher.compileRules()
	return &matcher, err
}

// Match finds the first rule that matches.
// For mimeType to match a rule's mimetype it has to contain the
// exact string as defined in the rule. No wildcard matching
// is done.
func (matcher *Matcher) Match(path, mimeType string) *Rule {
	for _, rule := range matcher.Rules {
		if rule.Regexc.Match([]byte(path)) {
			if rule.MimeType == "" || strings.Contains(mimeType, rule.MimeType) {
				return rule
			}
		}
	}
	return matcher.Rules[len(matcher.Rules)-1]
}

func (matcher *Matcher) compileRules() {
	for _, rule := range matcher.Rules {
		rule.Regexc = regexp.MustCompile(rule.Regex)
	}
}
