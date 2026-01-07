package forms_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstSuggestionForm", func() {
	var suggestion burstfact.BurstSuggestion

	BeforeEach(func() {
		suggestion = burstfact.BurstSuggestion{
			Name:            "Project Alpha Work",
			Description:     "Work related to Project Alpha",
			EventIDs:        []string{"event1", "event2", "event3"},
			ConfidenceScore: 0.85,
		}
	})

	Describe("NewBurstSuggestionEditForm", func() {
		It("should create a form with name and description fields", func() {
			form := forms.NewBurstSuggestionEditForm(suggestion)
			Expect(form).NotTo(BeNil())
		})

		It("should initialize form with suggestion data", func() {
			formData := forms.GetBurstSuggestionFormData(suggestion)
			Expect(formData.Name).To(Equal("Project Alpha Work"))
			Expect(formData.Description).To(Equal("Work related to Project Alpha"))
		})
	})

	Describe("NewBurstSuggestionEditFormWithData", func() {
		It("should create a form with provided data", func() {
			data := &forms.BurstSuggestionFormData{
				Name:        "Custom Name",
				Description: "Custom Description",
			}
			form := forms.NewBurstSuggestionEditFormWithData(data)
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("GetBurstSuggestionFormData", func() {
		It("should extract form data from suggestion", func() {
			data := forms.GetBurstSuggestionFormData(suggestion)
			Expect(data).NotTo(BeNil())
			Expect(data.Name).To(Equal("Project Alpha Work"))
			Expect(data.Description).To(Equal("Work related to Project Alpha"))
		})

		It("should handle empty name", func() {
			suggestion.Name = ""
			data := forms.GetBurstSuggestionFormData(suggestion)
			Expect(data.Name).To(BeEmpty())
			Expect(data.Description).To(Equal("Work related to Project Alpha"))
		})

		It("should handle empty description", func() {
			suggestion.Description = ""
			data := forms.GetBurstSuggestionFormData(suggestion)
			Expect(data.Name).To(Equal("Project Alpha Work"))
			Expect(data.Description).To(BeEmpty())
		})

		It("should handle both fields empty", func() {
			suggestion.Name = ""
			suggestion.Description = ""
			data := forms.GetBurstSuggestionFormData(suggestion)
			Expect(data.Name).To(BeEmpty())
			Expect(data.Description).To(BeEmpty())
		})
	})

	Describe("ApplyBurstSuggestionFormData", func() {
		It("should apply form data to suggestion", func() {
			data := &forms.BurstSuggestionFormData{
				Name:        "Updated Name",
				Description: "Updated Description",
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)

			Expect(suggestion.Name).To(Equal("Updated Name"))
			Expect(suggestion.Description).To(Equal("Updated Description"))
		})

		It("should overwrite existing data", func() {
			data := &forms.BurstSuggestionFormData{
				Name:        "New Name",
				Description: "New Description",
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)

			Expect(suggestion.Name).To(Equal("New Name"))
			Expect(suggestion.Description).To(Equal("New Description"))
		})

		It("should allow clearing name", func() {
			data := &forms.BurstSuggestionFormData{
				Name:        "",
				Description: "Keep description",
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)

			Expect(suggestion.Name).To(BeEmpty())
			Expect(suggestion.Description).To(Equal("Keep description"))
		})

		It("should allow clearing description", func() {
			data := &forms.BurstSuggestionFormData{
				Name:        "Keep name",
				Description: "",
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)

			Expect(suggestion.Name).To(Equal("Keep name"))
			Expect(suggestion.Description).To(BeEmpty())
		})

		It("should preserve EventIDs and ConfidenceScore", func() {
			originalEventIDs := suggestion.EventIDs
			originalConfidence := suggestion.ConfidenceScore

			data := &forms.BurstSuggestionFormData{
				Name:        "New Name",
				Description: "New Description",
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)

			Expect(suggestion.EventIDs).To(Equal(originalEventIDs))
			Expect(suggestion.ConfidenceScore).To(Equal(originalConfidence))
		})
	})

	Describe("Roundtrip data conversion", func() {
		It("should preserve data through get and apply cycle", func() {
			// Get data from suggestion
			data := forms.GetBurstSuggestionFormData(suggestion)

			// Modify data
			data.Name = "Modified Name"
			data.Description = "Modified Description"

			// Apply back to suggestion
			forms.ApplyBurstSuggestionFormData(&suggestion, data)

			// Verify changes
			Expect(suggestion.Name).To(Equal("Modified Name"))
			Expect(suggestion.Description).To(Equal("Modified Description"))
		})

		It("should handle empty values in roundtrip", func() {
			suggestion.Name = ""
			suggestion.Description = ""

			data := forms.GetBurstSuggestionFormData(suggestion)
			Expect(data.Name).To(BeEmpty())
			Expect(data.Description).To(BeEmpty())

			data.Name = "Filled Name"
			data.Description = "Filled Description"

			forms.ApplyBurstSuggestionFormData(&suggestion, data)
			Expect(suggestion.Name).To(Equal("Filled Name"))
			Expect(suggestion.Description).To(Equal("Filled Description"))
		})
	})

	Describe("Form behavior expectations", func() {
		It("should accept names up to 100 characters", func() {
			longName := string(make([]byte, 100))
			for i := range longName {
				longName = longName[:i] + "a" + longName[i+1:]
			}

			data := &forms.BurstSuggestionFormData{
				Name:        longName,
				Description: "Description",
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)
			Expect(len(suggestion.Name)).To(Equal(100))
		})

		It("should accept descriptions up to 500 characters", func() {
			longDesc := string(make([]byte, 500))
			for i := range longDesc {
				longDesc = longDesc[:i] + "a" + longDesc[i+1:]
			}

			data := &forms.BurstSuggestionFormData{
				Name:        "Name",
				Description: longDesc,
			}

			forms.ApplyBurstSuggestionFormData(&suggestion, data)
			Expect(len(suggestion.Description)).To(Equal(500))
		})
	})
})
