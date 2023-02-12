package toc

import "go.abhg.dev/goldmark/toc"

// Transformer is a Goldmark AST transformer adds a TOC to the top of a
// Markdown document.
//
// To use this, either install the Extender on the goldmark.Markdown object,
// or install the AST transformer on the Markdown parser like so.
//
//	markdown := goldmark.New(...)
//	markdown.Parser().AddOptions(
//	  parser.WithAutoHeadingID(),
//	  parser.WithASTTransformers(
//	    util.Prioritized(&toc.Transformer{}, 100),
//	  ),
//	)
//
// NOTE: Unless you've supplied your own parser.IDs implementation, you'll
// need to enable the WithAutoHeadingID option on the parser to generate IDs
// and links for headings.
type Transformer = toc.Transformer
