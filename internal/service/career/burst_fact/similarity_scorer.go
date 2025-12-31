package burst_fact

import (
	"strings"
	"unicode"
)

// EventSimilarityInput represents event data used for similarity scoring
type EventSimilarityInput struct {
	Text     string
	Keywords []string
	Company  string
	Project  string
}

// SimilarityScorer provides methods to compute event similarity for burst detection
type SimilarityScorer struct{}

// NewSimilarityScorer creates a new SimilarityScorer instance
func NewSimilarityScorer() *SimilarityScorer {
	return &SimilarityScorer{}
}

// TextSimilarity computes similarity between two event texts using token overlap
// Returns a score from 0.0 (completely different) to 1.0 (identical)
// Algorithm: Jaccard similarity on word tokens
func (s *SimilarityScorer) TextSimilarity(text1, text2 string) float64 {
	if text1 == text2 {
		return 1.0
	}

	if text1 == "" || text2 == "" {
		if text1 == "" && text2 == "" {
			return 1.0
		}
		return 0.0
	}

	tokens1 := s.tokenize(text1)
	tokens2 := s.tokenize(text2)

	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	intersection := s.countIntersection(tokens1, tokens2)
	union := len(tokens1) + len(tokens2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// KeywordOverlapScore computes similarity based on shared keywords/tags
// Returns a score from 0.0 (no overlap) to 1.0 (identical keywords)
// Algorithm: Jaccard similarity on keyword sets
func (s *SimilarityScorer) KeywordOverlapScore(keywords1, keywords2 []string) float64 {
	if len(keywords1) == 0 && len(keywords2) == 0 {
		return 1.0
	}

	if len(keywords1) == 0 || len(keywords2) == 0 {
		return 0.0
	}

	// Normalize keywords to lowercase for comparison
	normalized1 := s.normalizeStringSlice(keywords1)
	normalized2 := s.normalizeStringSlice(keywords2)

	intersection := 0
	for _, kw1 := range normalized1 {
		for _, kw2 := range normalized2 {
			if kw1 == kw2 {
				intersection++
				break
			}
		}
	}

	union := len(normalized1) + len(normalized2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// CompanyMatchScore returns 1.0 if companies match (case-insensitive), 0.0 otherwise
// Treats empty companies as matching (no company specified)
func (s *SimilarityScorer) CompanyMatchScore(company1, company2 string) float64 {
	if company1 == "" && company2 == "" {
		return 1.0
	}

	if company1 == "" || company2 == "" {
		return 0.0
	}

	if strings.EqualFold(company1, company2) {
		return 1.0
	}

	return 0.0
}

// ProjectMatchScore returns 1.0 if projects match (case-insensitive), 0.0 otherwise
// Treats empty projects as matching (no project specified)
func (s *SimilarityScorer) ProjectMatchScore(project1, project2 string) float64 {
	if project1 == "" && project2 == "" {
		return 1.0
	}

	if project1 == "" || project2 == "" {
		return 0.0
	}

	if strings.EqualFold(project1, project2) {
		return 1.0
	}

	return 0.0
}

// CombinedSimilarityScore computes weighted similarity across all event components
// Weights:
//   - Text similarity: 50% (most important for semantic relevance)
//   - Keyword overlap: 20% (tag-based classification)
//   - Company match: 15% (organizational context)
//   - Project match: 15% (project scope)
//
// Returns a score from 0.0 to 1.0
func (s *SimilarityScorer) CombinedSimilarityScore(event1, event2 EventSimilarityInput) float64 {
	textScore := s.TextSimilarity(event1.Text, event2.Text)
	keywordScore := s.KeywordOverlapScore(event1.Keywords, event2.Keywords)
	companyScore := s.CompanyMatchScore(event1.Company, event2.Company)
	projectScore := s.ProjectMatchScore(event1.Project, event2.Project)

	// Apply weights: text=50%, keywords=20%, company=15%, project=15%
	const (
		textWeight    = 0.50
		keywordWeight = 0.20
		companyWeight = 0.15
		projectWeight = 0.15
	)

	combined := (textScore * textWeight) +
		(keywordScore * keywordWeight) +
		(companyScore * companyWeight) +
		(projectScore * projectWeight)

	return combined
}

// tokenize splits text into lowercase word tokens, filtering out empty strings
func (s *SimilarityScorer) tokenize(text string) []string {
	text = strings.ToLower(text)

	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Filter out very short tokens (noise)
	var filtered []string
	for _, token := range tokens {
		if len(token) > 2 {
			filtered = append(filtered, token)
		}
	}

	return filtered
}

// countIntersection counts matching tokens between two slices
func (s *SimilarityScorer) countIntersection(tokens1, tokens2 []string) int {
	count := 0
	for _, t1 := range tokens1 {
		for _, t2 := range tokens2 {
			if t1 == t2 {
				count++
				break
			}
		}
	}
	return count
}

// normalizeStringSlice converts all strings to lowercase
func (s *SimilarityScorer) normalizeStringSlice(strs []string) []string {
	normalized := make([]string, len(strs))
	for i, str := range strs {
		normalized[i] = strings.ToLower(str)
	}
	return normalized
}
