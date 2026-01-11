package career

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Skill", func() {
	Describe("Validation", func() {
		Describe("Name validation", func() {
			It("should fail when name is empty", func() {
				skill := &Skill{
					Name:     "",
					Category: "backend",
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
			})

			It("should fail when name is only whitespace", func() {
				skill := &Skill{
					Name:     "   ",
					Category: "backend",
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
			})

			It("should fail when name exceeds 100 characters", func() {
				longName := string(make([]byte, 101))
				for i := range longName {
					longName = string(append([]byte(longName[:i]), 'a'))
				}
				skill := &Skill{
					Name:     longName,
					Category: "backend",
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name cannot exceed 100 characters"))
			})

			It("should pass when name is valid", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("Category validation", func() {
			It("should fail when category is empty", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "",
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("category cannot be empty"))
			})

			It("should fail when category is only whitespace", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "   ",
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("category cannot be empty"))
			})

			It("should fail when category exceeds 50 characters", func() {
				longCategory := string(make([]byte, 51))
				for i := range longCategory {
					longCategory = string(append([]byte(longCategory[:i]), 'a'))
				}
				skill := &Skill{
					Name:     "Ruby",
					Category: longCategory,
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("category cannot exceed 50 characters"))
			})

			It("should pass when category is valid", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("Level validation", func() {
			It("should pass when level is empty (optional)", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
					Level:    "",
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass when level is valid", func() {
				validLevels := []string{"beginner", "intermediate", "advanced", "expert"}
				for _, level := range validLevels {
					skill := &Skill{
						Name:     "Ruby",
						Category: "backend",
						Level:    level,
					}
					err := skill.Validate()
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("should fail when level is invalid", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
					Level:    "master",
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("level must be one of"))
			})
		})

		Describe("YearsUsed validation", func() {
			It("should pass when years_used is nil (optional)", func() {
				skill := &Skill{
					Name:      "Ruby",
					Category:  "backend",
					YearsUsed: nil,
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass when years_used is within range", func() {
				years := 5
				skill := &Skill{
					Name:      "Ruby",
					Category:  "backend",
					YearsUsed: &years,
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail when years_used is negative", func() {
				years := -1
				skill := &Skill{
					Name:      "Ruby",
					Category:  "backend",
					YearsUsed: &years,
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("years_used must be between 0 and 50"))
			})

			It("should fail when years_used exceeds 50", func() {
				years := 51
				skill := &Skill{
					Name:      "Ruby",
					Category:  "backend",
					YearsUsed: &years,
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("years_used must be between 0 and 50"))
			})
		})

		Describe("LastUsed validation", func() {
			It("should pass when last_used is nil (optional)", func() {
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
					LastUsed: nil,
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass when last_used is in the past", func() {
				pastDate := time.Now().AddDate(0, 0, -10)
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
					LastUsed: &pastDate,
				}
				err := skill.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail when last_used is in the future", func() {
				futureDate := time.Now().AddDate(0, 0, 10)
				skill := &Skill{
					Name:     "Ruby",
					Category: "backend",
					LastUsed: &futureDate,
				}
				err := skill.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("last_used cannot be in the future"))
			})
		})
	})
})
