package career

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fact", func() {
	Context("when creating a fact", func() {
		It("should have all required fields", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Led team to deliver critical project",
				CompetencyCategories: []string{"leadership", "technical"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager", "peer"},
				StrengthSignal:       "Strong leadership and technical execution",
				SourceEventID:        "event-1",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			Expect(fact.ID).To(Equal("fact-123"))
			Expect(fact.Text).To(Equal("Led team to deliver critical project"))
			Expect(fact.CompetencyCategories).To(HaveLen(2))
			Expect(fact.RoleFit).To(Equal(RoleFitPrincipal))
			Expect(fact.AudienceRelevance).To(HaveLen(2))
		})
	})

	Context("when validating a fact", func() {
		It("should validate successfully with all required fields", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Led team to deliver critical project",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should validate with burst source instead of event", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Led team to deliver critical project",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceBurstID:        "burst-1",
			}

			err := fact.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject fact with empty ID", func() {
			fact := &Fact{
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ID cannot be empty"))
		})

		It("should reject fact with empty text", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("text cannot be empty"))
		})

		It("should reject fact with text exceeding 2000 characters", func() {
			longText := ""
			for i := 0; i < 2001; i++ {
				longText += "a"
			}

			fact := &Fact{
				ID:                   "fact-123",
				Text:                 longText,
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot exceed 2000 characters"))
		})

		It("should reject fact with no competency categories", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least one competency category"))
		})

		It("should reject fact with invalid competency category", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"invalid-category"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid competency category"))
		})

		It("should reject fact with duplicate competency categories", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership", "leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate competency categories"))
		})

		It("should reject fact with empty role fit", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              "",
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("role fit cannot be empty"))
		})

		It("should reject fact with invalid role fit", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              "invalid-role",
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid role fit"))
		})

		It("should accept all valid role fits", func() {
			roleFits := []RoleFit{RoleFitPrincipal, RoleFitEM, RoleFitStaff, RoleFitSeniorIC}

			for _, roleFit := range roleFits {
				fact := &Fact{
					ID:                   "fact-123",
					Text:                 "Test",
					CompetencyCategories: []string{"leadership"},
					RoleFit:              roleFit,
					AudienceRelevance:    []string{"hiring_manager"},
					SourceEventID:        "event-1",
				}

				err := fact.Validate()
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should reject fact with no audience relevance", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least one audience relevance type"))
		})

		It("should reject fact with invalid audience relevance", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"invalid-audience"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid audience relevance"))
		})

		It("should reject fact with duplicate audience relevance", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager", "hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate audience relevance types"))
		})

		It("should reject fact with no source reference", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Test",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "",
				SourceBurstID:        "",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least one source"))
		})

		It("should reject fact with aspirational language 'will'", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "I will lead the team",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aspirational language"))
		})

		It("should reject fact with aspirational language 'should'", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "You should know that I managed the project",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aspirational language"))
		})

		It("should reject fact with aspirational language 'could'", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "I could have done better",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aspirational language"))
		})

		It("should accept fact with grounded language", func() {
			fact := &Fact{
				ID:                   "fact-123",
				Text:                 "Led cross-functional team of 8 engineers to deliver critical platform migration",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              RoleFitPrincipal,
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        "event-1",
			}

			err := fact.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept all valid audience types", func() {
			audiences := []string{"hiring_manager", "recruiter", "peer"}

			for _, audience := range audiences {
				fact := &Fact{
					ID:                   "fact-123",
					Text:                 "Test",
					CompetencyCategories: []string{"leadership"},
					RoleFit:              RoleFitPrincipal,
					AudienceRelevance:    []string{audience},
					SourceEventID:        "event-1",
				}

				err := fact.Validate()
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
