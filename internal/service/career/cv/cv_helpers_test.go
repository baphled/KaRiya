package cv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/service/career/cv"
)

var _ = Describe("CVHelpers", func() {
	Describe("DefaultNarrativeProfile", func() {
		It("should not contain hardcoded personal data", func() {
			profile := cv.DefaultNarrativeProfile()

			// Should NOT contain specific person's data
			Expect(profile.Name).NotTo(ContainSubstring("Yomi"))
			Expect(profile.Name).NotTo(ContainSubstring("Colledge"))
			Expect(profile.Email).NotTo(ContainSubstring("boodah"))
			Expect(profile.GitHub).NotTo(ContainSubstring("baphled"))
			Expect(profile.Portfolio).NotTo(ContainSubstring("boodah"))
		})

		It("should return empty strings for personal fields", func() {
			profile := cv.DefaultNarrativeProfile()

			// Personal data should be empty (user must provide via onboarding)
			Expect(profile.Name).To(BeEmpty())
			Expect(profile.Email).To(BeEmpty())
			Expect(profile.GitHub).To(BeEmpty())
			Expect(profile.Portfolio).To(BeEmpty())
			Expect(profile.Location).To(BeEmpty())
			Expect(profile.Role).To(BeEmpty())
		})

		It("should return empty slices for inferred fields", func() {
			profile := cv.DefaultNarrativeProfile()

			// Inferred fields should be empty (will be populated by inference service)
			Expect(profile.CoreStrengths).To(BeEmpty())
			Expect(profile.ValuePropositions).To(BeEmpty())
		})

		It("should return empty slices for technology fields", func() {
			profile := cv.DefaultNarrativeProfile()

			// Technology fields should be empty slices (will be populated from skills)
			Expect(profile.Languages).To(BeEmpty())
			Expect(profile.Frontend).To(BeEmpty())
			Expect(profile.Systems).To(BeEmpty())
		})
	})
})
