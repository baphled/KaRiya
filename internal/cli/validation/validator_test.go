package validation_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/validation"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}

var _ = Describe("ValidationError", func() {
	Describe("Error message formatting", func() {
		It("should include explanation and suggestion", func() {
			verr := validation.NewValidationError(
				"text",
				"Event text is required",
				"Please enter a description of your career event",
			)

			Expect(verr.Field).To(Equal("text"))
			Expect(verr.Message).To(Equal("Event text is required"))
			Expect(verr.Suggestion).To(Equal("Please enter a description of your career event"))
			Expect(verr.Error()).To(ContainSubstring("Event text is required"))
		})
	})

	Describe("Formatted display", func() {
		It("should format error with explanation and suggestion", func() {
			verr := validation.NewValidationError(
				"date",
				"Date cannot be in the future",
				"Please enter a past date or use 'today'",
			)

			formatted := verr.Formatted()
			Expect(formatted).To(ContainSubstring("Date cannot be in the future"))
			Expect(formatted).To(ContainSubstring("Please enter a past date or use 'today'"))
		})
	})
})

var _ = Describe("EventValidator", func() {
	var validator *validation.EventValidator

	BeforeEach(func() {
		validator = validation.NewEventValidator()
	})

	Describe("ValidateText", func() {
		Context("when text is valid", func() {
			It("should return nil error", func() {
				err := validator.ValidateText("Valid event text")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when text is empty", func() {
			It("should return ValidationError with suggestion", func() {
				err := validator.ValidateText("")
				Expect(err).To(HaveOccurred())

				verr, ok := err.(*validation.ValidationError)
				Expect(ok).To(BeTrue())
				Expect(verr.Field).To(Equal("text"))
				Expect(verr.Message).To(ContainSubstring("required"))
				Expect(verr.Suggestion).ToNot(BeEmpty())
			})
		})

		Context("when text is whitespace only", func() {
			It("should return ValidationError with suggestion", func() {
				err := validator.ValidateText("   ")
				Expect(err).To(HaveOccurred())

				verr, ok := err.(*validation.ValidationError)
				Expect(ok).To(BeTrue())
				Expect(verr.Suggestion).ToNot(BeEmpty())
			})
		})

		Context("when text exceeds 2000 characters", func() {
			It("should return ValidationError with suggestion", func() {
				longText := make([]byte, 2001)
				for i := range longText {
					longText[i] = 'a'
				}
				err := validator.ValidateText(string(longText))
				Expect(err).To(HaveOccurred())

				verr, ok := err.(*validation.ValidationError)
				Expect(ok).To(BeTrue())
				Expect(verr.Message).To(ContainSubstring("2000 characters"))
				Expect(verr.Suggestion).ToNot(BeEmpty())
			})
		})
	})

	Describe("ValidateDate", func() {
		Context("when date is valid", func() {
			It("should return nil error", func() {
				pastDate := time.Now().Add(-24 * time.Hour)
				err := validator.ValidateDate(pastDate)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when date is in the future", func() {
			It("should return error", func() {
				futureDate := time.Now().Add(24 * time.Hour)
				err := validator.ValidateDate(futureDate)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("future"))
			})
		})
	})

	Describe("ValidateDateForMode", func() {
		Context("when using TimelineJournaling mode", func() {
			It("should allow dates within 30 days", func() {
				recentDate := time.Now().Add(-15 * 24 * time.Hour)
				err := validator.ValidateDateForMode(recentDate, careerservice.TimelineJournaling)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should reject dates older than 30 days", func() {
				oldDate := time.Now().Add(-31 * 24 * time.Hour)
				err := validator.ValidateDateForMode(oldDate, careerservice.TimelineJournaling)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("30 days"))
			})
		})

		Context("when using CVBackfill mode", func() {
			It("should allow any past date", func() {
				oldDate := time.Now().Add(-365 * 24 * time.Hour)
				err := validator.ValidateDateForMode(oldDate, careerservice.CVBackfill)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when using ManualEntry mode", func() {
			It("should allow any past date", func() {
				oldDate := time.Now().Add(-365 * 24 * time.Hour)
				err := validator.ValidateDateForMode(oldDate, careerservice.ManualEntry)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("ValidateTags", func() {
		Context("when tags are valid", func() {
			It("should return nil error", func() {
				tags := []string{"technical", "leadership"}
				err := validator.ValidateTags(tags)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when tags contain duplicates", func() {
			It("should return error", func() {
				tags := []string{"technical", "technical"}
				err := validator.ValidateTags(tags)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("duplicate"))
			})
		})

		Context("when tags exceed maximum of 8", func() {
			It("should return error", func() {
				tags := []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7", "tag8", "tag9"}
				err := validator.ValidateTags(tags)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("8 tags"))
			})
		})

		Context("when tags contain invalid tag", func() {
			It("should return error", func() {
				tags := []string{"invalid_tag"}
				err := validator.ValidateTags(tags)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid tag"))
			})
		})
	})
})
