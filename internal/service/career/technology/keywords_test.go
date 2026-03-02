package technology_test

import (
	"context"
	"strings"

	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/technology"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keywords", func() {
	Describe("Dictionary Structure", func() {
		It("should have at least 480 keyword entries", func() {
			Expect(len(technology.Keywords)).To(BeNumerically(">=", 480))
		})

		It("should only use canonical categories from constants.AllSkillCategories", func() {
			allCategories := constants.AllSkillCategories()
			validCategories := make(map[string]bool, len(allCategories))
			for _, cat := range allCategories {
				validCategories[string(cat)] = true
			}

			for _, kw := range technology.Keywords {
				Expect(validCategories).To(HaveKey(kw.Category),
					"keyword %q uses invalid category %q", kw.Keyword, kw.Category)
			}
		})

		It("should have no duplicate keywords", func() {
			seen := make(map[string]bool)
			for _, kw := range technology.Keywords {
				Expect(seen).NotTo(HaveKey(kw.Keyword),
					"duplicate keyword: %q", kw.Keyword)
				seen[kw.Keyword] = true
			}
		})

		It("should have all keywords in lowercase", func() {
			for _, kw := range technology.Keywords {
				Expect(kw.Keyword).To(Equal(strings.ToLower(kw.Keyword)),
					"keyword %q should be lowercase", kw.Keyword)
			}
		})

		It("should have non-empty skill names", func() {
			for _, kw := range technology.Keywords {
				Expect(kw.Skill).NotTo(BeEmpty(),
					"keyword %q has empty skill name", kw.Keyword)
			}
		})
	})

	Describe("Keyword Coverage", func() {
		var categoryCounts map[string]int

		BeforeEach(func() {
			categoryCounts = make(map[string]int)
			for _, kw := range technology.Keywords {
				categoryCounts[kw.Category]++
			}
		})

		It("should have architecture keywords", func() {
			Expect(categoryCounts["architecture"]).To(BeNumerically(">=", 55))
		})

		It("should have security keywords", func() {
			Expect(categoryCounts["security"]).To(BeNumerically(">=", 12))
		})

		It("should have practices keywords", func() {
			Expect(categoryCounts["practices"]).To(BeNumerically(">=", 80))
		})

		It("should have backend keywords", func() {
			Expect(categoryCounts["backend"]).To(BeNumerically(">=", 68))
		})

		It("should have frontend keywords", func() {
			Expect(categoryCounts["frontend"]).To(BeNumerically(">=", 50))
		})

		It("should have devops keywords", func() {
			Expect(categoryCounts["devops"]).To(BeNumerically(">=", 38))
		})

		It("should have database keywords", func() {
			Expect(categoryCounts["database"]).To(BeNumerically(">=", 15))
		})

		It("should have testing keywords", func() {
			Expect(categoryCounts["testing"]).To(BeNumerically(">=", 16))
		})

		It("should have monitoring keywords", func() {
			Expect(categoryCounts["monitoring"]).To(BeNumerically(">=", 16))
		})

		It("should have data keywords", func() {
			Expect(categoryCounts["data"]).To(BeNumerically(">=", 12))
		})

		It("should have ml keywords", func() {
			Expect(categoryCounts["ml"]).To(BeNumerically(">=", 8))
		})

		It("should have tooling keywords", func() {
			Expect(categoryCounts["tooling"]).To(BeNumerically(">=", 34))
		})
	})

	Describe("RecategorizeSkills", func() {
		var skillRepo careerRepo.SkillRepository

		BeforeEach(func() {
			skillRepo = careermemory.NewSkillRepository()
		})

		It("should return zero changes for an empty repository", func() {
			result, err := technology.RecategorizeSkills(context.Background(), skillRepo)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(0))
			Expect(result.Updated).To(Equal(0))
		})

		It("should recategorize skills with matching keywords", func() {
			ctx := context.Background()
			skill := fixtures.SkillWith("s1", "Docker", "other", "intermediate")
			Expect(skillRepo.Create(ctx, skill)).To(Succeed())

			result, err := technology.RecategorizeSkills(ctx, skillRepo)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Updated).To(Equal(1))
			Expect(result.ByCategory).To(HaveKeyWithValue("devops", 1))

			updated, err := skillRepo.GetByID(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Category).To(Equal("devops"))
		})

		It("should skip skills that already have the correct category", func() {
			ctx := context.Background()
			skill := fixtures.SkillWith("s1", "Docker", "devops", "intermediate")
			Expect(skillRepo.Create(ctx, skill)).To(Succeed())

			result, err := technology.RecategorizeSkills(ctx, skillRepo)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Updated).To(Equal(0))
			Expect(result.Skipped).To(Equal(1))
		})

		It("should skip skills with no keyword match", func() {
			ctx := context.Background()
			skill := fixtures.SkillWith("s1", "UnknownTechnology12345", "other", "intermediate")
			Expect(skillRepo.Create(ctx, skill)).To(Succeed())

			result, err := technology.RecategorizeSkills(ctx, skillRepo)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Updated).To(Equal(0))
			Expect(result.NoMatch).To(Equal(1))
		})

		It("should handle multiple skills across categories", func() {
			ctx := context.Background()
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "React", "other", "advanced"),
				fixtures.SkillWith("s2", "PostgreSQL", "other", "advanced"),
				fixtures.SkillWith("s3", "TDD", "other", "advanced"),
				fixtures.SkillWith("s4", "Go", "backend", "expert"),
			}
			for _, s := range skills {
				Expect(skillRepo.Create(ctx, s)).To(Succeed())
			}

			result, err := technology.RecategorizeSkills(ctx, skillRepo)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Total).To(Equal(4))
			Expect(result.Updated).To(Equal(3))
			Expect(result.Skipped).To(Equal(1))
			Expect(result.ByCategory).To(HaveKeyWithValue("frontend", 1))
			Expect(result.ByCategory).To(HaveKeyWithValue("database", 1))
			Expect(result.ByCategory).To(HaveKeyWithValue("practices", 1))
		})

		It("should move OAuth from tooling to security", func() {
			ctx := context.Background()
			skill := fixtures.SkillWith("s1", "OAuth", "tooling", "intermediate")
			Expect(skillRepo.Create(ctx, skill)).To(Succeed())

			result, err := technology.RecategorizeSkills(ctx, skillRepo)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Updated).To(Equal(1))

			updated, err := skillRepo.GetByID(ctx, skill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Category).To(Equal("security"))
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

		Context("compound phrase lookups", func() {
			It("should categorize architecture compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Backend Architecture")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("Service Architecture")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("Component Architecture")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("Real-time Systems")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("Scalable Systems")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("Repository Pattern")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("Builder Pattern")).To(Equal("architecture"))
				Expect(technology.GetCategoryForSkillName("ADR")).To(Equal("architecture"))
			})

			It("should categorize practices compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Agile Delivery")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Project Management")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Software Engineering")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Engineering Practices")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Code Standards")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Technical Leadership")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Knowledge Transfer")).To(Equal("practices"))
				Expect(technology.GetCategoryForSkillName("Collaboration")).To(Equal("practices"))
			})

			It("should categorize backend compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Backend Development")).To(Equal("backend"))
				Expect(technology.GetCategoryForSkillName("Full-stack Development")).To(Equal("backend"))
				Expect(technology.GetCategoryForSkillName("Web Development")).To(Equal("backend"))
				Expect(technology.GetCategoryForSkillName("API Development")).To(Equal("backend"))
				Expect(technology.GetCategoryForSkillName("CRUD")).To(Equal("backend"))
				Expect(technology.GetCategoryForSkillName("Error Handling")).To(Equal("backend"))
			})

			It("should categorize frontend compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Component Library")).To(Equal("frontend"))
				Expect(technology.GetCategoryForSkillName("Keyboard Navigation")).To(Equal("frontend"))
				Expect(technology.GetCategoryForSkillName("Layout Design")).To(Equal("frontend"))
				Expect(technology.GetCategoryForSkillName("UX")).To(Equal("frontend"))
				Expect(technology.GetCategoryForSkillName("User Experience")).To(Equal("frontend"))
				Expect(technology.GetCategoryForSkillName("Bubble Tea")).To(Equal("frontend"))
				Expect(technology.GetCategoryForSkillName("Vite")).To(Equal("frontend"))
			})

			It("should categorize devops compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Incident Response")).To(Equal("devops"))
				Expect(technology.GetCategoryForSkillName("Production Support")).To(Equal("devops"))
				Expect(technology.GetCategoryForSkillName("Embedded Systems")).To(Equal("devops"))
				Expect(technology.GetCategoryForSkillName("Reliability Engineering")).To(Equal("devops"))
				Expect(technology.GetCategoryForSkillName("Firmware Development")).To(Equal("devops"))
			})

			It("should categorize monitoring compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Kibana")).To(Equal("monitoring"))
				Expect(technology.GetCategoryForSkillName("Logstash")).To(Equal("monitoring"))
			})

			It("should categorize database compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("TimescaleDB")).To(Equal("database"))
			})

			It("should categorize tooling compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Technical Writing")).To(Equal("tooling"))
				Expect(technology.GetCategoryForSkillName("Diagramming")).To(Equal("tooling"))
			})

			It("should categorize data compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Streaming")).To(Equal("data"))
			})

			It("should categorize testing compound phrases", func() {
				Expect(technology.GetCategoryForSkillName("Playwright")).To(Equal("testing"))
			})
		})

		Context("substring matching fallback", func() {
			It("should match when keyword is a substring of skill name", func() {
				Expect(technology.GetCategoryForSkillName("ELK Stack")).To(Equal("monitoring"))
			})

			It("should match plural/suffix variants", func() {
				Expect(technology.GetCategoryForSkillName("Code Reviews")).To(Equal("practices"))
			})

			It("should match keyword prefix in compound names", func() {
				Expect(technology.GetCategoryForSkillName("Zend Framework")).To(Equal("backend"))
			})

			It("should match compound skill names with existing keywords", func() {
				Expect(technology.GetCategoryForSkillName("REST API")).To(Equal("backend"))
			})

			It("should match hyphenated compound names", func() {
				Expect(technology.GetCategoryForSkillName("Event-driven Architecture")).To(Equal("architecture"))
			})

			It("should not false-match short keywords inside unrelated words", func() {
				Expect(technology.GetCategoryForSkillName("Communication")).To(BeEmpty())
				Expect(technology.GetCategoryForSkillName("Logistics Systems")).To(BeEmpty())
				Expect(technology.GetCategoryForSkillName("Career Development")).To(BeEmpty())
			})

			It("should still prefer exact match over substring match", func() {
				Expect(technology.GetCategoryForSkillName("Docker")).To(Equal("devops"))
			})

			It("should match keyword at second occurrence when first fails boundary", func() {
				Expect(technology.GetCategoryForSkillName("norest rest")).To(Equal("backend"))
			})
		})

		Context("short keyword safety for ai", func() {
			It("should not false-match ai inside Maintaining", func() {
				Expect(technology.GetCategoryForSkillName("Maintaining")).To(BeEmpty())
			})

			It("should not false-match ai inside Email", func() {
				Expect(technology.GetCategoryForSkillName("Email")).To(BeEmpty())
			})

			It("should match AI as standalone skill", func() {
				Expect(technology.GetCategoryForSkillName("AI")).To(Equal("ml"))
			})
		})

		Context("short keyword safety for r", func() {
			It("should match R as standalone skill", func() {
				Expect(technology.GetCategoryForSkillName("R")).To(Equal("backend"))
			})

			It("should not false-match r inside Career Development", func() {
				Expect(technology.GetCategoryForSkillName("Career Development")).To(BeEmpty())
			})
		})

		Context("broad keyword priority", func() {
			It("should categorize Database Design as database not practices", func() {
				Expect(technology.GetCategoryForSkillName("Database Design")).NotTo(Equal("practices"))
			})

			It("should categorize Performance Testing as testing not practices", func() {
				Expect(technology.GetCategoryForSkillName("Performance Testing")).To(Equal("testing"))
			})

			It("should categorize Data Processing as data not practices", func() {
				Expect(technology.GetCategoryForSkillName("Data Processing")).To(Equal("data"))
			})
		})

		Context("Swift category", func() {
			It("should categorize Swift as mobile", func() {
				Expect(technology.GetCategoryForSkillName("Swift")).To(Equal("mobile"))
			})
		})
	})
})

var _ = Describe("GetKeywordMap", func() {
	It("should return all keywords as a map keyed by keyword string", func() {
		keywordMap := technology.GetKeywordMap()
		Expect(keywordMap).NotTo(BeEmpty())
		Expect(keywordMap).To(HaveLen(len(technology.Keywords)))
	})

	It("should allow lookup of keyword entries by keyword string", func() {
		keywordMap := technology.GetKeywordMap()
		// Pick a known keyword from the Keywords slice
		if len(technology.Keywords) > 0 {
			expectedEntry := technology.Keywords[0]
			Expect(keywordMap).To(HaveKey(expectedEntry.Keyword))
			Expect(keywordMap[expectedEntry.Keyword]).To(Equal(expectedEntry))
		}
	})
})
