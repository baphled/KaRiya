# Task 44: Wire EnhancedBulletGenerator

## Overview
- **Goal**: Replace basic `BulletGenerator` with existing `EnhancedBulletGenerator` in CV generation for improved bullet quality
- **Time Estimate**: 2-3 hours
- **Prerequisites**: Understanding of CV generation service, bullet generator interfaces, data processing service

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - `internal/service/career/cv/cv_generation_service.go` (current service)
  - `internal/service/career/cv/enhanced_bullet_generator.go` (target generator)
  - `internal/service/career/cv/bullet_generator.go` (current basic generator)
  - `internal/service/career/cv/data_processing_service.go` (achievement extraction)
- [ ] Confirmed this is ONE atomic task (wire enhanced generator)
- [ ] Identified which test files will be created/modified

## Current Status

**TASK 44 NOT STARTED** - Ready to Begin

## Context

The `EnhancedBulletGenerator` already exists with advanced features that improve CV bullet quality:
- **Multi-factor scoring**: role relevance, audience fit, metrics, impact
- **Verb enhancement**: replaces weak verbs ("worked on" → "led", "helped with" → "contributed to")
- **Metric detection**: identifies and scores quantifiable achievements
- **Impact level assessment**: categorizes bullets by scope (personal, team, company, industry)

However, `CVGenerationService` currently uses the basic `BulletGenerator` which only provides:
- Simple confidence scoring
- Basic inclusion/exclusion filtering
- No sophisticated ranking

This task wires the enhanced generator for immediate CV quality improvement without adding new features.

## Files to Modify

### Service Layer
- [ ] `internal/service/career/cv/cv_generation_service.go` - Add enhanced generator and data processor fields
- [ ] `internal/service/career/cv/cv_generation_service_test.go` - Update tests for new fields

### Intent Layer
- [ ] `internal/cli/intents/generate_cv_context.go` - Pass DataProcessingService to service

### Application Wiring
- [ ] `cmd/cli/main.go` - Wire enhanced generator and data processor in initialization

## Implementation Plan

### Phase 1: Update CVGenerationService Structure

**Goal**: Add enhanced bullet generator and data processor to the service

#### Current State
```go
// internal/service/career/cv/cv_generation_service.go

type DefaultCVGenerationService struct {
    eventRepo       careerrepo.Repository
    factRepo        careerrepo.FactRepository
    configManager   ConfigManager
    bulletGenerator BulletGenerator      // Basic generator
    sectionBuilder  SectionBuilder
    logger          *logger.Logger
}

func NewCVGenerationService(
    eventRepo careerrepo.Repository,
    factRepo careerrepo.FactRepository,
    configManager ConfigManager,
    bulletGenerator BulletGenerator,
    sectionBuilder SectionBuilder,
    log *logger.Logger,
) *DefaultCVGenerationService {
    return &DefaultCVGenerationService{
        eventRepo:       eventRepo,
        factRepo:        factRepo,
        configManager:   configManager,
        bulletGenerator: bulletGenerator,
        sectionBuilder:  sectionBuilder,
        logger:          log,
    }
}
```

#### Target State
```go
// internal/service/career/cv/cv_generation_service.go

type DefaultCVGenerationService struct {
    eventRepo               careerrepo.Repository
    factRepo                careerrepo.FactRepository
    configManager           ConfigManager
    bulletGenerator         BulletGenerator              // Keep for backward compat
    enhancedBulletGenerator EnhancedBulletGenerator      // NEW - Advanced scoring
    dataProcessor           DataProcessingService        // NEW - Achievement extraction
    sectionBuilder          SectionBuilder
    logger                  *logger.Logger
}

func NewCVGenerationService(
    eventRepo careerrepo.Repository,
    factRepo careerrepo.FactRepository,
    configManager ConfigManager,
    bulletGenerator BulletGenerator,
    enhancedBulletGenerator EnhancedBulletGenerator,      // NEW
    dataProcessor DataProcessingService,                  // NEW
    sectionBuilder SectionBuilder,
    log *logger.Logger,
) *DefaultCVGenerationService {
    return &DefaultCVGenerationService{
        eventRepo:               eventRepo,
        factRepo:                factRepo,
        configManager:           configManager,
        bulletGenerator:         bulletGenerator,
        enhancedBulletGenerator: enhancedBulletGenerator,  // NEW
        dataProcessor:           dataProcessor,            // NEW
        sectionBuilder:          sectionBuilder,
        logger:                  log,
    }
}
```

**TDD Checklist - Phase 1:**
- [ ] Write failing test: NewCVGenerationService accepts enhanced generator parameter
- [ ] Test fails with compilation error (parameter doesn't exist)
- [ ] Add `enhancedBulletGenerator` parameter to constructor
- [ ] Test passes
- [ ] Write failing test: NewCVGenerationService accepts data processor parameter
- [ ] Test fails with compilation error (parameter doesn't exist)
- [ ] Add `dataProcessor` parameter to constructor
- [ ] Test passes
- [ ] Write failing test: Service stores enhanced generator in field
- [ ] Test fails (field doesn't exist)
- [ ] Add fields to struct
- [ ] Test passes
- [ ] Commit: `test(cv): add enhanced generator constructor tests`
- [ ] Commit: `feat(cv): add enhanced generator to service`

### Phase 2: Use Enhanced Generator in CV Generation

**Goal**: Replace basic bullet generation with enhanced generation

#### Current Implementation
```go
// internal/service/career/cv/cv_generation_service.go
// In GenerateCVFromConfig()

// Generate bullets using BulletGenerator
bullets, err := svc.bulletGenerator.GenerateBullets(ctx, events, facts, config.TargetRole, config.TargetAudience)
if err != nil {
    svc.logger.Error("Failed to generate bullets: %v", err)
    return nil, fmt.Errorf("failed to generate bullets: %w", err)
}
```

#### Target Implementation
```go
// internal/service/career/cv/cv_generation_service.go
// In GenerateCVFromConfig()

// Extract achievements from events for enhanced generation
companyGroups, err := svc.dataProcessor.GroupEventsByCompany(ctx, events)
if err != nil {
    svc.logger.Warn("Failed to group events by company: %v", err)
    companyGroups = make(map[string][]*CareerEvent)
}

achievements, err := svc.dataProcessor.ExtractAchievements(ctx, companyGroups)
if err != nil {
    svc.logger.Warn("Failed to extract achievements: %v", err)
    achievements = []*Achievement{}
}

svc.logger.Info("Extracted %d achievements for enhanced bullet generation", len(achievements))

// Use enhanced generator with achievements
enhancedBullets, err := svc.enhancedBulletGenerator.GenerateBullets(
    ctx,
    events,
    facts,
    achievements,
    config.TargetRole,
    config.TargetAudience,
)
if err != nil {
    svc.logger.Error("Failed to generate enhanced bullets: %v", err)
    return nil, fmt.Errorf("failed to generate bullets: %w", err)
}

// Convert EnhancedBullet to CVBullet for section builder compatibility
bullets := make([]*career.CVBullet, len(enhancedBullets))
for i, eb := range enhancedBullets {
    bullets[i] = &career.CVBullet{
        ID:              eb.ID,
        Text:            eb.EnhancedText,  // Use enhanced text if available
        SourceEventIDs:  eb.SourceEventIDs,
        SourceFactIDs:   eb.SourceFactIDs,
        Rank:            eb.FinalScore,     // Use final score as rank
        InclusionReason: eb.InclusionReason,
        Confidence:      eb.Confidence,
    }

    // Fallback to original text if enhanced is empty
    if bullets[i].Text == "" {
        bullets[i].Text = eb.Text
    }
}

svc.logger.Info("Generated %d enhanced bullets from %d events and %d facts",
    len(bullets), len(events), len(facts))
```

**TDD Checklist - Phase 2:**
- [ ] Write failing test: GenerateCVFromConfig calls dataProcessor.GroupEventsByCompany
- [ ] Test fails (method not called)
- [ ] Add company grouping logic
- [ ] Test passes
- [ ] Write failing test: GenerateCVFromConfig calls dataProcessor.ExtractAchievements
- [ ] Test fails (method not called)
- [ ] Add achievement extraction
- [ ] Test passes
- [ ] Write failing test: GenerateCVFromConfig calls enhancedBulletGenerator.GenerateBullets with achievements
- [ ] Test fails (enhanced generator not called)
- [ ] Add enhanced bullet generation
- [ ] Test passes
- [ ] Write failing test: EnhancedBullets are converted to CVBullets correctly
- [ ] Test fails (conversion not implemented)
- [ ] Add conversion logic
- [ ] Test passes
- [ ] Write failing test: Enhanced text is preferred over original text
- [ ] Test fails (not using enhanced text)
- [ ] Update to use EnhancedText with fallback
- [ ] Test passes
- [ ] Write failing test: FinalScore is used as bullet rank
- [ ] Test fails (not using score)
- [ ] Map FinalScore to Rank
- [ ] Test passes
- [ ] Commit: `test(cv): add enhanced bullet generation tests`
- [ ] Commit: `feat(cv): use enhanced bullet generator in CV generation`

### Phase 3: Wire in Application

**Goal**: Initialize and pass enhanced generator through the application stack

#### Update generate_cv_context.go
```go
// internal/cli/intents/generate_cv_context.go

type GenerateCVContext struct {
    Context        context.Context
    Service        *career.Service
    CVService      cv.CVGenerationService
    DataProcessor  cv.DataProcessingService   // NEW - Add this field
    // ... other fields
}

func (ctx *GenerateCVContext) Validate() error {
    if ctx.Context == nil {
        return fmt.Errorf("context cannot be nil")
    }
    if ctx.Service == nil {
        return fmt.Errorf("service cannot be nil")
    }
    if ctx.CVService == nil {
        return fmt.Errorf("CV service cannot be nil")
    }
    if ctx.DataProcessor == nil {                      // NEW
        return fmt.Errorf("data processor cannot be nil")  // NEW
    }                                                  // NEW
    return nil
}
```

#### Update cmd/cli/main.go
```go
// cmd/cli/main.go - In CV service initialization

// Create data processing service
dataProcessor := cv.NewDataProcessingService(logger)

// Create bullet generators
bulletGenerator := cv.NewDefaultBulletGenerator(logger)
enhancedBulletGenerator := cv.NewEnhancedBulletGenerator(logger)  // NEW

// Create section builder
sectionBuilder := cv.NewDefaultSectionBuilder(logger)

// Create CV generation service
cvService := cv.NewCVGenerationService(
    eventRepo,
    factRepo,
    configManager,
    bulletGenerator,
    enhancedBulletGenerator,  // NEW - Pass enhanced generator
    dataProcessor,            // NEW - Pass data processor
    sectionBuilder,
    logger,
)

// ... later when creating GenerateCVContext ...

cvContext := &intents.GenerateCVContext{
    Context:       ctx,
    Service:       careerService,
    CVService:     cvService,
    DataProcessor: dataProcessor,  // NEW - Pass to context
}
```

**TDD Checklist - Phase 3:**
- [ ] Write failing test: GenerateCVContext validates DataProcessor is not nil
- [ ] Test fails (no validation)
- [ ] Add DataProcessor field and validation
- [ ] Test passes
- [ ] Write failing integration test: Full CV generation uses enhanced bullets
- [ ] Test fails (wiring not complete)
- [ ] Update main.go to create enhanced generator and data processor
- [ ] Update context initialization
- [ ] Test passes
- [ ] Verify all existing CV generation tests still pass
- [ ] Commit: `test(intent): add data processor validation`
- [ ] Commit: `feat(app): wire enhanced bullet generator`

### Phase 4: Backward Compatibility & Fallback

**Goal**: Ensure system works if enhanced generator is nil (defensive programming)

```go
// internal/service/career/cv/cv_generation_service.go

func (svc *DefaultCVGenerationService) GenerateCVFromConfig(ctx context.Context, config *career.CVConfig) (*career.CVView, error) {
    // ... existing code ...

    var bullets []*career.CVBullet
    var err error

    // Use enhanced generator if available, fallback to basic
    if svc.enhancedBulletGenerator != nil && svc.dataProcessor != nil {
        svc.logger.Info("Using enhanced bullet generator")

        // Enhanced generation path (from Phase 2)
        companyGroups, _ := svc.dataProcessor.GroupEventsByCompany(ctx, events)
        achievements, _ := svc.dataProcessor.ExtractAchievements(ctx, companyGroups)

        enhancedBullets, err := svc.enhancedBulletGenerator.GenerateBullets(
            ctx, events, facts, achievements, config.TargetRole, config.TargetAudience,
        )
        if err != nil {
            svc.logger.Warn("Enhanced generator failed, falling back to basic: %v", err)
            // Fall through to basic generator
        } else {
            // Convert and use enhanced bullets
            bullets = convertEnhancedBullets(enhancedBullets)
        }
    }

    // Fallback to basic generator if enhanced failed or unavailable
    if bullets == nil {
        svc.logger.Info("Using basic bullet generator")
        bullets, err = svc.bulletGenerator.GenerateBullets(
            ctx, events, facts, config.TargetRole, config.TargetAudience,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to generate bullets: %w", err)
        }
    }

    // ... continue with bullets ...
}
```

**TDD Checklist - Phase 4:**
- [ ] Write failing test: Falls back to basic generator if enhanced is nil
- [ ] Test fails (no fallback)
- [ ] Add nil check and fallback logic
- [ ] Test passes
- [ ] Write failing test: Falls back to basic generator if data processor is nil
- [ ] Test fails (no check for data processor)
- [ ] Add data processor nil check
- [ ] Test passes
- [ ] Write failing test: Falls back to basic generator if enhanced generation fails
- [ ] Test fails (no error handling)
- [ ] Add error handling with fallback
- [ ] Test passes
- [ ] Verify all existing tests pass with new code paths
- [ ] Commit: `test(cv): add fallback generator tests`
- [ ] Commit: `feat(cv): add fallback to basic generator`

## TDD Checklist (MUST COMPLETE IN ORDER)

### RED Phase (Write Tests First)
- [ ] Phase 1: Write constructor tests (enhanced generator, data processor parameters)
- [ ] Phase 2: Write enhanced bullet generation tests (grouping, extraction, conversion)
- [ ] Phase 3: Write integration tests (full CV generation with enhanced bullets)
- [ ] Phase 4: Write fallback tests (nil checks, error handling)

### GREEN Phase (Implement Features)
- [ ] Phase 1: Update constructor and struct (add new fields)
- [ ] Phase 2: Implement enhanced bullet generation (in GenerateCVFromConfig)
- [ ] Phase 3: Wire in application (main.go, context)
- [ ] Phase 4: Add fallback logic (defensive programming)

### REFACTOR Phase (if needed)
- [ ] Extract conversion logic to helper method if complex
- [ ] Simplify fallback logic if nested conditionals become deep
- [ ] Consider builder pattern if constructor has too many parameters

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes (REQUIRED)
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)
- [ ] All tests pass locally

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All tests pass (including race detector: `go test -race ./...`)
- [ ] Code coverage maintained ≥ 80%
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria
- [ ] `CVGenerationService` uses `EnhancedBulletGenerator` when available
- [ ] Achievements are extracted via `DataProcessingService`
- [ ] CV bullets have `FinalScore` from multi-factor scoring
- [ ] Enhanced text is used (verb improvements, metric context)
- [ ] Fallback to basic generator works if enhanced unavailable
- [ ] All existing CV generation tests pass (zero regressions)
- [ ] New tests cover enhanced generation path

## Expected Impact

### Before (Basic Generator)
```
• Worked on API improvements
• Helped with migration to new database
• Involved in team meetings
```

### After (Enhanced Generator)
```
• Led API performance improvements, reducing latency by 40%
• Architected migration to PostgreSQL for 5 critical services
• Mentored team of 3 engineers on database optimization
```

### Scoring Example
```
Bullet: "Led API performance improvements, reducing latency by 40%"

RoleScore:     0.85  (leadership verb, technical content)
AudienceScore: 0.90  (high impact, quantified)
MetricScore:   0.80  (percentage metric present)
ImpactScore:   0.75  (team-level impact)
Confidence:    0.80  (from fact extraction)

FinalScore = (0.25 * 0.85) + (0.20 * 0.90) + (0.20 * 0.80) + (0.20 * 0.75) + (0.15 * 0.80)
           = 0.8125 (ranks near top of CV bullets)
```

## Rollback Plan
If this change causes issues:
1. Revert to commit before Phase 1
2. Remove enhanced generator parameters from constructor
3. Remove enhanced generation logic from `GenerateCVFromConfig`
4. Remove wiring from `main.go` and context
5. Run `make check-compliance` to verify rollback
6. All tests should pass with basic generator

## Dependencies
- None (this task is independent)
- All required components already exist

## Next Steps After Completion
- Task 45: Burst Auto-Detection UI (independent)
- Task 46: Skill Inference Service (depends on 45)
- Task 47: Enhanced CV Skills Section (depends on 46)
