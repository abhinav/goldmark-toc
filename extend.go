package toc

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Extender extends a Goldmark Markdown parser and renderer to always include
// a table of contents in the output.
//
// To use this, install it into your Goldmark Markdown object.
//
//	md := goldmark.New(
//	  // ...
//	  goldmark.WithParserOptions(parser.WithAutoHeadingID()),
//	  goldmark.WithExtensions(
//	    // ...
//	    &toc.Extender{
//	    },
//	  ),
//	)
//
// This will install the default Transformer. For more control, install the
// Transformer directly on the Markdown Parser.
//
// NOTE: Unless you've supplied your own parser.IDs implementation, you'll
// need to enable the WithAutoHeadingID option on the parser to generate IDs
// and links for headings.
type Extender struct {
	// Title is the title of the table of contents section.
	// Defaults to "Table of Contents" if unspecified.
	Title string

	// MinDepth is the minimum depth of the table of contents.
	// Headings with a level lower than the specified depth will be ignored.
	// See the documentation for MinDepth for more information.
	//
	// Defaults to 0 (no limit) if unspecified.
	MinDepth int

	// MaxDepth is the maximum depth of the table of contents.
	// Headings with a level greater than the specified depth will be ignored.
	// See the documentation for MaxDepth for more information.
	//
	// Defaults to 0 (no limit) if unspecified.
	MaxDepth int

	// ListID is the id for the list of TOC items rendered in the HTML.
	//
	// See the documentation for Transformer.ListID for more information.
	ListID string

	// Compact controls whether empty items should be removed
	// from the table of contents.
	//
	// See the documentation for Compact for more information.
	Compact bool
}

// Extend adds support for rendering a table of contents to the provided
// Markdown parser/renderer.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				Title:    e.Title,
				MinDepth: e.MinDepth,
				MaxDepth: e.MaxDepth,
				ListID:   e.ListID,
				Compact:  e.Compact,
			}, 100),
		),
	)
}
