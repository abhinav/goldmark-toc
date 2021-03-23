package toc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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

			got, err := Inspect(doc, src)
			require.NoError(t, err, "inspect error")
			assert.Equal(t, &TOC{Items: tt.want}, got)
		})
	}
}
