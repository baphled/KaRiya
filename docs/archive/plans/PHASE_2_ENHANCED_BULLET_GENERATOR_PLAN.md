---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Phase 2: Enhanced BulletGenerator Implementation Plan

**Date**: January 4, 2026
**Phase**: 2 of 6
**Status**: In Progress
**Estimated Hours**: 35

---

## Overview

Phase 2 enhances the existing BulletGenerator to leverage the DataProcessingService from Phase 1. The enhanced BulletGenerator will:

1. Use structured data from DataProcessingService (CompanyGroup, Achievement, Metric, Skill)
2. Implement intelligent role/audience filtering
3. Implement comprehensive bullet ranking algorithm
4. Generate professional, impact-driven bullet text
5. Apply role-specific customization

---

## Current BulletGenerator Analysis

### Existing Strengths
- ✅ Basic bullet generation from events and facts
- ✅ Filtering by inclusion criteria
- ✅ Basic ranking algorithm
- ✅ Role-specific bullet caps
- ✅ Category-based role filtering

### Areas for Enhancement
- ❌ Uses raw event/fact text without processing
- ❌ Limited role-specific filtering
- ❌ Simplistic ranking algorithm
- ❌ No metric-based scoring
- ❌ No professional wording enhancement
- ❌ No audience-specific customization

---

## Enhancement Strategy

### 1. Integration with DataProcessingService

**Current Flow**:
```
Events/Facts → generateInitialBullets → filterByInclusionCriteria → rankBullets → compress
```

**Enhanced Flow**:
```
Events/Facts → DataProcessingService (extract achievements, metrics, skills)
            ↓
         CompanyGroups, Achievements, Metrics, Skills
            ↓
    Enhanced BulletGenerator
            ├─ generateAchievementBullets
            ├─ filterByRole
            ├─ filterByAudience
            ├─ rankByRelevance
            ├─ enhanceWording
            └─ compress
```

### 2. Enhanced Filtering

#### Role-Based Filtering
```go
type RoleFilter struct {
    PrimaryCategories   []string  // Most relevant categories
    SecondaryCategories []string  // Secondary categories
    MinConfidence       float64   // Minimum confidence threshold
    PreferredMetrics    []string  // Preferred metric types
}

// Role filters:
Principal: {
    Primary: ["leadership", "strategy", "architecture"],
    Secondary: ["technical", "mentoring"],
    MinConfidence: 0.8,
    PreferredMetrics: ["percentage", "count", "currency"]
}

Staff: {
    Primary: ["technical", "architecture"],
    Secondary: ["leadership", "mentoring"],
    MinConfidence: 0.75,
    PreferredMetrics: ["percentage", "count"]
}

EM: {
    Primary: ["leadership", "mentoring"],
    Secondary: ["strategy", "product"],
    MinConfidence: 0.75,
    PreferredMetrics: ["count", "percentage"]
}

SeniorIC: {
    Primary: ["technical", "architecture"],
    Secondary: ["leadership", "strategy"],
    MinConfidence: 0.75,
    PreferredMetrics: ["percentage", "count"]
}
```

#### Audience-Based Filtering
```go
type AudienceFilter struct {
    Name                string
    FocusAreas         []string  // What this audience cares about
    PreferredMetrics   []string  // Metrics they value
    MinImpactLevel     string    // "low", "medium", "high"
}

// Audience filters:
HiringManager: {
    Focus: ["impact", "results", "team_fit"],
    Metrics: ["percentage", "count", "currency"],
    MinImpact: "medium"
}

Recruiter: {
    Focus: ["growth", "progression", "skills"],
    Metrics: ["count", "percentage"],
    MinImpact: "low"
}

Peer: {
    Focus: ["technical_depth", "collaboration", "innovation"],
    Metrics: ["percentage", "count"],
    MinImpact: "medium"
}
```

### 3. Enhanced Ranking Algorithm

**Scoring Components**:

```go
type BulletScore struct {
    BaseScore       float64  // 0.0-1.0
    RoleRelevance   float64  // 0.0-1.0 (how relevant to target role)
    AudienceFit     float64  // 0.0-1.0 (how relevant to audience)
    MetricScore     float64  // 0.0-1.0 (based on metrics present)
    ImpactScore     float64  // 0.0-1.0 (based on achievement impact)
    ConfidenceScore float64  // 0.0-1.0 (extraction confidence)
    FinalScore      float64  // Weighted combination
}

// Scoring formula:
FinalScore = (
    0.25 * RoleRelevance +
    0.20 * AudienceF it +
    0.20 * MetricScore +
    0.20 * ImpactScore +
    0.15 * ConfidenceScore
)
```

**Ranking Factors**:

1. **Role Relevance** (25%)
   - Category match with role preferences
   - Skill alignment with role requirements
   - Achievement type relevance

2. **Audience Fit** (20%)
   - Alignment with audience focus areas
   - Impact level match
   - Metric type preference

3. **Metric Score** (20%)
   - Presence of quantifiable metrics
   - Multiple metric types (higher is better)
   - Metric magnitude (larger is better)

4. **Impact Score** (20%)
   - Number of metrics
   - Breadth of impact (team, org, customer)
   - Outcome clarity

5. **Confidence Score** (15%)
   - Source quality (fact > event)
   - Extraction confidence
   - Verification status

### 4. Professional Wording Enhancement

**Enhancement Rules**:

```go
type WordingEnhancer interface {
    // EnhanceBullet improves bullet text for professional CV use
    EnhanceBullet(bullet *CVBullet, role string) (*EnhancedBullet, error)

    // ApplyActionVerbs replaces weak verbs with strong action verbs
    ApplyActionVerbs(text string) string

    // AddMetricContext adds context around metrics
    AddMetricContext(text string, metrics []*Metric) string

    // StructureBullet applies "action + context + result" structure
    StructureBullet(text string, metrics []*Metric) string
}
```

**Wording Improvements**:

1. **Action Verb Enhancement**
   - Weak: "worked on", "helped with", "involved in"
   - Strong: "led", "architected", "delivered", "optimized"
   - Role-specific: Principal → "architected", "drove", "established"

2. **Metric Integration**
   - Before: "Improved performance"
   - After: "Improved performance by 25% for 12-person team"

3. **Impact Emphasis**
   - Before: "Led team"
   - After: "Led team of 12 engineers, delivering microservices architecture"

4. **Result Focus**
   - Before: "Implemented system"
   - After: "Implemented real-time monitoring system, reducing MTTR by 60%"

### 5. Role-Specific Customization

**Principal-Level Customization**:
- Emphasize strategic impact
- Highlight organizational influence
- Focus on architecture and vision
- Include mentoring and growth
- Metrics: % improvement, scale (people/customers)

**Staff-Level Customization**:
- Emphasize technical depth
- Highlight architectural decisions
- Focus on complexity and scale
- Include ownership and autonomy
- Metrics: % improvement, complexity metrics

**EM-Level Customization**:
- Emphasize team development
- Highlight people impact
- Focus on growth and retention
- Include organizational contributions
- Metrics: team size, retention, growth

**SeniorIC-Level Customization**:
- Emphasize technical excellence
- Highlight innovation
- Focus on quality and reliability
- Include mentoring junior engineers
- Metrics: % improvement, scale, quality metrics

---

## Implementation Approach

### Step 1: Create Enhanced Data Structures

```go
// EnhancedBullet with additional metadata
type EnhancedBullet struct {
    *CVBullet
    RoleScore         float64
    AudienceScore     float64
    EnhancedText      string
    Metrics           []*Metric
    ImpactLevel       string  // "low", "medium", "high"
    KeywordMatches    []string
}

// RoleAudienceConfig for filtering
type RoleAudienceConfig struct {
    Role               string
    Audiences          []string
    RoleFilter         *RoleFilter
    AudienceFilters    []*AudienceFilter
}
```

### Step 2: Enhance Existing Methods

```go
// Enhanced method signatures
func (bg *EnhancedBulletGenerator) GenerateBullets(
    ctx context.Context,
    events []*CareerEvent,
    facts []*Fact,
    targetRole string,
    targetAudiences []string,
    dataProcessor DataProcessingService) ([]*EnhancedBullet, error)

func (bg *EnhancedBulletGenerator) FilterByRole(
    bullets []*EnhancedBullet,
    role string) []*EnhancedBullet

func (bg *EnhancedBulletGenerator) FilterByAudience(
    bullets []*EnhancedBullet,
    audiences []string) []*EnhancedBullet

func (bg *EnhancedBulletGenerator) RankByRelevance(
    bullets []*EnhancedBullet,
    role string,
    audiences []string) []*EnhancedBullet

func (bg *EnhancedBulletGenerator) EnhanceBulletWording(
    bullet *EnhancedBullet,
    role string) (*EnhancedBullet, error)
```

### Step 3: Implement New Filtering Logic

```go
// Role-based filtering
func (bg *EnhancedBulletGenerator) FilterByRole(bullets, role) {
    filter := bg.getRoleFilter(role)
    var filtered []*EnhancedBullet

    for _, bullet := range bullets {
        // Check category match
        if bg.categoriesMatch(bullet.Categories, filter.PrimaryCategories) {
            bullet.RoleScore += 0.5
        } else if bg.categoriesMatch(bullet.Categories, filter.SecondaryCategories) {
            bullet.RoleScore += 0.3
        }

        // Check confidence threshold
        if bullet.Confidence >= filter.MinConfidence {
            filtered = append(filtered, bullet)
        }
    }

    return filtered
}

// Audience-based filtering
func (bg *EnhancedBulletGenerator) FilterByAudience(bullets, audiences) {
    var filtered []*EnhancedBullet

    for _, bullet := range bullets {
        for _, audience := range audiences {
            filter := bg.getAudienceFilter(audience)

            // Check focus area match
            if bg.hasRelevantFocusArea(bullet, filter.FocusAreas) {
                bullet.AudienceScore += 0.3
            }

            // Check impact level
            if bg.meetsImpactLevel(bullet, filter.MinImpactLevel) {
                filtered = append(filtered, bullet)
                break
            }
        }
    }

    return filtered
}
```

### Step 4: Implement Enhanced Ranking

```go
func (bg *EnhancedBulletGenerator) RankByRelevance(bullets, role, audiences) {
    for _, bullet := range bullets {
        bullet.Rank = bg.calculateEnhancedScore(bullet, role, audiences)
    }

    sort.Slice(bullets, func(i, j int) bool {
        return bullets[i].Rank > bullets[j].Rank
    })

    return bullets
}

func (bg *EnhancedBulletGenerator) calculateEnhancedScore(bullet, role, audiences) float64 {
    roleScore := bullet.RoleScore * 0.25
    audienceScore := bg.calculateAudienceScore(bullet, audiences) * 0.20
    metricScore := bg.calculateMetricScore(bullet) * 0.20
    impactScore := bg.calculateImpactScore(bullet) * 0.20
    confidenceScore := bullet.Confidence * 0.15

    return roleScore + audienceScore + metricScore + impactScore + confidenceScore
}
```

### Step 5: Implement Wording Enhancement

```go
func (bg *EnhancedBulletGenerator) EnhanceBulletWording(bullet, role) {
    // 1. Apply action verb enhancement
    enhancedText := bg.enhanceActionVerb(bullet.Text, role)

    // 2. Add metric context
    enhancedText = bg.addMetricContext(enhancedText, bullet.Metrics)

    // 3. Structure for impact
    enhancedText = bg.structureForImpact(enhancedText, bullet.Metrics)

    // 4. Role-specific customization
    enhancedText = bg.customizeForRole(enhancedText, role)

    return &EnhancedBullet{
        CVBullet:     bullet,
        EnhancedText: enhancedText,
    }
}
```

---

## Testing Strategy

### Unit Tests (15+ tests)

1. **Filtering Tests**
   - Role-based filtering
   - Audience-based filtering
   - Combined filtering
   - Edge cases

2. **Ranking Tests**
   - Score calculation
   - Ranking order
   - Tie-breaking
   - Role-specific caps

3. **Wording Tests**
   - Action verb enhancement
   - Metric integration
   - Impact structure
   - Role customization

4. **Integration Tests**
   - End-to-end generation
   - Multiple roles/audiences
   - Large event sets
   - Performance

### Performance Tests

- Generate bullets for 100 events: < 100ms
- Generate bullets for 1000 events: < 1s
- Rank 500 bullets: < 50ms

---

## Success Criteria

✅ Enhanced filtering by role and audience
✅ Improved ranking algorithm with multiple factors
✅ Professional wording enhancement
✅ Role-specific customization
✅ 90%+ test coverage
✅ < 1 second for 1000 events
✅ Zero breaking changes

---

## Next Steps

1. Implement EnhancedBulletGenerator
2. Create comprehensive test suite
3. Integrate with DataProcessingService
4. Performance testing and optimization
5. Documentation and examples


