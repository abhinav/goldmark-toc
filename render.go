package toc

import (
	"github.com/yuin/goldmark/ast"
	"go.abhg.dev/goldmark/toc"
)

// RenderList renders a table of contents as a nested list with a sane,
// default configuration for the ListRenderer.
func RenderList(t *TOC) ast.Node {
	return toc.RenderList(t)
}

// ListRenderer builds a nested list from a table of contents.
//
// For example,
//
//	# Foo
//	## Bar
//	## Baz
//	# Qux
//
// Becomes,
//
//   - Foo
//   - Bar
//   - Baz
//   - Qux
type ListRenderer = toc.ListRenderer
