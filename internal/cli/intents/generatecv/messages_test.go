package generatecv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Messages", func() {
	Describe("CVGenerationCompleteMsg", func() {
		It("should carry a CV view on success", func() {
			cv := fixtures.CVView("cv-1")
			msg := generatecv.CVGenerationCompleteMsg{CV: cv, Error: nil}

			Expect(msg.CV).To(Equal(cv))
			Expect(msg.Error).ToNot(HaveOccurred())
		})

		It("should carry an error on failure", func() {
			err := errors.New("generation failed")
			msg := generatecv.CVGenerationCompleteMsg{CV: nil, Error: err}

			Expect(msg.CV).To(BeNil())
			Expect(msg.Error).To(MatchError("generation failed"))
		})
	})

	Describe("TechnologiesExtractedMsg", func() {
		It("should carry technologies on success", func() {
			msg := generatecv.TechnologiesExtractedMsg{
				Technologies: []*generatecv.ExtractedTechnology{},
				Suggestion:   &generatecv.FocusAreaSuggestion{},
				Error:        nil,
			}

			Expect(msg.Technologies).To(BeEmpty())
			Expect(msg.Suggestion).NotTo(BeNil())
			Expect(msg.Error).ToNot(HaveOccurred())
		})

		It("should carry an error on failure", func() {
			err := errors.New("extraction failed")
			msg := generatecv.TechnologiesExtractedMsg{Error: err}

			Expect(msg.Error).To(MatchError("extraction failed"))
		})
	})

	Describe("WizardCompleteMsg", func() {
		It("should carry all wizard configuration fields", func() {
			msg := generatecv.WizardCompleteMsg{
				ProfileID:    "profile-1",
				Audience:     "hiring_manager",
				TechFocus:    "all",
				Technologies: []string{"Go", "Python"},
				FocusArea:    "backend",
				SkillsFormat: "grouped",
				SkillsLimit:  10,
				CVLength:     "standard",
			}

			Expect(msg.ProfileID).To(Equal("profile-1"))
			Expect(msg.Audience).To(Equal("hiring_manager"))
			Expect(msg.TechFocus).To(Equal("all"))
			Expect(msg.Technologies).To(ConsistOf("Go", "Python"))
			Expect(msg.FocusArea).To(Equal("backend"))
			Expect(msg.SkillsFormat).To(Equal("grouped"))
			Expect(msg.SkillsLimit).To(Equal(10))
			Expect(msg.CVLength).To(Equal("standard"))
		})
	})

	Describe("ExportCompleteMsg", func() {
		It("should carry a path on success", func() {
			msg := generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil}

			Expect(msg.Path).To(Equal("/tmp/cv.md"))
			Expect(msg.Error).ToNot(HaveOccurred())
		})

		It("should carry an error on failure", func() {
			err := errors.New("export failed")
			msg := generatecv.ExportCompleteMsg{Path: "", Error: err}

			Expect(msg.Path).To(BeEmpty())
			Expect(msg.Error).To(MatchError("export failed"))
		})
	})
})
