package toc_test

import (
	"testing"

	toc "github.com/abhinav/goldmark-toc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	goldtestutil "github.com/yuin/goldmark/testutil"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	md := goldmark.New(
		goldmark.WithExtensions(&toc.Extender{}),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	goldtestutil.DoTestCaseFile(md, "testdata/tests.txt", t)
}
