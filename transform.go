package toc

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const _defaultTitle = "Table of Contents"

// Transformer is a Goldmark AST transformer adds a TOC to the top of a
// Markdown document.
//
// To use this, either install the Extender on the goldmark.Markdown object,
// or install the AST transformer on the Markdown parser like so.
//
//   markdown := goldmark.New(...)
//   markdown.Parser().AddOptions(
//     parser.WithAutoHeadingID(),
//     parser.WithASTTransformers(
//       util.Prioritized(&toc.Transformer{}, 100),
//     ),
//   )
//
// NOTE: Unless you've supplied your own parser.IDs implementation, you'll
// need to enable the WithAutoHeadingID option on the parser to generate IDs
// and links for headings.
type Transformer struct {
	// Title is the title of the table of contents section.
	// Defaults to "Table of Contents" if unspecified.
	Title string
}

var _ parser.ASTTransformer = (*Transformer)(nil) // interface compliance

// Transform adds a table of contents to the provided Markdown document.
//
// Errors encountered while transforming are ignored. For more fine-grained
// control, use Inspect and transform the document manually.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	toc, err := Inspect(doc, reader.Source())
	if err != nil {
		// There are currently no scenarios under which Inspect
		// returns an error but we have to account for it anyway.
		return
	}

	// Don't add anything for documents with no headings.
	if len(toc.Items) == 0 {
		return
	}

	doc.InsertBefore(doc, doc.FirstChild(), RenderList(toc))

	title := t.Title
	if len(title) == 0 {
		title = _defaultTitle
	}

	heading := ast.NewHeading(1)
	heading.AppendChild(heading, ast.NewString([]byte(title)))

	doc.InsertBefore(doc, doc.FirstChild(), heading)
}
