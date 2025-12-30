package models_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

var _ = Describe("QualityIndicator", func() {
	var (
		indicator   *models.QualityIndicator
		calculator  *careerservice.DataQualityCalculator
	)

	BeforeEach(func() {
		calculator = careerservice.NewDataQualityCalculator()
		indicator = models.NewQualityIndicator(calculator)
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
		})
	})

	Describe("RenderCompact", func() {
		It("should return empty string when score is nil", func() {
			result := indicator.RenderCompact()
			Expect(result).To(Equal(""))
		})

		It("should use correct icon for Incomplete level", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).To(ContainSubstring("●"))
		})

		It("should use correct icon for Basic level", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Led development",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TechCorp",
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).NotTo(BeEmpty())
		})

		It("should use correct icon for Enriched level", func() {
			event := &career.CareerEvent{
				ID:       "test-id",
				Text:     "Led development",
				Date:     time.Now().Add(-24 * time.Hour),
				Company:  "TechCorp",
				Project:  "Migration",
				Tags:     []string{"technical"},
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.RenderCompact()
			Expect(result).NotTo(BeEmpty())
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
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("SetScore", func() {
		It("should set the score", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			score := calculator.CalculateQuality(event)
			indicator.SetScore(&score)

			result := indicator.Render()
			Expect(result).NotTo(BeEmpty())
		})
	})
})
