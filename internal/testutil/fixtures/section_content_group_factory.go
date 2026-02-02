package fixtures

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// ContentGroup creates a minimal SectionContentGroup with the given header.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.SectionContentGroup ready for use.
//
// Side effects:
//   - None.
func ContentGroup(header string) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header: header,
	}
}

// ContentGroupWith creates a SectionContentGroup with dates.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized career.SectionContentGroup ready for use.
//
// Side effects:
//   - None.
func ContentGroupWith(header, startDate, endDate string) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header:    header,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// ContentGroupWithBullets creates a SectionContentGroup with bullets attached.
//
// Expected:
//   - Must be a valid string.
//   - cvbullet must be valid.
//
// Returns:
//   - A fully initialized career.SectionContentGroup ready for use.
//
// Side effects:
//   - None.
func ContentGroupWithBullets(header string, bullets []*career.CVBullet) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header:  header,
		Bullets: bullets,
	}
}

// ContentGroupFull creates a fully-populated SectionContentGroup.
//
// Expected:
//   - Must be a valid string.
//   - cvbullet must be valid.
//
// Returns:
//   - A fully initialized career.SectionContentGroup ready for use.
//
// Side effects:
//   - None.
func ContentGroupFull(header, startDate, endDate string, bullets []*career.CVBullet) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header:    header,
		StartDate: startDate,
		EndDate:   endDate,
		Bullets:   bullets,
	}
}
