package main

import (
	"testing"

	"ibfd.org/d-way/rule"

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
		Test{"/document/", "OTHER"},
		Test{"/document/vatst_de", "LINKRESOLVER"},
		Test{"/linkresolver/static", "OTHER"},
		Test{"/linkresolver/static/", "OTHER"},
		Test{"/linkresolver/static/tt_al-de_50_ger_2015_tt__td1", "LINKRESOLVER"},
		Test{"/linkresolver/static/some/extra/path/tt_al-de_50_ger_2015_tt__td1", "LINKRESOLVER"},

		Test{"/data/tns/docs/html/tns_1988-01-15_au_1.html", "HTML"},
		Test{"/archive/ggd/html/ggd_2016_ar_1.html", "HTML"},
		Test{"/data/tns/html/.html", "HTML"},
		Test{"/data//html/.html", "HTML"},
		Test{"/archive/ggd/html/ggd_2016_ar_1.htm", "OTHER"},
		Test{"/archive/knorhaan/banana/html/kh_299.html", "HTML"},

		Test{"/collections/eulaw/pdf/doc1.pdf", "PDF1"},
		Test{"/collections/ftn/pdf/doc1.pdf", "PDF1"},
		Test{"/collections/oecd/pdf/doc1.pdf", "PDF1"},
		Test{"/collections/tni/pdf/doc1.pdf", "PDF1"},
		Test{"/collections/vatst/pdf/doc1.pdf", "PDF2"},
		Test{"/collections/eufake/pdf/doc1.pdf", "PDF2"},

		Test{"/data/tns/docs/pdf/tns_1988-01-15_au_1.pdf", "PDF3"},
		Test{"/data/treaty/docs/pdf/tt_al-de_50_ger_2015_tt__td1.pdf", "PDF3"},
		Test{"/data/kf/docs/pdf/kf_nl.pdf", "OTHER"},
		Test{"/collections/vatst/docs/printversion/pdf/vatst_nl.pdf", "PDF4"},
		Test{"/data/tns/docs/printversion/pdf/tns_1988-01-15_au_1.pdf", "PDF5"},
		Test{"/data/tns/docs/printversion/pdf/some/extra/path/here/tns_1997-02-20_nl_1.pdf", "PDF5"},
		Test{"/banana/tns/docs/printversion/pdf/tns_2019-01-25_ar_1.pdf", "OTHER"},
		Test{"/data/tns/docs/printversion/pdf/tns_2019-01-25_ar_1.pdf", "PDF5"},

		Test{"/collections/kf/excel/kf_dz.xls", "EXCEL"},
		Test{"/collections//excel/kf_dz.xls", "EXCEL"}}
}

func TestRules(t *testing.T) {
	matcher := cfg.GetMatcher()
	for _, test := range tests {
		rule := matcher.Match(test.path, "")
		assert(t, rule, test.ruleName, test.path)
	}
}

func assert(t *testing.T, rule *rule.Rule, expectedName, what string) {
	if rule.Name == expectedName {
		t.Logf("\tFor path %s the rule should be \"%s\" %s", what, rule.Name, checkMark)
	} else {
		t.Fatalf("\tFor path %s the rule should be \"%s\" but it was \"%s\" %s", what, expectedName, rule.Name, ballotX)
	}
}
