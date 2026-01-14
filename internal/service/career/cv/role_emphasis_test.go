package cv_test

import (
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Role Emphasis Configuration", func() {
	Describe("RoleEmphasisConfig", func() {
		It("should define config for senior_backend role emphasis", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisSeniorBackend)
			Expect(config).NotTo(BeNil())
			Expect(config.ID).To(Equal(cv.RoleEmphasisSeniorBackend))
			Expect(config.PrimaryCategories).To(ContainElements("technical", "architecture"))
			Expect(config.SecondaryCategories).To(ContainElements("product", "delivery"))
			Expect(config.BulletFocus).To(ContainSubstring("Technical"))
		})

		It("should define config for staff_principal role emphasis", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisStaffPrincipal)
			Expect(config).NotTo(BeNil())
			Expect(config.ID).To(Equal(cv.RoleEmphasisStaffPrincipal))
			Expect(config.PrimaryCategories).To(ContainElements("leadership", "strategy", "architecture"))
			Expect(config.SecondaryCategories).To(ContainElements("technical", "mentoring"))
			Expect(config.BulletFocus).To(ContainSubstring("Architecture"))
		})

		It("should define config for consulting role emphasis", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisConsulting)
			Expect(config).NotTo(BeNil())
			Expect(config.ID).To(Equal(cv.RoleEmphasisConsulting))
			Expect(config.PrimaryCategories).To(ContainElements("strategy", "delivery"))
			Expect(config.SecondaryCategories).To(ContainElements("technical", "leadership"))
			Expect(config.BulletFocus).To(ContainSubstring("Client"))
		})

		It("should define config for language_agnostic role emphasis", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisLanguageAgnostic)
			Expect(config).NotTo(BeNil())
			Expect(config.ID).To(Equal(cv.RoleEmphasisLanguageAgnostic))
			Expect(config.PrimaryCategories).To(ContainElements("technical", "architecture"))
			Expect(config.BulletFocus).To(ContainSubstring("Multi-language"))
		})

		It("should return default config for unknown role emphasis", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasis("unknown"))
			Expect(config).NotTo(BeNil())
			Expect(config.PrimaryCategories).NotTo(BeEmpty())
		})
	})

	Describe("Role Emphasis Scoring", func() {
		It("should score bullets higher when categories match primary", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisSeniorBackend)
			// Primary categories: technical, architecture
			technicalScore := config.ScoreBulletCategory("technical")
			architectureScore := config.ScoreBulletCategory("architecture")
			otherScore := config.ScoreBulletCategory("leadership")

			Expect(technicalScore).To(BeNumerically(">", otherScore))
			Expect(architectureScore).To(BeNumerically(">", otherScore))
		})

		It("should score bullets medium when categories match secondary", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisSeniorBackend)
			// Secondary categories: product, delivery
			productScore := config.ScoreBulletCategory("product")
			deliveryScore := config.ScoreBulletCategory("delivery")
			unrelatedScore := config.ScoreBulletCategory("mentoring")

			Expect(productScore).To(BeNumerically(">", unrelatedScore))
			Expect(deliveryScore).To(BeNumerically(">", unrelatedScore))
		})

		It("should score staff_principal higher for leadership bullets", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisStaffPrincipal)
			leadershipScore := config.ScoreBulletCategory("leadership")
			strategyScore := config.ScoreBulletCategory("strategy")

			Expect(leadershipScore).To(BeNumerically(">=", 0.8))
			Expect(strategyScore).To(BeNumerically(">=", 0.8))
		})

		It("should score consulting higher for delivery bullets", func() {
			config := cv.GetRoleEmphasisConfig(cv.RoleEmphasisConsulting)
			deliveryScore := config.ScoreBulletCategory("delivery")
			strategyScore := config.ScoreBulletCategory("strategy")

			Expect(deliveryScore).To(BeNumerically(">=", 0.8))
			Expect(strategyScore).To(BeNumerically(">=", 0.8))
		})
	})

	Describe("ListRoleEmphasisConfigs", func() {
		It("should return all 4 role emphasis configs", func() {
			configs := cv.ListRoleEmphasisConfigs()
			Expect(configs).To(HaveLen(4))
		})

		It("should include all role emphasis types", func() {
			configs := cv.ListRoleEmphasisConfigs()
			ids := make([]cv.RoleEmphasis, len(configs))
			for i, c := range configs {
				ids[i] = c.ID
			}
			Expect(ids).To(ContainElements(
				cv.RoleEmphasisSeniorBackend,
				cv.RoleEmphasisStaffPrincipal,
				cv.RoleEmphasisConsulting,
				cv.RoleEmphasisLanguageAgnostic,
			))
		})
	})
})
