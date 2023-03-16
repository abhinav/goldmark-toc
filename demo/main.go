// demo implements a WASM module that can be used to format markdown
// with the goldmark-toc extension.
package main

import (
	"bytes"
	"syscall/js"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/toc"
)

func main() {
	js.Global().Set("formatMarkdown", js.FuncOf(func(this js.Value, args []js.Value) any {
		var req request
		req.Decode(args[0])
		return formatMarkdown(req)
	}))
	select {}
}

type request struct {
	Markdown string
	Title    string
	MaxDepth int
	Compact  bool
}

func (r *request) Decode(v js.Value) {
	r.Markdown = v.Get("markdown").String()
	r.Title = v.Get("title").String()
	r.MaxDepth = v.Get("maxDepth").Int()
	r.Compact = v.Get("compact").Bool()
}

func formatMarkdown(req request) string {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&toc.Extender{
				Title:    req.Title,
				MaxDepth: req.MaxDepth,
				Compact:  req.Compact,
			},
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(req.Markdown), &buf); err != nil {
		return err.Error()
	}
	return buf.String()
}
