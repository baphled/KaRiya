package shared_test

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventDetailCard", func() {
	var (
		testTheme themes.Theme
		testEvent *career.Event
	)

	BeforeEach(func() {
		testTheme = themes.NewDefaultTheme()
		testEvent = fixtures.EventWith("test-123", "Implemented new feature for the product", "Acme Corp", "Project Alpha")
		testEvent.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		testEvent.Tags = []string{"go", "backend", "api"}
		testEvent.Categories = []string{"development", "feature"}
		testEvent.Skills = []string{"skill-1", "skill-2"}
	})

	Describe("Construction", func() {
		It("creates a new EventDetailCard", func() {
			card := shared.NewEventDetailCard(testEvent, testTheme)
			Expect(card).NotTo(BeNil())
		})

		It("creates card with nil theme", func() {
			card := shared.NewEventDetailCard(testEvent, nil)
			Expect(card).NotTo(BeNil())
			// Should still render without panic
			result := card.Render()
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("Rendering", func() {
		It("renders 'No event selected' for nil event", func() {
			card := shared.NewEventDetailCard(nil, testTheme)
			result := card.Render()
			Expect(result).To(ContainSubstring("No event selected"))
		})

		It("renders the event title", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Event Details"))
		})

		It("renders the date", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Date:"))
			Expect(result).To(ContainSubstring("2024-01-15"))
		})

		It("renders the company when present", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Company:"))
			Expect(result).To(ContainSubstring("Acme Corp"))
		})

		It("renders the project when present", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Project:"))
			Expect(result).To(ContainSubstring("Project Alpha"))
		})

		It("renders the event text", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Text:"))
			Expect(result).To(ContainSubstring("Implemented new feature"))
		})

		It("renders tags as comma-separated list", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Tags:"))
			Expect(result).To(ContainSubstring("go, backend, api"))
			// Should NOT have Go array syntax
			Expect(result).NotTo(ContainSubstring("[go backend api]"))
		})

		It("renders categories as comma-separated list", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Categories:"))
			Expect(result).To(ContainSubstring("development, feature"))
		})

		It("renders skill count", func() {
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Skills:"))
			Expect(result).To(ContainSubstring("2 associated"))
		})
	})

	Describe("Optional Fields", func() {
		It("omits company when empty", func() {
			testEvent.Company = ""
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).NotTo(ContainSubstring("Company:"))
		})

		It("omits project when empty", func() {
			testEvent.Project = ""
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).NotTo(ContainSubstring("Project:"))
		})

		It("omits tags when empty", func() {
			testEvent.Tags = nil
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).NotTo(ContainSubstring("Tags:"))
		})

		It("omits categories when empty", func() {
			testEvent.Categories = nil
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).NotTo(ContainSubstring("Categories:"))
		})

		It("omits skills when empty", func() {
			testEvent.Skills = nil
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).NotTo(ContainSubstring("Skills:"))
		})
	})

	Describe("Long Text Handling", func() {
		It("renders long event text", func() {
			longText := "This is a very long event description that contains multiple sentences. " +
				"It describes a complex project with many technical details. " +
				"The implementation involved multiple teams and technologies. " +
				"We achieved significant results and learned valuable lessons."

			testEvent.Text = longText
			result := shared.RenderEventDetailCard(testEvent, testTheme)

			// The text should be present
			Expect(result).To(ContainSubstring("This is a very long event description"))
			Expect(result).To(ContainSubstring("achieved significant results"))
		})

		It("handles multi-line event text", func() {
			multiLineText := "First paragraph of the event.\n\nSecond paragraph with more details.\n\nThird paragraph with conclusions."

			testEvent.Text = multiLineText
			result := shared.RenderEventDetailCard(testEvent, testTheme)

			Expect(result).To(ContainSubstring("First paragraph"))
			Expect(result).To(ContainSubstring("Second paragraph"))
			Expect(result).To(ContainSubstring("Third paragraph"))
		})

		It("handles long tag lists", func() {
			testEvent.Tags = []string{
				"go", "python", "javascript", "typescript", "rust",
				"kubernetes", "docker", "aws", "gcp", "terraform",
			}

			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Tags:"))
			// Should have all tags
			Expect(result).To(ContainSubstring("go"))
			Expect(result).To(ContainSubstring("terraform"))
		})

		It("handles long category lists", func() {
			testEvent.Categories = []string{
				"development", "architecture", "leadership",
				"mentoring", "process-improvement", "cost-optimization",
			}

			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("Categories:"))
			Expect(result).To(ContainSubstring("development"))
			Expect(result).To(ContainSubstring("cost-optimization"))
		})
	})

	Describe("Special Characters", func() {
		It("handles unicode in event text", func() {
			testEvent.Text = "Implemented internationalization: 日本語, 中文, العربية"
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("日本語"))
			Expect(result).To(ContainSubstring("中文"))
		})

		It("handles emojis in event text", func() {
			testEvent.Text = "Launched new feature 🚀 with great success ✅"
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("🚀"))
			Expect(result).To(ContainSubstring("✅"))
		})

		It("handles special characters in company name", func() {
			testEvent.Company = "O'Reilly & Associates"
			result := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(result).To(ContainSubstring("O'Reilly & Associates"))
		})
	})

	Describe("Helper Function", func() {
		It("RenderEventDetailCard produces same output as Render()", func() {
			card := shared.NewEventDetailCard(testEvent, testTheme)
			expected := card.Render()
			actual := shared.RenderEventDetailCard(testEvent, testTheme)
			Expect(actual).To(Equal(expected))
		})
	})
})
