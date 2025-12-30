package importer_test

import (
	"bytes"
	"time"

	"github.com/baphled/kariya/internal/cli/importer"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)


var _ = Describe("CSV Parser", func() {
	var (
		parser *importer.CSVParser
		
	)

	BeforeEach(func() {
		parser = importer.NewCSVParser([]*career.CareerEvent{})
		
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
			existingEvent := &career.CareerEvent{
				ID:        "123",
				Text:      "Existing event",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Company:   "MyCompany",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			parserWithExisting := importer.NewCSVParser([]*career.CareerEvent{existingEvent})

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
})
