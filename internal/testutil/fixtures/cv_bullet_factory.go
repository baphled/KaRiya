package fixtures

import (
	"fmt"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

// CVBulletFactory creates CVBullet fixtures with realistic fake data.
var CVBulletFactory = factory.NewFactory(
	&career.CVBullet{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("bullet-%d", n), nil
}).Attr("SectionID", func(args factory.Args) (interface{}, error) {
	return fmt.Sprintf("section-%d", gofakeit.Number(1, 10)), nil
}).Attr("Text", func(args factory.Args) (interface{}, error) {
	templates := []string{
		"Reduced %s latency by %d%% through system redesign",
		"Led cross-functional team of %d to deliver %s on schedule",
		"Implemented %s pipeline reducing deployment time by %d%%",
		"Designed %s architecture serving %dK daily active users",
	}
	template := templates[gofakeit.Number(0, len(templates)-1)]
	return fmt.Sprintf(template, gofakeit.BuzzWord(), gofakeit.Number(20, 80)), nil
}).Attr("SourceEventIDs", func(args factory.Args) (interface{}, error) {
	return []string{fmt.Sprintf("event-%d", gofakeit.Number(1, 100))}, nil
}).Attr("Rank", func(args factory.Args) (interface{}, error) {
	return gofakeit.Float64Range(0.3, 0.95), nil
}).Attr("InclusionReason", func(args factory.Args) (interface{}, error) {
	reasons := []string{"ownership", "contribution", "strategy", "execution", "outcome", "fact_extraction"}
	return reasons[gofakeit.Number(0, len(reasons)-1)], nil
}).Attr("Confidence", func(args factory.Args) (interface{}, error) {
	return gofakeit.Float64Range(0.5, 1.0), nil
})

// CVBullet creates a minimal valid CVBullet with the given ID and section ID.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.CVBullet ready for use.
//
// Side effects:
//   - None.
func CVBullet(id, sectionID string) *career.CVBullet {
	return &career.CVBullet{
		ID:              id,
		SectionID:       sectionID,
		Text:            "Test bullet " + id,
		SourceEventIDs:  []string{"event-1"},
		Rank:            0.8,
		InclusionReason: "ownership",
		Confidence:      0.9,
	}
}

// CVBulletWith creates a CVBullet with custom text.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.CVBullet ready for use.
//
// Side effects:
//   - None.
func CVBulletWith(id, sectionID, text string) *career.CVBullet {
	return &career.CVBullet{
		ID:              id,
		SectionID:       sectionID,
		Text:            text,
		SourceEventIDs:  []string{"event-1"},
		Rank:            0.8,
		InclusionReason: "ownership",
		Confidence:      0.9,
	}
}

// CVBulletWithScores creates a CVBullet with all scoring fields populated.
//
// Expected:
//   - Must be a valid string.
//   - float64 must be valid.
//
// Returns:
//   - A fully initialized career.CVBullet ready for use.
//
// Side effects:
//   - None.
func CVBulletWithScores(id, sectionID, text string, rank, confidence, roleScore, audienceScore float64) *career.CVBullet {
	return &career.CVBullet{
		ID:              id,
		SectionID:       sectionID,
		Text:            text,
		SourceEventIDs:  []string{"event-1"},
		Rank:            rank,
		InclusionReason: "ownership",
		Confidence:      confidence,
		RoleScore:       roleScore,
		AudienceScore:   audienceScore,
	}
}

// CVBulletWithSources creates a CVBullet with specific source event and fact IDs.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.CVBullet ready for use.
//
// Side effects:
//   - None.
func CVBulletWithSources(id, sectionID, text string, eventIDs, factIDs []string) *career.CVBullet {
	return &career.CVBullet{
		ID:              id,
		SectionID:       sectionID,
		Text:            text,
		SourceEventIDs:  eventIDs,
		SourceFactIDs:   factIDs,
		Rank:            0.8,
		InclusionReason: "fact_extraction",
		Confidence:      0.9,
	}
}

// CVBullets creates n bullets for a given section.
//
// Expected:
//   - int must be valid.
//   - Must be a valid string.
//
// Returns:
//   - A []*career.CVBullet value.
//
// Side effects:
//   - None.
func CVBullets(n int, sectionID string) []*career.CVBullet {
	bullets := make([]*career.CVBullet, n)
	for i := range n {
		bullet, ok := CVBulletFactory.MustCreate().(*career.CVBullet)
		if !ok {
			continue
		}
		bullet.SectionID = sectionID
		bullets[i] = bullet
	}
	return bullets
}

// CVBulletVal creates a CVBullet value (not pointer) with the given ID and section ID.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A career.CVBullet value.
//
// Side effects:
//   - None.
func CVBulletVal(id, sectionID string) career.CVBullet {
	return *CVBullet(id, sectionID)
}
