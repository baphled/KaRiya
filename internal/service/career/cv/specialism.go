package cv

// DeriveSpecialism returns the CV specialism based on technology, position, and sector.
//
// Expected:
//   - technology must be a valid technology name (e.g. "Go", "Ruby", "JavaScript").
//   - position must be a valid position type (e.g. "backend", "frontend", "devops").
//   - sector must be a valid sector name (e.g. "startup", "enterprise", "public-sector").
//
// Returns:
//   - A string representing the derived specialism, or "General Engineering" if no mapping exists.
//
// Side effects:
//   - None.
func DeriveSpecialism(technology, position, sector string) string {
	mapping := map[string]string{
		"Go-backend-startup":                "Platform Engineering",
		"Ruby-backend-enterprise":           "Enterprise Engineering",
		"JavaScript-frontend-public-sector": "Public Sector Engineering",
		"DevOps-ai/ml-startup":              "AI/ML Engineering",
		"Fullstack-backend-startup":         "Platform Engineering",
		"Go-frontend-enterprise":            "Enterprise Engineering",
		"Ruby-devops-public-sector":         "Public Sector Engineering",
	}
	key := technology + "-" + position + "-" + sector
	if val, ok := mapping[key]; ok {
		return val
	}
	return "General Engineering"
}
