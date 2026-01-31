package forms_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillForm", func() {
	Describe("NewSkillForm", func() {
		Context("creating form for new skill", func() {
			It("should create form with empty fields", func() {
				form := forms.NewSkillForm(nil)
				Expect(form).NotTo(BeNil())
			})
		})

		Context("creating form for existing skill", func() {
			It("should pre-populate fields with skill data", func() {
				skill := fixtures.SkillWith("skill-ruby", "Ruby", "backend", "advanced")
				form := forms.NewSkillForm(skill)
				Expect(form).NotTo(BeNil())
			})
		})
	})

	Describe("NewSkillFormWithData", func() {
		It("should create form with provided data", func() {
			data := &forms.SkillFormData{
				Name:      "Go",
				Category:  "backend",
				Level:     "expert",
				YearsUsed: "5",
			}
			form := forms.NewSkillFormWithData(data)
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("GetSkillFormData", func() {
		It("should extract form data from skill", func() {
			skill := fixtures.SkillWithYears("skill-react", "React", "frontend", 3)
			skill.Level = "intermediate"

			data := forms.GetSkillFormData(skill)
			Expect(data.Name).To(Equal("React"))
			Expect(data.Category).To(Equal("frontend"))
			Expect(data.Level).To(Equal("intermediate"))
			Expect(data.YearsUsed).To(Equal("3"))
		})

		It("should handle nil years", func() {
			skill := fixtures.SkillWith("skill-docker", "Docker", "devops", "")

			data := forms.GetSkillFormData(skill)
			Expect(data.YearsUsed).To(Equal(""))
		})
	})

	Describe("ApplySkillFormData", func() {
		Context("with valid data", func() {
			It("should apply data to new skill", func() {
				skill := fixtures.SkillWith("", "", "", "")
				data := &forms.SkillFormData{
					Name:      "Kubernetes",
					Category:  "devops",
					Level:     "advanced",
					YearsUsed: "2",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.Name).To(Equal("Kubernetes"))
				Expect(skill.Category).To(Equal("devops"))
				Expect(skill.Level).To(Equal("advanced"))
				Expect(skill.YearsUsed).NotTo(BeNil())
				Expect(*skill.YearsUsed).To(Equal(2))
			})

			It("should update existing skill", func() {
				skill := fixtures.SkillWithYears("skill-ruby", "Ruby", "backend", 1)
				data := &forms.SkillFormData{
					Name:      "Ruby",
					Category:  "backend",
					Level:     "expert",
					YearsUsed: "5",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.Level).To(Equal("expert"))
				Expect(*skill.YearsUsed).To(Equal(5))
			})

			It("should handle empty optional fields", func() {
				skill := fixtures.SkillWith("", "", "", "")
				data := &forms.SkillFormData{
					Name:      "Python",
					Category:  "backend",
					Level:     "",
					YearsUsed: "",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.Name).To(Equal("Python"))
				Expect(skill.Category).To(Equal("backend"))
				Expect(skill.Level).To(Equal(""))
				Expect(skill.YearsUsed).To(BeNil())
			})

			It("should trim whitespace from name and category", func() {
				skill := fixtures.SkillWith("", "", "", "")
				data := &forms.SkillFormData{
					Name:     "  PostgreSQL  ",
					Category: "  database  ",
					Level:    "intermediate",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.Name).To(Equal("PostgreSQL"))
				Expect(skill.Category).To(Equal("database"))
			})
		})

		Context("with invalid data", func() {
			It("should ignore invalid years format", func() {
				skill := fixtures.SkillWith("", "", "", "")
				data := &forms.SkillFormData{
					Name:      "Go",
					Category:  "backend",
					YearsUsed: "invalid",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.YearsUsed).To(BeNil())
			})

			It("should ignore negative years", func() {
				skill := fixtures.SkillWith("", "", "", "")
				data := &forms.SkillFormData{
					Name:      "Go",
					Category:  "backend",
					YearsUsed: "-1",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.YearsUsed).To(BeNil())
			})

			It("should ignore years over 50", func() {
				skill := fixtures.SkillWith("", "", "", "")
				data := &forms.SkillFormData{
					Name:      "Go",
					Category:  "backend",
					YearsUsed: "51",
				}

				forms.ApplySkillFormData(skill, data)
				Expect(skill.YearsUsed).To(BeNil())
			})
		})
	})

	Describe("Skill validators", func() {
		Describe("SkillName", func() {
			It("should require non-empty name", func() {
				err := forms.SkillName("")
				Expect(err).To(HaveOccurred())
			})

			It("should reject names that are too short", func() {
				err := forms.SkillName("")
				Expect(err).To(HaveOccurred())
			})

			It("should reject names that are too long", func() {
				longName := string(make([]byte, 101))
				for i := range longName {
					longName = string(append([]byte(longName[:i]), 'a'))
				}
				err := forms.SkillName(longName)
				Expect(err).To(HaveOccurred())
			})

			It("should accept valid names", func() {
				err := forms.SkillName("Ruby")
				Expect(err).NotTo(HaveOccurred())

				err = forms.SkillName("React.js")
				Expect(err).NotTo(HaveOccurred())

				err = forms.SkillName("AWS Lambda")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should trim whitespace", func() {
				err := forms.SkillName("  Go  ")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("SkillCategory", func() {
			It("should require non-empty category", func() {
				err := forms.SkillCategory("")
				Expect(err).To(HaveOccurred())
			})

			It("should reject categories that are too long", func() {
				longCategory := string(make([]byte, 51))
				for i := range longCategory {
					longCategory = string(append([]byte(longCategory[:i]), 'a'))
				}
				err := forms.SkillCategory(longCategory)
				Expect(err).To(HaveOccurred())
			})

			It("should accept valid categories", func() {
				err := forms.SkillCategory("backend")
				Expect(err).NotTo(HaveOccurred())

				err = forms.SkillCategory("frontend")
				Expect(err).NotTo(HaveOccurred())

				err = forms.SkillCategory("custom-category")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("SkillLevel", func() {
			It("should allow empty level", func() {
				err := forms.SkillLevel("")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should accept valid levels", func() {
				validLevels := []string{"beginner", "intermediate", "advanced", "expert"}
				for _, level := range validLevels {
					err := forms.SkillLevel(level)
					Expect(err).NotTo(HaveOccurred(), "Expected %s to be valid", level)
				}
			})

			It("should reject invalid levels", func() {
				err := forms.SkillLevel("master")
				Expect(err).To(HaveOccurred())

				err = forms.SkillLevel("novice")
				Expect(err).To(HaveOccurred())
			})

			It("should be case-sensitive", func() {
				err := forms.SkillLevel("Beginner")
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("SkillYearsUsed", func() {
			It("should allow empty years", func() {
				err := forms.SkillYearsUsed("")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should accept valid years", func() {
				err := forms.SkillYearsUsed("0")
				Expect(err).NotTo(HaveOccurred())

				err = forms.SkillYearsUsed("5")
				Expect(err).NotTo(HaveOccurred())

				err = forms.SkillYearsUsed("50")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject non-numeric values", func() {
				err := forms.SkillYearsUsed("five")
				Expect(err).To(HaveOccurred())

				err = forms.SkillYearsUsed("5.5")
				Expect(err).To(HaveOccurred())
			})

			It("should reject negative values", func() {
				err := forms.SkillYearsUsed("-1")
				Expect(err).To(HaveOccurred())
			})

			It("should reject values over 50", func() {
				err := forms.SkillYearsUsed("51")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
