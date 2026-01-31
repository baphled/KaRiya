package constants_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/constants"
)

func TestConstants(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Constants Suite")
}

var _ = Describe("Constants", func() {
	Describe("RoleFit", func() {
		Describe("AllRoleFits", func() {
			It("returns all defined role fits", func() {
				roles := constants.AllRoleFits()
				Expect(roles).To(HaveLen(4))
				Expect(roles).To(ContainElements(
					constants.RoleFitPrincipal,
					constants.RoleFitEM,
					constants.RoleFitStaff,
					constants.RoleFitSeniorIC,
				))
			})
		})

		Describe("IsValidRoleFit", func() {
			It("returns true for valid role fits", func() {
				Expect(constants.IsValidRoleFit("principal")).To(BeTrue())
				Expect(constants.IsValidRoleFit("em")).To(BeTrue())
				Expect(constants.IsValidRoleFit("staff")).To(BeTrue())
				Expect(constants.IsValidRoleFit("senior_ic")).To(BeTrue())
			})

			It("returns false for invalid role fits", func() {
				Expect(constants.IsValidRoleFit("invalid")).To(BeFalse())
				Expect(constants.IsValidRoleFit("")).To(BeFalse())
				Expect(constants.IsValidRoleFit("PRINCIPAL")).To(BeFalse())
			})
		})

		Describe("CVTargetRoleFits", func() {
			It("returns role fits suitable for CV targeting", func() {
				roles := constants.CVTargetRoleFits()
				Expect(roles).To(HaveLen(4))
				Expect(roles).To(ContainElements(
					constants.RoleFitPrincipal,
					constants.RoleFitStaff,
					constants.RoleFitEM,
					constants.RoleFitSeniorIC,
				))
			})
		})
	})

	Describe("CompetencyCategory", func() {
		Describe("AllCompetencyCategories", func() {
			It("returns all defined competency categories including soft skills", func() {
				categories := constants.AllCompetencyCategories()
				Expect(categories).To(HaveLen(11))
				Expect(categories).To(ContainElements(
					constants.CompetencyTechnical,
					constants.CompetencyLeadership,
					constants.CompetencyProduct,
					constants.CompetencyConsulting,
					constants.CompetencyResearch,
					constants.CompetencyMentoring,
					constants.CompetencyCommunication,
					constants.CompetencyCollaboration,
					constants.CompetencyProblemSolving,
					constants.CompetencyProjectManagement,
					constants.CompetencyArchitecture,
				))
			})
		})

		Describe("IsValidCompetencyCategory", func() {
			It("returns true for original categories", func() {
				Expect(constants.IsValidCompetencyCategory("technical")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("leadership")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("product")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("consulting")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("research")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("mentoring")).To(BeTrue())
			})

			It("returns true for soft skill categories", func() {
				Expect(constants.IsValidCompetencyCategory("communication")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("collaboration")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("problem-solving")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("project-management")).To(BeTrue())
				Expect(constants.IsValidCompetencyCategory("architecture")).To(BeTrue())
			})

			It("returns false for invalid categories", func() {
				Expect(constants.IsValidCompetencyCategory("invalid")).To(BeFalse())
				Expect(constants.IsValidCompetencyCategory("")).To(BeFalse())
				Expect(constants.IsValidCompetencyCategory("TECHNICAL")).To(BeFalse())
			})
		})
	})

	Describe("SkillLevel", func() {
		Describe("AllSkillLevels", func() {
			It("returns all defined skill levels", func() {
				levels := constants.AllSkillLevels()
				Expect(levels).To(HaveLen(4))
				Expect(levels).To(ContainElements(
					constants.SkillLevelBeginner,
					constants.SkillLevelIntermediate,
					constants.SkillLevelAdvanced,
					constants.SkillLevelExpert,
				))
			})
		})

		Describe("IsValidSkillLevel", func() {
			It("returns true for valid skill levels", func() {
				Expect(constants.IsValidSkillLevel("beginner")).To(BeTrue())
				Expect(constants.IsValidSkillLevel("intermediate")).To(BeTrue())
				Expect(constants.IsValidSkillLevel("advanced")).To(BeTrue())
				Expect(constants.IsValidSkillLevel("expert")).To(BeTrue())
			})

			It("returns false for invalid skill levels", func() {
				Expect(constants.IsValidSkillLevel("invalid")).To(BeFalse())
				Expect(constants.IsValidSkillLevel("")).To(BeFalse())
				Expect(constants.IsValidSkillLevel("BEGINNER")).To(BeFalse())
			})
		})
	})

	Describe("EventTag", func() {
		Describe("AllEventTags", func() {
			It("returns all defined event tags", func() {
				tags := constants.AllEventTags()
				Expect(tags).To(HaveLen(8))
				Expect(tags).To(ContainElements(
					constants.EventTagProject,
					constants.EventTagAchievement,
					constants.EventTagLeadership,
					constants.EventTagTechnical,
					constants.EventTagConsulting,
					constants.EventTagResearch,
					constants.EventTagProduct,
					constants.EventTagMentoring,
				))
			})
		})

		Describe("IsValidEventTag", func() {
			It("returns true for valid event tags", func() {
				Expect(constants.IsValidEventTag("project")).To(BeTrue())
				Expect(constants.IsValidEventTag("achievement")).To(BeTrue())
				Expect(constants.IsValidEventTag("leadership")).To(BeTrue())
				Expect(constants.IsValidEventTag("technical")).To(BeTrue())
				Expect(constants.IsValidEventTag("consulting")).To(BeTrue())
				Expect(constants.IsValidEventTag("research")).To(BeTrue())
				Expect(constants.IsValidEventTag("product")).To(BeTrue())
				Expect(constants.IsValidEventTag("mentoring")).To(BeTrue())
			})

			It("returns false for invalid event tags", func() {
				Expect(constants.IsValidEventTag("invalid")).To(BeFalse())
				Expect(constants.IsValidEventTag("")).To(BeFalse())
				Expect(constants.IsValidEventTag("PROJECT")).To(BeFalse())
			})
		})
	})

	Describe("Audience", func() {
		Describe("AllAudiences", func() {
			It("returns all defined audiences", func() {
				audiences := constants.AllAudiences()
				Expect(audiences).To(HaveLen(3))
				Expect(audiences).To(ContainElements(
					constants.AudienceHiringManager,
					constants.AudienceRecruiter,
					constants.AudiencePeer,
				))
			})
		})

		Describe("IsValidAudience", func() {
			It("returns true for valid audiences", func() {
				Expect(constants.IsValidAudience("hiring_manager")).To(BeTrue())
				Expect(constants.IsValidAudience("recruiter")).To(BeTrue())
				Expect(constants.IsValidAudience("peer")).To(BeTrue())
			})

			It("returns false for invalid audiences", func() {
				Expect(constants.IsValidAudience("invalid")).To(BeFalse())
				Expect(constants.IsValidAudience("")).To(BeFalse())
				Expect(constants.IsValidAudience("HIRING_MANAGER")).To(BeFalse())
			})
		})
	})

	Describe("SkillCategory", func() {
		Describe("AllSkillCategories", func() {
			It("returns all 12 defined skill categories", func() {
				categories := constants.AllSkillCategories()
				Expect(categories).To(HaveLen(12))
				Expect(categories).To(ContainElements(
					constants.SkillCategoryBackend,
					constants.SkillCategoryFrontend,
					constants.SkillCategoryDevOps,
					constants.SkillCategoryDatabase,
					constants.SkillCategoryCloud,
					constants.SkillCategoryMobile,
					constants.SkillCategoryTooling,
					constants.SkillCategoryTesting,
					constants.SkillCategoryData,
					constants.SkillCategoryML,
					constants.SkillCategoryMonitoring,
					constants.SkillCategoryOther,
				))
			})
		})

		Describe("SuggestedSkillCategories", func() {
			It("returns all skill categories for form dropdowns", func() {
				categories := constants.SuggestedSkillCategories()
				Expect(categories).To(HaveLen(12))
				Expect(categories).To(ContainElements(
					constants.SkillCategoryBackend,
					constants.SkillCategoryFrontend,
					constants.SkillCategoryDevOps,
					constants.SkillCategoryDatabase,
					constants.SkillCategoryCloud,
					constants.SkillCategoryMobile,
					constants.SkillCategoryTooling,
					constants.SkillCategoryTesting,
					constants.SkillCategoryData,
					constants.SkillCategoryML,
					constants.SkillCategoryMonitoring,
					constants.SkillCategoryOther,
				))
			})
		})

		Describe("IsValidSkillCategory", func() {
			It("returns true for all canonical skill categories", func() {
				Expect(constants.IsValidSkillCategory("backend")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("frontend")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("devops")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("database")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("cloud")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("mobile")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("tooling")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("testing")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("data")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("ml")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("monitoring")).To(BeTrue())
				Expect(constants.IsValidSkillCategory("other")).To(BeTrue())
			})

			It("returns false for invalid categories", func() {
				Expect(constants.IsValidSkillCategory("invalid")).To(BeFalse())
				Expect(constants.IsValidSkillCategory("")).To(BeFalse())
				Expect(constants.IsValidSkillCategory("BACKEND")).To(BeFalse())
			})

			It("returns false for categories that were merged into existing ones", func() {
				Expect(constants.IsValidSkillCategory("build")).To(BeFalse())
				Expect(constants.IsValidSkillCategory("documentation")).To(BeFalse())
				Expect(constants.IsValidSkillCategory("os")).To(BeFalse())
			})
		})
	})

	Describe("SectionType", func() {
		Describe("AllSectionTypes", func() {
			It("returns all defined section types", func() {
				types := constants.AllSectionTypes()
				Expect(types).To(HaveLen(4))
				Expect(types).To(ContainElements(
					constants.SectionTypeExperience,
					constants.SectionTypeProjects,
					constants.SectionTypeSkills,
					constants.SectionTypeSummary,
				))
			})
		})

		Describe("IsValidSectionType", func() {
			It("returns true for valid section types", func() {
				Expect(constants.IsValidSectionType("experience")).To(BeTrue())
				Expect(constants.IsValidSectionType("projects")).To(BeTrue())
				Expect(constants.IsValidSectionType("skills")).To(BeTrue())
				Expect(constants.IsValidSectionType("summary")).To(BeTrue())
			})

			It("returns false for invalid section types", func() {
				Expect(constants.IsValidSectionType("invalid")).To(BeFalse())
				Expect(constants.IsValidSectionType("")).To(BeFalse())
			})
		})
	})

	Describe("ImpactLevel", func() {
		Describe("AllImpactLevels", func() {
			It("returns all defined impact levels", func() {
				levels := constants.AllImpactLevels()
				Expect(levels).To(HaveLen(3))
				Expect(levels).To(ContainElements(
					constants.ImpactLevelLow,
					constants.ImpactLevelMedium,
					constants.ImpactLevelHigh,
				))
			})
		})

		Describe("IsValidImpactLevel", func() {
			It("returns true for valid impact levels", func() {
				Expect(constants.IsValidImpactLevel("low")).To(BeTrue())
				Expect(constants.IsValidImpactLevel("medium")).To(BeTrue())
				Expect(constants.IsValidImpactLevel("high")).To(BeTrue())
			})

			It("returns true for empty string (not set)", func() {
				Expect(constants.IsValidImpactLevel("")).To(BeTrue())
			})

			It("returns false for invalid impact levels", func() {
				Expect(constants.IsValidImpactLevel("invalid")).To(BeFalse())
				Expect(constants.IsValidImpactLevel("LOW")).To(BeFalse())
			})
		})
	})

	Describe("InclusionReason", func() {
		Describe("AllInclusionReasons", func() {
			It("returns all defined inclusion reasons", func() {
				reasons := constants.AllInclusionReasons()
				Expect(reasons).To(HaveLen(9))
				// Semantic reasons
				Expect(reasons).To(ContainElements(
					constants.InclusionReasonOwnership,
					constants.InclusionReasonContribution,
					constants.InclusionReasonStrategy,
					constants.InclusionReasonExecution,
					constants.InclusionReasonOutcome,
					constants.InclusionReasonActivity,
				))
				// Source-based reasons
				Expect(reasons).To(ContainElements(
					constants.InclusionReasonFactExtraction,
					constants.InclusionReasonEventDirect,
					constants.InclusionReasonAchievementExtraction,
				))
			})
		})

		Describe("IsValidInclusionReason", func() {
			It("returns true for valid semantic reasons", func() {
				Expect(constants.IsValidInclusionReason("ownership")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("contribution")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("strategy")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("execution")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("outcome")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("activity")).To(BeTrue())
			})

			It("returns true for valid source-based reasons", func() {
				Expect(constants.IsValidInclusionReason("fact_extraction")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("event_direct")).To(BeTrue())
				Expect(constants.IsValidInclusionReason("achievement_extraction")).To(BeTrue())
			})

			It("returns false for invalid reasons", func() {
				Expect(constants.IsValidInclusionReason("invalid")).To(BeFalse())
				Expect(constants.IsValidInclusionReason("")).To(BeFalse())
			})
		})
	})

	Describe("FocusArea", func() {
		Describe("AllFocusAreas", func() {
			It("returns all defined focus areas", func() {
				areas := constants.AllFocusAreas()
				Expect(areas).To(HaveLen(4))
				Expect(areas).To(ContainElements(
					constants.FocusAreaBackend,
					constants.FocusAreaFrontend,
					constants.FocusAreaFullstack,
					constants.FocusAreaDevOps,
				))
			})
		})

		Describe("IsValidFocusArea", func() {
			It("returns true for valid focus areas", func() {
				Expect(constants.IsValidFocusArea("backend")).To(BeTrue())
				Expect(constants.IsValidFocusArea("frontend")).To(BeTrue())
				Expect(constants.IsValidFocusArea("fullstack")).To(BeTrue())
				Expect(constants.IsValidFocusArea("devops")).To(BeTrue())
			})

			It("returns false for invalid focus areas", func() {
				Expect(constants.IsValidFocusArea("invalid")).To(BeFalse())
				Expect(constants.IsValidFocusArea("")).To(BeFalse())
			})
		})
	})

	Describe("TechnologyFocus", func() {
		Describe("AllTechnologyFocuses", func() {
			It("returns all defined technology focuses", func() {
				focuses := constants.AllTechnologyFocuses()
				Expect(focuses).To(HaveLen(3))
				Expect(focuses).To(ContainElements(
					constants.TechnologyFocusLanguageAgnostic,
					constants.TechnologyFocusGeneralist,
					constants.TechnologyFocusSpecialist,
				))
			})
		})

		Describe("IsValidTechnologyFocus", func() {
			It("returns true for valid technology focuses", func() {
				Expect(constants.IsValidTechnologyFocus("language_agnostic")).To(BeTrue())
				Expect(constants.IsValidTechnologyFocus("generalist")).To(BeTrue())
				Expect(constants.IsValidTechnologyFocus("specialist")).To(BeTrue())
			})

			It("returns false for invalid technology focuses", func() {
				Expect(constants.IsValidTechnologyFocus("invalid")).To(BeFalse())
				Expect(constants.IsValidTechnologyFocus("")).To(BeFalse())
			})
		})
	})
})
