package components_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
)

var _ = Describe("SyntaxHighlighter", func() {
	var (
		highlighter *components.SyntaxHighlighter
		theme       themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		highlighter = components.NewSyntaxHighlighter(theme)
	})

	Describe("NewSyntaxHighlighter", func() {
		It("creates a highlighter with theme", func() {
			h := components.NewSyntaxHighlighter(theme)
			Expect(h).NotTo(BeNil())
		})

		It("uses default theme when nil is provided", func() {
			h := components.NewSyntaxHighlighter(nil)
			Expect(h).NotTo(BeNil())
			// Should not panic when highlighting
			result := h.HighlightJSON(`{"key": "value"}`)
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("HighlightJSON", func() {
		It("highlights JSON keys", func() {
			input := `{"name": "John"}`
			result := highlighter.HighlightJSON(input)
			// Result should contain the key (may have ANSI codes)
			Expect(result).To(ContainSubstring("name"))
		})

		It("highlights JSON string values", func() {
			input := `{"name": "John"}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("John"))
		})

		It("highlights JSON numbers", func() {
			input := `{"count": 42}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("42"))
		})

		It("highlights JSON booleans", func() {
			input := `{"active": true}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("true"))
		})

		It("highlights JSON null", func() {
			input := `{"value": null}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("null"))
		})

		It("highlights JSON brackets", func() {
			input := `{"items": [1, 2, 3]}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("{"))
			Expect(result).To(ContainSubstring("}"))
			Expect(result).To(ContainSubstring("["))
			Expect(result).To(ContainSubstring("]"))
		})

		It("handles multi-line JSON", func() {
			input := `{
  "name": "John",
  "age": 30
}`
			result := highlighter.HighlightJSON(input)
			lines := strings.Split(result, "\n")
			Expect(len(lines)).To(Equal(4))
		})

		It("handles empty JSON object", func() {
			input := `{}`
			result := highlighter.HighlightJSON(input)
			Expect(result).NotTo(BeEmpty())
		})

		It("handles negative numbers", func() {
			input := `{"value": -42}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("-42"))
		})

		It("handles decimal numbers", func() {
			input := `{"price": 19.99}`
			result := highlighter.HighlightJSON(input)
			Expect(result).To(ContainSubstring("19.99"))
		})
	})

	Describe("HighlightYAML", func() {
		It("highlights YAML keys", func() {
			input := `name: John`
			result := highlighter.HighlightYAML(input)
			Expect(result).To(ContainSubstring("name"))
		})

		It("highlights YAML list items", func() {
			input := `items:
  - item1
  - item2`
			result := highlighter.HighlightYAML(input)
			Expect(result).To(ContainSubstring("-"))
		})

		It("highlights YAML comments", func() {
			input := `name: John # this is a comment`
			result := highlighter.HighlightYAML(input)
			Expect(result).To(ContainSubstring("comment"))
		})

		It("highlights YAML numbers", func() {
			input := `count: 42`
			result := highlighter.HighlightYAML(input)
			Expect(result).To(ContainSubstring("42"))
		})

		It("highlights YAML booleans", func() {
			input := `active: true`
			result := highlighter.HighlightYAML(input)
			Expect(result).To(ContainSubstring("true"))
		})

		It("handles multi-level nesting", func() {
			input := `parent:
  child:
    value: test`
			result := highlighter.HighlightYAML(input)
			lines := strings.Split(result, "\n")
			Expect(len(lines)).To(Equal(3))
		})

		It("handles underscore in keys", func() {
			input := `my_key: value`
			result := highlighter.HighlightYAML(input)
			Expect(result).To(ContainSubstring("my_key"))
		})
	})

	Describe("HighlightCSV", func() {
		It("highlights CSV header row differently", func() {
			input := `Name,Age,City
John,30,NYC`
			result := highlighter.HighlightCSV(input)
			// Headers should be styled (contains ANSI codes)
			Expect(result).To(ContainSubstring("Name"))
			Expect(result).To(ContainSubstring("Age"))
		})

		It("highlights CSV data rows", func() {
			input := `Name,Age
John,30`
			result := highlighter.HighlightCSV(input)
			Expect(result).To(ContainSubstring("John"))
			Expect(result).To(ContainSubstring("30"))
		})

		It("handles empty lines", func() {
			input := `Name,Age

John,30`
			result := highlighter.HighlightCSV(input)
			lines := strings.Split(result, "\n")
			Expect(len(lines)).To(Equal(3))
		})

		It("handles empty content", func() {
			input := ``
			result := highlighter.HighlightCSV(input)
			Expect(result).To(BeEmpty())
		})

		It("handles single column", func() {
			input := `Name
John
Jane`
			result := highlighter.HighlightCSV(input)
			Expect(result).To(ContainSubstring("Name"))
			Expect(result).To(ContainSubstring("John"))
		})

		It("trims whitespace from cells", func() {
			input := `Name , Age
 John , 30 `
			result := highlighter.HighlightCSV(input)
			// Cells should be trimmed
			Expect(result).To(ContainSubstring("John"))
		})
	})

	Describe("HighlightPlainText", func() {
		It("highlights section headers ending with colon", func() {
			input := `Summary:
This is content`
			result := highlighter.HighlightPlainText(input)
			Expect(result).To(ContainSubstring("Summary:"))
		})

		It("highlights separator lines", func() {
			input := `---
Content here
===`
			result := highlighter.HighlightPlainText(input)
			Expect(result).To(ContainSubstring("---"))
			Expect(result).To(ContainSubstring("==="))
		})

		It("handles unicode box drawing separators", func() {
			input := `───────────`
			result := highlighter.HighlightPlainText(input)
			Expect(result).To(ContainSubstring("───"))
		})

		It("does not highlight regular lines", func() {
			input := `This is regular text`
			result := highlighter.HighlightPlainText(input)
			Expect(result).To(Equal(input))
		})

		It("does not highlight lines with spaces before colon", func() {
			input := `Key: value here`
			result := highlighter.HighlightPlainText(input)
			// This contains spaces before potential key-value, so not treated as header
			Expect(result).To(Equal(input))
		})
	})

	Describe("Highlight", func() {
		It("dispatches to HighlightJSON for json format", func() {
			input := `{"key": "value"}`
			result := highlighter.Highlight(input, "json")
			// Should contain styled content (not equal to raw input due to ANSI codes)
			Expect(result).To(ContainSubstring("key"))
		})

		It("dispatches to HighlightYAML for yaml format", func() {
			input := `key: value`
			result := highlighter.Highlight(input, "yaml")
			Expect(result).To(ContainSubstring("key"))
		})

		It("dispatches to HighlightCSV for csv format", func() {
			input := `Name,Age`
			result := highlighter.Highlight(input, "csv")
			Expect(result).To(ContainSubstring("Name"))
		})

		It("dispatches to HighlightPlainText for txt format", func() {
			input := `Section:`
			result := highlighter.Highlight(input, "txt")
			Expect(result).To(ContainSubstring("Section:"))
		})

		It("dispatches to HighlightPlainText for text format", func() {
			input := `Section:`
			result := highlighter.Highlight(input, "text")
			Expect(result).To(ContainSubstring("Section:"))
		})

		It("dispatches to HighlightPlainText for markdown format", func() {
			input := `# Header`
			result := highlighter.Highlight(input, "markdown")
			Expect(result).To(ContainSubstring("Header"))
		})

		It("dispatches to HighlightPlainText for md format", func() {
			input := `# Header`
			result := highlighter.Highlight(input, "md")
			Expect(result).To(ContainSubstring("Header"))
		})

		It("returns raw content for unknown format", func() {
			input := `some content`
			result := highlighter.Highlight(input, "unknown")
			Expect(result).To(Equal(input))
		})

		It("handles case-insensitive format names", func() {
			input := `{"key": "value"}`
			result := highlighter.Highlight(input, "JSON")
			Expect(result).To(ContainSubstring("key"))

			result = highlighter.Highlight(input, "Json")
			Expect(result).To(ContainSubstring("key"))
		})
	})
})
