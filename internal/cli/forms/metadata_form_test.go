package forms_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("MetadataForm", func() {
	var testEvent *career.CareerEvent

	BeforeEach(func() {
		testEvent = fixtures.Event("event-123")
		testEvent.Text = "Test event"
		testEvent.Date = time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)
		testEvent.Company = "Test Company"
		testEvent.Project = "Test Project"
		testEvent.Tags = []string{"go", "testing"}
		testEvent.Categories = []string{"technical", "leadership"}
		testEvent.Skills = []string{"skill-id-1", "skill-id-2"}
	})

	Describe("NewMetadataEditorForm", func() {
		It("should create a form with event data", func() {
			availableTags := []string{"go", "python", "testing"}
			availableCategories := []string{"technical", "leadership", "mentoring"}
			availableSkills := []*career.Skill{}

			form := forms.NewMetadataEditorForm(testEvent, availableTags, availableCategories, availableSkills)

			Expect(form).NotTo(BeNil())
		})
	})

	Describe("MetadataFormData", func() {
		It("should extract data from event", func() {
			data := forms.GetMetadataFormData(testEvent)

			Expect(data.Date).To(Equal("2024-01-07"))
			Expect(data.Company).To(Equal("Test Company"))
			Expect(data.Project).To(Equal("Test Project"))
			Expect(data.Tags).To(Equal([]string{"go", "testing"}))
			Expect(data.Categories).To(Equal([]string{"technical", "leadership"}))
			Expect(data.Skills).To(Equal([]string{"skill-id-1", "skill-id-2"}))
		})

		It("should apply data to event", func() {
			newEvent := &career.CareerEvent{}
			data := &forms.MetadataFormData{
				Date:       "2024-01-15",
				Company:    "New Company",
				Project:    "New Project",
				Tags:       []string{"rust", "performance"},
				Categories: []string{"research"},
				Skills:     []string{"skill-id-3", "skill-id-4"},
			}

			err := forms.ApplyMetadataFormData(newEvent, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(newEvent.Date.Format("2006-01-02")).To(Equal("2024-01-15"))
			Expect(newEvent.Company).To(Equal("New Company"))
			Expect(newEvent.Project).To(Equal("New Project"))
			Expect(newEvent.Tags).To(Equal([]string{"rust", "performance"}))
			Expect(newEvent.Categories).To(Equal([]string{"research"}))
			Expect(newEvent.Skills).To(Equal([]string{"skill-id-3", "skill-id-4"}))
		})
	})

	Describe("ParseDateString", func() {
		It("should parse YYYY-MM-DD format", func() {
			date, err := forms.ParseDateString("2024-01-07")

			Expect(err).NotTo(HaveOccurred())
			Expect(date.Year()).To(Equal(2024))
			Expect(date.Month()).To(Equal(time.January))
			Expect(date.Day()).To(Equal(7))
		})

		It("should parse 'today'", func() {
			date, err := forms.ParseDateString("today")

			Expect(err).NotTo(HaveOccurred())
			Expect(date.Year()).To(Equal(time.Now().Year()))
			Expect(date.Month()).To(Equal(time.Now().Month()))
			Expect(date.Day()).To(Equal(time.Now().Day()))
		})

		It("should parse relative dates - days", func() {
			date, err := forms.ParseDateString("7 days ago")

			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -7)
			Expect(date.Year()).To(Equal(expected.Year()))
			Expect(date.Month()).To(Equal(expected.Month()))
			Expect(date.Day()).To(Equal(expected.Day()))
		})

		It("should parse relative dates - weeks", func() {
			date, err := forms.ParseDateString("2 weeks ago")

			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -14)
			Expect(date.Year()).To(Equal(expected.Year()))
			Expect(date.Month()).To(Equal(expected.Month()))
			Expect(date.Day()).To(Equal(expected.Day()))
		})

		It("should parse relative dates - months", func() {
			date, err := forms.ParseDateString("1 month ago")

			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, -1, 0)
			Expect(date.Year()).To(Equal(expected.Year()))
			Expect(date.Month()).To(Equal(expected.Month()))
		})

		It("should handle singular units", func() {
			date, err := forms.ParseDateString("1 day ago")

			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -1)
			Expect(date.Day()).To(Equal(expected.Day()))
		})

		It("should fail on invalid format", func() {
			_, err := forms.ParseDateString("invalid")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid date format"))
		})

		It("should fail on invalid relative format", func() {
			_, err := forms.ParseDateString("ago 1 week")

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("NewMetadataEditorFormWithData", func() {
		It("should create a form with initial data", func() {
			data := &forms.MetadataFormData{
				Date:       "2024-01-15",
				Company:    "Initial Company",
				Project:    "Initial Project",
				Tags:       []string{"tag1"},
				Categories: []string{"cat1"},
			}

			availableTags := []string{"tag1", "tag2"}
			availableCategories := []string{"cat1", "cat2"}
			availableSkills := []*career.Skill{}

			form := forms.NewMetadataEditorFormWithData(data, availableTags, availableCategories, availableSkills)

			Expect(form).NotTo(BeNil())
		})
	})
})
