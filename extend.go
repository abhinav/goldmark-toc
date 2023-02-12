package toc

import "go.abhg.dev/goldmark/toc"

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
type Extender = toc.Extender
