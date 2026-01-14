package cv_test

import (
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CV Variants", func() {
	Describe("RoleEmphasis", func() {
		It("should define 4 role emphases", func() {
			emphases := []cv.RoleEmphasis{
				cv.RoleEmphasisSeniorBackend,
				cv.RoleEmphasisStaffPrincipal,
				cv.RoleEmphasisConsulting,
				cv.RoleEmphasisLanguageAgnostic,
			}
			Expect(len(emphases)).To(Equal(4))
		})

		It("should have valid string values", func() {
			Expect(string(cv.RoleEmphasisSeniorBackend)).To(Equal("senior_backend"))
			Expect(string(cv.RoleEmphasisStaffPrincipal)).To(Equal("staff_principal"))
			Expect(string(cv.RoleEmphasisConsulting)).To(Equal("consulting"))
			Expect(string(cv.RoleEmphasisLanguageAgnostic)).To(Equal("language_agnostic"))
		})
	})

	Describe("LengthFormat", func() {
		It("should define 4 length formats", func() {
			formats := []cv.LengthFormat{
				cv.LengthFull,
				cv.LengthStandard,
				cv.LengthShort,
				cv.LengthUltraShort,
			}
			Expect(len(formats)).To(Equal(4))
		})

		It("should have valid string values", func() {
			Expect(string(cv.LengthFull)).To(Equal("full"))
			Expect(string(cv.LengthStandard)).To(Equal("standard"))
			Expect(string(cv.LengthShort)).To(Equal("short"))
			Expect(string(cv.LengthUltraShort)).To(Equal("ultra_short"))
		})
	})

	Describe("CVVariant", func() {
		It("should have required fields", func() {
			variant := &cv.CVVariant{
				ID:            "test_variant",
				Name:          "Test Variant",
				Description:   "A test variant",
				RoleEmphasis:  cv.RoleEmphasisSeniorBackend,
				LengthFormat:  cv.LengthStandard,
				BaseStructure: cv.CVStructureStandard,
				IsBuiltIn:     true,
			}

			Expect(variant.ID).To(Equal("test_variant"))
			Expect(variant.Name).To(Equal("Test Variant"))
			Expect(variant.RoleEmphasis).To(Equal(cv.RoleEmphasisSeniorBackend))
			Expect(variant.LengthFormat).To(Equal(cv.LengthStandard))
			Expect(variant.BaseStructure).To(Equal(cv.CVStructureStandard))
			Expect(variant.IsBuiltIn).To(BeTrue())
		})
	})

	Describe("BulletConfig", func() {
		It("should support nil values for defaults", func() {
			config := &cv.BulletConfig{}
			Expect(config.MaxBulletsPerCompany).To(BeNil())
			Expect(config.MinConfidence).To(BeNil())
			Expect(config.MaxYearsHistory).To(BeNil())
			Expect(config.MaxCompanies).To(BeNil())
		})

		It("should support explicit values", func() {
			maxBullets := 5
			minConf := 0.75
			maxYears := 10
			maxCompanies := 3

			config := &cv.BulletConfig{
				MaxBulletsPerCompany: &maxBullets,
				MinConfidence:        &minConf,
				MaxYearsHistory:      &maxYears,
				MaxCompanies:         &maxCompanies,
			}

			Expect(*config.MaxBulletsPerCompany).To(Equal(5))
			Expect(*config.MinConfidence).To(Equal(0.75))
			Expect(*config.MaxYearsHistory).To(Equal(10))
			Expect(*config.MaxCompanies).To(Equal(3))
		})
	})

	Describe("SectionConfig", func() {
		It("should define section configuration", func() {
			config := cv.SectionConfig{
				SectionType:   "experience",
				Title:         "Professional Experience",
				Order:         1,
				Enabled:       true,
				MinConfidence: 0.7,
			}

			Expect(config.SectionType).To(Equal("experience"))
			Expect(config.Title).To(Equal("Professional Experience"))
			Expect(config.Order).To(Equal(1))
			Expect(config.Enabled).To(BeTrue())
			Expect(config.MinConfidence).To(Equal(0.7))
		})
	})

	Describe("ProfileOverride", func() {
		It("should support optional fields", func() {
			override := &cv.ProfileOverride{}
			Expect(override.ProfessionalTitle).To(BeNil())
			Expect(override.CoreStrengths).To(BeNil())
			Expect(override.CareerDifferentiators).To(BeNil())
			Expect(override.CareerPositioning).To(BeNil())
		})

		It("should support explicit values", func() {
			title := "Senior Consulting Engineer"
			positioning := "Building scalable systems"

			override := &cv.ProfileOverride{
				ProfessionalTitle:     &title,
				CoreStrengths:         []string{"Go", "Architecture", "Leadership"},
				CareerDifferentiators: []string{"20 years experience"},
				CareerPositioning:     &positioning,
			}

			Expect(*override.ProfessionalTitle).To(Equal("Senior Consulting Engineer"))
			Expect(override.CoreStrengths).To(HaveLen(3))
			Expect(override.CareerDifferentiators).To(HaveLen(1))
			Expect(*override.CareerPositioning).To(Equal("Building scalable systems"))
		})
	})

	Describe("BuiltInVariants", func() {
		It("should define exactly 16 built-in variants", func() {
			Expect(cv.BuiltInVariants).To(HaveLen(16))
		})

		It("should have unique IDs", func() {
			ids := make(map[string]bool)
			for _, v := range cv.BuiltInVariants {
				Expect(ids[v.ID]).To(BeFalse(), "Duplicate ID: "+v.ID)
				ids[v.ID] = true
			}
		})

		It("should all be marked as built-in", func() {
			for _, v := range cv.BuiltInVariants {
				Expect(v.IsBuiltIn).To(BeTrue(), "Variant not marked built-in: "+v.ID)
			}
		})

		It("should have valid base structures", func() {
			validStructures := map[cv.CVStructure]bool{
				cv.CVStructureStandard:   true,
				cv.CVStructureNarrative:  true,
				cv.CVStructureConsulting: true,
				cv.CVStructureHighlights: true,
			}
			for _, v := range cv.BuiltInVariants {
				Expect(validStructures[v.BaseStructure]).To(BeTrue(),
					"Invalid structure for "+v.ID+": "+string(v.BaseStructure))
			}
		})

		Context("senior_backend variants", func() {
			It("should have 4 variants for senior_backend role", func() {
				count := 0
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisSeniorBackend {
						count++
					}
				}
				Expect(count).To(Equal(4))
			})

			It("should use standard structure for full/standard/short", func() {
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisSeniorBackend {
						if v.LengthFormat == cv.LengthUltraShort {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureHighlights))
						} else {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureStandard))
						}
					}
				}
			})
		})

		Context("staff_principal variants", func() {
			It("should have 4 variants for staff_principal role", func() {
				count := 0
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisStaffPrincipal {
						count++
					}
				}
				Expect(count).To(Equal(4))
			})

			It("should use standard structure for full/standard/short", func() {
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisStaffPrincipal {
						if v.LengthFormat == cv.LengthUltraShort {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureHighlights))
						} else {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureStandard))
						}
					}
				}
			})
		})

		Context("consulting variants", func() {
			It("should have 4 variants for consulting role", func() {
				count := 0
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisConsulting {
						count++
					}
				}
				Expect(count).To(Equal(4))
			})

			It("should use consulting structure for full/standard/short", func() {
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisConsulting {
						if v.LengthFormat == cv.LengthUltraShort {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureHighlights))
						} else {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureConsulting))
						}
					}
				}
			})
		})

		Context("language_agnostic variants", func() {
			It("should have 4 variants for language_agnostic role", func() {
				count := 0
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisLanguageAgnostic {
						count++
					}
				}
				Expect(count).To(Equal(4))
			})

			It("should use narrative structure for full/standard/short", func() {
				for _, v := range cv.BuiltInVariants {
					if v.RoleEmphasis == cv.RoleEmphasisLanguageAgnostic {
						if v.LengthFormat == cv.LengthUltraShort {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureHighlights))
						} else {
							Expect(v.BaseStructure).To(Equal(cv.CVStructureNarrative))
						}
					}
				}
			})
		})
	})

	Describe("VariantService", func() {
		var service cv.VariantService

		BeforeEach(func() {
			service = cv.NewVariantService()
		})

		Describe("ListVariants", func() {
			It("should return all 16 variants", func() {
				variants := service.ListVariants()
				Expect(variants).To(HaveLen(16))
			})
		})

		Describe("GetVariant", func() {
			It("should return variant by ID", func() {
				variant, err := service.GetVariant("senior_backend_standard")
				Expect(err).NotTo(HaveOccurred())
				Expect(variant).NotTo(BeNil())
				Expect(variant.ID).To(Equal("senior_backend_standard"))
				Expect(variant.RoleEmphasis).To(Equal(cv.RoleEmphasisSeniorBackend))
				Expect(variant.LengthFormat).To(Equal(cv.LengthStandard))
			})

			It("should return error for unknown ID", func() {
				variant, err := service.GetVariant("unknown_variant")
				Expect(err).To(HaveOccurred())
				Expect(variant).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})

		Describe("GetVariantByDimensions", func() {
			It("should return variant for role and length combination", func() {
				variant, err := service.GetVariantByDimensions(
					cv.RoleEmphasisConsulting,
					cv.LengthShort,
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(variant).NotTo(BeNil())
				Expect(variant.ID).To(Equal("consulting_short"))
				Expect(variant.RoleEmphasis).To(Equal(cv.RoleEmphasisConsulting))
				Expect(variant.LengthFormat).To(Equal(cv.LengthShort))
			})

			It("should return correct variant for all combinations", func() {
				roles := []cv.RoleEmphasis{
					cv.RoleEmphasisSeniorBackend,
					cv.RoleEmphasisStaffPrincipal,
					cv.RoleEmphasisConsulting,
					cv.RoleEmphasisLanguageAgnostic,
				}
				lengths := []cv.LengthFormat{
					cv.LengthFull,
					cv.LengthStandard,
					cv.LengthShort,
					cv.LengthUltraShort,
				}

				for _, role := range roles {
					for _, length := range lengths {
						variant, err := service.GetVariantByDimensions(role, length)
						Expect(err).NotTo(HaveOccurred(), "Failed for %s/%s", role, length)
						Expect(variant.RoleEmphasis).To(Equal(role))
						Expect(variant.LengthFormat).To(Equal(length))
					}
				}
			})
		})

		Describe("ListRoleEmphases", func() {
			It("should return 4 role emphases with metadata", func() {
				emphases := service.ListRoleEmphases()
				Expect(emphases).To(HaveLen(4))

				// Check each has required fields
				for _, e := range emphases {
					Expect(e.ID).NotTo(BeEmpty())
					Expect(e.Name).NotTo(BeEmpty())
					Expect(e.Description).NotTo(BeEmpty())
				}
			})

			It("should include all role emphasis types", func() {
				emphases := service.ListRoleEmphases()
				ids := make([]cv.RoleEmphasis, len(emphases))
				for i, e := range emphases {
					ids[i] = e.ID
				}

				Expect(ids).To(ContainElements(
					cv.RoleEmphasisSeniorBackend,
					cv.RoleEmphasisStaffPrincipal,
					cv.RoleEmphasisConsulting,
					cv.RoleEmphasisLanguageAgnostic,
				))
			})
		})

		Describe("ListLengthFormats", func() {
			It("should return 4 length formats with metadata", func() {
				formats := service.ListLengthFormats()
				Expect(formats).To(HaveLen(4))

				// Check each has required fields
				for _, f := range formats {
					Expect(f.ID).NotTo(BeEmpty())
					Expect(f.Name).NotTo(BeEmpty())
					Expect(f.Description).NotTo(BeEmpty())
				}
			})

			It("should include all length format types", func() {
				formats := service.ListLengthFormats()
				ids := make([]cv.LengthFormat, len(formats))
				for i, f := range formats {
					ids[i] = f.ID
				}

				Expect(ids).To(ContainElements(
					cv.LengthFull,
					cv.LengthStandard,
					cv.LengthShort,
					cv.LengthUltraShort,
				))
			})
		})
	})
})
