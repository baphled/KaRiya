package cv_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
				ev1 := fixtures.EventWithCategories("ev-1", "Led migration of monolith to microservices", []string{"technical"})
				ev1.Date = time.Now().AddDate(-1, 0, 0)
				ev1.Tags = []string{"architecture", "leadership"}

				ev2 := fixtures.EventWithCategories("ev-2", "Designed and implemented distributed caching layer", []string{"technical"})
				ev2.Date = time.Now().AddDate(0, -6, 0)
				ev2.Tags = []string{"technical"}

				events = []*career.Event{ev1, ev2}
			})

			It("should infer technical-related core strengths", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				Expect(strengths).To(ContainElement(ContainSubstring("technical")))
			})
		})

		Context("when events have leadership categories", func() {
			BeforeEach(func() {
				ev1 := fixtures.EventWithCategories("ev-1", "Mentored 5 junior engineers to senior level", []string{"leadership", "mentoring"})
				ev1.Date = time.Now().AddDate(-1, 0, 0)
				ev1.Tags = []string{"leadership", "mentoring"}

				ev2 := fixtures.EventWithCategories("ev-2", "Led cross-functional team of 8 engineers", []string{"leadership"})
				ev2.Date = time.Now().AddDate(0, -3, 0)
				ev2.Tags = []string{"leadership"}

				events = []*career.Event{ev1, ev2}
			})

			It("should infer leadership-related core strengths", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				Expect(strengths).To(ContainElement(ContainSubstring("leadership")))
			})
		})

		Context("when facts indicate consulting competencies", func() {
			BeforeEach(func() {
				fact1 := fixtures.FactWithCategories("fact-1", "Delivered architecture review for enterprise client", "ev-1", []string{"consulting"}, []string{"hiring_manager"})

				fact2 := fixtures.FactWithCategories("fact-2", "Led technical due diligence for acquisition", "ev-2", []string{"consulting", "technical"}, []string{"hiring_manager"})
				fact2.RoleFit = career.RoleFitPrincipal

				facts = []*career.Fact{fact1, fact2}
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
					fixtures.SkillWith("skill-1", "Go", "backend", "expert"),
					fixtures.SkillWith("skill-2", "Kubernetes", "devops", "advanced"),
					fixtures.SkillWith("skill-3", "PostgreSQL", "database", "advanced"),
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

		Context("when skills indicate architecture expertise", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					fixtures.SkillWith("skill-1", "Microservices", "architecture", "expert"),
					fixtures.SkillWith("skill-2", "DDD", "architecture", "advanced"),
				}
			})

			It("should infer architecture-related core strength", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				hasArchStrength := false
				for _, s := range strengths {
					if ContainsAny(s, []string{"architecture", "system design"}) {
						hasArchStrength = true
						break
					}
				}
				Expect(hasArchStrength).To(BeTrue(), "Should have architecture-related strength")
			})
		})

		Context("when skills indicate security expertise", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					fixtures.SkillWith("skill-1", "OAuth", "security", "expert"),
					fixtures.SkillWith("skill-2", "TLS", "security", "advanced"),
				}
			})

			It("should infer security-related core strength", func() {
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).NotTo(BeEmpty())
				hasSecurityStrength := false
				for _, s := range strengths {
					if ContainsAny(s, []string{"security", "compliance"}) {
						hasSecurityStrength = true
						break
					}
				}
				Expect(hasSecurityStrength).To(BeTrue(), "Should have security-related strength")
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

		Context("when facts indicate soft skill competencies", func() {
			It("should infer communication strength", func() {
				facts = []*career.Fact{
					fixtures.FactForValidation("fact-1", "Presented quarterly results to stakeholders", career.RoleFitStaff, []string{"communication"}, []string{"peer"}, "ev-1"),
				}
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).To(ContainElement(ContainSubstring("communication")))
			})

			It("should infer collaboration strength", func() {
				facts = []*career.Fact{
					fixtures.FactForValidation("fact-1", "Worked with cross-functional teams", career.RoleFitStaff, []string{"collaboration"}, []string{"peer"}, "ev-1"),
				}
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).To(ContainElement(ContainSubstring("collaboration")))
			})

			It("should infer problem-solving strength", func() {
				facts = []*career.Fact{
					fixtures.FactForValidation("fact-1", "Debugged complex production issues", career.RoleFitSeniorIC, []string{"problem-solving"}, []string{"peer"}, "ev-1"),
				}
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).To(ContainElement(ContainSubstring("problem-solving")))
			})

			It("should infer project-management strength", func() {
				facts = []*career.Fact{
					fixtures.FactForValidation("fact-1", "Planned sprint milestones and delivery schedule", career.RoleFitEM, []string{"project-management"}, []string{"hiring_manager"}, "ev-1"),
				}
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).To(ContainElement(ContainSubstring("planning")))
			})

			It("should infer architecture strength", func() {
				facts = []*career.Fact{
					fixtures.FactForValidation("fact-1", "Designed scalable distributed system architecture", career.RoleFitPrincipal, []string{"architecture"}, []string{"peer"}, "ev-1"),
				}
				strengths := service.InferCoreStrengths(events, facts, skills)
				Expect(strengths).To(ContainElement(ContainSubstring("architecture")))
			})
		})
	})

	Describe("InferValuePropositions", func() {
		Context("when events show cross-functional work", func() {
			BeforeEach(func() {
				ev1 := fixtures.EventWithCategories("ev-1", "Collaborated with product, design, and engineering teams", []string{"leadership", "product"})
				ev1.Date = time.Now().AddDate(-1, 0, 0)
				ev1.Tags = []string{"leadership"}

				events = []*career.Event{ev1}
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
					fixtures.FactWithCategories("fact-1", "Established engineering onboarding program", "ev-1", []string{"mentoring"}, []string{"hiring_manager"}),
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
					fixtures.SkillWith("skill-1", "Ruby", "backend", "expert"),
					fixtures.SkillWith("skill-2", "Go", "backend", "advanced"),
					fixtures.SkillWith("skill-3", "Python", "backend", "intermediate"),
					fixtures.SkillWith("skill-4", "React", "frontend", "intermediate"),
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
					fixtures.SkillWith("skill-1", "Go", "backend", ""),
					fixtures.SkillWith("skill-2", "Ruby", "backend", ""),
					fixtures.SkillWith("skill-3", "React", "frontend", ""),
					fixtures.SkillWith("skill-4", "PostgreSQL", "database", ""),
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

		Context("when skills include architecture category", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					fixtures.SkillWith("skill-1", "Microservices", "architecture", "expert"),
					fixtures.SkillWith("skill-2", "CQRS", "architecture", "advanced"),
				}
			})

			It("should route architecture skills to systems", func() {
				result := service.InferTechnologies(events, skills)
				Expect(result.Systems).To(ContainElement("Microservices"))
				Expect(result.Systems).To(ContainElement("CQRS"))
			})
		})

		Context("when skills include security category", func() {
			BeforeEach(func() {
				skills = []*career.Skill{
					fixtures.SkillWith("skill-1", "OAuth", "security", "advanced"),
					fixtures.SkillWith("skill-2", "TLS", "security", "intermediate"),
				}
			})

			It("should route security skills to systems", func() {
				result := service.InferTechnologies(events, skills)
				Expect(result.Systems).To(ContainElement("OAuth"))
				Expect(result.Systems).To(ContainElement("TLS"))
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
