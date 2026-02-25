package cv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
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

	Describe("NarrativeProfileFromConfig", func() {
		var cfg *config.ProfileConfig

		BeforeEach(func() {
			cfg = &config.ProfileConfig{}
		})

		Context("when FirstName and LastName are provided", func() {
			It("should use FirstName and LastName combined", func() {
				cfg.FirstName = "John"
				cfg.LastName = "Doe"

				profile := cv.NarrativeProfileFromConfig(cfg)

				Expect(profile.Name).To(Equal("John Doe"))
			})

			It("should handle empty LastName gracefully", func() {
				cfg.FirstName = "John"
				cfg.LastName = ""

				profile := cv.NarrativeProfileFromConfig(cfg)

				Expect(profile.Name).To(Equal("John"))
			})

			It("should handle only FirstName", func() {
				cfg.FirstName = "Madonna"
				cfg.LastName = ""

				profile := cv.NarrativeProfileFromConfig(cfg)

				Expect(profile.Name).To(Equal("Madonna"))
			})
		})

		Context("when FirstName is empty but Name is set (legacy)", func() {
			It("should fall back to Name field", func() {
				cfg.Name = "John Doe"
				cfg.FirstName = ""
				cfg.LastName = ""

				profile := cv.NarrativeProfileFromConfig(cfg)

				Expect(profile.Name).To(Equal("John Doe"))
			})
		})

		Context("when both FirstName and Name are empty", func() {
			It("should return empty string", func() {
				cfg.Name = ""
				cfg.FirstName = ""
				cfg.LastName = ""

				profile := cv.NarrativeProfileFromConfig(cfg)

				Expect(profile.Name).To(BeEmpty())
			})
		})

		Context("when FirstName is set but Name is also set", func() {
			It("should prefer FirstName/LastName over Name", func() {
				cfg.Name = "Legacy Name"
				cfg.FirstName = "Jane"
				cfg.LastName = "Smith"

				profile := cv.NarrativeProfileFromConfig(cfg)

				Expect(profile.Name).To(Equal("Jane Smith"))
				Expect(profile.Name).NotTo(Equal("Legacy Name"))
			})
		})
	})
})
