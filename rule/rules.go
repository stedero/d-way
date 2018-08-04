package rule

import (
	"encoding/json"
	"regexp"

	"ibfd.org/d-way/doc"
)

// Rule defines process rules
type Rule struct {
	Regex     string `json:"regex"`
	Regexc    *regexp.Regexp
	Processes []string `json:"process"`
}

// Matcher defines the rules to match paths
type Matcher struct {
	Comment string  `json:"comment"`
	Rules   []*Rule `json:"rules"`
}

// NewMatcher creates a Matcher
func NewMatcher(data []byte) (*Matcher, error) {
	var matcher Matcher
	err := json.Unmarshal(data, &matcher)
	if err != nil {
		return nil, err
	}
	matcher.compileRules()
	return &matcher, err
}

func (matcher *Matcher) compileRules() {
	for _, rule := range matcher.Rules {
		rule.Regexc = regexp.MustCompile(rule.Regex)
	}
}

// Match finds a rule that matches path
func (matcher *Matcher) Match(d *doc.Document) *Rule {
	for _, rule := range matcher.Rules {
		if rule.Regexc.Match([]byte(d.Path)) {
			return rule
		}
	}
	return matcher.Rules[len(matcher.Rules)-1]
}
