package skills

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewSkillFromInput", func() {
	Context("with valid input", func() {
		It("creates a skill with all fields", func() {
			input := SkillInput{
				Name:      "Go",
				Category:  "backend",
				Level:     "advanced",
				YearsUsed: "10",
			}
			skill, err := NewSkillFromInput(input, "skill-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill).NotTo(BeNil())
			Expect(skill.ID).To(Equal("skill-1"))
			Expect(skill.Name).To(Equal("Go"))
			Expect(skill.Category).To(Equal("backend"))
			Expect(skill.Level).To(Equal("advanced"))
			Expect(skill.YearsUsed).NotTo(BeNil())
			Expect(*skill.YearsUsed).To(Equal(10))
		})

		It("creates a skill with empty YearsUsed", func() {
			input := SkillInput{
				Name:      "Python",
				Category:  "backend",
				Level:     "intermediate",
				YearsUsed: "",
			}
			skill, err := NewSkillFromInput(input, "skill-2")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.YearsUsed).To(BeNil())
		})

		It("creates a skill with zero years", func() {
			input := SkillInput{
				Name:      "Rust",
				Category:  "backend",
				Level:     "beginner",
				YearsUsed: "0",
			}
			skill, err := NewSkillFromInput(input, "skill-3")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.YearsUsed).NotTo(BeNil())
			Expect(*skill.YearsUsed).To(Equal(0))
		})

		It("creates a skill at boundary 50 years", func() {
			input := SkillInput{
				Name:      "COBOL",
				Category:  "backend",
				Level:     "expert",
				YearsUsed: "50",
			}
			skill, err := NewSkillFromInput(input, "skill-4")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.YearsUsed).NotTo(BeNil())
			Expect(*skill.YearsUsed).To(Equal(50))
		})

		It("creates a skill with empty skillID", func() {
			input := SkillInput{
				Name:      "JavaScript",
				Category:  "frontend",
				Level:     "advanced",
				YearsUsed: "8",
			}
			skill, err := NewSkillFromInput(input, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).To(Equal(""))
		})

		It("creates a skill with empty level", func() {
			input := SkillInput{
				Name:      "Docker",
				Category:  "devops",
				Level:     "",
				YearsUsed: "5",
			}
			skill, err := NewSkillFromInput(input, "skill-5")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Level).To(Equal(""))
		})

		It("trims whitespace from all fields", func() {
			input := SkillInput{
				Name:      "  Go  ",
				Category:  "  backend  ",
				Level:     "  advanced  ",
				YearsUsed: "  10  ",
			}
			skill, err := NewSkillFromInput(input, "skill-6")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Go"))
			Expect(skill.Category).To(Equal("backend"))
			Expect(skill.Level).To(Equal("advanced"))
			Expect(skill.YearsUsed).NotTo(BeNil())
			Expect(*skill.YearsUsed).To(Equal(10))
		})
	})

	Context("with invalid years", func() {
		It("returns YearsParseError for non-numeric years", func() {
			input := SkillInput{
				Name:      "Go",
				Category:  "backend",
				Level:     "advanced",
				YearsUsed: "abc",
			}
			_, err := NewSkillFromInput(input, "skill-7")
			Expect(err).To(HaveOccurred())
			var parseErr *YearsParseError
			Expect(errors.As(err, &parseErr)).To(BeTrue())
		})

		It("returns YearsParseError for decimal years", func() {
			input := SkillInput{
				Name:      "Python",
				Category:  "backend",
				Level:     "intermediate",
				YearsUsed: "3.5",
			}
			_, err := NewSkillFromInput(input, "skill-8")
			Expect(err).To(HaveOccurred())
			var parseErr *YearsParseError
			Expect(errors.As(err, &parseErr)).To(BeTrue())
		})
	})

	Context("with invalid fields that fail Validate", func() {
		It("returns error for empty name", func() {
			input := SkillInput{
				Name:      "",
				Category:  "backend",
				Level:     "advanced",
				YearsUsed: "5",
			}
			_, err := NewSkillFromInput(input, "skill-9")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
		})

		It("returns error for invalid category", func() {
			input := SkillInput{
				Name:      "Go",
				Category:  "invalid",
				Level:     "advanced",
				YearsUsed: "5",
			}
			_, err := NewSkillFromInput(input, "skill-10")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("category must be one of"))
		})

		It("returns error for invalid level", func() {
			input := SkillInput{
				Name:      "Go",
				Category:  "backend",
				Level:     "master",
				YearsUsed: "5",
			}
			_, err := NewSkillFromInput(input, "skill-11")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("level must be one of"))
		})

		It("returns error for years below range", func() {
			input := SkillInput{
				Name:      "Go",
				Category:  "backend",
				Level:     "advanced",
				YearsUsed: "-1",
			}
			_, err := NewSkillFromInput(input, "skill-12")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("years_used must be between 0 and 50"))
		})

		It("returns error for years above range", func() {
			input := SkillInput{
				Name:      "Go",
				Category:  "backend",
				Level:     "advanced",
				YearsUsed: "51",
			}
			_, err := NewSkillFromInput(input, "skill-13")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("years_used must be between 0 and 50"))
		})
	})
})

var _ = Describe("YearsParseError", func() {
	Context("Error method", func() {
		It("returns formatted error message with value", func() {
			err := &YearsParseError{Value: "abc"}
			Expect(err.Error()).To(Equal("unable to parse years: abc"))
		})

		It("returns formatted error message with empty value", func() {
			err := &YearsParseError{Value: ""}
			Expect(err.Error()).To(Equal("unable to parse years: "))
		})

		It("returns formatted error message with numeric value", func() {
			err := &YearsParseError{Value: "3.5"}
			Expect(err.Error()).To(Equal("unable to parse years: 3.5"))
		})
	})

	Context("Unwrap method", func() {
		It("returns the wrapped error", func() {
			wrappedErr := errors.New("strconv.Atoi: parsing \"abc\": invalid syntax")
			err := &YearsParseError{Value: "abc", Err: wrappedErr}
			Expect(err.Unwrap()).To(Equal(wrappedErr))
		})

		It("returns nil when no error is wrapped", func() {
			err := &YearsParseError{Value: "abc", Err: nil}
			Expect(err.Unwrap()).To(Succeed())
		})
	})
})
