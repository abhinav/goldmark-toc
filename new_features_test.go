package toc_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/toc"
)

func TestHideTitle(t *testing.T) {
	t.Parallel()

	src := []byte(`
# Section 1
## Subsection 1.1
# Section 2
`)

	tests := []struct {
		name      string
		hideTitle bool
		wantTitle bool
	}{
		{
			name:      "with title (default)",
			hideTitle: false,
			wantTitle: true,
		},
		{
			name:      "without title",
			hideTitle: true,
			wantTitle: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			md := goldmark.New(
				goldmark.WithParserOptions(parser.WithAutoHeadingID()),
				goldmark.WithExtensions(&toc.Extender{
					HideTitle: tt.hideTitle,
				}),
			)

			var buf bytes.Buffer
			if err := md.Convert(src, &buf); err != nil {
				t.Fatalf("Convert failed: %v", err)
			}

			output := buf.String()
			hasTitle := strings.Contains(output, "Table of Contents")

			if hasTitle != tt.wantTitle {
				t.Errorf("HideTitle=%v: got title present=%v, want %v\nOutput:\n%s",
					tt.hideTitle, hasTitle, tt.wantTitle, output)
			}
		})
	}
}

func TestContainerElement(t *testing.T) {
	t.Parallel()

	src := []byte(`
# Section 1
## Subsection 1.1
`)

	tests := []struct {
		name             string
		containerElement string
		containerClass   string
		containerID      string
		wantContains     []string
		wantNotContains  []string
	}{
		{
			name:             "no container",
			containerElement: "",
			wantNotContains:  []string{"<nav", "</nav>", "<div", "</div>"},
		},
		{
			name:             "nav container",
			containerElement: "nav",
			wantContains:     []string{"<nav>", "</nav>"},
		},
		{
			name:             "nav container with class",
			containerElement: "nav",
			containerClass:   "toc-nav",
			wantContains:     []string{`<nav class="toc-nav">`, "</nav>"},
		},
		{
			name:             "nav container with id",
			containerElement: "nav",
			containerID:      "table-of-contents",
			wantContains:     []string{`<nav id="table-of-contents">`, "</nav>"},
		},
		{
			name:             "nav container with class and id",
			containerElement: "nav",
			containerClass:   "toc-nav",
			containerID:      "table-of-contents",
			wantContains:     []string{`id="table-of-contents"`, `class="toc-nav"`, "</nav>"},
		},
		{
			name:             "div container",
			containerElement: "div",
			containerClass:   "toc-wrapper",
			wantContains:     []string{`<div class="toc-wrapper">`, "</div>"},
		},
		{
			name:             "aside container",
			containerElement: "aside",
			wantContains:     []string{"<aside>", "</aside>"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			md := goldmark.New(
				goldmark.WithParserOptions(parser.WithAutoHeadingID()),
				goldmark.WithExtensions(&toc.Extender{
					ContainerElement: tt.containerElement,
					ContainerClass:   tt.containerClass,
					ContainerID:      tt.containerID,
				}),
			)

			var buf bytes.Buffer
			if err := md.Convert(src, &buf); err != nil {
				t.Fatalf("Convert failed: %v", err)
			}

			output := buf.String()

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output should contain %q but doesn't.\nOutput:\n%s", want, output)
				}
			}

			for _, notWant := range tt.wantNotContains {
				if strings.Contains(output, notWant) {
					t.Errorf("Output should not contain %q but does.\nOutput:\n%s", notWant, output)
				}
			}
		})
	}
}

func TestHideTitleWithContainer(t *testing.T) {
	t.Parallel()

	src := []byte(`
# Section 1
## Subsection 1.1
`)

	md := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(&toc.Extender{
			HideTitle:        true,
			ContainerElement: "nav",
			ContainerClass:   "toc",
		}),
	)

	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	output := buf.String()

	// Should have nav container
	if !strings.Contains(output, `<nav class="toc">`) {
		t.Errorf("Output should contain nav container.\nOutput:\n%s", output)
	}

	// Should NOT have "Table of Contents" title
	if strings.Contains(output, "Table of Contents") {
		t.Errorf("Output should not contain title when HideTitle=true.\nOutput:\n%s", output)
	}

	// Should still have the TOC list
	if !strings.Contains(output, "<ul>") {
		t.Errorf("Output should still contain the TOC list.\nOutput:\n%s", output)
	}
}

func TestContainerWithCustomTitle(t *testing.T) {
	t.Parallel()

	src := []byte(`
# Section 1
## Subsection 1.1
`)

	md := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(&toc.Extender{
			Title:            "Contents",
			TitleDepth:       2,
			ContainerElement: "nav",
			ContainerID:      "toc",
		}),
	)

	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	output := buf.String()

	// Should have nav container with ID
	if !strings.Contains(output, `<nav id="toc">`) {
		t.Errorf("Output should contain nav container with id.\nOutput:\n%s", output)
	}

	// Should have h2 title (TitleDepth=2)
	if !strings.Contains(output, "<h2") {
		t.Errorf("Output should contain h2 heading for title.\nOutput:\n%s", output)
	}

	// Should have custom title text
	if !strings.Contains(output, "Contents") {
		t.Errorf("Output should contain custom title 'Contents'.\nOutput:\n%s", output)
	}
}
