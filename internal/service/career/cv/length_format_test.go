package cv_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/service/career/cv"
)

var _ = Describe("LengthFormat", func() {
	Describe("GetLengthFormatConfig", func() {
		It("should return config for LengthFull", func() {
			config := cv.GetLengthFormatConfig(cv.LengthFull)

			Expect(config.ID).To(Equal(cv.LengthFull))
			Expect(config.Name).To(Equal("Full"))
			Expect(config.MaxYearsHistory).To(BeNil())
			Expect(config.MaxCompanies).To(BeNil())
			Expect(config.MinConfidence).To(Equal(0.50))
			Expect(config.TargetPages).To(Equal("3+"))
			Expect(config.IncludeSummary).To(BeTrue())
			Expect(config.IncludeEducation).To(BeTrue())
			Expect(config.IncludeSkills).To(BeTrue())
		})

		It("should return config for LengthStandard", func() {
			config := cv.GetLengthFormatConfig(cv.LengthStandard)

			Expect(config.ID).To(Equal(cv.LengthStandard))
			Expect(config.Name).To(Equal("Standard"))
			Expect(*config.MaxYearsHistory).To(Equal(10))
			Expect(config.MaxCompanies).To(BeNil())
			Expect(*config.MaxBulletsPerJob).To(Equal(5))
			Expect(config.MinConfidence).To(Equal(0.65))
			Expect(config.TargetPages).To(Equal("2-3"))
			Expect(config.IncludeSummary).To(BeTrue())
			Expect(config.IncludeEducation).To(BeTrue())
			Expect(config.IncludeSkills).To(BeTrue())
		})

		It("should return config for LengthShort", func() {
			config := cv.GetLengthFormatConfig(cv.LengthShort)

			Expect(config.ID).To(Equal(cv.LengthShort))
			Expect(config.Name).To(Equal("Short"))
			Expect(*config.MaxYearsHistory).To(Equal(5))
			Expect(*config.MaxCompanies).To(Equal(5))
			Expect(*config.MaxBulletsPerJob).To(Equal(3))
			Expect(config.MinConfidence).To(Equal(0.75))
			Expect(config.TargetPages).To(Equal("1-2"))
			Expect(config.IncludeSummary).To(BeTrue())
			Expect(config.IncludeEducation).To(BeFalse())
			Expect(config.IncludeSkills).To(BeTrue())
		})

		It("should return config for LengthUltraShort", func() {
			config := cv.GetLengthFormatConfig(cv.LengthUltraShort)

			Expect(config.ID).To(Equal(cv.LengthUltraShort))
			Expect(config.Name).To(Equal("1-Page"))
			Expect(*config.MaxYearsHistory).To(Equal(3))
			Expect(*config.MaxCompanies).To(Equal(3))
			Expect(*config.MaxBulletsPerJob).To(Equal(2))
			Expect(config.MinConfidence).To(Equal(0.85))
			Expect(config.TargetPages).To(Equal("1"))
			Expect(config.IncludeSummary).To(BeTrue())
			Expect(config.IncludeEducation).To(BeFalse())
			Expect(config.IncludeSkills).To(BeFalse())
		})

		It("should return default config for unknown format", func() {
			config := cv.GetLengthFormatConfig(cv.LengthFormat("unknown"))

			Expect(config.ID).To(Equal(cv.LengthFormat("unknown")))
			Expect(config.Name).To(Equal("unknown"))
			Expect(config.MinConfidence).To(Equal(0.65))
			Expect(config.TargetPages).To(Equal("2-3"))
		})
	})

	Describe("ListLengthFormatConfigs", func() {
		It("should return all 4 length format configurations", func() {
			configs := cv.ListLengthFormatConfigs()

			Expect(configs).To(HaveLen(4))
			Expect(configs[0].ID).To(Equal(cv.LengthFull))
			Expect(configs[1].ID).To(Equal(cv.LengthStandard))
			Expect(configs[2].ID).To(Equal(cv.LengthShort))
			Expect(configs[3].ID).To(Equal(cv.LengthUltraShort))
		})
	})

	Describe("ShouldIncludeEvent", func() {
		Context("with no year limit (Full format)", func() {
			It("should include events from any date", func() {
				config := cv.GetLengthFormatConfig(cv.LengthFull)

				// Event from 20 years ago
				oldEvent := time.Now().AddDate(-20, 0, 0)
				Expect(config.ShouldIncludeEvent(oldEvent)).To(BeTrue())

				// Recent event
				recentEvent := time.Now().AddDate(0, -1, 0)
				Expect(config.ShouldIncludeEvent(recentEvent)).To(BeTrue())
			})
		})

		Context("with 10-year limit (Standard format)", func() {
			It("should include events within 10 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthStandard)

				// Event from 5 years ago
				event := time.Now().AddDate(-5, 0, 0)
				Expect(config.ShouldIncludeEvent(event)).To(BeTrue())
			})

			It("should exclude events older than 10 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthStandard)

				// Event from 15 years ago
				event := time.Now().AddDate(-15, 0, 0)
				Expect(config.ShouldIncludeEvent(event)).To(BeFalse())
			})
		})

		Context("with 5-year limit (Short format)", func() {
			It("should include events within 5 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthShort)

				// Event from 3 years ago
				event := time.Now().AddDate(-3, 0, 0)
				Expect(config.ShouldIncludeEvent(event)).To(BeTrue())
			})

			It("should exclude events older than 5 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthShort)

				// Event from 7 years ago
				event := time.Now().AddDate(-7, 0, 0)
				Expect(config.ShouldIncludeEvent(event)).To(BeFalse())
			})
		})

		Context("with 3-year limit (UltraShort format)", func() {
			It("should include events within 3 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthUltraShort)

				// Event from 2 years ago
				event := time.Now().AddDate(-2, 0, 0)
				Expect(config.ShouldIncludeEvent(event)).To(BeTrue())
			})

			It("should exclude events older than 3 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthUltraShort)

				// Event from 5 years ago
				event := time.Now().AddDate(-5, 0, 0)
				Expect(config.ShouldIncludeEvent(event)).To(BeFalse())
			})
		})
	})

	Describe("ShouldIncludeEventByYear", func() {
		currentYear := time.Now().Year()

		Context("with no year limit (Full format)", func() {
			It("should include events from any year", func() {
				config := cv.GetLengthFormatConfig(cv.LengthFull)

				Expect(config.ShouldIncludeEventByYear(currentYear - 20)).To(BeTrue())
				Expect(config.ShouldIncludeEventByYear(currentYear)).To(BeTrue())
			})
		})

		Context("with 10-year limit (Standard format)", func() {
			It("should include events from current year minus 10", func() {
				config := cv.GetLengthFormatConfig(cv.LengthStandard)

				Expect(config.ShouldIncludeEventByYear(currentYear)).To(BeTrue())
				Expect(config.ShouldIncludeEventByYear(currentYear - 10)).To(BeTrue())
			})

			It("should exclude events older than 10 years", func() {
				config := cv.GetLengthFormatConfig(cv.LengthStandard)

				Expect(config.ShouldIncludeEventByYear(currentYear - 11)).To(BeFalse())
				Expect(config.ShouldIncludeEventByYear(currentYear - 15)).To(BeFalse())
			})
		})
	})

	Describe("FilterCompaniesByLimit", func() {
		companies := []string{"Company A", "Company B", "Company C", "Company D", "Company E", "Company F"}

		Context("with no company limit (Full format)", func() {
			It("should return all companies", func() {
				config := cv.GetLengthFormatConfig(cv.LengthFull)
				result := config.FilterCompaniesByLimit(companies)

				Expect(result).To(Equal(companies))
			})
		})

		Context("with 5 company limit (Short format)", func() {
			It("should return first 5 companies", func() {
				config := cv.GetLengthFormatConfig(cv.LengthShort)
				result := config.FilterCompaniesByLimit(companies)

				Expect(result).To(HaveLen(5))
				Expect(result).To(Equal([]string{"Company A", "Company B", "Company C", "Company D", "Company E"}))
			})

			It("should return all companies if less than limit", func() {
				config := cv.GetLengthFormatConfig(cv.LengthShort)
				fewCompanies := []string{"Company A", "Company B"}
				result := config.FilterCompaniesByLimit(fewCompanies)

				Expect(result).To(Equal(fewCompanies))
			})
		})

		Context("with 3 company limit (UltraShort format)", func() {
			It("should return first 3 companies", func() {
				config := cv.GetLengthFormatConfig(cv.LengthUltraShort)
				result := config.FilterCompaniesByLimit(companies)

				Expect(result).To(HaveLen(3))
				Expect(result).To(Equal([]string{"Company A", "Company B", "Company C"}))
			})
		})
	})

	Describe("GetEffectiveBulletLimit", func() {
		It("should return default limit when MaxBulletsPerJob is nil", func() {
			config := cv.GetLengthFormatConfig(cv.LengthFull)

			Expect(config.GetEffectiveBulletLimit(10)).To(Equal(10))
			Expect(config.GetEffectiveBulletLimit(5)).To(Equal(5))
		})

		It("should return configured limit when MaxBulletsPerJob is set", func() {
			config := cv.GetLengthFormatConfig(cv.LengthStandard)

			Expect(config.GetEffectiveBulletLimit(10)).To(Equal(5)) // Configured as 5
		})

		It("should return 3 for Short format", func() {
			config := cv.GetLengthFormatConfig(cv.LengthShort)

			Expect(config.GetEffectiveBulletLimit(10)).To(Equal(3))
		})

		It("should return 2 for UltraShort format", func() {
			config := cv.GetLengthFormatConfig(cv.LengthUltraShort)

			Expect(config.GetEffectiveBulletLimit(10)).To(Equal(2))
		})
	})

	Describe("MeetsConfidenceThreshold", func() {
		It("should accept bullets meeting Full format threshold (0.50)", func() {
			config := cv.GetLengthFormatConfig(cv.LengthFull)

			Expect(config.MeetsConfidenceThreshold(0.50)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.75)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.49)).To(BeFalse())
		})

		It("should accept bullets meeting Standard format threshold (0.65)", func() {
			config := cv.GetLengthFormatConfig(cv.LengthStandard)

			Expect(config.MeetsConfidenceThreshold(0.65)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.80)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.64)).To(BeFalse())
		})

		It("should accept bullets meeting Short format threshold (0.75)", func() {
			config := cv.GetLengthFormatConfig(cv.LengthShort)

			Expect(config.MeetsConfidenceThreshold(0.75)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.90)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.74)).To(BeFalse())
		})

		It("should accept bullets meeting UltraShort format threshold (0.85)", func() {
			config := cv.GetLengthFormatConfig(cv.LengthUltraShort)

			Expect(config.MeetsConfidenceThreshold(0.85)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.95)).To(BeTrue())
			Expect(config.MeetsConfidenceThreshold(0.84)).To(BeFalse())
		})
	})

	Describe("LengthFormat Progression", func() {
		It("should have progressively stricter confidence thresholds", func() {
			full := cv.GetLengthFormatConfig(cv.LengthFull)
			standard := cv.GetLengthFormatConfig(cv.LengthStandard)
			short := cv.GetLengthFormatConfig(cv.LengthShort)
			ultraShort := cv.GetLengthFormatConfig(cv.LengthUltraShort)

			Expect(full.MinConfidence).To(BeNumerically("<", standard.MinConfidence))
			Expect(standard.MinConfidence).To(BeNumerically("<", short.MinConfidence))
			Expect(short.MinConfidence).To(BeNumerically("<", ultraShort.MinConfidence))
		})

		It("should have progressively stricter year limits", func() {
			// Full has no limit
			full := cv.GetLengthFormatConfig(cv.LengthFull)
			Expect(full.MaxYearsHistory).To(BeNil())

			// Others have progressively stricter limits
			standard := cv.GetLengthFormatConfig(cv.LengthStandard)
			short := cv.GetLengthFormatConfig(cv.LengthShort)
			ultraShort := cv.GetLengthFormatConfig(cv.LengthUltraShort)

			Expect(*standard.MaxYearsHistory).To(Equal(10))
			Expect(*short.MaxYearsHistory).To(Equal(5))
			Expect(*ultraShort.MaxYearsHistory).To(Equal(3))
		})

		It("should have progressively stricter bullet limits", func() {
			// Full has no limit
			full := cv.GetLengthFormatConfig(cv.LengthFull)
			Expect(full.MaxBulletsPerJob).To(BeNil())

			// Others have progressively stricter limits
			standard := cv.GetLengthFormatConfig(cv.LengthStandard)
			short := cv.GetLengthFormatConfig(cv.LengthShort)
			ultraShort := cv.GetLengthFormatConfig(cv.LengthUltraShort)

			Expect(*standard.MaxBulletsPerJob).To(Equal(5))
			Expect(*short.MaxBulletsPerJob).To(Equal(3))
			Expect(*ultraShort.MaxBulletsPerJob).To(Equal(2))
		})
	})

	Describe("MapUILengthToFormat", func() {
		It("should map '1_page' to LengthUltraShort", func() {
			result := cv.MapUILengthToFormat("1_page")
			Expect(result).To(Equal(cv.LengthUltraShort))
		})

		It("should map '2_page' to LengthShort", func() {
			result := cv.MapUILengthToFormat("2_page")
			Expect(result).To(Equal(cv.LengthShort))
		})

		It("should map 'standard' to LengthStandard", func() {
			result := cv.MapUILengthToFormat("standard")
			Expect(result).To(Equal(cv.LengthStandard))
		})

		It("should map 'detailed' to LengthFull", func() {
			result := cv.MapUILengthToFormat("detailed")
			Expect(result).To(Equal(cv.LengthFull))
		})

		It("should return LengthStandard as default for unknown values", func() {
			result := cv.MapUILengthToFormat("unknown")
			Expect(result).To(Equal(cv.LengthStandard))

			result = cv.MapUILengthToFormat("")
			Expect(result).To(Equal(cv.LengthStandard))
		})
	})
})
