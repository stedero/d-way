package rule

// Rule defines the mapping of document names to processes
type Rule struct {
	Regex   string   `json:"regex"`
	Process []string `json:"process"`
}

type engine struct {
	Rules []Rule `json:"rules"`
}

func match(path string) *Rule {
	return &Rule{}
}
