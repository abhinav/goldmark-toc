package toc

import (
	"fmt"
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

	type testCase struct {
		desc      string
		giveDepth int
		wantDepth int
	}

	tests := []testCase{
		{
			desc:      "default",
			wantDepth: _defaultTitleDepth,
		},
		{
			desc:      "< 1",
			giveDepth: -1,
			wantDepth: 1,
		},
		{
			desc:      "> 6",
			giveDepth: 7,
			wantDepth: 6,
		},
		{
			desc:      "absurd",
			giveDepth: 130931,
			wantDepth: 6,
		},
	}

	for i := _defaultTitleDepth; i <= _maxTitleDepth; i++ {
		tests = append(tests, testCase{
			desc:      fmt.Sprintf("valid/%d", i),
			giveDepth: i,
			wantDepth: i,
		})
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
						TitleDepth: tt.giveDepth,
					}, 100),
				),
			).Parse(text.NewReader(src))

			// Should definitely still be a heading
			heading, ok := doc.FirstChild().(*ast.Heading)

			require.True(t, ok, "first child must be a heading, got %T", doc.FirstChild())
			assert.Equal(t, tt.wantDepth, heading.Level, "level mismatch")
		})
	}
}
