package intents

// NarrativeProfileData contains hardcoded profile data for narrative CV structure.
// In Phase 5, this will be configurable via Configure System → Profile.
type NarrativeProfileData struct {
	Name              string
	Role              string
	Location          string
	Email             string
	GitHub            string
	Portfolio         string
	CoreStrengths     []string
	Languages         string
	Frontend          string
	Systems           string
	ValuePropositions []string
}

// DefaultNarrativeProfile returns the default narrative profile data.
// This is hardcoded for Phase 3 and will be made configurable in Phase 5.
func DefaultNarrativeProfile() *NarrativeProfileData {
	return &NarrativeProfileData{
		Name:      "Yomi Colledge",
		Role:      "Senior Software Engineer / Technical Consultant",
		Location:  "Remote (UK)",
		Email:     "yomi@boodah.net",
		GitHub:    "https://github.com/baphled",
		Portfolio: "http://boodah.net",
		CoreStrengths: []string{
			"Language-agnostic backend and systems engineering",
			"System design and architectural ownership",
			"Cross-functional collaboration and mentorship",
			"Pragmatic problem-solving with production focus",
			"Technical leadership without formal authority",
		},
		Languages: "Ruby, Go, PHP, C/C++, JavaScript, Shell",
		Frontend:  "Vue.js",
		Systems:   "Linux, SQL, APIs, CI/CD, automation",
		ValuePropositions: []string{
			"Languages as tools, not identity",
			"Calm handling of complexity",
			"Ownership without ego",
			"Clear technical communication",
			"Production-first mindset",
		},
	}
}

// MinConfidenceForNarrative is the minimum confidence score for bullets to appear in narrative CV.
const MinConfidenceForNarrative = 0.75
