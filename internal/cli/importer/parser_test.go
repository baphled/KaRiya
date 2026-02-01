package importer_test

import (
	"bytes"
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/importer"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSV Parser", func() {
	var (
		parser *importer.CSVParser
	)

	BeforeEach(func() {
		parser = importer.NewCSVParser([]*career.Event{}, nil, context.Background())

	})

	Describe("Basic CSV Parsing", func() {
		It("should parse valid CSV with required fields", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test event,2024-01,Technical,technical,MyProject,MyCompany`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Text).To(Equal("Test event"))
			Expect(rows[0].Event.Company).To(Equal("MyCompany"))
			Expect(rows[0].Event.Project).To(Equal("MyProject"))
		})

		It("should parse multiple rows", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Event 1,2024-01,Technical,technical,Project1,Company1
Event 2,2024-02,Leadership,leadership,Project2,Company2
Event 3,2024-03,Product,product,Project3,Company3`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(3))
			Expect(rows[0].Event.Text).To(Equal("Event 1"))
			Expect(rows[1].Event.Text).To(Equal("Event 2"))
			Expect(rows[2].Event.Text).To(Equal("Event 3"))
		})

		It("should handle optional fields", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Minimal event,2024-01,,,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Company).To(Equal(""))
			Expect(rows[0].Event.Project).To(Equal(""))
		})
	})

	Describe("Date Parsing", func() {
		It("should parse YYYY-MM format", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Date.Year()).To(Equal(2024))
			Expect(rows[0].Event.Date.Month()).To(Equal(time.January))
		})

		It("should parse YYYY-MM-DD format", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01-15,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeTrue())
		})

		It("should reject invalid date format", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,invalid-date,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement(ContainSubstring("Invalid date format")))
		})

		It("should reject future dates", func() {
			futureDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
			csv := "Text,Date,Categories,Tags,Project,Company\nTest," + futureDate + ",Technical,technical,,"

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement(ContainSubstring("future")))
		})
	})

	Describe("Tag Validation", func() {
		It("should parse valid tags", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01,Technical,technical;leadership;mentoring,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Tags).To(HaveLen(3))
			Expect(rows[0].Event.Tags).To(ContainElement("technical"))
			Expect(rows[0].Event.Tags).To(ContainElement("leadership"))
		})

		It("should reject invalid tags", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01,Technical,invalid-tag,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement(ContainSubstring("Invalid tag")))
		})

		It("should normalize tags to lowercase", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01,Technical,TECHNICAL;Leadership,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Tags).To(ContainElement("technical"))
			Expect(rows[0].Event.Tags).To(ContainElement("leadership"))
		})
	})

	Describe("Category Validation", func() {
		It("should parse valid categories", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01,Technical;Leadership,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Categories).To(HaveLen(2))
		})

		It("should reject invalid categories", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test,2024-01,InvalidCategory,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement(ContainSubstring("Invalid category")))
		})
	})

	Describe("Validation Errors", func() {
		It("should report missing text", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
,2024-01,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement("Text is required"))
		})

		It("should report missing date", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test event,,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement("Date is required"))
		})

		It("should report text exceeding character limit", func() {
			longText := ""
			for i := 0; i < 2001; i++ {
				longText += "a"
			}
			csv := "Text,Date,Categories,Tags,Project,Company\n" + longText + ",2024-01,Technical,technical,,"

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsValid).To(BeFalse())
			Expect(rows[0].ValidationErrors).To(ContainElement(ContainSubstring("exceed")))
		})
	})

	Describe("Duplicate Detection", func() {
		It("should detect duplicates with existing events", func() {
			// Create an existing event
			existingEvent := fixtures.EventWith("123", "Existing event", "MyCompany", "")
			existingEvent.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			parserWithExisting := importer.NewCSVParser([]*career.Event{existingEvent}, nil, context.Background())

			csv := `Text,Date,Categories,Tags,Project,Company
Existing event,2024-01,Technical,technical,,MyCompany`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parserWithExisting.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsDuplicate).To(BeTrue())
			Expect(rows[0].DuplicateOf).To(Equal("123"))
		})

		It("should detect duplicates within parsed rows", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Duplicate event,2024-01,Technical,technical,,Company1
Duplicate event,2024-01,Technical,technical,,Company1
Different event,2024-02,Leadership,leadership,,Company2`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsDuplicate).To(BeFalse())
			Expect(rows[1].IsDuplicate).To(BeTrue())
			Expect(rows[2].IsDuplicate).To(BeFalse())
		})

		It("should not mark valid events as duplicates", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Event 1,2024-01,Technical,technical,,Company1
Event 2,2024-02,Leadership,leadership,,Company2`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].IsDuplicate).To(BeFalse())
			Expect(rows[1].IsDuplicate).To(BeFalse())
		})
	})

	Describe("Header Validation", func() {
		It("should reject CSV with missing Text column", func() {
			csv := `Date,Categories,Tags,Project,Company
2024-01,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			_, err := parser.Parse(reader)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required column: Text"))
		})

		It("should reject CSV with missing Date column", func() {
			csv := `Text,Categories,Tags,Project,Company
Test event,Technical,technical,,`

			reader := bytes.NewReader([]byte(csv))
			_, err := parser.Parse(reader)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required column: Date"))
		})

		It("should handle extra columns", func() {
			csv := `Text,Date,Categories,Tags,Project,Company,ExtraColumn
Test event,2024-01,Technical,technical,MyProject,MyCompany,Extra`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
		})
	})

	Describe("Field origin detection", func() {
		It("should detect when text field came from CSV", func() {
			// Arrange: Create a parsed row with text in RawData
			rawData := map[string]string{
				"Text": "Event from CSV",
				"Date": "2024-01-15",
			}

			// Act & Assert: Text should be detected as from CSV
			Expect(rawData["Text"]).To(Equal("Event from CSV"))
		})

		It("should detect when company field came from CSV", func() {
			// Arrange
			rawData := map[string]string{
				"Text":    "Event text",
				"Date":    "2024-01-15",
				"Company": "TechCorp",
			}

			// Act & Assert
			Expect(rawData["Company"]).To(Equal("TechCorp"))
		})

		It("should detect when project field came from CSV", func() {
			// Arrange
			rawData := map[string]string{
				"Text":    "Event text",
				"Date":    "2024-01-15",
				"Project": "Platform Migration",
			}

			// Act & Assert
			Expect(rawData["Project"]).To(Equal("Platform Migration"))
		})

		It("should detect when tags field came from CSV", func() {
			// Arrange
			rawData := map[string]string{
				"Text": "Event text",
				"Date": "2024-01-15",
				"Tags": "technical,leadership",
			}

			// Act & Assert
			Expect(rawData["Tags"]).To(Equal("technical,leadership"))
		})

		It("should detect missing fields (not from CSV)", func() {
			// Arrange: RawData without Company field
			rawData := map[string]string{
				"Text": "Event text",
				"Date": "2024-01-15",
			}

			// Act & Assert: Company field not in RawData
			_, hasCompany := rawData["Company"]
			Expect(hasCompany).To(BeFalse())
		})
	})

	Describe("Skills Column Handling", func() {
		It("should handle CSV without Skills column", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Test event,2024-01,Technical,technical,MyProject,MyCompany`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(BeEmpty())
		})

		It("should handle empty Skills column", func() {
			csv := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(BeEmpty())
		})

		It("should handle whitespace-only Skills column", func() {
			csv := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,   `

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(BeEmpty())
		})

		It("should handle empty skill names after splitting", func() {
			csv := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,;;`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(BeEmpty())
		})

		It("should handle mixed empty and whitespace skill names", func() {
			csv := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,; ;  ;`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(BeEmpty())
		})

		It("should skip processing skills when skillRepository is nil", func() {
			// Parser already initialized with nil skillRepository in BeforeEach
			csv := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,Go;Ruby`

			reader := bytes.NewReader([]byte(csv))
			rows, err := parser.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			// Skills should be empty because skillRepository is nil
			Expect(rows[0].Event.Skills).To(BeEmpty())
		})
	})

	Describe("Skills Import with memory.SkillRepository", func() {
		var (
			skillRepo       *careermemory.SkillRepository
			parserWithSkill *importer.CSVParser
			ctx             context.Context
		)

		BeforeEach(func() {
			ctx = context.Background()
			skillRepo = careermemory.NewSkillRepository()
			parserWithSkill = importer.NewCSVParser([]*career.Event{}, skillRepo, ctx)
		})

		It("should create new skills when not found in repository", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,Go;Ruby`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(HaveLen(2))

			goSkill, err := skillRepo.GetByName(ctx, "Go")
			Expect(err).NotTo(HaveOccurred())
			Expect(goSkill).NotTo(BeNil())
			Expect(goSkill.Category).To(Equal("backend"))

			rubySkill, err := skillRepo.GetByName(ctx, "Ruby")
			Expect(err).NotTo(HaveOccurred())
			Expect(rubySkill).NotTo(BeNil())
			Expect(rubySkill.Category).To(Equal("backend"))
		})

		It("should reuse existing skills from repository", func() {
			// Pre-create a skill
			existingSkill := fixtures.SkillWith("", "Go", "backend", "expert")
			err := skillRepo.Create(ctx, existingSkill)
			Expect(err).NotTo(HaveOccurred())

			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,Go;Ruby`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[0].Event.Skills).To(HaveLen(2))

			// Verify the existing skill ID was used
			Expect(rows[0].Event.Skills).To(ContainElement(existingSkill.ID))

			rubySkill, err := skillRepo.GetByName(ctx, "Ruby")
			Expect(err).NotTo(HaveOccurred())
			Expect(rubySkill.Category).To(Equal("backend"))
		})

		It("should handle multiple events with shared skills", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Event 1,2024-01,Technical,technical,Project1,Company1,Go;Docker
Event 2,2024-02,Technical,technical,Project2,Company2,Go;Kubernetes`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(2))
			Expect(rows[0].IsValid).To(BeTrue(), "Row 0 validation errors: %v", rows[0].ValidationErrors)
			Expect(rows[1].IsValid).To(BeTrue(), "Row 1 validation errors: %v", rows[1].ValidationErrors)
			Expect(rows[0].Event).NotTo(BeNil())
			Expect(rows[1].Event).NotTo(BeNil())

			// Both events should have Go skill
			goSkill, err := skillRepo.GetByName(ctx, "Go")
			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].Event.Skills).To(ContainElement(goSkill.ID))
			Expect(rows[1].Event.Skills).To(ContainElement(goSkill.ID))

			// Check other skills exist
			dockerSkill, err := skillRepo.GetByName(ctx, "Docker")
			Expect(err).NotTo(HaveOccurred())
			Expect(rows[0].Event.Skills).To(ContainElement(dockerSkill.ID))

			k8sSkill, err := skillRepo.GetByName(ctx, "Kubernetes")
			Expect(err).NotTo(HaveOccurred())
			Expect(rows[1].Event.Skills).To(ContainElement(k8sSkill.ID))
		})

		It("should trim whitespace from skill names", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,  Go  ;  Ruby  `

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].Event.Skills).To(HaveLen(2))

			// Verify skills were created with trimmed names
			goSkill, err := skillRepo.GetByName(ctx, "Go")
			Expect(err).NotTo(HaveOccurred())
			Expect(goSkill).NotTo(BeNil())

			rubySkill, err := skillRepo.GetByName(ctx, "Ruby")
			Expect(err).NotTo(HaveOccurred())
			Expect(rubySkill).NotTo(BeNil())
		})

		It("should handle case-insensitive skill matching", func() {
			existingSkill := fixtures.SkillWith("", "Go", "backend", "")
			err := skillRepo.Create(ctx, existingSkill)
			Expect(err).NotTo(HaveOccurred())

			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,go`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].Event.Skills).To(HaveLen(1))

			foundSkill, err := skillRepo.GetByName(ctx, "go")
			Expect(err).NotTo(HaveOccurred())
			Expect(foundSkill).NotTo(BeNil())
			Expect(foundSkill.ID).To(Equal(existingSkill.ID),
				"should reuse existing 'Go' skill when import has 'go'")
		})

		It("should skip empty skill names in semicolon-separated list", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,Go;;Ruby;  ;Python`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].Event.Skills).To(HaveLen(3))

			// Verify only non-empty skills were created
			skills, err := skillRepo.List(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(3))
		})

		It("should auto-categorize known skills from technology dictionary", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,PostgreSQL;Docker;React`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].Event.Skills).To(HaveLen(3))

			pgSkill, err := skillRepo.GetByName(ctx, "PostgreSQL")
			Expect(err).NotTo(HaveOccurred())
			Expect(pgSkill.Category).To(Equal("database"))

			dockerSkill, err := skillRepo.GetByName(ctx, "Docker")
			Expect(err).NotTo(HaveOccurred())
			Expect(dockerSkill.Category).To(Equal("devops"))

			reactSkill, err := skillRepo.GetByName(ctx, "React")
			Expect(err).NotTo(HaveOccurred())
			Expect(reactSkill.Category).To(Equal("frontend"))
		})

		It("should fall back to other for unknown skills", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,CustomFramework`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].Event.Skills).To(HaveLen(1))

			skill, err := skillRepo.GetByName(ctx, "CustomFramework")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Category).To(Equal("other"))
		})

		It("should auto-categorize case-insensitively", func() {
			csvData := `Text,Date,Categories,Tags,Project,Company,Skills
Test event,2024-01,Technical,technical,MyProject,MyCompany,postgresql;docker`

			reader := bytes.NewReader([]byte(csvData))
			rows, err := parserWithSkill.Parse(reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(1))
			Expect(rows[0].Event.Skills).To(HaveLen(2))

			pgSkill, err := skillRepo.GetByName(ctx, "postgresql")
			Expect(err).NotTo(HaveOccurred())
			Expect(pgSkill.Category).To(Equal("database"))

			dockerSkill, err := skillRepo.GetByName(ctx, "docker")
			Expect(err).NotTo(HaveOccurred())
			Expect(dockerSkill.Category).To(Equal("devops"))
		})
	})

})
