package cv

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CV Variants", func() {
	Describe("TechnologyFocus", func() {
		It("should have Language Agnostic constant", func() {
			Expect(TechnologyFocusLanguageAgnostic).To(Equal(TechnologyFocus("language_agnostic")))
		})

		It("should have Generalist constant", func() {
			Expect(TechnologyFocusGeneralist).To(Equal(TechnologyFocus("generalist")))
		})

		It("should have Specialist constant", func() {
			Expect(TechnologyFocusSpecialist).To(Equal(TechnologyFocus("specialist")))
		})
	})

	Describe("FocusArea", func() {
		It("should have Backend constant", func() {
			Expect(FocusAreaBackend).To(Equal(FocusArea("backend")))
		})

		It("should have Frontend constant", func() {
			Expect(FocusAreaFrontend).To(Equal(FocusArea("frontend")))
		})

		It("should have Fullstack constant", func() {
			Expect(FocusAreaFullstack).To(Equal(FocusArea("fullstack")))
		})

		It("should have DevOps constant", func() {
			Expect(FocusAreaDevOps).To(Equal(FocusArea("devops")))
		})
	})

	Describe("LengthFormat", func() {
		It("should have UltraShort constant", func() {
			Expect(LengthUltraShort).To(Equal(LengthFormat("ultra_short")))
		})

		It("should have Short constant", func() {
			Expect(LengthShort).To(Equal(LengthFormat("short")))
		})

		It("should have Standard constant", func() {
			Expect(LengthStandard).To(Equal(LengthFormat("standard")))
		})

		It("should have Full constant", func() {
			Expect(LengthFull).To(Equal(LengthFormat("full")))
		})
	})

	Describe("CVStructure", func() {
		It("should have Highlights constant", func() {
			Expect(CVStructureHighlights).To(Equal(CVStructure("highlights")))
		})

		It("should have Standard constant", func() {
			Expect(CVStructureStandard).To(Equal(CVStructure("standard")))
		})

		It("should have Narrative constant", func() {
			Expect(CVStructureNarrative).To(Equal(CVStructure("narrative")))
		})
	})

	Describe("CVVariant", func() {
		It("should have all required fields", func() {
			variant := &CVVariant{
				ID:              "agnostic_backend_full",
				Name:            "Language Agnostic - Backend - Full",
				Description:     "Language agnostic CV for backend roles",
				TechnologyFocus: TechnologyFocusLanguageAgnostic,
				FocusArea:       FocusAreaBackend,
				LengthFormat:    LengthFull,
				Technologies:    []string{},
				BaseStructure:   CVStructureNarrative,
			}

			Expect(variant.ID).To(Equal("agnostic_backend_full"))
			Expect(variant.TechnologyFocus).To(Equal(TechnologyFocusLanguageAgnostic))
			Expect(variant.FocusArea).To(Equal(FocusAreaBackend))
			Expect(variant.LengthFormat).To(Equal(LengthFull))
			Expect(variant.Technologies).To(BeEmpty())
			Expect(variant.BaseStructure).To(Equal(CVStructureNarrative))
		})

		It("should support Generalist with multiple technologies", func() {
			variant := &CVVariant{
				ID:              "generalist_fullstack_standard",
				Name:            "Generalist - Fullstack - Standard",
				TechnologyFocus: TechnologyFocusGeneralist,
				FocusArea:       FocusAreaFullstack,
				LengthFormat:    LengthStandard,
				Technologies:    []string{"skill-1", "skill-2", "skill-3"},
				BaseStructure:   CVStructureStandard,
			}

			Expect(variant.Technologies).To(HaveLen(3))
			Expect(variant.BaseStructure).To(Equal(CVStructureStandard))
		})

		It("should support Specialist with single technology", func() {
			variant := &CVVariant{
				ID:              "specialist_ruby_backend_full",
				Name:            "Specialist - Ruby - Backend - Full",
				TechnologyFocus: TechnologyFocusSpecialist,
				FocusArea:       FocusAreaBackend,
				LengthFormat:    LengthFull,
				Technologies:    []string{"ruby"},
				BaseStructure:   CVStructureStandard,
			}

			Expect(variant.Technologies).To(HaveLen(1))
			Expect(variant.Technologies[0]).To(Equal("ruby"))
		})
	})

	Describe("BuildVariantID", func() {
		It("should build Language Agnostic variant ID", func() {
			id := BuildVariantID(
				TechnologyFocusLanguageAgnostic,
				FocusAreaBackend,
				LengthFull,
				[]string{},
			)

			Expect(id).To(Equal("agnostic_backend_full"))
		})

		It("should build Generalist variant ID", func() {
			id := BuildVariantID(
				TechnologyFocusGeneralist,
				FocusAreaFullstack,
				LengthStandard,
				[]string{"ruby", "react"},
			)

			Expect(id).To(Equal("generalist_fullstack_standard"))
		})

		It("should build Specialist variant ID with technology", func() {
			id := BuildVariantID(
				TechnologyFocusSpecialist,
				FocusAreaBackend,
				LengthFull,
				[]string{"ruby"},
			)

			Expect(id).To(Equal("specialist_ruby_backend_full"))
		})
	})

	Describe("DetermineStructure", func() {
		It("should use Highlights for UltraShort length", func() {
			structure := DetermineStructure(
				TechnologyFocusLanguageAgnostic,
				LengthUltraShort,
			)

			Expect(structure).To(Equal(CVStructureHighlights))
		})

		It("should use Narrative for Language Agnostic (non-UltraShort)", func() {
			structure := DetermineStructure(
				TechnologyFocusLanguageAgnostic,
				LengthFull,
			)

			Expect(structure).To(Equal(CVStructureNarrative))
		})

		It("should use Standard for Generalist", func() {
			structure := DetermineStructure(
				TechnologyFocusGeneralist,
				LengthStandard,
			)

			Expect(structure).To(Equal(CVStructureStandard))
		})

		It("should use Standard for Specialist", func() {
			structure := DetermineStructure(
				TechnologyFocusSpecialist,
				LengthFull,
			)

			Expect(structure).To(Equal(CVStructureStandard))
		})
	})
})
