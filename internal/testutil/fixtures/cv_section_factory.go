package fixtures

import (
	"fmt"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

// CVSectionFactory creates CVSection fixtures with realistic fake data.
var CVSectionFactory = factory.NewFactory(
	&career.CVSection{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("section-%d", n), nil
}).Attr("CVViewID", func(args factory.Args) (interface{}, error) {
	return fmt.Sprintf("cv-%d", gofakeit.Number(1, 10)), nil
}).Attr("SectionType", func(args factory.Args) (interface{}, error) {
	types := []string{"experience", "projects", "skills", "summary"}
	return types[gofakeit.Number(0, len(types)-1)], nil
}).Attr("Title", func(args factory.Args) (interface{}, error) {
	titles := []string{"Professional Experience", "Key Projects", "Technical Skills", "Executive Summary"}
	return titles[gofakeit.Number(0, len(titles)-1)], nil
}).SeqInt("Order", func(n int) (interface{}, error) {
	return n, nil
})

// CVSection creates a minimal valid CVSection with the given ID and CV view ID.
func CVSection(id, cvViewID string) *career.CVSection {
	return &career.CVSection{
		ID:          id,
		CVViewID:    cvViewID,
		SectionType: "experience",
		Title:       "Professional Experience",
		Order:       0,
	}
}

// CVSectionWith creates a CVSection with custom fields.
func CVSectionWith(id, cvViewID, sectionType, title string, order int) *career.CVSection {
	return &career.CVSection{
		ID:          id,
		CVViewID:    cvViewID,
		SectionType: sectionType,
		Title:       title,
		Order:       order,
	}
}

// CVSectionWithContent creates a CVSection with content groups.
func CVSectionWithContent(id, cvViewID string, content []*career.SectionContentGroup) *career.CVSection {
	section := CVSection(id, cvViewID)
	section.Content = content
	return section
}

// CVSectionWithSummary creates a summary-type CVSection.
func CVSectionWithSummary(id, cvViewID, summary string) *career.CVSection {
	return &career.CVSection{
		ID:          id,
		CVViewID:    cvViewID,
		SectionType: "summary",
		Title:       "Summary",
		Order:       0,
		Summary:     summary,
	}
}

// CVSections creates n sections for a given CV view.
func CVSections(n int, cvViewID string) []*career.CVSection {
	sectionTypes := []string{"experience", "projects", "skills", "summary"}
	sections := make([]*career.CVSection, n)
	for i := 0; i < n; i++ {
		sectionType := sectionTypes[i%len(sectionTypes)]
		sections[i] = CVSectionWith(
			fmt.Sprintf("section-%d", i+1),
			cvViewID,
			sectionType,
			fmt.Sprintf("Section %d", i+1),
			i,
		)
	}
	return sections
}
