package career

// Event represents a career event (minimal test double).
type Event struct {
	ID   string
	Text string
}

// Fact represents a career fact (minimal test double).
type Fact struct {
	ID   string
	Text string
}

// Burst represents a career burst (minimal test double).
type Burst struct {
	ID   string
	Name string
}

// Skill represents a career skill (minimal test double).
type Skill struct {
	ID   string
	Name string
}

// CVView represents a CV view (minimal test double).
type CVView struct {
	ID   string
	Name string
}

// CVSection represents a CV section (minimal test double).
type CVSection struct {
	ID string
}

// CVBullet represents a CV bullet (minimal test double).
type CVBullet struct {
	ID string
}

// CVConfig represents a CV config (minimal test double).
type CVConfig struct {
	Name string
}

// SectionContentGroup represents section content (minimal test double).
type SectionContentGroup struct {
	Header string
}

// unexportedTestDouble simulates an unexported test type (should be skipped by analyzer).
type unexportedTestDouble struct {
	name string
}
