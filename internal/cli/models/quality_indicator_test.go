package models

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

func TestQualityIndicator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Quality Indicator Suite")
}

var _ = Describe("QualityIndicator", func() {
	var (
		indicator   *QualityIndicator
		calculator  *careerservice.DataQualityCalculator
	)

	BeforeEach(func() {
		calculator = careerservice.NewDataQualityCalculator()
		indicator = NewQualityIndicator(calculator)
	})

	Describe("Render", func() {
		It("should return empty string when score is nil", func() {
			result := indicator.Render()
			Expect(result).To(Equal(""))
		})

		It("should render score bar for incomplete event", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.Render()
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("40%"))
		})

		It("should render score bar for complete event", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Led development of backend system",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "TechCorp",
				Project:    "Migration",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical"},
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.Render()
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("100%"))
		})

		It("should include quality level in output", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.Render()
			Expect(result).To(ContainSubstring("Basic"))
		})

		It("should include suggestion for missing fields", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.Render()
			Expect(result).To(ContainSubstring("Add"))
		})
	})

	Describe("RenderCompact", func() {
		It("should return empty string when score is nil", func() {
			result := indicator.RenderCompact()
			Expect(result).To(Equal(""))
		})

		It("should render single line for incomplete event", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Basic"))
			Expect(result).To(ContainSubstring("40%"))
		})

		It("should render single line for complete event", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Led development of backend system",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "TechCorp",
				Project:    "Migration",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical"},
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Complete"))
			Expect(result).To(ContainSubstring("100%"))
		})
	})

	Describe("RenderWithDetails", func() {
		It("should return empty string when score is nil", func() {
			result := indicator.RenderWithDetails()
			Expect(result).To(Equal(""))
		})

		It("should include field scores for event with all fields", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Led development of backend system",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "TechCorp",
				Project:    "Migration",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical"},
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderWithDetails()
			Expect(result).To(ContainSubstring("Text:"))
			Expect(result).To(ContainSubstring("Date:"))
			Expect(result).To(ContainSubstring("Company:"))
			Expect(result).To(ContainSubstring("Project:"))
			Expect(result).To(ContainSubstring("Tags:"))
			Expect(result).To(ContainSubstring("Categories:"))
		})

		It("should not include zero-score fields", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderWithDetails()
			Expect(result).To(ContainSubstring("Text:"))
			Expect(result).To(ContainSubstring("Date:"))
			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})
	})

	Describe("SetWidth", func() {
		It("should set width for rendering", func() {
			indicator.SetWidth(80)
			Expect(indicator.width).To(Equal(80))
		})
	})

	Describe("Quality Level Icons", func() {
		It("should use correct icon for Incomplete level", func() {
			event := &career.CareerEvent{
				ID: "test-id",
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).To(ContainSubstring("◯"))
		})

		It("should use correct icon for Basic level", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).To(ContainSubstring("◑"))
		})

		It("should use correct icon for Enriched level", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TechCorp",
				Project: "Migration",
				Tags:    []string{"technical"},
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).To(ContainSubstring("◐"))
		})

		It("should use correct icon for Complete level", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Led development of backend system",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "TechCorp",
				Project:    "Migration",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical"},
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).To(ContainSubstring("✓"))
		})
	})
})

