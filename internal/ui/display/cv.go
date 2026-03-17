package display

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// CVView is a presentation-only view of a generated CV.
type CVView struct {
	ID               string
	Name             string
	TargetRole       string
	TargetAudience   string
	EventFilters     map[string]interface{}
	GeneratedAt      time.Time
	SourceEventCount int
	SourceFactCount  int
	Sections         []CVSection
}

// CVSection is a presentation-only view of a CV section.
type CVSection struct {
	ID          string
	CVViewID    string
	SectionType string
	Title       string
	Order       int
	Content     []SectionContentGroup
	Summary     string
}

// SectionContentGroup is a presentation-only view of grouped CV content.
type SectionContentGroup struct {
	Header    string
	StartDate string
	EndDate   string
	Bullets   []CVBullet
}

// CVBullet is a presentation-only view of a CV bullet.
type CVBullet struct {
	ID                string
	SectionID         string
	Text              string
	SourceEventIDs    []string
	SourceFactIDs     []string
	Rank              float64
	InclusionReason   string
	Confidence        float64
	EnhancedText      string
	Category          string
	RoleScore         float64
	AudienceScore     float64
	MetricScore       float64
	ImpactScore       float64
	ImpactLevel       string
	KeywordMatches    []string
	AudienceRelevance map[string]float64
}

// CVConfig is a presentation-only view of CV generation configuration.
type CVConfig struct {
	Name                 string
	TargetRole           string
	TargetAudience       string
	EventFilters         map[string]interface{}
	TechnologyFocus      string
	SelectedTechnologies []string
	FocusArea            string
	LengthFormat         string
	SkillsFormat         string
	SkillsLimit          int
	SummaryHeading       string
	ProfileTitle         string
	WhatIBring           []string
	CoreStrengths        []string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// CVViewFromDomain converts a domain CV view to a display CV view.
//
// Expected:
//   - cvview must be valid.
//
// Returns:
//   - A CVView value.
//
// Side effects:
//   - None.
func CVViewFromDomain(view *career.CVView) CVView {
	if view == nil {
		return CVView{}
	}

	return CVView{
		ID:               view.ID,
		Name:             view.Name,
		TargetRole:       view.TargetRole,
		TargetAudience:   view.TargetAudience,
		EventFilters:     cloneStringInterfaceMap(view.EventFilters),
		GeneratedAt:      view.GeneratedAt,
		SourceEventCount: view.SourceEventCount,
		SourceFactCount:  view.SourceFactCount,
		Sections:         CVSectionsFromDomain(view.Sections),
	}
}

// CVViewsFromDomain converts domain CV views to display CV views.
//
// Expected:
//   - cvview must be valid.
//
// Returns:
//   - A []CVView value.
//
// Side effects:
//   - None.
func CVViewsFromDomain(views []*career.CVView) []CVView {
	if views == nil {
		return nil
	}

	result := make([]CVView, len(views))
	for i, view := range views {
		result[i] = CVViewFromDomain(view)
	}

	return result
}

// CVSectionFromDomain converts a domain CV section to a display CV section.
//
// Expected:
//   - cvsection must be valid.
//
// Returns:
//   - A CVSection value.
//
// Side effects:
//   - None.
func CVSectionFromDomain(section *career.CVSection) CVSection {
	if section == nil {
		return CVSection{}
	}

	return CVSection{
		ID:          section.ID,
		CVViewID:    section.CVViewID,
		SectionType: section.SectionType,
		Title:       section.Title,
		Order:       section.Order,
		Content:     SectionContentGroupsFromDomain(section.Content),
		Summary:     section.Summary,
	}
}

// CVSectionsFromDomain converts domain CV sections to display CV sections.
//
// Expected:
//   - cvsection must be valid.
//
// Returns:
//   - A []CVSection value.
//
// Side effects:
//   - None.
func CVSectionsFromDomain(sections []*career.CVSection) []CVSection {
	if sections == nil {
		return nil
	}

	result := make([]CVSection, len(sections))
	for i, section := range sections {
		result[i] = CVSectionFromDomain(section)
	}

	return result
}

// SectionContentGroupFromDomain converts a domain content group to a display content group.
//
// Expected:
//   - sectioncontentgroup must be valid.
//
// Returns:
//   - A SectionContentGroup value.
//
// Side effects:
//   - None.
func SectionContentGroupFromDomain(group *career.SectionContentGroup) SectionContentGroup {
	if group == nil {
		return SectionContentGroup{}
	}

	return SectionContentGroup{
		Header:    group.Header,
		StartDate: group.StartDate,
		EndDate:   group.EndDate,
		Bullets:   CVBulletsFromDomain(group.Bullets),
	}
}

// SectionContentGroupsFromDomain converts domain content groups to display content groups.
//
// Expected:
//   - sectioncontentgroup must be valid.
//
// Returns:
//   - A []SectionContentGroup value.
//
// Side effects:
//   - None.
func SectionContentGroupsFromDomain(groups []*career.SectionContentGroup) []SectionContentGroup {
	if groups == nil {
		return nil
	}

	result := make([]SectionContentGroup, len(groups))
	for i, group := range groups {
		result[i] = SectionContentGroupFromDomain(group)
	}

	return result
}

// CVBulletFromDomain converts a domain bullet to a display bullet.
//
// Expected:
//   - cvbullet must be valid.
//
// Returns:
//   - A CVBullet value.
//
// Side effects:
//   - None.
func CVBulletFromDomain(bullet *career.CVBullet) CVBullet {
	if bullet == nil {
		return CVBullet{}
	}

	return CVBullet{
		ID:                bullet.ID,
		SectionID:         bullet.SectionID,
		Text:              bullet.Text,
		SourceEventIDs:    append([]string(nil), bullet.SourceEventIDs...),
		SourceFactIDs:     append([]string(nil), bullet.SourceFactIDs...),
		Rank:              bullet.Rank,
		InclusionReason:   bullet.InclusionReason,
		Confidence:        bullet.Confidence,
		EnhancedText:      bullet.EnhancedText,
		Category:          string(bullet.Category),
		RoleScore:         bullet.RoleScore,
		AudienceScore:     bullet.AudienceScore,
		MetricScore:       bullet.MetricScore,
		ImpactScore:       bullet.ImpactScore,
		ImpactLevel:       bullet.ImpactLevel,
		KeywordMatches:    append([]string(nil), bullet.KeywordMatches...),
		AudienceRelevance: cloneStringFloatMap(bullet.AudienceRelevance),
	}
}

// CVBulletsFromDomain converts domain bullets to display bullets.
//
// Expected:
//   - cvbullet must be valid.
//
// Returns:
//   - A []CVBullet value.
//
// Side effects:
//   - None.
func CVBulletsFromDomain(bullets []*career.CVBullet) []CVBullet {
	if bullets == nil {
		return nil
	}

	result := make([]CVBullet, len(bullets))
	for i, bullet := range bullets {
		result[i] = CVBulletFromDomain(bullet)
	}

	return result
}

// CVConfigFromDomain converts a domain CV config to a display CV config.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A CVConfig value.
//
// Side effects:
//   - None.
func CVConfigFromDomain(config *career.CVConfig) CVConfig {
	if config == nil {
		return CVConfig{}
	}

	return CVConfig{
		Name:                 config.Name,
		TargetRole:           config.TargetRole,
		TargetAudience:       config.TargetAudience,
		EventFilters:         cloneStringInterfaceMap(config.EventFilters),
		TechnologyFocus:      config.TechnologyFocus,
		SelectedTechnologies: append([]string(nil), config.SelectedTechnologies...),
		FocusArea:            config.FocusArea,
		LengthFormat:         config.LengthFormat,
		SkillsFormat:         config.SkillsFormat,
		SkillsLimit:          config.SkillsLimit,
		SummaryHeading:       config.SummaryHeading,
		ProfileTitle:         config.ProfileTitle,
		WhatIBring:           append([]string(nil), config.WhatIBring...),
		CoreStrengths:        append([]string(nil), config.CoreStrengths...),
		CreatedAt:            config.CreatedAt,
		UpdatedAt:            config.UpdatedAt,
	}
}

// CVConfigsFromDomain converts domain CV configs to display CV configs.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A []CVConfig value.
//
// Side effects:
//   - None.
func CVConfigsFromDomain(configs []*career.CVConfig) []CVConfig {
	if configs == nil {
		return nil
	}

	result := make([]CVConfig, len(configs))
	for i, config := range configs {
		result[i] = CVConfigFromDomain(config)
	}

	return result
}

func cloneStringInterfaceMap(input map[string]interface{}) map[string]interface{} {
	if input == nil {
		return nil
	}

	result := make(map[string]interface{}, len(input))
	for key, value := range input {
		result[key] = value
	}

	return result
}

func cloneStringFloatMap(input map[string]float64) map[string]float64 {
	if input == nil {
		return nil
	}

	result := make(map[string]float64, len(input))
	for key, value := range input {
		result[key] = value
	}

	return result
}

func cloneTimePointer(input *time.Time) *time.Time {
	if input == nil {
		return nil
	}

	value := *input
	return &value
}

func cloneInt(input *int) *int {
	if input == nil {
		return nil
	}

	value := *input
	return &value
}
