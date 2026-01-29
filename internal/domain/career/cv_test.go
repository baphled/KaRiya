//nolint:errcheck // Test file - error handling for test setup is not relevant.
package career

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVView", func() {
	Describe("Validate", func() {
		var cvView *CVView

		BeforeEach(func() {
			cvView = &CVView{
				ID:             "test-id-1",
				Name:           "Senior IC CV",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
				EventFilters: map[string]interface{}{
					"date_from": "2020-01-01",
					"date_to":   "2024-01-01",
				},
				GeneratedAt:      time.Now(),
				SourceEventCount: 15,
				SourceFactCount:  8,
			}
		})

		Context("with valid CV", func() {
			It("should pass validation", func() {
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with empty ID", func() {
			It("should fail validation", func() {
				cvView.ID = ""
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ID cannot be empty"))
			})
		})

		Context("with whitespace-only ID", func() {
			It("should fail validation", func() {
				cvView.ID = "   "
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with empty name", func() {
			It("should fail validation", func() {
				cvView.Name = ""
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
			})
		})

		Context("with name exceeding 200 characters", func() {
			It("should fail validation", func() {
				cvView.Name = string(make([]byte, 201))
				for i := range cvView.Name {
					cvView.Name = cvView.Name[:i] + "a" + cvView.Name[i+1:]
				}
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot exceed 200 characters"))
			})
		})

		Context("with empty target role", func() {
			It("should fail validation", func() {
				cvView.TargetRole = ""
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid target role", func() {
			It("should fail validation", func() {
				cvView.TargetRole = "invalid_role"
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid target role"))
			})
		})

		Context("with all valid target roles", func() {
			It("should pass validation for principal", func() {
				cvView.TargetRole = "principal"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation for staff", func() {
				cvView.TargetRole = "staff"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation for em", func() {
				cvView.TargetRole = "em"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation for senior_ic", func() {
				cvView.TargetRole = "senior_ic"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with uppercase target role", func() {
			It("should pass validation (case-insensitive)", func() {
				cvView.TargetRole = "PRINCIPAL"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with empty audience", func() {
			It("should fail validation", func() {
				cvView.TargetAudience = ""
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with whitespace-only audience", func() {
			It("should fail validation", func() {
				cvView.TargetAudience = "   "
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid audience", func() {
			It("should fail validation", func() {
				cvView.TargetAudience = "invalid_audience"
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid target audience"))
			})
		})

		Context("with all valid audiences", func() {
			It("should pass validation with hiring_manager", func() {
				cvView.TargetAudience = "hiring_manager"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation with recruiter", func() {
				cvView.TargetAudience = "recruiter"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation with peer", func() {
				cvView.TargetAudience = "peer"
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with negative source event count", func() {
			It("should fail validation", func() {
				cvView.SourceEventCount = -1
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("source event count cannot be negative"))
			})
		})

		Context("with negative source fact count", func() {
			It("should fail validation", func() {
				cvView.SourceFactCount = -1
				err := cvView.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("source fact count cannot be negative"))
			})
		})

		Context("with zero source counts", func() {
			It("should pass validation", func() {
				cvView.SourceEventCount = 0
				cvView.SourceFactCount = 0
				err := cvView.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

var _ = Describe("CVSection", func() {
	Describe("Validate", func() {
		var cvSection *CVSection

		BeforeEach(func() {
			cvSection = &CVSection{
				ID:          "section-id-1",
				CVViewID:    "cv-id-1",
				SectionType: "experience",
				Title:       "Professional Experience",
				Order:       0,
				Content: []*SectionContentGroup{
					{
						Header: "Test Company",
						Bullets: []*CVBullet{
							{ID: "bullet-1", Text: "Test bullet"},
						},
					},
				},
			}
		})

		Context("with valid section", func() {
			It("should pass validation", func() {
				err := cvSection.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with empty ID", func() {
			It("should fail validation", func() {
				cvSection.ID = ""
				err := cvSection.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ID cannot be empty"))
			})
		})

		Context("with empty CVViewID", func() {
			It("should fail validation", func() {
				cvSection.CVViewID = ""
				err := cvSection.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("CV view ID cannot be empty"))
			})
		})

		Context("with invalid section type", func() {
			It("should fail validation", func() {
				cvSection.SectionType = "invalid_type"
				err := cvSection.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid section type"))
			})
		})

		Context("with all valid section types", func() {
			It("should pass for experience", func() {
				cvSection.SectionType = "experience"
				err := cvSection.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass for skills", func() {
				cvSection.SectionType = "skills"
				err := cvSection.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass for summary", func() {
				cvSection.SectionType = "summary"
				err := cvSection.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with empty title", func() {
			It("should fail validation", func() {
				cvSection.Title = ""
				err := cvSection.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("title cannot be empty"))
			})
		})

		Context("with title exceeding 200 characters", func() {
			It("should fail validation", func() {
				cvSection.Title = string(make([]byte, 201))
				err := cvSection.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot exceed 200 characters"))
			})
		})

		Context("with negative order", func() {
			It("should fail validation", func() {
				cvSection.Order = -1
				err := cvSection.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("order cannot be negative"))
			})
		})

		Context("with zero order", func() {
			It("should pass validation", func() {
				cvSection.Order = 0
				err := cvSection.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

var _ = Describe("CVBullet", func() {
	Describe("Validate", func() {
		var cvBullet *CVBullet

		BeforeEach(func() {
			cvBullet = &CVBullet{
				ID:              "bullet-id-1",
				SectionID:       "section-id-1",
				Text:            "Led the design of a distributed system architecture",
				SourceEventIDs:  []string{"event-1", "event-2"},
				SourceFactIDs:   []string{"fact-1"},
				Rank:            0.85,
				InclusionReason: "ownership",
				Confidence:      0.9,
			}
		})

		Context("with valid bullet", func() {
			It("should pass validation", func() {
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with empty ID", func() {
			It("should fail validation", func() {
				cvBullet.ID = ""
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ID cannot be empty"))
			})
		})

		Context("with empty SectionID", func() {
			It("should fail validation", func() {
				cvBullet.SectionID = ""
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("section ID cannot be empty"))
			})
		})

		Context("with empty text", func() {
			It("should fail validation", func() {
				cvBullet.Text = ""
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("text cannot be empty"))
			})
		})

		Context("with text exceeding 500 characters", func() {
			It("should fail validation", func() {
				cvBullet.Text = string(make([]byte, 501))
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot exceed 500 characters"))
			})
		})

		Context("with no source events", func() {
			It("should fail validation", func() {
				cvBullet.SourceEventIDs = []string{}
				err := cvBullet.Validate()
				Expect(err).To(Equal(ErrNoSourceEvents))
			})
		})

		Context("with nil source events", func() {
			It("should fail validation", func() {
				cvBullet.SourceEventIDs = nil
				err := cvBullet.Validate()
				Expect(err).To(Equal(ErrNoSourceEvents))
			})
		})

		Context("with empty event ID in source list", func() {
			It("should fail validation", func() {
				cvBullet.SourceEventIDs = []string{"event-1", ""}
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be empty"))
			})
		})

		Context("with duplicate event IDs", func() {
			It("should fail validation", func() {
				cvBullet.SourceEventIDs = []string{"event-1", "event-1"}
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("duplicate"))
			})
		})

		Context("with rank below 0.0", func() {
			It("should fail validation", func() {
				cvBullet.Rank = -0.1
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("between 0.0 and 1.0"))
			})
		})

		Context("with rank above 1.0", func() {
			It("should fail validation", func() {
				cvBullet.Rank = 1.1
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("between 0.0 and 1.0"))
			})
		})

		Context("with rank at boundaries", func() {
			It("should pass with rank 0.0", func() {
				cvBullet.Rank = 0.0
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass with rank 1.0", func() {
				cvBullet.Rank = 1.0
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with confidence below 0.0", func() {
			It("should fail validation", func() {
				cvBullet.Confidence = -0.1
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with confidence above 1.0", func() {
			It("should fail validation", func() {
				cvBullet.Confidence = 1.1
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid inclusion reason", func() {
			It("should fail validation", func() {
				cvBullet.InclusionReason = "invalid_reason"
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid inclusion reason"))
			})
		})

		Context("with all valid inclusion reasons", func() {
			validReasons := []string{"ownership", "contribution", "strategy", "execution", "outcome", "activity"}

			for _, reason := range validReasons {
				It("should pass for "+reason, func() {
					cvBullet.InclusionReason = reason
					err := cvBullet.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			}
		})

		Context("with generator-based inclusion reasons", func() {
			// These reasons are used by BulletGenerator
			// to indicate the source of the bullet (Task 44)
			generatorReasons := []string{"fact_extraction", "event_direct", "achievement_extraction"}

			for _, reason := range generatorReasons {
				It("should pass for "+reason, func() {
					cvBullet.InclusionReason = reason
					err := cvBullet.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			}
		})

		// Task 44: Enhanced fields from BulletGenerator
		Context("with enhanced fields", func() {
			It("should pass validation with EnhancedText set", func() {
				cvBullet.EnhancedText = "Architected and led implementation of distributed system serving 10M users"
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation with all score fields set", func() {
				cvBullet.RoleScore = 0.9
				cvBullet.AudienceScore = 0.85
				cvBullet.MetricScore = 0.75
				cvBullet.ImpactScore = 0.95
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation with ImpactLevel set", func() {
				cvBullet.ImpactLevel = "high"
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation with KeywordMatches set", func() {
				cvBullet.KeywordMatches = []string{"distributed", "architecture", "leadership"}
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass validation with all enhanced fields set", func() {
				cvBullet.EnhancedText = "Architected distributed system at scale"
				cvBullet.RoleScore = 0.9
				cvBullet.AudienceScore = 0.85
				cvBullet.MetricScore = 0.75
				cvBullet.ImpactScore = 0.95
				cvBullet.ImpactLevel = "high"
				cvBullet.KeywordMatches = []string{"architecture", "scale"}
				err := cvBullet.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		// Task 44: Enhanced field validation
		Context("with invalid enhanced score fields", func() {
			It("should fail validation when RoleScore is negative", func() {
				cvBullet.RoleScore = -0.1
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("role score"))
			})

			It("should fail validation when RoleScore exceeds 1.0", func() {
				cvBullet.RoleScore = 1.1
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("role score"))
			})

			It("should fail validation when AudienceScore is negative", func() {
				cvBullet.AudienceScore = -0.5
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("audience score"))
			})

			It("should fail validation when AudienceScore exceeds 1.0", func() {
				cvBullet.AudienceScore = 2.0
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("audience score"))
			})

			It("should fail validation when MetricScore is negative", func() {
				cvBullet.MetricScore = -0.01
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("metric score"))
			})

			It("should fail validation when MetricScore exceeds 1.0", func() {
				cvBullet.MetricScore = 1.5
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("metric score"))
			})

			It("should fail validation when ImpactScore is negative", func() {
				cvBullet.ImpactScore = -1.0
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("impact score"))
			})

			It("should fail validation when ImpactScore exceeds 1.0", func() {
				cvBullet.ImpactScore = 100.0
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("impact score"))
			})
		})

		Context("with invalid ImpactLevel", func() {
			It("should fail validation with invalid impact level", func() {
				cvBullet.ImpactLevel = "invalid"
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("impact level"))
			})

			It("should fail validation with uppercase impact level", func() {
				cvBullet.ImpactLevel = "HIGH"
				err := cvBullet.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("impact level"))
			})

			It("should pass validation with valid impact levels", func() {
				for _, level := range []string{"", "low", "medium", "high"} {
					cvBullet.ImpactLevel = level
					err := cvBullet.Validate()
					Expect(err).NotTo(HaveOccurred())
				}
			})
		})
	})
})

var _ = Describe("CVConfig", func() {
	Describe("Validate", func() {
		var cvConfig *CVConfig

		BeforeEach(func() {
			cvConfig = &CVConfig{
				Name:           "Staff Engineer CV",
				TargetRole:     "staff",
				TargetAudience: "hiring_manager",
				EventFilters: map[string]interface{}{
					"date_from": "2020-01-01",
					"tags":      []string{"technical", "leadership"},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		})

		Context("with valid config", func() {
			It("should pass validation", func() {
				err := cvConfig.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with empty name", func() {
			It("should fail validation", func() {
				cvConfig.Name = ""
				err := cvConfig.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
			})
		})

		Context("with name exceeding 200 characters", func() {
			It("should fail validation", func() {
				cvConfig.Name = string(make([]byte, 201))
				err := cvConfig.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid filename characters", func() {
			invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}

			for _, char := range invalidChars {
				It("should fail with character "+char, func() {
					cvConfig.Name = "test" + char + "config"
					err := cvConfig.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("invalid character"))
				})
			}
		})

		Context("with invalid target role", func() {
			It("should fail validation", func() {
				cvConfig.TargetRole = "invalid_role"
				err := cvConfig.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with empty audience", func() {
			It("should fail validation", func() {
				cvConfig.TargetAudience = ""
				err := cvConfig.Validate()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ToJSON and FromJSON", func() {
		var cvConfig *CVConfig

		BeforeEach(func() {
			cvConfig = &CVConfig{
				Name:           "Principal CV",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters: map[string]interface{}{
					"date_from": "2015-01-01",
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
		})

		Context("ToJSON", func() {
			It("should convert to JSON", func() {
				jsonBytes, err := cvConfig.ToJSON()
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonBytes).NotTo(BeEmpty())
			})

			It("should produce valid JSON", func() {
				jsonBytes, err := cvConfig.ToJSON()
				Expect(err).NotTo(HaveOccurred())

				var result map[string]interface{}
				err = json.Unmarshal(jsonBytes, &result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result["name"]).To(Equal("Principal CV"))
			})
		})

		Context("FromJSON", func() {
			It("should populate from JSON", func() {
				jsonBytes, _ := cvConfig.ToJSON()

				newConfig := &CVConfig{}
				err := newConfig.FromJSON(jsonBytes)
				Expect(err).NotTo(HaveOccurred())
				Expect(newConfig.Name).To(Equal(cvConfig.Name))
			})

			It("should handle invalid JSON", func() {
				newConfig := &CVConfig{}
				err := newConfig.FromJSON([]byte("invalid json"))
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("Validation Helper Functions", func() {
	Describe("IsAspirationLanguage", func() {
		It("should detect 'will' as aspirational", func() {
			Expect(IsAspirationLanguage("I will build a new system")).To(BeTrue())
		})

		It("should detect 'hoping to' as aspirational", func() {
			Expect(IsAspirationLanguage("I'm hoping to improve performance")).To(BeTrue())
		})

		It("should detect 'aiming to' as aspirational", func() {
			Expect(IsAspirationLanguage("Aiming to reduce latency")).To(BeTrue())
		})

		It("should detect 'in progress' as aspirational", func() {
			Expect(IsAspirationLanguage("Work is in progress")).To(BeTrue())
		})

		It("should not flag factual language", func() {
			Expect(IsAspirationLanguage("Built a distributed system")).To(BeFalse())
		})

		It("should be case-insensitive", func() {
			Expect(IsAspirationLanguage("I WILL build it")).To(BeTrue())
		})
	})

	Describe("IsSingleClaimBullet", func() {
		It("should pass single claim", func() {
			Expect(IsSingleClaimBullet("Led the API redesign")).To(BeTrue())
		})

		It("should allow single 'and' conjunction", func() {
			text := "Led the API redesign and implemented caching"
			Expect(IsSingleClaimBullet(text)).To(BeTrue())
		})

		It("should allow single 'while' conjunction", func() {
			text := "Built the system while maintaining uptime"
			Expect(IsSingleClaimBullet(text)).To(BeTrue())
		})

		It("should detect multiple claims with multiple conjunctions", func() {
			text := "Designed architecture; implemented performance improvements and optimized queries"
			Expect(IsSingleClaimBullet(text)).To(BeFalse())
		})

		It("should allow single conjunction in simple contexts", func() {
			text := "Designed and built the data pipeline"
			Expect(IsSingleClaimBullet(text)).To(BeTrue())
		})
	})

	Describe("HasInferredMetrics", func() {
		It("should detect 'thousands of' as inferred", func() {
			Expect(HasInferredMetrics("Impacted thousands of users")).To(BeTrue())
		})

		It("should detect 'millions of' as inferred", func() {
			Expect(HasInferredMetrics("Millions of customers benefited")).To(BeTrue())
		})

		It("should detect 'improved' as potentially inferred", func() {
			Expect(HasInferredMetrics("Improved performance significantly")).To(BeTrue())
		})

		It("should detect 'increased' as inferred", func() {
			Expect(HasInferredMetrics("Increased throughput by x%")).To(BeTrue())
		})

		It("should not flag specific metrics", func() {
			Expect(HasInferredMetrics("Reduced latency from 500ms to 100ms")).To(BeFalse())
		})
	})

	Describe("IsRoleInflation", func() {
		It("should detect inflation for Senior IC role", func() {
			Expect(IsRoleInflation("Managed the team for 5 years", "senior_ic")).To(BeTrue())
		})

		It("should detect inflation for Principal role", func() {
			Expect(IsRoleInflation("Single-handedly built the entire platform", "principal")).To(BeTrue())
		})

		It("should not flag appropriate Senior IC claims", func() {
			Expect(IsRoleInflation("Architected the service layer", "senior_ic")).To(BeFalse())
		})

		It("should be case-insensitive for role", func() {
			Expect(IsRoleInflation("Managed the team", "SENIOR_IC")).To(BeTrue())

		})
	})
})
