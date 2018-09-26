package rule

import (
	"encoding/json"
	"regexp"
)

// Rule defines a process rule.
type Rule struct {
	Regex  string `json:"regex"`
	Regexc *regexp.Regexp
	Steps  []string `json:"steps"`
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
	ResolveURL           string     `json:"resolve_url"`
	SdrmURL              string     `json:"sdrm_url"`
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

// Match finds the first rule that matches a path.
func (matcher *Matcher) Match(path string) *Rule {
	for _, rule := range matcher.Rules {
		if rule.Regexc.Match([]byte(path)) {
			return rule
		}
	}
	return matcher.Rules[len(matcher.Rules)-1]
}

func (matcher *Matcher) compileRules() {
	for _, rule := range matcher.Rules {
		rule.Regexc = regexp.MustCompile(rule.Regex)
	}
}
