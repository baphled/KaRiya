package fixtures

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// ContentGroup creates a minimal SectionContentGroup with the given header.
func ContentGroup(header string) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header: header,
	}
}

// ContentGroupWith creates a SectionContentGroup with dates.
func ContentGroupWith(header, startDate, endDate string) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header:    header,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// ContentGroupWithBullets creates a SectionContentGroup with bullets attached.
func ContentGroupWithBullets(header string, bullets []*career.CVBullet) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header:  header,
		Bullets: bullets,
	}
}

// ContentGroupFull creates a fully-populated SectionContentGroup.
func ContentGroupFull(header, startDate, endDate string, bullets []*career.CVBullet) *career.SectionContentGroup {
	return &career.SectionContentGroup{
		Header:    header,
		StartDate: startDate,
		EndDate:   endDate,
		Bullets:   bullets,
	}
}
