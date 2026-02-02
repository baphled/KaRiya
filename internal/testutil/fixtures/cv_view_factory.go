package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

// CVViewFactory creates CVView fixtures with realistic fake data.
var CVViewFactory = factory.NewFactory(
	&career.CVView{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("cv-%d", n), nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	roles := []string{"Staff Engineer", "Principal Engineer", "Engineering Manager", "Senior IC"}
	companies := []string{"at FAANG", "for Startup", "at Scale-up", "for Enterprise"}
	return fmt.Sprintf("%s %s", roles[gofakeit.Number(0, len(roles)-1)], companies[gofakeit.Number(0, len(companies)-1)]), nil
}).Attr("TargetRole", func(args factory.Args) (interface{}, error) {
	roles := []string{"staff", "principal", "em", "senior_ic"}
	return roles[gofakeit.Number(0, len(roles)-1)], nil
}).Attr("TargetAudience", func(args factory.Args) (interface{}, error) {
	audiences := []string{"hiring_manager", "recruiter", "peer"}
	return audiences[gofakeit.Number(0, len(audiences)-1)], nil
}).Attr("GeneratedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("SourceEventCount", func(args factory.Args) (interface{}, error) {
	return gofakeit.Number(5, 50), nil
}).Attr("SourceFactCount", func(args factory.Args) (interface{}, error) {
	return gofakeit.Number(3, 30), nil
})

// CVView creates a minimal valid CVView with the given ID.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func CVView(id string) *career.CVView {
	return &career.CVView{
		ID:               id,
		Name:             "Test CV " + id,
		TargetRole:       "staff",
		TargetAudience:   "hiring_manager",
		GeneratedAt:      time.Now(),
		SourceEventCount: 10,
		SourceFactCount:  5,
	}
}

// CVViewWith creates a CVView with custom key fields.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func CVViewWith(id, name, targetRole, targetAudience string) *career.CVView {
	return &career.CVView{
		ID:               id,
		Name:             name,
		TargetRole:       targetRole,
		TargetAudience:   targetAudience,
		GeneratedAt:      time.Now(),
		SourceEventCount: 10,
		SourceFactCount:  5,
	}
}

// CVViewWithSections creates a CVView with attached sections.
//
// Expected:
//   - Must be a valid string.
//   - cvsection must be valid.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func CVViewWithSections(id string, sections []*career.CVSection) *career.CVView {
	cv := CVView(id)
	cv.Sections = sections
	return cv
}

// CVViewWithCounts creates a CVView with specific source counts.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func CVViewWithCounts(id string, eventCount, factCount int) *career.CVView {
	cv := CVView(id)
	cv.SourceEventCount = eventCount
	cv.SourceFactCount = factCount
	return cv
}

// CVViewWithFilters creates a CVView with event filters.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func CVViewWithFilters(id string, filters map[string]interface{}) *career.CVView {
	cv := CVView(id)
	cv.EventFilters = filters
	return cv
}
