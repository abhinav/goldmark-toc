package toc

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const (
	_defaultTitle = "Table of Contents"
)

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
type Transformer struct {
	// Title is the title of the table of contents section.
	// Defaults to "Table of Contents" if unspecified.
	// The title is rendered as a <p> element.
	Title string

	// MinDepth is the minimum depth of the table of contents.
	// See the documentation for MinDepth for more information.
	MinDepth int

	// MaxDepth is the maximum depth of the table of contents.
	// See the documentation for MaxDepth for more information.
	MaxDepth int

	// ListID is the id for the list of TOC items rendered in the HTML.
	//
	// For example, if ListID is "toc", the table of contents will be
	// rendered as:
	//
	//	<ul id="toc">
	//	  ...
	//	</ul>
	//
	// The HTML element does not have an ID if ListID is empty.
	ListID string

	// TitleID is the id for the Title element rendered in the HTML.
	//
	// For example, if TitleID is "toc-title",
	// the title will be rendered as:
	//
	//	<p id="toc-title">Table of Contents</p>
	//
	// If TitleID is empty, a value will be requested
	// from the Goldmark Parser.
	TitleID string

	// Compact controls whether empty items should be removed
	// from the table of contents.
	// See the documentation for Compact for more information.
	Compact bool

	// HideTitle controls whether the title is rendered.
	// When set to true, the title (e.g., <p>Table of Contents</p>) is not rendered,
	// and only the TOC list is output.
	//
	// When HideTitle is true and ContainerElement is set,
	// an aria-label attribute with the title text is added to the container
	// for accessibility.
	//
	// Defaults to false (title is shown).
	HideTitle bool

	// ContainerElement specifies the HTML element to wrap the TOC in.
	// Common values are "nav", "div", "aside", etc.
	//
	// For example, if ContainerElement is "nav", the table of contents
	// will be rendered as:
	//
	//	<nav>
	//	  <p>Table of Contents</p>
	//	  <ul>...</ul>
	//	</nav>
	//
	// If ContainerElement is empty, no wrapper element is added.
	ContainerElement string

	// ContainerClass specifies the CSS class(es) for the container element.
	// This is only used when ContainerElement is set.
	//
	// For example, if ContainerElement is "nav" and ContainerClass is "toc-nav",
	// the table of contents will be rendered as:
	//
	//	<nav class="toc-nav">
	//	  ...
	//	</nav>
	//
	// Multiple classes can be specified separated by spaces.
	ContainerClass string

	// ContainerID specifies the ID for the container element.
	// This is only used when ContainerElement is set.
	//
	// For example, if ContainerElement is "nav" and ContainerID is "table-of-contents",
	// the table of contents will be rendered as:
	//
	//	<nav id="table-of-contents">
	//	  ...
	//	</nav>
	ContainerID string
}

var _ parser.ASTTransformer = (*Transformer)(nil) // interface compliance

// Transform adds a table of contents to the provided Markdown document.
//
// Errors encountered while transforming are ignored. For more fine-grained
// control, use Inspect and transform the document manually.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, ctx parser.Context) {
	toc, err := Inspect(doc, reader.Source(), MinDepth(t.MinDepth), MaxDepth(t.MaxDepth), Compact(t.Compact))
	if err != nil {
		// There are currently no scenarios under which Inspect
		// returns an error but we have to account for it anyway.
		return
	}

	// Don't add anything for documents with no headings.
	if len(toc.Items) == 0 {
		return
	}

	listNode := RenderList(toc)
	if id := t.ListID; len(id) > 0 {
		listNode.SetAttributeString("id", []byte(id))
	}

	// Determine the title text
	title := t.Title
	if len(title) == 0 {
		title = _defaultTitle
	}

	// Build the title paragraph if not hidden
	var titleNode *ast.Paragraph
	if !t.HideTitle {
		titleBytes := []byte(title)
		titleNode = ast.NewParagraph()
		titleNode.AppendChild(titleNode, ast.NewString(titleBytes))
		if id := t.TitleID; len(id) > 0 {
			titleNode.SetAttributeString("id", []byte(id))
		} else if ids := ctx.IDs(); ids != nil {
			id := ids.Generate(titleBytes, titleNode.Kind())
			titleNode.SetAttributeString("id", id)
		}
	}

	// If container element is specified, wrap everything in it
	if t.ContainerElement != "" {
		container := &containerNode{
			element: t.ContainerElement,
			class:   t.ContainerClass,
			id:      t.ContainerID,
		}
		// Add aria-label when title is hidden for accessibility
		if t.HideTitle {
			container.ariaLabel = title
		}
		if titleNode != nil {
			container.AppendChild(container, titleNode)
		}
		container.AppendChild(container, listNode)
		doc.InsertBefore(doc, doc.FirstChild(), container)
	} else {
		// Insert without container (original behavior)
		doc.InsertBefore(doc, doc.FirstChild(), listNode)
		if titleNode != nil {
			doc.InsertBefore(doc, doc.FirstChild(), titleNode)
		}
	}
}

// containerNode is a custom AST node that renders as a specified HTML element.
type containerNode struct {
	ast.BaseBlock
	element   string
	class     string
	id        string
	ariaLabel string
}

// KindContainerNode is the kind of the container node.
var KindContainerNode = ast.NewNodeKind("TOCContainer")

// Kind returns the kind of the container node.
func (n *containerNode) Kind() ast.NodeKind {
	return KindContainerNode
}

// Dump dumps the container node to the given writer.
func (n *containerNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
