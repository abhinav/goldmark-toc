package toc

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"pgregory.net/rapid"
)

func item(title, id string, items ...*Item) *Item {
	n := new(Item)
	if len(title) > 0 {
		n.Title = []byte(title)
	}
	if len(id) > 0 {
		n.ID = []byte(id)
	}
	for _, item := range items {
		n.Items = append(n.Items, item)
	}
	return n
}

func TestInspect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give []string // lines of a doc
		opts []InspectOption
		want Items
	}{
		{
			desc: "empty",
			give: nil,
		},
		{
			desc: "single level",
			give: []string{
				"# Foo",
				"# Bar",
				"# Baz",
			},
			want: Items{
				item("Foo", "foo"),
				item("Bar", "bar"),
				item("Baz", "baz"),
			},
		},
		{
			desc: "subitems",
			give: []string{
				"# Foo",
				"## Bar",
				"## Baz",
			},
			want: Items{
				item("Foo", "foo",
					item("Bar", "bar"),
					item("Baz", "baz"),
				),
			},
		},
		{
			desc: "decrease level",
			give: []string{
				"# Foo",
				"## Bar",
				"# Baz",
				"# Qux",
			},
			want: Items{
				item("Foo", "foo",
					item("Bar", "bar"),
				),
				item("Baz", "baz"),
				item("Qux", "qux"),
			},
		},
		{
			desc: "alternating levels",
			give: []string{
				"# Foo",
				"## Bar",
				"# Baz",
				"## Qux",
				"# Quux",
			},
			want: Items{
				item("Foo", "foo",
					item("Bar", "bar"),
				),
				item("Baz", "baz",
					item("Qux", "qux"),
				),
				item("Quux", "quux"),
			},
		},
		{
			desc: "several levels offset",
			give: []string{
				"# A",
				"###### B",
				"### C",
				"##### D",
				"## E",
				"# F",
				"# G",
			},
			// Levels:
			// 	1	2	3	4	5	6
			want: Items{
				item("A", "a",
					item("", "",
						item("", "",
							item("", "",
								item("", "",
									item("B", "b"),
								),
							),
						),
						item("C", "c",
							item("", "",
								item("D", "d"),
							),
						),
					),
					item("E", "e"),
				),
				item("F", "f"),
				item("G", "g"),
			},
		},
		{
			desc: "escaped punctuation in title",
			give: []string{
				`# Foo\-Bar`,
				`## Bar\-Baz`,
			},
			want: Items{
				item("Foo-Bar", "foo-bar",
					item("Bar-Baz", "bar-baz"),
				),
			},
		},
		{
			desc: "depth",
			give: []string{
				"# A",
				"###### B",
				"### C",
				"##### D",
				"## E",
				"# F",
				"# G",
			},
			opts: []InspectOption{MaxDepth(3)},
			want: Items{
				item("A", "a",
					item("", "",
						item("C", "c")),
					item("E", "e")),
				item("F", "f"),
				item("G", "g"),
			},
		},
		{
			desc: "compact",
			give: []string{
				"# A",
				"### B",
				"#### C",
				"# D",
				"#### E",
			},
			opts: []InspectOption{Compact(true)},
			want: Items{
				item("A", "a",
					item("B", "b",
						item("C", "c"),
					),
				),
				item("D", "d",
					item("E", "e"),
				),
			},
		},
		{
			desc: "compact complex",
			give: []string{
				"## A",
				"##### B",
				"###### C",
				"## D",
				"# E",
				"### F",
				"# G",
				"#### H",
				"### I",
				"## J",
			},
			opts: []InspectOption{Compact(true)},
			want: Items{
				item("A", "a",
					item("B", "b",
						item("C", "c"),
					),
				),
				item("D", "d"),
				item("E", "e",
					item("F", "f"),
				),
				item("G", "g",
					item("H", "h"),
					item("I", "i"),
					item("J", "j"),
				),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			src := []byte(strings.Join(tt.give, "\n") + "\n")
			doc := parser.NewParser(
				parser.WithInlineParsers(parser.DefaultInlineParsers()...),
				parser.WithBlockParsers(parser.DefaultBlockParsers()...),
				parser.WithAutoHeadingID(),
			).Parse(text.NewReader(src))

			got, err := Inspect(doc, src, tt.opts...)
			require.NoError(t, err, "inspect error")
			assert.Equal(t, &TOC{Items: tt.want}, got)
		})
	}
}

func TestInspectOption_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give InspectOption
		want string
	}{
		{give: MaxDepth(3), want: "MaxDepth(3)"},
		{give: MaxDepth(0), want: "MaxDepth(0)"},
		{give: MaxDepth(-1), want: "MaxDepth(-1)"},
		{give: Compact(true), want: "Compact(true)"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, fmt.Sprint(tt.give))
		})
	}
}

func TestInspectRandomHeadings(t *testing.T) {
	t.Parallel()

	rapid.Check(t, testInspectRandomHeadings)
}

func FuzzInspectRandomHeadings(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(testInspectRandomHeadings))
}

func testInspectRandomHeadings(t *rapid.T) {
	// Generate a random hierarchy.
	levels := rapid.SliceOf(rapid.IntRange(1, 6)).Draw(t, "levels")
	var buf bytes.Buffer
	for i, level := range levels {
		buf.WriteString(strings.Repeat("#", level))
		buf.WriteString(" Heading ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteByte('\n')
	}

	src := buf.Bytes()
	doc := parser.NewParser(
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithAutoHeadingID(),
	).Parse(text.NewReader(src))

	toc, err := Inspect(doc, src)
	require.NoError(t, err, "inspect error")

	// Verify that the number of items in the TOC is the same as the number
	// of headings in the document.
	assert.Equal(t, len(levels), nonEmptyItems(toc.Items),
		"number of non-empty items in TOC "+
			"does not match number of headings in document:\n%s", src)
}

func TestInspectCompactRandomHeadings(t *testing.T) {
	t.Parallel()

	rapid.Check(t, testInspectCompactRandomHeadings)
}

func FuzzInspectCompactRandomHeadings(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(testInspectCompactRandomHeadings))
}

func testInspectCompactRandomHeadings(t *rapid.T) {
	// Generate a random hierarchy.
	levels := rapid.SliceOf(rapid.IntRange(1, 6)).Draw(t, "levels")
	var buf bytes.Buffer
	for i, level := range levels {
		buf.WriteString(strings.Repeat("#", level))
		buf.WriteString(" Heading ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteByte('\n')
	}

	src := buf.Bytes()
	doc := parser.NewParser(
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithAutoHeadingID(),
	).Parse(text.NewReader(src))

	toc, err := Inspect(doc, src, Compact(true))
	require.NoError(t, err, "inspect error")

	// There must be no empty items in the TOC.
	assert.Equal(t, nonEmptyItems(toc.Items), totalItems(toc.Items),
		"number of non-empty items in TOC "+
			"does not match number of items in TOC:\n%s", src)
	assert.Equal(t, len(levels), totalItems(toc.Items),
		"number of items in TOC "+
			"does not match number of headings in document:\n%s", src)
}

func totalItems(items Items) (total int) {
	for _, item := range items {
		total++
		total += totalItems(item.Items)
	}
	return total
}

func nonEmptyItems(items Items) (total int) {
	for _, item := range items {
		if len(item.Title) > 0 {
			total++
		}
		total += nonEmptyItems(item.Items)
	}
	return total
}
