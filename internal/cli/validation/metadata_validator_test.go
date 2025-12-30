package validation_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/validation"
)

var _ = Describe("MetadataValidator", func() {
	var (
		validator *validation.MetadataValidator
	)

	BeforeEach(func() {
		validator = validation.NewMetadataValidator()
	})

	Describe("ValidateDate", func() {
		It("should accept valid past date", func() {
			pastDate := time.Now().Add(-24 * time.Hour)
			err := validator.ValidateDate(pastDate)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept today's date", func() {
			today := time.Now()
			err := validator.ValidateDate(today)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject future date", func() {
			futureDate := time.Now().Add(24 * time.Hour)
			err := validator.ValidateDate(futureDate)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot be in the future"))
		})

		It("should accept date in 1900", func() {
			date := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
			err := validator.ValidateDate(date)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject date before 1900", func() {
			date := time.Date(1899, 12, 31, 0, 0, 0, 0, time.UTC)
			err := validator.ValidateDate(date)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("before 1900"))
		})

		It("should accept date in 2000s", func() {
			date := time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC)
			err := validator.ValidateDate(date)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ValidateCompany", func() {
		It("should accept valid company name", func() {
			err := validator.ValidateCompany("TechCorp Inc.")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept empty company (optional field)", func() {
			err := validator.ValidateCompany("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept company at max length (200)", func() {
			company := strings.Repeat("a", 200)
			err := validator.ValidateCompany(company)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject company exceeding 200 characters", func() {
			company := strings.Repeat("a", 201)
			err := validator.ValidateCompany(company)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot exceed 200"))
		})

		It("should trim whitespace when validating length", func() {
			company := "  " + strings.Repeat("a", 200) + "  "
			err := validator.ValidateCompany(company)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept whitespace-only as empty (optional)", func() {
			err := validator.ValidateCompany("   ")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ValidateProject", func() {
		It("should accept valid project name", func() {
			err := validator.ValidateProject("Platform Migration")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept empty project (optional field)", func() {
			err := validator.ValidateProject("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept project at max length (200)", func() {
			project := strings.Repeat("a", 200)
			err := validator.ValidateProject(project)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject project exceeding 200 characters", func() {
			project := strings.Repeat("a", 201)
			err := validator.ValidateProject(project)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot exceed 200"))
		})

		It("should trim whitespace when validating length", func() {
			project := "  " + strings.Repeat("a", 200) + "  "
			err := validator.ValidateProject(project)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ValidateTags", func() {
		It("should accept empty tags (optional field)", func() {
			err := validator.ValidateTags([]string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept valid single tag", func() {
			err := validator.ValidateTags([]string{"technical"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept valid multiple tags", func() {
			err := validator.ValidateTags([]string{"technical", "leadership", "achievement"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept tags at max count (8)", func() {
			tags := []string{"technical", "leadership", "achievement", "project", "consulting", "research", "product", "mentoring"}
			err := validator.ValidateTags(tags)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject tags exceeding max count (8)", func() {
			tags := []string{"technical", "leadership", "achievement", "project", "consulting", "research", "product", "mentoring", "extra"}
			err := validator.ValidateTags(tags)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("more than 8"))
		})

		It("should reject duplicate tags (case-insensitive)", func() {
			err := validator.ValidateTags([]string{"technical", "Technical"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Duplicate"))
		})

		It("should reject invalid tag", func() {
			err := validator.ValidateTags([]string{"invalid-tag"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Invalid tag"))
		})

		It("should accept case-insensitive tags", func() {
			err := validator.ValidateTags([]string{"TECHNICAL", "Leadership", "ACHIEVEMENT"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should allow all valid tags", func() {
			validTags := []string{"project", "achievement", "leadership", "technical", "consulting", "research", "product", "mentoring"}
			err := validator.ValidateTags(validTags)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ValidateCategories", func() {
		It("should accept empty categories (optional field)", func() {
			err := validator.ValidateCategories([]string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept valid single category", func() {
			err := validator.ValidateCategories([]string{"technical"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept valid multiple categories", func() {
			err := validator.ValidateCategories([]string{"technical", "leadership", "product"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept case-insensitive categories", func() {
			err := validator.ValidateCategories([]string{"TECHNICAL", "Leadership", "PRODUCT"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject invalid category", func() {
			err := validator.ValidateCategories([]string{"invalid-category"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Invalid category"))
		})

		It("should allow all valid categories", func() {
			validCategories := []string{"technical", "leadership", "product", "consulting", "research", "mentoring"}
			err := validator.ValidateCategories(validCategories)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept duplicate categories (allowed for categories)", func() {
			err := validator.ValidateCategories([]string{"technical", "technical"})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
