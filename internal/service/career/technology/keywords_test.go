package technology_test

import (
	"github.com/baphled/kariya/internal/service/career/technology"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TechnologyKeywords", func() {
	Describe("Dictionary Structure", func() {
		It("should have at least 200 keyword entries", func() {
			keywords := technology.GetTechnologyKeywords()
			Expect(len(keywords)).To(BeNumerically(">=", 200),
				"Expected at least 200 comprehensive keywords")
		})

		It("should have all keywords in lowercase", func() {
			keywords := technology.GetTechnologyKeywords()
			for _, kw := range keywords {
				// Check each character for uppercase
				for _, ch := range kw.Keyword {
					if ch >= 'A' && ch <= 'Z' {
						Fail("Keyword \"" + kw.Keyword + "\" contains uppercase character")
					}
				}
			}
		})

		It("should have no duplicate keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			seen := make(map[string]bool)

			for _, kw := range keywords {
				if seen[kw.Keyword] {
					Fail("Duplicate keyword found: " + kw.Keyword)
				}
				seen[kw.Keyword] = true
			}
		})

		It("should only use canonical categories from constants.AllSkillCategories", func() {
			keywords := technology.GetTechnologyKeywords()
			validCategories := map[string]bool{
				"backend":    true,
				"frontend":   true,
				"database":   true,
				"devops":     true,
				"cloud":      true,
				"mobile":     true,
				"tooling":    true,
				"testing":    true,
				"ml":         true,
				"data":       true,
				"monitoring": true,
				"other":      true,
			}

			for _, kw := range keywords {
				if !validCategories[kw.Category] {
					Fail("Invalid category \"" + kw.Category + "\" for keyword \"" + kw.Keyword + "\". Must use canonical categories.")
				}
			}
		})

		It("should have non-empty skill names", func() {
			keywords := technology.GetTechnologyKeywords()
			for _, kw := range keywords {
				Expect(kw.Skill).NotTo(BeEmpty(),
					"Keyword %q has empty Skill field", kw.Keyword)
			}
		})

		It("should have non-empty keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			for _, kw := range keywords {
				Expect(kw.Keyword).NotTo(BeEmpty(),
					"Found keyword with empty Keyword field")
			}
		})
	})

	Describe("GetKeywordMap", func() {
		It("should return a map with correct count", func() {
			keywords := technology.GetTechnologyKeywords()
			keywordMap := technology.GetKeywordMap()

			Expect(len(keywordMap)).To(Equal(len(keywords)),
				"Map should have same count as slice")
		})

		It("should allow O(1) keyword lookup", func() {
			keywordMap := technology.GetKeywordMap()

			// Test a few expected keywords
			expectedKeywords := []string{"go", "golang", "python", "react", "docker"}

			for _, keyword := range expectedKeywords {
				_, exists := keywordMap[keyword]
				Expect(exists).To(BeTrue(),
					"Expected keyword %q to exist in map", keyword)
			}
		})

		It("should map keyword to correct TechnologyKeyword struct", func() {
			keywordMap := technology.GetKeywordMap()

			// Test specific mapping
			goLang := keywordMap["go"]
			Expect(goLang.Skill).To(Equal("Go"))
			Expect(goLang.Category).To(Equal("backend"))

			golangAlias := keywordMap["golang"]
			Expect(golangAlias.Skill).To(Equal("Go"))
			Expect(golangAlias.Category).To(Equal("backend"))
		})
	})

	Describe("Keyword Coverage", func() {
		It("should have backend keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			backendCount := 0
			for _, kw := range keywords {
				if kw.Category == "backend" {
					backendCount++
				}
			}
			Expect(backendCount).To(BeNumerically(">=", 20),
				"Expected at least 20 backend keywords, got %d", backendCount)
		})

		It("should have frontend keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			frontendCount := 0
			for _, kw := range keywords {
				if kw.Category == "frontend" {
					frontendCount++
				}
			}
			Expect(frontendCount).To(BeNumerically(">=", 15),
				"Expected at least 15 frontend keywords, got %d", frontendCount)
		})

		It("should have database keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			databaseCount := 0
			for _, kw := range keywords {
				if kw.Category == "database" {
					databaseCount++
				}
			}
			Expect(databaseCount).To(BeNumerically(">=", 10),
				"Expected at least 10 database keywords, got %d", databaseCount)
		})

		It("should have devops keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			devopsCount := 0
			for _, kw := range keywords {
				if kw.Category == "devops" {
					devopsCount++
				}
			}
			Expect(devopsCount).To(BeNumerically(">=", 10),
				"Expected at least 10 devops keywords, got %d", devopsCount)
		})

		It("should have cloud keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			cloudCount := 0
			for _, kw := range keywords {
				if kw.Category == "cloud" {
					cloudCount++
				}
			}
			Expect(cloudCount).To(BeNumerically(">=", 10),
				"Expected at least 10 cloud keywords, got %d", cloudCount)
		})

		It("should have testing keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			count := 0
			for _, kw := range keywords {
				if kw.Category == "testing" {
					count++
				}
			}
			Expect(count).To(BeNumerically(">=", 10),
				"Expected at least 10 testing keywords, got %d", count)
		})

		It("should have tooling keywords including build tools and documentation", func() {
			keywords := technology.GetTechnologyKeywords()
			count := 0
			for _, kw := range keywords {
				if kw.Category == "tooling" {
					count++
				}
			}
			Expect(count).To(BeNumerically(">=", 30),
				"Expected at least 30 tooling keywords (includes build tools and docs), got %d", count)
		})

		It("should have ML/data keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			count := 0
			for _, kw := range keywords {
				if kw.Category == "ml" || kw.Category == "data" {
					count++
				}
			}
			Expect(count).To(BeNumerically(">=", 15),
				"Expected at least 15 ML/data keywords, got %d", count)
		})

		It("should have monitoring keywords", func() {
			keywords := technology.GetTechnologyKeywords()
			count := 0
			for _, kw := range keywords {
				if kw.Category == "monitoring" {
					count++
				}
			}
			Expect(count).To(BeNumerically(">=", 5),
				"Expected at least 5 monitoring keywords, got %d", count)
		})
	})

	Describe("GetCategoryForSkillName", func() {
		It("should return correct category for known canonical skill names", func() {
			Expect(technology.GetCategoryForSkillName("Go")).To(Equal("backend"))
			Expect(technology.GetCategoryForSkillName("PostgreSQL")).To(Equal("database"))
			Expect(technology.GetCategoryForSkillName("Docker")).To(Equal("devops"))
			Expect(technology.GetCategoryForSkillName("React")).To(Equal("frontend"))
			Expect(technology.GetCategoryForSkillName("AWS")).To(Equal("cloud"))
			Expect(technology.GetCategoryForSkillName("Jest")).To(Equal("testing"))
			Expect(technology.GetCategoryForSkillName("Git")).To(Equal("tooling"))
			Expect(technology.GetCategoryForSkillName("TensorFlow")).To(Equal("ml"))
			Expect(technology.GetCategoryForSkillName("Apache Spark")).To(Equal("data"))
			Expect(technology.GetCategoryForSkillName("Splunk")).To(Equal("monitoring"))
			Expect(technology.GetCategoryForSkillName("iOS")).To(Equal("mobile"))
		})

		It("should match case-insensitively", func() {
			Expect(technology.GetCategoryForSkillName("go")).To(Equal("backend"))
			Expect(technology.GetCategoryForSkillName("POSTGRESQL")).To(Equal("database"))
			Expect(technology.GetCategoryForSkillName("docker")).To(Equal("devops"))
			Expect(technology.GetCategoryForSkillName("react")).To(Equal("frontend"))
		})

		It("should match keywords as well as canonical names", func() {
			Expect(technology.GetCategoryForSkillName("golang")).To(Equal("backend"))
			Expect(technology.GetCategoryForSkillName("postgres")).To(Equal("database"))
			Expect(technology.GetCategoryForSkillName("k8s")).To(Equal("devops"))
			Expect(technology.GetCategoryForSkillName("nodejs")).To(Equal("backend"))
		})

		It("should return empty string for unknown skill names", func() {
			Expect(technology.GetCategoryForSkillName("UnknownTechnology")).To(BeEmpty())
			Expect(technology.GetCategoryForSkillName("FooBar")).To(BeEmpty())
			Expect(technology.GetCategoryForSkillName("")).To(BeEmpty())
		})
	})

	Describe("Aliases", func() {
		It("should support golang -> Go alias", func() {
			keywordMap := technology.GetKeywordMap()

			golangKw := keywordMap["golang"]
			goKw := keywordMap["go"]

			Expect(golangKw.Skill).To(Equal("Go"))
			Expect(goKw.Skill).To(Equal("Go"))
			Expect(golangKw.Skill).To(Equal(goKw.Skill))
		})

		It("should support k8s -> Kubernetes alias", func() {
			keywordMap := technology.GetKeywordMap()

			k8s := keywordMap["k8s"]
			kubernetes := keywordMap["kubernetes"]

			Expect(k8s.Skill).To(Equal("Kubernetes"))
			Expect(kubernetes.Skill).To(Equal("Kubernetes"))
			Expect(k8s.Skill).To(Equal(kubernetes.Skill))
		})

		It("should support postgres -> PostgreSQL alias", func() {
			keywordMap := technology.GetKeywordMap()

			postgres := keywordMap["postgres"]
			postgresql := keywordMap["postgresql"]

			Expect(postgres.Skill).To(Equal("PostgreSQL"))
			Expect(postgresql.Skill).To(Equal("PostgreSQL"))
		})
	})
})
