package arff

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("quote",
	func(s, exp string) {
		Expect(quote(s)).To(BeIdenticalTo(exp))
	},

	Entry("plain", "plain", "plain"),
	Entry("spaces", `with space`, "'with space'"),
	Entry("question mark", `?`, "'?'"),
	Entry("comments", `with % comment`, "'with % comment'"),
	Entry("brackets", `with{a}`, "'with{a}'"),
	Entry("quotes", `with "quoted"`, `'with "quoted"'`),
	Entry("single quotes", `with 'quoted'`, `'with \'quoted\''`),
	Entry("tabs", "with\ttab", "'with\\ttab'"),
	Entry("new lines", "line\r\nbreak", "'line\\r\\nbreak'"),
	Entry("whitespace", ` with whitespace`, `'with whitespace'`),

	Entry("backslashes", "back\\slash", "back\\slash"),
	Entry("backslashes escaped", "with back\\slash", "'with back\\\\slash'"),
	Entry("unicode", "日本", "日本"),
)

var _ = DescribeTable("unquote",
	func(s, exp string) {
		Expect(unquote(s)).To(BeIdenticalTo(exp))
	},

	Entry("blank", "", ""),
	Entry("single-char", "x", "x"),
	Entry("one-sided", "'one quote", "'one quote"),

	Entry("plain", "plain", "plain"),
	Entry("quoted", "'with space'", `with space`),
	Entry("question mark", "'?'", `?`),

	Entry("quotes", `'with "quoted"'`, `with "quoted"`),
	Entry("single quotes", `'with \'quoted\''`, `with 'quoted'`),
	Entry("tabs", "'with\\ttab'", "with\ttab"),
	Entry("new lines", "'line\\r\\nbreak'", "line\r\nbreak"),
	Entry("whitespace", `'with whitespace '`, `with whitespace`),

	Entry("backslashes", "back\\slash", "back\\slash"),
	Entry("backslashes trailing", "'back\\\\'", "back\\"),
	Entry("backslashes escaped", "'with back\\\\slash'", "with back\\slash"),
	Entry("unicode", "日本", "日本"),
)
