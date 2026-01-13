package technology_test

import (
	"github.com/baphled/kariya/internal/service/career/technology"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analyzer", func() {
	var analyzer *technology.Analyzer

	BeforeEach(func() {
		analyzer = &technology.Analyzer{}
	})

	Describe("AnalyzeSkills", func() {
		Context("when 70%+ skills are backend", func() {
			It("should suggest Backend focus area", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "Ruby", Category: "backend", EventCount: 10},
					{ID: "2", Name: "Go", Category: "backend", EventCount: 8},
					{ID: "3", Name: "PostgreSQL", Category: "database", EventCount: 7},
					{ID: "4", Name: "Redis", Category: "database", EventCount: 5},
					{ID: "5", Name: "Docker", Category: "devops", EventCount: 3},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Area).To(Equal(technology.FocusAreaBackend))
				Expect(suggestion.Confidence).To(BeNumerically(">", 0.5))
				Expect(suggestion.Evidence).To(HaveKey("backend"))
				Expect(suggestion.Evidence["backend"]).To(Equal(2))
			})
		})

		Context("when 70%+ skills are frontend", func() {
			It("should suggest Frontend focus area", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "React", Category: "frontend", EventCount: 10},
					{ID: "2", Name: "TypeScript", Category: "frontend", EventCount: 8},
					{ID: "3", Name: "Vue", Category: "frontend", EventCount: 7},
					{ID: "4", Name: "CSS", Category: "frontend", EventCount: 5},
					{ID: "5", Name: "Webpack", Category: "tooling", EventCount: 3},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Area).To(Equal(technology.FocusAreaFrontend))
				Expect(suggestion.Confidence).To(BeNumerically(">", 0.7))
				Expect(suggestion.Evidence["frontend"]).To(Equal(4))
			})
		})

		Context("when 70%+ skills are devops", func() {
			It("should suggest DevOps focus area", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "Kubernetes", Category: "devops", EventCount: 10},
					{ID: "2", Name: "Docker", Category: "devops", EventCount: 9},
					{ID: "3", Name: "Terraform", Category: "devops", EventCount: 8},
					{ID: "4", Name: "Jenkins", Category: "devops", EventCount: 7},
					{ID: "5", Name: "Ansible", Category: "devops", EventCount: 6},
					{ID: "6", Name: "PostgreSQL", Category: "database", EventCount: 3},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Area).To(Equal(technology.FocusAreaDevOps))
				Expect(suggestion.Confidence).To(BeNumerically(">", 0.7))
				Expect(suggestion.Evidence["devops"]).To(Equal(5))
			})
		})

		Context("when skills are mixed between backend and frontend", func() {
			It("should suggest Fullstack focus area", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "Ruby", Category: "backend", EventCount: 10},
					{ID: "2", Name: "React", Category: "frontend", EventCount: 9},
					{ID: "3", Name: "PostgreSQL", Category: "database", EventCount: 7},
					{ID: "4", Name: "TypeScript", Category: "frontend", EventCount: 6},
					{ID: "5", Name: "Node.js", Category: "backend", EventCount: 5},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Area).To(Equal(technology.FocusAreaFullstack))
				Expect(suggestion.Evidence["backend"]).To(Equal(2))
				Expect(suggestion.Evidence["frontend"]).To(Equal(2))
			})
		})

		Context("when analyzing diverse skill set", func() {
			It("should include all category counts in evidence", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "Ruby", Category: "backend", EventCount: 10},
					{ID: "2", Name: "React", Category: "frontend", EventCount: 8},
					{ID: "3", Name: "Docker", Category: "devops", EventCount: 6},
					{ID: "4", Name: "PostgreSQL", Category: "database", EventCount: 4},
					{ID: "5", Name: "AWS", Category: "cloud", EventCount: 3},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Evidence).To(HaveLen(5))
				Expect(suggestion.Evidence["backend"]).To(Equal(1))
				Expect(suggestion.Evidence["frontend"]).To(Equal(1))
				Expect(suggestion.Evidence["devops"]).To(Equal(1))
				Expect(suggestion.Evidence["database"]).To(Equal(1))
				Expect(suggestion.Evidence["cloud"]).To(Equal(1))
			})
		})

		Context("when confidence is calculated", func() {
			It("should be high when one category dominates (>80%)", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "Ruby", Category: "backend", EventCount: 10},
					{ID: "2", Name: "Go", Category: "backend", EventCount: 9},
					{ID: "3", Name: "Python", Category: "backend", EventCount: 8},
					{ID: "4", Name: "Java", Category: "backend", EventCount: 7},
					{ID: "5", Name: "Docker", Category: "devops", EventCount: 1},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion.Confidence).To(BeNumerically(">", 0.8))
				Expect(suggestion.Area).To(Equal(technology.FocusAreaBackend))
			})

			It("should be moderate when balanced (50-70%)", func() {
				techs := []*technology.ExtractedTechnology{
					{ID: "1", Name: "Ruby", Category: "backend", EventCount: 10},
					{ID: "2", Name: "React", Category: "frontend", EventCount: 9},
					{ID: "3", Name: "PostgreSQL", Category: "database", EventCount: 8},
				}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion.Confidence).To(BeNumerically(">=", 0.3))
				Expect(suggestion.Confidence).To(BeNumerically("<=", 0.7))
			})
		})

		Context("when skills list is empty", func() {
			It("should return Backend as default with zero confidence", func() {
				techs := []*technology.ExtractedTechnology{}

				suggestion := analyzer.AnalyzeSkills(techs)

				Expect(suggestion).NotTo(BeNil())
				Expect(suggestion.Area).To(Equal(technology.FocusAreaBackend))
				Expect(suggestion.Confidence).To(Equal(0.0))
				Expect(suggestion.Evidence).To(BeEmpty())
			})
		})
	})
})
