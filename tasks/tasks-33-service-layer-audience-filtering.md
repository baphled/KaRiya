---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 33: Service Layer Audience Filtering

**Created**: 2026-01-08
**Status**: Ready for Implementation
**Priority**: HIGH
**Estimated Time**: 2-3 hours
**Related**: Codebase Audit (2026-01-08)

---

## Overview

CV generation service has stub functions that always return true/all items instead of filtering by audience. This means CVs aren't tailored to target audience - all bullets appear for all audiences regardless of relevance.

**Stub Functions**:
1. `FilterByAudience()` - Returns all bullets for all audiences
2. `isEventRelevantToAudience()` - Always returns true
3. `isFactRelevantToAudience()` - Always returns true
4. `customizeForRole()` - Returns text unchanged

---

## Files to Modify

- [ ] `internal/service/career/cv/enhanced_bullet_generator.go`
- [ ] `internal/service/career/cv/bullet_generator.go`
- [ ] `internal/service/career/cv/enhanced_bullet_generator_test.go`
- [ ] `internal/service/career/cv/bullet_generator_test.go`

---

## Implementation Plan

### Phase 1: Implement FilterByAudience (1 hour)

**Location**: `enhanced_bullet_generator.go:163-171`

**Current (STUB)**:
```go
func (g *EnhancedBulletGenerator) FilterByAudience(bullets []*career.CVBullet, audience string) []*career.CVBullet {
    // For now, accept all bullets for all audiences.
    // Future: implement audience-specific filtering
    return bullets
}
```

**Implementation**:
```go
func (g *EnhancedBulletGenerator) FilterByAudience(bullets []*career.CVBullet, audience string) []*career.CVBullet {
    filtered := make([]*career.CVBullet, 0)
    
    for _, bullet := range bullets {
        if g.isBulletRelevantToAudience(bullet, audience) {
            filtered = append(filtered, bullet)
        }
    }
    
    return filtered
}

func (g *EnhancedBulletGenerator) isBulletRelevantToAudience(bullet *career.CVBullet, audience string) bool {
    // Check source facts for audience relevance
    if len(bullet.SourceFactIDs) > 0 {
        hasRelevantFact := false
        for _, factID := range bullet.SourceFactIDs {
            fact := g.getFactByID(factID)
            if fact != nil && g.isFactRelevantToAudience(fact, audience) {
                hasRelevantFact = true
                break
            }
        }
        if !hasRelevantFact {
            return false
        }
    }
    
    // Check source events for audience relevance
    if len(bullet.SourceEventIDs) > 0 {
        hasRelevantEvent := false
        for _, eventID := range bullet.SourceEventIDs {
            event := g.getEventByID(eventID)
            if event != nil && g.isEventRelevantToAudience(event, audience) {
                hasRelevantEvent = true
                break
            }
        }
        if !hasRelevantEvent {
            return false
        }
    }
    
    return true
}
```

**Tasks**:
- [ ] Implement FilterByAudience() with real filtering logic
- [ ] Add isBulletRelevantToAudience() helper
- [ ] Check source facts for audience match
- [ ] Check source events for audience match
- [ ] Add tests for each audience type
- [ ] Verify filtered bullets are relevant

---

### Phase 2: Implement isEventRelevantToAudience (30 min)

**Location**: `bullet_generator.go:287-295`

**Current (STUB)**:
```go
func (g *BulletGenerator) isEventRelevantToAudience(event *career.CareerEvent, audience string) bool {
    // For now, accept all events for all audiences.
    // In future, could implement audience-specific filtering
    return true
}
```

**Implementation**:
```go
func (g *BulletGenerator) isEventRelevantToAudience(event *career.CareerEvent, audience string) bool {
    // Map audiences to relevant categories/tags
    audienceKeywords := map[string][]string{
        "technical": {"engineering", "development", "architecture", "technical", "coding", "infrastructure"},
        "leadership": {"management", "leadership", "strategy", "team", "mentoring", "hiring"},
        "product": {"product", "roadmap", "feature", "user", "design", "requirements"},
        "executive": {"executive", "strategic", "business", "revenue", "growth", "transformation"},
    }
    
    keywords, exists := audienceKeywords[strings.ToLower(audience)]
    if !exists {
        // Unknown audience, accept all
        return true
    }
    
    // Check categories
    for _, category := range event.Categories {
        for _, keyword := range keywords {
            if strings.Contains(strings.ToLower(category), keyword) {
                return true
            }
        }
    }
    
    // Check tags
    for _, tag := range event.Tags {
        for _, keyword := range keywords {
            if strings.Contains(strings.ToLower(tag), keyword) {
                return true
            }
        }
    }
    
    // Check event text
    lowerText := strings.ToLower(event.Text)
    for _, keyword := range keywords {
        if strings.Contains(lowerText, keyword) {
            return true
        }
    }
    
    return false
}
```

**Tasks**:
- [ ] Replace stub with keyword-based matching
- [ ] Define audience keyword mappings
- [ ] Check event categories against keywords
- [ ] Check event tags against keywords
- [ ] Check event text for keywords
- [ ] Add tests for each audience type

---

### Phase 3: Implement isFactRelevantToAudience (30 min)

**Location**: `bullet_generator.go:305-312`

**Current (STUB)**:
```go
func (g *BulletGenerator) isFactRelevantToAudience(fact *career.Fact, audience string) bool {
    // For now, accept all facts for all audiences
    return true
}
```

**Implementation**:
```go
func (g *BulletGenerator) isFactRelevantToAudience(fact *career.Fact, audience string) bool {
    // Check fact's AudienceRelevance field
    if len(fact.AudienceRelevance) > 0 {
        for _, relevantAudience := range fact.AudienceRelevance {
            if strings.EqualFold(relevantAudience, audience) {
                return true
            }
        }
        return false  // Has relevance list but audience not in it
    }
    
    // No explicit relevance set, check competency categories
    if len(fact.CompetencyCategory) > 0 {
        // Map audiences to competency categories
        audienceCompetencies := map[string][]string{
            "technical": {"technical_expertise", "architecture", "development"},
            "leadership": {"leadership", "team_management", "mentoring"},
            "product": {"product_development", "user_focus"},
            "executive": {"strategic_thinking", "business_impact"},
        }
        
        competencies, exists := audienceCompetencies[strings.ToLower(audience)]
        if !exists {
            return true  // Unknown audience, accept
        }
        
        for _, factComp := range fact.CompetencyCategory {
            for _, audComp := range competencies {
                if strings.Contains(strings.ToLower(factComp), audComp) {
                    return true
                }
            }
        }
        
        return false
    }
    
    // No explicit data, accept all
    return true
}
```

**Tasks**:
- [ ] Check fact's AudienceRelevance field
- [ ] Map audiences to competency categories
- [ ] Check fact competency categories
- [ ] Return false if explicitly not relevant
- [ ] Add tests for audience matching

---

### Phase 4: Implement customizeForRole (30 min)

**Location**: `enhanced_bullet_generator.go:462-467`

**Current (STUB)**:
```go
func (g *EnhancedBulletGenerator) customizeForRole(text string, role career.TargetRole) string {
    // For now, return as-is. Future: apply role-specific wording preferences
    return text
}
```

**Implementation**:
```go
func (g *EnhancedBulletGenerator) customizeForRole(text string, role career.TargetRole) string {
    // Role-specific terminology mappings
    roleAdjustments := map[career.TargetRole]map[string]string{
        career.RoleSeniorIC: {
            "led team": "contributed to team",
            "managed": "worked with",
            "directed": "collaborated on",
        },
        career.RolePrincipal: {
            "worked on": "led development of",
            "helped": "drove",
            "assisted": "architected",
        },
        career.RoleStaff: {
            "worked on": "led strategy for",
            "implemented": "established",
            "developed": "defined",
        },
        career.RoleEM: {
            "implemented": "enabled team to implement",
            "developed": "guided development of",
            "created": "built team that created",
        },
    }
    
    adjustments, exists := roleAdjustments[role]
    if !exists {
        return text
    }
    
    adjusted := text
    for old, new := range adjustments {
        // Case-insensitive replacement while preserving original case
        adjusted = replacePreservingCase(adjusted, old, new)
    }
    
    return adjusted
}

func replacePreservingCase(text, old, new string) string {
    // Simple implementation - can be enhanced
    lowerText := strings.ToLower(text)
    lowerOld := strings.ToLower(old)
    
    if strings.Contains(lowerText, lowerOld) {
        // Find position and check original case
        idx := strings.Index(lowerText, lowerOld)
        if idx >= 0 {
            // Preserve case of first letter
            if unicode.IsUpper(rune(text[idx])) {
                new = strings.ToUpper(string(new[0])) + new[1:]
            }
            return text[:idx] + new + text[idx+len(old):]
        }
    }
    
    return text
}
```

**Tasks**:
- [ ] Define role-specific terminology mappings
- [ ] Implement text replacement with case preservation
- [ ] Add for each role (senior_ic, principal, staff, em)
- [ ] Test customization for each role

---

## Acceptance Criteria

### Must Have
- [ ] `FilterByAudience()` filters bullets by audience relevance
- [ ] `isEventRelevantToAudience()` checks categories, tags, text
- [ ] `isFactRelevantToAudience()` checks AudienceRelevance field
- [ ] `customizeForRole()` adapts language for role seniority
- [ ] Generated CVs are tailored to selected audience
- [ ] All 203 CV service tests still pass
- [ ] No regressions in CV generation

### Should Have
- [ ] Audience matching is case-insensitive
- [ ] Unknown audiences don't break filtering
- [ ] Empty/missing audience data handled gracefully

---

## Testing Strategy

```go
It("should filter bullets by technical audience", func() {
    bullets := createTestBullets()  // Mix of technical and non-technical
    filtered := generator.FilterByAudience(bullets, "technical")
    for _, bullet := range filtered {
        Expect(bullet.SourceFacts).To(HaveRelevantAudience("technical"))
    }
})

It("should identify technical events", func() {
    event := &career.CareerEvent{
        Categories: []string{"engineering", "architecture"},
    }
    Expect(generator.isEventRelevantToAudience(event, "technical")).To(BeTrue())
})

It("should customize text for principal role", func() {
    text := "worked on implementation of feature"
    customized := generator.customizeForRole(text, career.RolePrincipal)
    Expect(customized).To(ContainSubstring("led development"))
})
```

---

## References

- `internal/service/career/cv/enhanced_bullet_generator.go` - Implementation
- `internal/domain/career/fact.go` - AudienceRelevance field
- `internal/domain/career/cv.go` - Role definitions

---

**Last Updated**: 2026-01-08
**Status**: Ready for implementation
