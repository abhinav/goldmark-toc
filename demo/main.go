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
	js.Global().Set("formatMarkdown", js.FuncOf(formatMarkdown))
	select {}
}

func formatMarkdown(this js.Value, args []js.Value) any {
	input := args[0].String()
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&toc.Extender{},
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return err.Error()
	}
	return buf.String()
}
