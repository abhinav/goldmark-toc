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

// From: https://github.com/abhinav/goldmark-toc/issues/61
func TestTransformerWithTitleDepth(t *testing.T) {
	t.Parallel()

	src := []byte(strings.Join([]string{
		"# Hey",
		"## Now",
		"# Then",
		"### There",
		"## Now",
	}, "\n") + "\n")

	tests := []struct {
		desc           string
		giveTitleDepth int
		wantTitleDepth int
	}{
		{
			desc:           "default title depth",
			wantTitleDepth: _defaultTitleDepth,
		},
		{
			desc:           "title depth 0",
			giveTitleDepth: 0,
			wantTitleDepth: 1,
		},
		{
			desc:           "title depth 1",
			giveTitleDepth: 1,
			wantTitleDepth: 1,
		},
		{
			desc:           "title depth 2",
			giveTitleDepth: 2,
			wantTitleDepth: 2,
		},
		{
			desc:           "title depth 3",
			giveTitleDepth: 3,
			wantTitleDepth: 3,
		},
		{
			desc:           "title depth 4",
			giveTitleDepth: 4,
			wantTitleDepth: 4,
		},
		{
			desc:           "title depth 5",
			giveTitleDepth: 5,
			wantTitleDepth: 5,
		},
		{
			desc:           "title depth 6",
			giveTitleDepth: 6,
			wantTitleDepth: 6,
		},
		{
			desc:           "title depth 7",
			giveTitleDepth: 7,
			wantTitleDepth: 6,
		},
		{
			desc:           "title depth 255",
			giveTitleDepth: 255,
			wantTitleDepth: 6,
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
						TitleDepth: tt.giveTitleDepth,
					}, 100),
				),
			).Parse(text.NewReader(src))

			// Should definitely still be a heading
			heading, ok := doc.FirstChild().(*ast.Heading)

			require.True(t, ok, "first child must be a heading, got %T", doc.FirstChild())
			assert.Equal(t, tt.wantTitleDepth, heading.Level, "level mismatch")
		})
	}
}
