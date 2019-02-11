package main

import (
	"testing"

	"ibfd.org/d-way/cfg"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

type Test struct {
	path     string
	ruleName string
}

var tests []Test

func init() {
	tests = []Test{
		Test{"/banana", "OTHER"},
		Test{"/data/tns/docs/printversion/pdf/tns_2019-01-25_ar_1.pdf", "PDF"}}
}

func TestRules(t *testing.T) {
	matcher := cfg.GetMatcher()
	for _, test := range tests {
		rule := matcher.Match(test.path, "")
		assert(t, test.ruleName, rule.Name, test.path)
	}
}

func assert(t *testing.T, expected string, actual string, what string) {
	t.Logf("\tWhen checking path %s", what)
	if actual == expected {
		t.Logf("\t\tThe rule should be \"%s\" %s", actual, checkMark)
	} else {
		t.Fatalf("\t\tThe rule should be \"%s\" but it was \"%s\" %s", expected, actual, ballotX)
	}
}
