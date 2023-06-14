package toc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestTransformer(t *testing.T) {
	t.Parallel()

	src := []byte(strings.Join([]string{
		"# Foo",
		"## Bar",
		"# Baz",
		"### Qux",
		"## Quux",
	}, "\n") + "\n")

	tests := []struct {
		desc      string
		giveTitle string
		wantTitle string
	}{
		{
			desc:      "default title",
			wantTitle: _defaultTitle,
		},
		{
			desc:      "custom title",
			giveTitle: "Contents",
			wantTitle: "Contents",
		},
	}

	for _, tt := range tests {
		tt := tt // for t.Parallel
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			doc := parser.NewParser(
				parser.WithInlineParsers(parser.DefaultInlineParsers()...),
				parser.WithBlockParsers(parser.DefaultBlockParsers()...),
				parser.WithAutoHeadingID(),
				parser.WithASTTransformers(
					util.Prioritized(&Transformer{
						Title: tt.giveTitle,
					}, 100),
				),
			).Parse(text.NewReader(src))

			heading, ok := doc.FirstChild().(*ast.Heading)
			require.True(t, ok, "first child must be a heading, got %T", doc.FirstChild())
			assert.Equal(t, tt.wantTitle, string(heading.Text(src)), "title mismatch")
		})
	}
}
