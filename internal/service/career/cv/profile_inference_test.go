package cv_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

var _ = Describe("ProfileInferenceService", func() {
	var (
		service *cv.ProfileInferenceService
		events  []*career.Event
		facts   []*career.Fact
		skills  []*career.Skill
	)

	BeforeEach(func() {
		service = cv.NewProfileInferenceService()
		events = nil
		facts = nil
		skills = nil
	})

	Describe("InferCoreStrengths", func() {
		Context("when events have technical categories", func() {
			BeforeEach(func() {
				events = []*career.Event{
					{
						ID:         "ev-1",
						Text:       "Led migration of monolith to microservices",
						Date:       time.Now().AddDate(-1, 0, 0),
						Categories: []string{"technical"},
						Tags:       []string{"architecture", "leadership"},
					},
					{
						ID:         "ev-2",
						Text:       "Designed and implemented distributed caching layer",
						Date:       time.Now().AddDate(0, -6, 0),
						Categories: []string{"technical"},
						Tags:       []string{"technical"},
					},
				}
			})

			It("should infer technical-related core strengths", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				Expect(strengths).To(ContainElement(ContainSubstring("technical")))
			})
		})

		Context("when events have leadership categories", func() {
			BeforeEach(func() {
				events = []*career.Event{
					{
						ID:         "ev-1",
						Text:       "Mentored 5 junior engineers to senior level",
						Date:       time.Now().AddDate(-1, 0, 0),
						Categories: []string{"leadership", "mentoring"},
						Tags:       []string{"leadership", "mentoring"},
					},
					{
						ID:         "ev-2",
						Text:       "Led cross-functional team of 8 engineers",
						Date:       time.Now().AddDate(0, -3, 0),
						Categories: []string{"leadership"},
						Tags:       []string{"leadership"},
					},
				}
			})

			It("should infer leadership-related core strengths", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				Expect(strengths).To(ContainElement(ContainSubstring("leadership")))
			})
		})

		Context("when facts indicate consulting competencies", func() {
			BeforeEach(func() {
				facts = []*career.Fact{
					{
						ID:                   "fact-1",
						Text:                 "Delivered architecture review for enterprise client",
						CompetencyCategories: []string{"consulting"},
						RoleFit:              career.RoleFitStaff,
						AudienceRelevance:    []string{"hiring_manager"},
						SourceEventID:        "ev-1",
					},
					{
						ID:                   "fact-2",
						Text:                 "Led technical due diligence for acquisition",
						CompetencyCategories: []string{"consulting", "technical"},
						RoleFit:              career.RoleFitPrincipal,
						AudienceRelevance:    []string{"hiring_manager"},
						SourceEventID:        "ev-2",
					},
				}
			})

			It("should infer consulting-related core strengths", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				Expect(strengths).To(ContainElement(ContainSubstring("consulting")))
			})
		})

		Context("when skills indicate technology expertise", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					{ID: "skill-1", Name: "Go", Category: "backend", Level: "expert"},
					{ID: "skill-2", Name: "Kubernetes", Category: "devops", Level: "advanced"},
					{ID: "skill-3", Name: "PostgreSQL", Category: "database", Level: "advanced"},
				}
			})

			It("should infer backend and systems-related core strengths", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				// Should reflect backend expertise (case-insensitive check)
				hasBackendStrength := false
				for _, s := range strengths {
					if ContainsAny(s, []string{"backend", "Backend"}) {
						hasBackendStrength = true
						break
					}
				}
				Expect(hasBackendStrength).To(BeTrue(), "Should have backend-related strength")
			})
		})

		Context("when no career data is provided", func() {
			It("should return generic core strengths", func() {
				strengths := service.InferCoreStrengths(nil, nil, nil)
				Expect(strengths).NotTo(BeEmpty())
				// Should return sensible defaults, not hardcoded personal data
				Expect(strengths).NotTo(ContainElement(ContainSubstring("Yomi")))
			})
		})

		Context("when data is empty", func() {
			It("should return generic core strengths", func() {
				strengths := service.InferCoreStrengths([]*career.Event{}, []*career.Fact{}, []*career.Skill{})
				Expect(strengths).NotTo(BeEmpty())
			})
		})
	})

	Describe("InferValuePropositions", func() {
		Context("when events show cross-functional work", func() {
			BeforeEach(func() {
				events = []*career.Event{
					{
						ID:         "ev-1",
						Text:       "Collaborated with product, design, and engineering teams",
						Date:       time.Now().AddDate(-1, 0, 0),
						Categories: []string{"leadership", "product"},
						Tags:       []string{"leadership"},
					},
				}
			})

			It("should infer collaboration-related value propositions", func() {
				propositions := service.InferValuePropositions(events, facts, skills)
				Expect(propositions).NotTo(BeEmpty())
				Expect(propositions).To(ContainElement(ContainSubstring("collaboration")))
			})
		})

		Context("when facts show mentoring competencies", func() {
			BeforeEach(func() {
				facts = []*career.Fact{
					{
						ID:                   "fact-1",
						Text:                 "Established engineering onboarding program",
						CompetencyCategories: []string{"mentoring"},
						RoleFit:              career.RoleFitStaff,
						AudienceRelevance:    []string{"hiring_manager"},
						SourceEventID:        "ev-1",
					},
				}
			})

			It("should infer mentoring-related value propositions", func() {
				propositions := service.InferValuePropositions(events, facts, skills)
				Expect(propositions).NotTo(BeEmpty())
				Expect(propositions).To(ContainElement(ContainSubstring("mentor")))
			})
		})

		Context("when skills show diverse technology stack", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					{ID: "skill-1", Name: "Ruby", Category: "backend", Level: "expert"},
					{ID: "skill-2", Name: "Go", Category: "backend", Level: "advanced"},
					{ID: "skill-3", Name: "Python", Category: "backend", Level: "intermediate"},
					{ID: "skill-4", Name: "React", Category: "frontend", Level: "intermediate"},
				}
			})

			It("should infer language-agnostic value proposition", func() {
				propositions := service.InferValuePropositions(events, facts, skills)
				Expect(propositions).NotTo(BeEmpty())
				// Should reflect ability to work with multiple languages
				hasLanguageFlexibility := false
				for _, p := range propositions {
					if ContainsAny(p, []string{"language", "adaptable", "versatile", "polyglot"}) {
						hasLanguageFlexibility = true
						break
					}
				}
				Expect(hasLanguageFlexibility).To(BeTrue(), "Should have language flexibility proposition")
			})
		})

		Context("when no career data is provided", func() {
			It("should return generic value propositions", func() {
				propositions := service.InferValuePropositions(nil, nil, nil)
				Expect(propositions).NotTo(BeEmpty())
				// Should NOT contain hardcoded personal data
				for _, p := range propositions {
					Expect(p).NotTo(ContainSubstring("Yomi"))
				}
			})
		})
	})

	Describe("InferTechnologies", func() {
		Context("when skills are provided", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					{ID: "skill-1", Name: "Go", Category: "backend"},
					{ID: "skill-2", Name: "Ruby", Category: "backend"},
					{ID: "skill-3", Name: "React", Category: "frontend"},
					{ID: "skill-4", Name: "PostgreSQL", Category: "database"},
				}
			})

			It("should extract languages from skills", func() {
				result := service.InferTechnologies(events, skills)
				Expect(result.Languages).To(ContainElement("Go"))
				Expect(result.Languages).To(ContainElement("Ruby"))
			})

			It("should extract frontend from skills", func() {
				result := service.InferTechnologies(events, skills)
				Expect(result.Frontend).To(ContainElement("React"))
			})

			It("should extract systems/database from skills", func() {
				result := service.InferTechnologies(events, skills)
				Expect(result.Systems).To(ContainElement("PostgreSQL"))
			})
		})

		Context("when no skills are provided", func() {
			It("should return empty technology slices", func() {
				result := service.InferTechnologies(nil, nil)
				Expect(result.Languages).To(BeEmpty())
				Expect(result.Frontend).To(BeEmpty())
				Expect(result.Systems).To(BeEmpty())
			})
		})
	})
})

// ContainsAny checks if the string contains any of the substrings (case-insensitive).
func ContainsAny(s string, substrings []string) bool {
	lower := strings.ToLower(s)
	for _, sub := range substrings {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
