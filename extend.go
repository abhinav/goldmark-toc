package toc

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
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
	// The title is rendered as a <p> element.
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

	// TitleID is the id for the Title element rendered in the HTML.
	//
	// See the documentation for Transformer.TitleID for more information.
	TitleID string

	// Compact controls whether empty items should be removed
	// from the table of contents.
	//
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

// Extend adds support for rendering a table of contents to the provided
// Markdown parser/renderer.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				Title:            e.Title,
				TitleID:          e.TitleID,
				ListID:           e.ListID,
				MinDepth:         e.MinDepth,
				MaxDepth:         e.MaxDepth,
				Compact:          e.Compact,
				HideTitle:        e.HideTitle,
				ContainerElement: e.ContainerElement,
				ContainerClass:   e.ContainerClass,
				ContainerID:      e.ContainerID,
			}, 100),
		),
	)

	// Always register the container node renderer.
	// If ContainerElement is empty, the renderer simply won't be used,
	// but it needs to be registered in case Transformer uses it.
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newContainerNodeRenderer(), 100),
		),
	)
}

// containerNodeRenderer renders containerNode as HTML.
type containerNodeRenderer struct {
	html.Config
}

// newContainerNodeRenderer returns a new renderer for containerNode.
func newContainerNodeRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &containerNodeRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers the render functions for containerNode.
func (r *containerNodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindContainerNode, r.renderContainer)
}

func (r *containerNodeRenderer) renderContainer(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*containerNode)
	if entering {
		_, _ = w.WriteString("<")
		_, _ = w.WriteString(n.element)
		if n.id != "" {
			_, _ = w.WriteString(` id="`)
			_, _ = w.WriteString(n.id)
			_, _ = w.WriteString(`"`)
		}
		if n.class != "" {
			_, _ = w.WriteString(` class="`)
			_, _ = w.WriteString(n.class)
			_, _ = w.WriteString(`"`)
		}
		if n.ariaLabel != "" {
			_, _ = w.WriteString(` aria-label="`)
			_, _ = w.WriteString(n.ariaLabel)
			_, _ = w.WriteString(`"`)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(n.element)
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}
