package cv

import (
	"github.com/baphled/kariya/internal/constants"
)

// TechnologyFocus is an alias to constants.TechnologyFocus for backward compatibility.
type TechnologyFocus = constants.TechnologyFocus

// Technology focus constants for backward compatibility.
const (
	TechnologyFocusLanguageAgnostic = constants.TechnologyFocusLanguageAgnostic
	TechnologyFocusGeneralist       = constants.TechnologyFocusGeneralist
	TechnologyFocusSpecialist       = constants.TechnologyFocusSpecialist
)

// FocusArea is an alias to constants.FocusArea for backward compatibility.
type FocusArea = constants.FocusArea

// Focus area constants for backward compatibility.
const (
	FocusAreaBackend   = constants.FocusAreaBackend
	FocusAreaFrontend  = constants.FocusAreaFrontend
	FocusAreaFullstack = constants.FocusAreaFullstack
	FocusAreaDevOps    = constants.FocusAreaDevOps
)

// LengthFormat defines CV density/length
type LengthFormat string

const (
	// LengthFull includes complete career history with all relevant details
	LengthFull LengthFormat = "full"

	// LengthStandard is a balanced 2-3 page format with recent 10 years emphasis
	LengthStandard LengthFormat = "standard"

	// LengthShort is a concise 1-2 page format focusing on recent 5 years
	LengthShort LengthFormat = "short"

	// LengthUltraShort is a 1-page highlights format with only top achievements
	LengthUltraShort LengthFormat = "ultra_short"
)
