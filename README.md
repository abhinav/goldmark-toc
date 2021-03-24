[![Go Reference](https://pkg.go.dev/badge/github.com/abhinav/goldmark-toc.svg)](https://pkg.go.dev/github.com/abhinav/goldmark-toc)
[![Go](https://github.com/abhinav/goldmark-toc/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-toc/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-toc/branch/main/graph/badge.svg?token=OLXTVHEIOG)](https://codecov.io/gh/abhinav/goldmark-toc)

goldmark-toc is an add-on for the [goldmark] Markdown parser that adds support
for rendering a table-of-contents.

  [goldmark]: http://github.com/yuin/goldmark

Note that this is *not* an extension to the library: it doesn't alter the
parsing or rendering process. Instead, it inspects a parsed document and
provides a means of generating a TOC.

# Usage

To use goldmark-toc, import the `toc` package.

```go
import toc "github.com/abhinav/goldmark-toc"
```

## Parse Markdown

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

## Build a table of contents

After parsing a Markdown document, inspect it with `toc`.

```go
tree, err := toc.Inspect(doc, src)
if err != nil {
  // handle the error
}
```

## Generate a Markdown list

You can render the table of contents into a Markdown list with
`toc.RenderList`.

```go
list := toc.RenderList(tree)
```

This builds a list representation of the table of contents to be rendered as
Markdown or HTML.

## Render HTML

Finally, render this table of contents along with your Markdown document:

```go
markdown.Renderer().Render(output, src, list) // table of contents
markdown.Renderer().Render(output, src, doc)  // document
```
