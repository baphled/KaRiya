package technology_test

import (
	"strings"

	"github.com/baphled/kariya/internal/constants"
	"github.com/baphled/kariya/internal/service/career/technology"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keywords", func() {
	Describe("Dictionary Structure", func() {
		It("should have at least 350 keyword entries", func() {
			Expect(len(technology.TechnologyKeywords)).To(BeNumerically(">=", 350))
		})

		It("should only use canonical categories from constants.AllSkillCategories", func() {
			allCategories := constants.AllSkillCategories()
			validCategories := make(map[string]bool, len(allCategories))
			for _, cat := range allCategories {
				validCategories[string(cat)] = true
			}

			for _, kw := range technology.TechnologyKeywords {
				Expect(validCategories).To(HaveKey(kw.Category),
					"keyword %q uses invalid category %q", kw.Keyword, kw.Category)
			}
		})

		It("should have no duplicate keywords", func() {
			seen := make(map[string]bool)
			for _, kw := range technology.TechnologyKeywords {
				Expect(seen).NotTo(HaveKey(kw.Keyword),
					"duplicate keyword: %q", kw.Keyword)
				seen[kw.Keyword] = true
			}
		})

		It("should have all keywords in lowercase", func() {
			for _, kw := range technology.TechnologyKeywords {
				Expect(kw.Keyword).To(Equal(strings.ToLower(kw.Keyword)),
					"keyword %q should be lowercase", kw.Keyword)
			}
		})

		It("should have non-empty skill names", func() {
			for _, kw := range technology.TechnologyKeywords {
				Expect(kw.Skill).NotTo(BeEmpty(),
					"keyword %q has empty skill name", kw.Keyword)
			}
		})
	})

	Describe("Keyword Coverage", func() {
		var categoryCounts map[string]int

		BeforeEach(func() {
			categoryCounts = make(map[string]int)
			for _, kw := range technology.TechnologyKeywords {
				categoryCounts[kw.Category]++
			}
		})

		It("should have architecture keywords", func() {
			Expect(categoryCounts["architecture"]).To(BeNumerically(">=", 25))
		})

		It("should have security keywords", func() {
			Expect(categoryCounts["security"]).To(BeNumerically(">=", 12))
		})

		It("should have practices keywords", func() {
			Expect(categoryCounts["practices"]).To(BeNumerically(">=", 25))
		})

		It("should have backend keywords", func() {
			Expect(categoryCounts["backend"]).To(BeNumerically(">=", 30))
		})

		It("should have frontend keywords", func() {
			Expect(categoryCounts["frontend"]).To(BeNumerically(">=", 20))
		})

		It("should have devops keywords", func() {
			Expect(categoryCounts["devops"]).To(BeNumerically(">=", 20))
		})

		It("should have database keywords", func() {
			Expect(categoryCounts["database"]).To(BeNumerically(">=", 14))
		})

		It("should have testing keywords", func() {
			Expect(categoryCounts["testing"]).To(BeNumerically(">=", 15))
		})

		It("should have monitoring keywords", func() {
			Expect(categoryCounts["monitoring"]).To(BeNumerically(">=", 10))
		})

		It("should have data keywords", func() {
			Expect(categoryCounts["data"]).To(BeNumerically(">=", 7))
		})

		It("should have ml keywords", func() {
			Expect(categoryCounts["ml"]).To(BeNumerically(">=", 8))
		})

		It("should have tooling keywords", func() {
			Expect(categoryCounts["tooling"]).To(BeNumerically(">=", 30))
		})
	})

	Describe("GetCategoryForSkillName", func() {
		It("should return architecture for architectural concepts", func() {
			Expect(technology.GetCategoryForSkillName("Microservices")).To(Equal("architecture"))
		})

		It("should return security for security concepts", func() {
			Expect(technology.GetCategoryForSkillName("Authentication")).To(Equal("security"))
		})

		It("should return practices for engineering practices", func() {
			Expect(technology.GetCategoryForSkillName("TDD")).To(Equal("practices"))
		})

		It("should return security for OAuth (moved from tooling)", func() {
			Expect(technology.GetCategoryForSkillName("OAuth")).To(Equal("security"))
		})

		It("should match multi-word keyword phrases", func() {
			Expect(technology.GetCategoryForSkillName("domain-driven design")).To(Equal("architecture"))
		})

		It("should be case-insensitive", func() {
			Expect(technology.GetCategoryForSkillName("DOCKER")).To(Equal("devops"))
			Expect(technology.GetCategoryForSkillName("react")).To(Equal("frontend"))
		})

		It("should return empty string for unknown skills", func() {
			Expect(technology.GetCategoryForSkillName("UnknownTech12345")).To(BeEmpty())
		})

		It("should return backend for Go", func() {
			Expect(technology.GetCategoryForSkillName("Go")).To(Equal("backend"))
		})

		It("should return backend for C", func() {
			Expect(technology.GetCategoryForSkillName("C")).To(Equal("backend"))
		})

		It("should return database for PostgreSQL", func() {
			Expect(technology.GetCategoryForSkillName("PostgreSQL")).To(Equal("database"))
		})

		It("should return cloud for AWS", func() {
			Expect(technology.GetCategoryForSkillName("AWS")).To(Equal("cloud"))
		})
	})
})
