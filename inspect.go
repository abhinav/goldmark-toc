package toc

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// InspectOption customizes the behavior of Inspect.
type InspectOption interface {
	apply(*inspectOptions)
}

type inspectOptions struct {
	maxDepth int
}

// MaxDepth limits the depth of the table of contents.
// Headings with a level greater than the specified depth will be ignored.
//
// For example, given the following:
//
//	# Foo
//	## Bar
//	### Baz
//	# Quux
//	## Qux
//
// MaxDepth(1) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Foo", ID: "foo"}
//	 |
//	 +--- &Item{Title: "Quux", ID: "quux", Items: ...}
//
// Whereas, MaxDepth(2) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Foo", ID: "foo", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Bar", ID: "bar"}
//	 |
//	 +--- &Item{Title: "Quux", ID: "quux", Items: ...}
//	       |
//	       +--- &Item{Title: "Qux", ID: "qux"}
//
// A value of 0 or less will result in no limit.
//
// The default is no limit.
func MaxDepth(depth int) InspectOption {
	return maxDepthOption(depth)
}

type maxDepthOption int

func (d maxDepthOption) apply(opts *inspectOptions) {
	opts.maxDepth = int(d)
}

func (d maxDepthOption) String() string {
	return fmt.Sprintf("MaxDepth(%d)", int(d))
}

// Inspect builds a table of contents by inspecting the provided document.
//
// The table of contents is represents as a tree where each item represents a
// heading or a heading level with zero or more children.
// The returned TOC will be empty if there are no headings in the document.
//
// For example,
//
//	# Section 1
//	## Subsection 1.1
//	## Subsection 1.2
//	# Section 2
//	## Subsection 2.1
//	# Section 3
//
// Will result in the following items.
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Section 1", ID: "section-1", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Subsection 1.1", ID: "subsection-1-1"}
//	 |     |
//	 |     +--- &Item{Title: "Subsection 1.2", ID: "subsection-1-2"}
//	 |
//	 +--- &Item{Title: "Section 2", ID: "section-2", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Subsection 2.1", ID: "subsection-2-1"}
//	 |
//	 +--- &Item{Title: "Section 3", ID: "section-3"}
//
// You may analyze or manipulate the table of contents before rendering it.
func Inspect(n ast.Node, src []byte, options ...InspectOption) (*TOC, error) {
	var opts inspectOptions
	for _, opt := range options {
		opt.apply(&opts)
	}

	// Appends an empty subitem to the given node
	// and returns a reference to it.
	appendChild := func(n *Item) *Item {
		child := new(Item)
		n.Items = append(n.Items, child)
		return child
	}

	// Returns the last subitem of the given node,
	// creating it if necessary.
	lastChild := func(n *Item) *Item {
		if len(n.Items) > 0 {
			return n.Items[len(n.Items)-1]
		}
		return appendChild(n)
	}

	var root Item

	stack := []*Item{&root}
	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		if opts.maxDepth > 0 && heading.Level > opts.maxDepth {
			return ast.WalkSkipChildren, nil
		}

		for len(stack) < heading.Level {
			parent := stack[len(stack)-1]
			stack = append(stack, lastChild(parent))
		}

		for len(stack) > heading.Level {
			stack = stack[:len(stack)-1]
		}

		parent := stack[len(stack)-1]
		target := lastChild(parent)
		if len(target.Title) > 0 || len(target.Items) > 0 {
			target = appendChild(parent)
		}

		target.Title = util.UnescapePunctuations(heading.Text(src))
		if id, ok := n.AttributeString("id"); ok {
			target.ID, _ = id.([]byte)
		}

		return ast.WalkSkipChildren, nil
	})

	return &TOC{Items: root.Items}, err
}
