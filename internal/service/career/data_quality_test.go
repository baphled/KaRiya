package career

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

func TestDataQuality(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Quality Suite")
}

var _ = Describe("DataQualityCalculator", func() {
	var (
		calculator *DataQualityCalculator
	)

	BeforeEach(func() {
		calculator = NewDataQualityCalculator()
	})

	Describe("CalculateQuality", func() {
		It("should score 0 for event with only ID and timestamps", func() {
			event := &career.CareerEvent{
				ID:        "test-id",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.Score).To(Equal(0))
			Expect(score.Level).To(Equal(QualityIncomplete))
			Expect(score.TextScore).To(Equal(0))
			Expect(score.DateScore).To(Equal(0))
		})

		It("should award 20 points for valid text", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Led team on critical project",
				Date: time.Now().Add(-24 * time.Hour),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TextScore).To(Equal(20))
		})

		It("should award 0 points for empty text", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "",
				Date: time.Now().Add(-24 * time.Hour),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TextScore).To(Equal(0))
		})

		It("should award 0 points for whitespace-only text", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "   \n  \t  ",
				Date: time.Now().Add(-24 * time.Hour),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TextScore).To(Equal(0))
		})

		It("should award 20 points for valid date", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event text",
				Date: time.Now().Add(-24 * time.Hour),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.DateScore).To(Equal(20))
		})

		It("should award 0 points for zero date", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event text",
				Date: time.Time{},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.DateScore).To(Equal(0))
		})

		It("should award 10 points for company", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event text",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TechCorp",
			}

			score := calculator.CalculateQuality(event)

			Expect(score.CompanyScore).To(Equal(10))
		})

		It("should award 0 points for empty company", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event text",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "",
			}

			score := calculator.CalculateQuality(event)

			Expect(score.CompanyScore).To(Equal(0))
		})

		It("should award 0 points for whitespace company", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event text",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "   ",
			}

			score := calculator.CalculateQuality(event)

			Expect(score.CompanyScore).To(Equal(0))
		})

		It("should award 10 points for project", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event text",
				Date:    time.Now().Add(-24 * time.Hour),
				Project: "Migration",
			}

			score := calculator.CalculateQuality(event)

			Expect(score.ProjectScore).To(Equal(10))
		})

		It("should award 0 points for empty project", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event text",
				Date:    time.Now().Add(-24 * time.Hour),
				Project: "",
			}

			score := calculator.CalculateQuality(event)

			Expect(score.ProjectScore).To(Equal(0))
		})

		It("should award 15 points for 2+ tags", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event text",
				Date: time.Now().Add(-24 * time.Hour),
				Tags: []string{"technical", "leadership"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TagsScore).To(Equal(15))
		})

		It("should award 7 points for 1 tag", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event text",
				Date: time.Now().Add(-24 * time.Hour),
				Tags: []string{"technical"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TagsScore).To(Equal(7))
		})

		It("should award 0 points for no tags", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event text",
				Date: time.Now().Add(-24 * time.Hour),
				Tags: []string{},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TagsScore).To(Equal(0))
		})

		It("should award 15 points for categories", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Event text",
				Date:       time.Now().Add(-24 * time.Hour),
				Categories: []string{"technical"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.CategoriesScore).To(Equal(15))
		})

		It("should award 0 points for no categories", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Event text",
				Date:       time.Now().Add(-24 * time.Hour),
				Categories: []string{},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.CategoriesScore).To(Equal(0))
		})

		It("should award 10 points for matching text and categories", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Led development of backend system",
				Date:       time.Now().Add(-24 * time.Hour),
				Categories: []string{"technical"},
				Tags:       []string{"technical"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.MatchScore).To(Equal(10))
		})

		It("should award 5 points for non-matching text and categories", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Did something",
				Date:       time.Now().Add(-24 * time.Hour),
				Categories: []string{"technical"},
				Tags:       []string{"technical"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.MatchScore).To(Equal(5))
		})

		It("should calculate correct total score for complete event", func() {
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

			// 20 (text) + 20 (date) + 10 (company) + 10 (project) + 15 (tags) + 15 (categories) + 10 (match)
			expectedScore := 100
			Expect(score.Score).To(Equal(expectedScore))
		})

		It("should cap score at 100", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Led development of backend system",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "TechCorp",
				Project:    "Migration",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical", "leadership"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.Score).To(BeNumerically("<=", 100))
		})

		It("should identify missing fields", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Event text",
				Date: time.Now().Add(-24 * time.Hour),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.MissingFields).To(ContainElement("company"))
			Expect(score.MissingFields).To(ContainElement("project"))
			Expect(score.MissingFields).To(ContainElement("tags"))
			Expect(score.MissingFields).To(ContainElement("categories"))
		})

		It("should not list fields as missing when present", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Event text",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "TechCorp",
				Project:    "Migration",
				Tags:       []string{"technical"},
				Categories: []string{"technical"},
			}

			score := calculator.CalculateQuality(event)

			Expect(score.MissingFields).To(BeEmpty())
		})
	})

	Describe("GetQualityLevel", func() {
		It("should return Incomplete for score 0-25", func() {
			Expect(calculator.GetQualityLevel(0)).To(Equal(QualityIncomplete))
			Expect(calculator.GetQualityLevel(10)).To(Equal(QualityIncomplete))
			Expect(calculator.GetQualityLevel(25)).To(Equal(QualityIncomplete))
		})

		It("should return Basic for score 26-50", func() {
			Expect(calculator.GetQualityLevel(26)).To(Equal(QualityBasic))
			Expect(calculator.GetQualityLevel(38)).To(Equal(QualityBasic))
			Expect(calculator.GetQualityLevel(50)).To(Equal(QualityBasic))
		})

		It("should return Enriched for score 51-75", func() {
			Expect(calculator.GetQualityLevel(51)).To(Equal(QualityEnriched))
			Expect(calculator.GetQualityLevel(63)).To(Equal(QualityEnriched))
			Expect(calculator.GetQualityLevel(75)).To(Equal(QualityEnriched))
		})

		It("should return Complete for score 76-100", func() {
			Expect(calculator.GetQualityLevel(76)).To(Equal(QualityComplete))
			Expect(calculator.GetQualityLevel(88)).To(Equal(QualityComplete))
			Expect(calculator.GetQualityLevel(100)).To(Equal(QualityComplete))
		})
	})

	Describe("GetMissingFieldsSuggestion", func() {
		It("should return positive message when no fields missing", func() {
			score := QualityScore{
				MissingFields: []string{},
			}

			suggestion := calculator.GetMissingFieldsSuggestion(score)

			Expect(suggestion).To(Equal("All fields are complete!"))
		})

		It("should return suggestion for single missing field", func() {
			score := QualityScore{
				MissingFields: []string{"company"},
			}

			suggestion := calculator.GetMissingFieldsSuggestion(score)

			Expect(suggestion).To(Equal("Add company to improve quality"))
		})

		It("should return suggestion for two missing fields", func() {
			score := QualityScore{
				MissingFields: []string{"company", "project"},
			}

			suggestion := calculator.GetMissingFieldsSuggestion(score)

			Expect(suggestion).To(Equal("Add company and project to improve quality"))
		})

		It("should return suggestion for three missing fields", func() {
			score := QualityScore{
				MissingFields: []string{"company", "project", "tags"},
			}

			suggestion := calculator.GetMissingFieldsSuggestion(score)

			Expect(suggestion).To(Equal("Add company, project and tags to improve quality"))
		})

		It("should return suggestion for four missing fields", func() {
			score := QualityScore{
				MissingFields: []string{"company", "project", "tags", "categories"},
			}

			suggestion := calculator.GetMissingFieldsSuggestion(score)

			Expect(suggestion).To(Equal("Add company, project, tags and categories to improve quality"))
		})
	})

	Describe("Integration with CareerEvent", func() {
		It("should calculate quality for minimal valid event", func() {
			event := &career.CareerEvent{
				ID:        "test-id",
				Text:      "Completed project",
				Date:      time.Now().Add(-7 * 24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.Score).To(Equal(40)) // 20 (text) + 20 (date)
			Expect(score.Level).To(Equal(QualityBasic))
		})

		It("should calculate quality for event with all optional fields", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Architected microservices migration",
				Date:       time.Now().Add(-30 * 24 * time.Hour),
				Company:    "GlobalTech",
				Project:    "Platform Modernization",
				Tags:       []string{"technical", "leadership", "achievement"},
				Categories: []string{"technical"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.Score).To(BeNumerically(">=", 80))
			Expect(score.Level).To(Equal(QualityComplete))
			Expect(score.MissingFields).To(BeEmpty())
		})

		It("should handle very old dates", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Early career achievement",
				Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.DateScore).To(Equal(20))
		})

		It("should handle long text", func() {
			longText := ""
			for i := 0; i < 100; i++ {
				longText += "word "
			}

			event := &career.CareerEvent{
				ID:   "test-id",
				Text: longText,
				Date: time.Now().Add(-24 * time.Hour),
			}

			score := calculator.CalculateQuality(event)

			Expect(score.TextScore).To(Equal(20))
		})
	})
})
