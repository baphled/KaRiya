package themes_test

import (
	"github.com/baphled/kariya/internal/cli/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Glamour Theme Integration", func() {
	var theme themes.Theme

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("NewGlamourStyle", func() {
		It("should create a glamour style name based on theme", func() {
			styleName := themes.NewGlamourStyleName(theme)
			Expect(styleName).NotTo(BeEmpty())
		})

		It("should return dark style for dark themes", func() {
			styleName := themes.NewGlamourStyleName(theme)
			// Default theme is dark
			Expect(styleName).To(Equal("dark"))
		})

		It("should return light style for light themes", func() {
			// Create a light theme
			lightPalette := &themes.ColorPalette{
				Background: "#ffffff",
				Foreground: "#000000",
			}
			lightTheme := themes.NewBaseTheme("light", "Light Theme", "Test", false, lightPalette)
			styleName := themes.NewGlamourStyleName(lightTheme)
			Expect(styleName).To(Equal("light"))
		})

		It("should handle nil theme", func() {
			styleName := themes.NewGlamourStyleName(nil)
			Expect(styleName).To(Equal("dark"))
		})
	})

	Describe("RenderMarkdown", func() {
		It("should render markdown content", func() {
			markdown := "# Hello World\n\nThis is a test."
			result, err := themes.RenderMarkdown(theme, markdown, 80)
			Expect(err).NotTo(HaveOccurred())
			// Content is rendered with ANSI styling, so just check it's not empty
			// and contains the text (possibly with ANSI codes interspersed)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(MatchRegexp("Hello.*World"))
		})

		It("should render bullet lists", func() {
			markdown := "- Item 1\n- Item 2\n- Item 3"
			result, err := themes.RenderMarkdown(theme, markdown, 80)
			Expect(err).NotTo(HaveOccurred())
			// Check content is present (may have ANSI codes)
			Expect(result).To(MatchRegexp("Item.*1"))
			Expect(result).To(MatchRegexp("Item.*2"))
		})

		It("should render code blocks", func() {
			markdown := "```go\nfunc main() {}\n```"
			result, err := themes.RenderMarkdown(theme, markdown, 80)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("func"))
		})

		It("should handle empty content", func() {
			result, err := themes.RenderMarkdown(theme, "", 80)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should handle nil theme", func() {
			result, err := themes.RenderMarkdown(nil, "# Test", 80)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Test"))
		})

		It("should respect width parameter", func() {
			longText := "This is a very long line that should wrap according to the width parameter specified."
			result, err := themes.RenderMarkdown(theme, longText, 40)
			Expect(err).NotTo(HaveOccurred())
			// The result should contain the text (wrapped)
			Expect(result).To(ContainSubstring("long"))
		})
	})

	Describe("RenderCVPreview", func() {
		It("should render CV content as formatted markdown", func() {
			cvContent := `# John Doe
## Senior Software Engineer

### Experience
- Led development of microservices architecture
- Mentored junior developers

### Skills
- Go, Python, JavaScript
- Docker, Kubernetes
`
			result, err := themes.RenderCVPreview(theme, cvContent, 80)
			Expect(err).NotTo(HaveOccurred())
			// Content is rendered with ANSI styling
			Expect(result).To(MatchRegexp("John.*Doe"))
			Expect(result).To(MatchRegexp("Experience"))
			Expect(result).To(MatchRegexp("Skills"))
		})

		It("should handle empty CV content", func() {
			result, err := themes.RenderCVPreview(theme, "", 80)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})
})
