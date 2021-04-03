[![Go Reference](https://pkg.go.dev/badge/github.com/abhinav/goldmark-toc.svg)](https://pkg.go.dev/github.com/abhinav/goldmark-toc)
[![Go](https://github.com/abhinav/goldmark-toc/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-toc/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-toc/branch/main/graph/badge.svg?token=OLXTVHEIOG)](https://codecov.io/gh/abhinav/goldmark-toc)

goldmark-toc is an add-on for the [goldmark] Markdown parser that adds support
for rendering a table-of-contents.

  [goldmark]: http://github.com/yuin/goldmark

# Usage

To use goldmark-toc, import the `toc` package.

```go
import toc "github.com/abhinav/goldmark-toc"
```

Following that, you have three options for using this package:

- [Extension][]: This is the easiest way to get a table of contents into your
  document and provides very little control over the output.
- [Transformer][]: This is the next easiest option and provides more control
  over the output.
- [Manual][]: This option requires the most work but also provides the most
  control.

  [Extension]: #extension
  [Transformer]: #transformer
  [Manual]: #manual

## Extension

To use this package as a simple Goldmark extension, install the `Extender`
when constructing the `goldmark.Markdown` object.

```go
markdown := goldmark.New(
    // ...
    goldmark.WithParserOptions(parser.WithAutoHeadingID()),
    goldmark.WithExtensions(
        // ...
        &toc.Extender{},
    ),
)
```

This will add a "Table of Contents" section to the top of every Markdown
document parsed by this Markdown object.

> NOTE: The example above enables `parser.WithAutoHeadingID`. Without this or
> a custom implementation of `parser.IDs`, none of the headings in the
> document will have links generated for them.

## Transformer

Installing this package as an AST Transformer provides slightly more control
over the output. To use it, install the AST transformer on the Goldmark
Markdown parser.

```go
markdown := goldmark.New(...)
markdown.Parser().AddOptions(
    parser.WithAutoHeadingID(),
    parser.WithASTTransformers(
        util.Prioritized(&toc.Transformer{
            Title: "Contents",
        }, 100),
    ),
)
```

This will generate a "Contents" section at the top of all Markdown documents
parsed by this parser.

As with the previous example, this enables `parser.WithAutoHeadingID` to get
auto-generated heading IDs.

## Manual

If you use this package manually to generate Tables of Contents, you have a
lot more control over the behavior. This requires a few steps.

### Parse Markdown

Parse a Markdown document with goldmark.

```go
markdown := goldmark.New(...)
markdown.Parser().AddOptions(parser.WithAutoHeadingID())
doc := markdown.Parser().Parse(text.NewReader(src))
```

Note that the parser must be configured to generate IDs for headers or the
headers in the table of contents won't have anything to point to. This can be
accomplished by adding the `parser.WithAutoHeadingID` option as in the example
above, or with a custom implementation of [`goldmark/parser.IDs`] by using the
snippet below.

  [`goldmark/parser.IDs`]: https://pkg.go.dev/github.com/yuin/goldmark/parser#IDs

```go
markdown := goldmark.New(...)
pctx := parser.NewContext(parser.WithIDs(ids))
doc := parser.Parse(text.NewReader(src), parser.WithContext(pctx))
```

### Build a table of contents

After parsing a Markdown document, inspect it with `toc`.

```go
tree, err := toc.Inspect(doc, src)
if err != nil {
  // handle the error
}
```

### Generate a Markdown list

You can render the table of contents into a Markdown list with
`toc.RenderList`.

```go
list := toc.RenderList(tree)
```

This builds a list representation of the table of contents to be rendered as
Markdown or HTML.

You may manipulate the `tree` before rendering the list.

### Render HTML

Finally, render this table of contents along with your Markdown document:

```go
markdown.Renderer().Render(output, src, list) // table of contents
markdown.Renderer().Render(output, src, doc)  // document
```

Alternatively, include the table of contents into your Markdown document in
your desired position and render it using your Markdown renderer.

```go
// Prepend table of contents to the front of the document.
doc.InsertBefore(doc, doc.FirstChild(), list)

// Render the document.
markdown.Renderer().Render(output, src, doc)
```
