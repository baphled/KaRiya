# Task 46: Skill Inference Service

## Overview
- **Goal**: Create skill inference service that detects technologies in event text and suggests skills to users
- **Time Estimate**: 8-10 hours
- **Prerequisites**: Task 45 (burst detection UI for pattern familiarity), understanding of skill domain model, text analysis basics

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - `internal/domain/career/skill.go` (skill model)
  - `internal/repository/career/skill_repository.go` (persistence)
  - `internal/service/career/burst_fact/detector.go` (similar detection pattern)
  - `internal/cli/intents/burst_management_intent.go` (suggestion UI pattern)
  - `internal/cli/intents/manage_skills_intent.go` (existing skills UI)
- [ ] Confirmed this is ONE atomic task (skill inference system)
- [ ] Identified which test files will be created/modified

## Current Status

**TASK 46 NOT STARTED** - Ready to Begin

## Context

Currently, skills are **manually managed** through the ManageSkills intent. Users must:
1. Manually create each skill
2. Manually link skills to events
3. Remember which technologies they've used

This task creates **automatic skill inference** that:
1. Analyzes event text for technology mentions (e.g., "built API with Go and PostgreSQL")
2. Suggests skills with confidence scores (0.0-1.0)
3. Auto-links inferred skills to source events via junction table
4. Integrates into burst confirmation flow (after fact extraction)
5. Provides standalone inference via "i" key in ManageSkills

### Why This Matters for CV Generation

From our analysis:
- CV "Core Competencies" currently uses generic categories from Facts
- The Skills table has rich data (Level, YearsUsed) but isn't used in CV generation
- Task 47 will use Skills table for professional CV skills sections
- **This task populates the Skills table automatically**

## Technology Dictionary Requirements

Based on user preference:
- **~100 comprehensive technology keywords** from the start
- Organized by category for easy maintenance
- File-based dictionary that can be internally updated
- Supports aliases (e.g., "golang" → "Go", "k8s" → "Kubernetes")

## Files to Create

### Service Layer
- [ ] `internal/service/career/technology_keywords.go` - Technology dictionary (~100 entries)
- [ ] `internal/service/career/technology_keywords_test.go` - Dictionary tests
- [ ] `internal/service/career/skill_inference.go` - Service interface
- [ ] `internal/service/career/skill_inference_impl.go` - Implementation
- [ ] `internal/service/career/skill_inference_test.go` - Service tests

## Files to Modify

### Intent Layer
- [ ] `internal/cli/intents/burst_management.go` - Add skill suggestion state
- [ ] `internal/cli/intents/burst_management_intent.go` - Integrate into confirmation flow
- [ ] `internal/cli/intents/manage_skills.go` - Add inference states
- [ ] `internal/cli/intents/manage_skills_intent.go` - Add "i" key handler for standalone inference

## Implementation Plan

### Phase 1: Create Technology Dictionary

**Goal**: Create comprehensive, maintainable dictionary of ~100 technology keywords

#### Technology Dictionary Design

```go
// internal/service/career/technology_keywords.go

package career

// TechnologyKeyword maps a search keyword to its canonical skill
type TechnologyKeyword struct {
    Keyword  string // Lowercase search term (for matching)
    Skill    string // Canonical skill name (for display)
    Category string // Skill category (backend, frontend, etc.)
}

// TechnologyKeywords is the master dictionary
// Organized by category for easy maintenance and updates
var TechnologyKeywords = []TechnologyKeyword{
    // ========================================
    // BACKEND LANGUAGES & FRAMEWORKS
    // ========================================
    {"go", "Go", "backend"},
    {"golang", "Go", "backend"},
    {"python", "Python", "backend"},
    {"ruby", "Ruby", "backend"},
    {"rails", "Ruby on Rails", "backend"},
    {"ruby on rails", "Ruby on Rails", "backend"},
    {"java", "Java", "backend"},
    {"spring", "Spring", "backend"},
    {"spring boot", "Spring Boot", "backend"},
    {"scala", "Scala", "backend"},
    {"rust", "Rust", "backend"},
    {"c#", "C#", "backend"},
    {"csharp", "C#", "backend"},
    {".net", ".NET", "backend"},
    {"dotnet", ".NET", "backend"},
    {"node", "Node.js", "backend"},
    {"nodejs", "Node.js", "backend"},
    {"node.js", "Node.js", "backend"},
    {"express", "Express.js", "backend"},
    {"expressjs", "Express.js", "backend"},
    {"php", "PHP", "backend"},
    {"laravel", "Laravel", "backend"},
    {"symfony", "Symfony", "backend"},
    {"elixir", "Elixir", "backend"},
    {"phoenix", "Phoenix", "backend"},
    {"clojure", "Clojure", "backend"},

    // ========================================
    // FRONTEND FRAMEWORKS & LIBRARIES
    // ========================================
    {"react", "React", "frontend"},
    {"reactjs", "React", "frontend"},
    {"react.js", "React", "frontend"},
    {"vue", "Vue.js", "frontend"},
    {"vuejs", "Vue.js", "frontend"},
    {"vue.js", "Vue.js", "frontend"},
    {"angular", "Angular", "frontend"},
    {"angularjs", "Angular", "frontend"},
    {"svelte", "Svelte", "frontend"},
    {"typescript", "TypeScript", "frontend"},
    {"javascript", "JavaScript", "frontend"},
    {"nextjs", "Next.js", "frontend"},
    {"next.js", "Next.js", "frontend"},
    {"nuxt", "Nuxt.js", "frontend"},
    {"gatsby", "Gatsby", "frontend"},
    {"ember", "Ember.js", "frontend"},
    {"backbone", "Backbone.js", "frontend"},
    {"jquery", "jQuery", "frontend"},
    {"html", "HTML", "frontend"},
    {"css", "CSS", "frontend"},
    {"tailwind", "Tailwind CSS", "frontend"},
    {"tailwindcss", "Tailwind CSS", "frontend"},
    {"bootstrap", "Bootstrap", "frontend"},
    {"sass", "Sass", "frontend"},
    {"scss", "Sass", "frontend"},
    {"less", "Less", "frontend"},
    {"webpack", "Webpack", "frontend"},
    {"vite", "Vite", "frontend"},
    {"rollup", "Rollup", "frontend"},
    {"parcel", "Parcel", "frontend"},

    // ========================================
    // DATABASES
    // ========================================
    {"postgresql", "PostgreSQL", "database"},
    {"postgres", "PostgreSQL", "database"},
    {"mysql", "MySQL", "database"},
    {"mariadb", "MariaDB", "database"},
    {"mongodb", "MongoDB", "database"},
    {"mongo", "MongoDB", "database"},
    {"redis", "Redis", "database"},
    {"elasticsearch", "Elasticsearch", "database"},
    {"elastic", "Elasticsearch", "database"},
    {"dynamodb", "DynamoDB", "database"},
    {"cassandra", "Cassandra", "database"},
    {"couchdb", "CouchDB", "database"},
    {"sqlite", "SQLite", "database"},
    {"oracle", "Oracle DB", "database"},
    {"sql server", "SQL Server", "database"},
    {"mssql", "SQL Server", "database"},
    {"neo4j", "Neo4j", "database"},
    {"influxdb", "InfluxDB", "database"},
    {"timescaledb", "TimescaleDB", "database"},

    // ========================================
    // DEVOPS & INFRASTRUCTURE
    // ========================================
    {"kubernetes", "Kubernetes", "devops"},
    {"k8s", "Kubernetes", "devops"},
    {"docker", "Docker", "devops"},
    {"terraform", "Terraform", "devops"},
    {"ansible", "Ansible", "devops"},
    {"puppet", "Puppet", "devops"},
    {"chef", "Chef", "devops"},
    {"jenkins", "Jenkins", "devops"},
    {"circleci", "CircleCI", "devops"},
    {"circle ci", "CircleCI", "devops"},
    {"github actions", "GitHub Actions", "devops"},
    {"gitlab ci", "GitLab CI", "devops"},
    {"travis", "Travis CI", "devops"},
    {"helm", "Helm", "devops"},
    {"prometheus", "Prometheus", "devops"},
    {"grafana", "Grafana", "devops"},
    {"datadog", "Datadog", "devops"},
    {"new relic", "New Relic", "devops"},
    {"nginx", "Nginx", "devops"},
    {"apache", "Apache", "devops"},
    {"istio", "Istio", "devops"},
    {"envoy", "Envoy", "devops"},
    {"consul", "Consul", "devops"},
    {"vault", "Vault", "devops"},

    // ========================================
    // CLOUD PLATFORMS & SERVICES
    // ========================================
    {"aws", "AWS", "cloud"},
    {"amazon web services", "AWS", "cloud"},
    {"ec2", "AWS EC2", "cloud"},
    {"s3", "AWS S3", "cloud"},
    {"lambda", "AWS Lambda", "cloud"},
    {"rds", "AWS RDS", "cloud"},
    {"ecs", "AWS ECS", "cloud"},
    {"eks", "AWS EKS", "cloud"},
    {"sqs", "AWS SQS", "cloud"},
    {"sns", "AWS SNS", "cloud"},
    {"cloudformation", "CloudFormation", "cloud"},
    {"gcp", "Google Cloud", "cloud"},
    {"google cloud", "Google Cloud", "cloud"},
    {"bigquery", "BigQuery", "cloud"},
    {"cloud run", "Cloud Run", "cloud"},
    {"gke", "Google Kubernetes Engine", "cloud"},
    {"azure", "Azure", "cloud"},
    {"heroku", "Heroku", "cloud"},
    {"vercel", "Vercel", "cloud"},
    {"netlify", "Netlify", "cloud"},
    {"cloudflare", "Cloudflare", "cloud"},
    {"digitalocean", "DigitalOcean", "cloud"},
    {"linode", "Linode", "cloud"},

    // ========================================
    // MOBILE DEVELOPMENT
    // ========================================
    {"ios", "iOS", "mobile"},
    {"android", "Android", "mobile"},
    {"swift", "Swift", "mobile"},
    {"objective-c", "Objective-C", "mobile"},
    {"objective c", "Objective-C", "mobile"},
    {"kotlin", "Kotlin", "mobile"},
    {"react native", "React Native", "mobile"},
    {"flutter", "Flutter", "mobile"},
    {"xamarin", "Xamarin", "mobile"},
    {"ionic", "Ionic", "mobile"},

    // ========================================
    // TOOLING & PROTOCOLS
    // ========================================
    {"git", "Git", "tooling"},
    {"github", "GitHub", "tooling"},
    {"gitlab", "GitLab", "tooling"},
    {"bitbucket", "Bitbucket", "tooling"},
    {"jira", "Jira", "tooling"},
    {"confluence", "Confluence", "tooling"},
    {"slack", "Slack", "tooling"},
    {"figma", "Figma", "tooling"},
    {"sketch", "Sketch", "tooling"},
    {"graphql", "GraphQL", "tooling"},
    {"rest", "REST API", "tooling"},
    {"rest api", "REST API", "tooling"},
    {"grpc", "gRPC", "tooling"},
    {"rabbitmq", "RabbitMQ", "tooling"},
    {"kafka", "Kafka", "tooling"},
    {"oauth", "OAuth", "tooling"},
    {"jwt", "JWT", "tooling"},
    {"websocket", "WebSocket", "tooling"},
    {"soap", "SOAP", "tooling"},
    {"protobuf", "Protocol Buffers", "tooling"},
}

// GetKeywordMap returns a map for O(1) keyword lookup
// This is used by the inference service for fast detection
func GetKeywordMap() map[string]TechnologyKeyword {
    m := make(map[string]TechnologyKeyword, len(TechnologyKeywords))
    for _, kw := range TechnologyKeywords {
        m[kw.Keyword] = kw
    }
    return m
}

// GetAllSkillNames returns unique canonical skill names
// Useful for deduplication and reporting
func GetAllSkillNames() []string {
    seen := make(map[string]bool)
    names := []string{}

    for _, kw := range TechnologyKeywords {
        if !seen[kw.Skill] {
            seen[kw.Skill] = true
            names = append(names, kw.Skill)
        }
    }

    return names
}
```

#### Dictionary Maintenance Guidelines

**Adding New Technologies**:
1. Find the appropriate category section
2. Add keyword entry: `{"keyword", "Canonical Name", "category"}`
3. Add common aliases (e.g., both "nodejs" and "node.js")
4. Keep keywords lowercase for case-insensitive matching
5. Run tests to ensure no duplicates

**Example Addition**:
```go
// Adding Deno to backend section
{"deno", "Deno", "backend"},
```

**TDD Checklist - Phase 1:**
- [ ] Write failing test: TechnologyKeywords slice is not empty
- [ ] Test fails (file doesn't exist)
- [ ] Create technology_keywords.go with ~100 entries
- [ ] Test passes
- [ ] Write failing test: All keywords are lowercase
- [ ] Test fails (some uppercase)
- [ ] Ensure all keywords are lowercase
- [ ] Test passes
- [ ] Write failing test: No duplicate keywords exist
- [ ] Test fails (duplicates found)
- [ ] Remove duplicates
- [ ] Test passes
- [ ] Write failing test: All categories are valid (from CommonSkillCategories)
- [ ] Test fails (invalid category)
- [ ] Fix category names
- [ ] Test passes
- [ ] Write failing test: GetKeywordMap returns correct count
- [ ] Test fails (function doesn't exist)
- [ ] Implement GetKeywordMap
- [ ] Test passes
- [ ] Write failing test: Map allows O(1) lookup
- [ ] Test fails (not a map)
- [ ] Verify map structure
- [ ] Test passes
- [ ] Commit: `test(service): add technology dictionary tests`
- [ ] Commit: `feat(service): add comprehensive technology dictionary`

### Phase 2: Create Skill Inference Service Interface

**Goal**: Define service interface and data types

#### Interface Design

```go
// internal/service/career/skill_inference.go

package career

import (
    "context"
    domain "github.com/baphled/kariya/internal/domain/career"
)

// SkillSuggestion represents a detected skill with metadata
type SkillSuggestion struct {
    Name       string   // Canonical skill name (e.g., "Go", "PostgreSQL")
    Category   string   // Skill category (backend, frontend, etc.)
    Confidence float64  // 0.0-1.0 confidence score
    EventIDs   []string // Events where this skill was detected
    Contexts   []string // Text snippets showing usage (max 3)
}

// SkillInferenceService detects skills from event text
type SkillInferenceService interface {
    // InferSkillsFromEvents analyzes all events for technology mentions
    InferSkillsFromEvents(ctx context.Context, events []*domain.CareerEvent) ([]SkillSuggestion, error)

    // InferSkillsFromBurst analyzes burst events for common skills
    InferSkillsFromBurst(ctx context.Context, burst *domain.Burst, events []*domain.CareerEvent) ([]SkillSuggestion, error)

    // CreateSkillsFromSuggestions persists accepted suggestions and links to events
    CreateSkillsFromSuggestions(ctx context.Context, suggestions []SkillSuggestion) ([]*domain.Skill, error)
}
```

**TDD Checklist - Phase 2:**
- [ ] Write failing test: SkillSuggestion type exists
- [ ] Test fails (type not defined)
- [ ] Define SkillSuggestion struct
- [ ] Test passes
- [ ] Write failing test: SkillInferenceService interface exists
- [ ] Test fails (interface not defined)
- [ ] Define interface with 3 methods
- [ ] Test passes
- [ ] Write failing test: Interface can be mocked
- [ ] Test fails (no implementation)
- [ ] Create mock implementation for tests
- [ ] Test passes
- [ ] Commit: `test(service): add skill inference interface tests`
- [ ] Commit: `feat(service): add skill inference interface`

### Phase 3: Implement Keyword Detection with Word Boundaries

**Goal**: Detect technology keywords in event text with word boundary matching

#### Implementation

```go
// internal/service/career/skill_inference_impl.go

package career

import (
    "context"
    "regexp"
    "strings"
    domain "github.com/baphled/kariya/internal/domain/career"
)

type DefaultSkillInferenceService struct {
    skillRepo   careerrepo.SkillRepository
    keywordMap  map[string]TechnologyKeyword
}

func NewSkillInferenceService(skillRepo careerrepo.SkillRepository) SkillInferenceService {
    return &DefaultSkillInferenceService{
        skillRepo:  skillRepo,
        keywordMap: GetKeywordMap(),
    }
}

func (svc *DefaultSkillInferenceService) InferSkillsFromEvents(
    ctx context.Context,
    events []*domain.CareerEvent,
) ([]SkillSuggestion, error) {
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }

    // Map to collect suggestions by skill name
    suggestionMap := make(map[string]*SkillSuggestion)

    for _, event := range events {
        // Detect skills in this event
        detectedSkills := svc.detectSkillsInText(event.Text, event.ID)

        // Merge into suggestion map
        for _, detected := range detectedSkills {
            if existing, found := suggestionMap[detected.Name]; found {
                // Merge with existing suggestion
                existing.EventIDs = append(existing.EventIDs, detected.EventIDs...)
                existing.Contexts = append(existing.Contexts, detected.Contexts...)

                // Keep highest confidence
                if detected.Confidence > existing.Confidence {
                    existing.Confidence = detected.Confidence
                }

                // Limit contexts to 3
                if len(existing.Contexts) > 3 {
                    existing.Contexts = existing.Contexts[:3]
                }
            } else {
                // New suggestion
                suggestionMap[detected.Name] = detected
            }
        }
    }

    // Convert map to slice
    suggestions := make([]SkillSuggestion, 0, len(suggestionMap))
    for _, suggestion := range suggestionMap {
        suggestions = append(suggestions, *suggestion)
    }

    return suggestions, nil
}

func (svc *DefaultSkillInferenceService) detectSkillsInText(
    text string,
    eventID string,
) []*SkillSuggestion {
    lowerText := strings.ToLower(text)
    detected := []*SkillSuggestion{}

    // Check each keyword in dictionary
    for keyword, tech := range svc.keywordMap {
        // Use word boundary regex to avoid partial matches
        // e.g., "goal" won't match "go"
        pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
        re := regexp.MustCompile(pattern)

        if re.MatchString(lowerText) {
            // Extract context (up to 80 chars around match)
            context := svc.extractContext(text, keyword)

            // Calculate confidence based on usage pattern
            confidence := svc.calculateConfidence(text, keyword)

            detected = append(detected, &SkillSuggestion{
                Name:       tech.Skill,
                Category:   tech.Category,
                Confidence: confidence,
                EventIDs:   []string{eventID},
                Contexts:   []string{context},
            })
        }
    }

    return detected
}

func (svc *DefaultSkillInferenceService) extractContext(text string, keyword string) string {
    lowerText := strings.ToLower(text)
    keywordIndex := strings.Index(lowerText, keyword)

    if keywordIndex == -1 {
        return ""
    }

    // Extract 40 chars before and after keyword
    start := keywordIndex - 40
    if start < 0 {
        start = 0
    }

    end := keywordIndex + len(keyword) + 40
    if end > len(text) {
        end = len(text)
    }

    context := text[start:end]

    // Trim to complete words
    context = strings.TrimSpace(context)

    // Add ellipsis if truncated
    if start > 0 {
        context = "..." + context
    }
    if end < len(text) {
        context = context + "..."
    }

    return context
}
```

**TDD Checklist - Phase 3:**
- [ ] Write failing test: detectSkillsInText finds exact keyword match
- [ ] Test fails (method doesn't exist)
- [ ] Implement detectSkillsInText with regex
- [ ] Test passes
- [ ] Write failing test: Word boundaries prevent partial matches ("goal" != "go")
- [ ] Test fails (partial match occurs)
- [ ] Add `\b` word boundary to regex
- [ ] Test passes
- [ ] Write failing test: Case-insensitive matching works
- [ ] Test fails (case sensitive)
- [ ] Convert text to lowercase before matching
- [ ] Test passes
- [ ] Write failing test: extractContext returns snippet around keyword
- [ ] Test fails (context empty)
- [ ] Implement extractContext
- [ ] Test passes
- [ ] Write failing test: Context limited to ~80 chars
- [ ] Test fails (too long)
- [ ] Add length limits
- [ ] Test passes
- [ ] Write failing test: InferSkillsFromEvents deduplicates same skill
- [ ] Test fails (duplicates present)
- [ ] Add deduplication logic with map
- [ ] Test passes
- [ ] Commit: `test(service): add keyword detection tests`
- [ ] Commit: `feat(service): implement keyword detection`

### Phase 4: Implement Confidence Scoring

**Goal**: Score skill detections based on usage patterns

#### Confidence Scoring Algorithm

```go
// internal/service/career/skill_inference_impl.go

func (svc *DefaultSkillInferenceService) calculateConfidence(text string, keyword string) float64 {
    lowerText := strings.ToLower(text)

    // High confidence patterns (0.9+)
    highPatterns := []string{
        "built with " + keyword,
        "built using " + keyword,
        "developed in " + keyword,
        "developed using " + keyword,
        "implemented in " + keyword,
        "implemented using " + keyword,
        "wrote " + keyword,
        "using " + keyword,
        keyword + " developer",
        keyword + " engineer",
        "expert in " + keyword,
        "proficient in " + keyword,
    }

    for _, pattern := range highPatterns {
        if strings.Contains(lowerText, pattern) {
            return 0.95
        }
    }

    // Medium confidence patterns (0.7)
    mediumPatterns := []string{
        "worked with " + keyword,
        "experience with " + keyword,
        keyword + " project",
        keyword + " system",
        keyword + " application",
        keyword + " service",
        "migrated to " + keyword,
        "integrated " + keyword,
    }

    for _, pattern := range mediumPatterns {
        if strings.Contains(lowerText, pattern) {
            return 0.75
        }
    }

    // Low confidence - simple presence (0.5)
    return 0.5
}
```

**TDD Checklist - Phase 4:**
- [ ] Write failing test: "built with X" gives high confidence (0.9+)
- [ ] Test fails (returns low confidence)
- [ ] Implement high confidence pattern detection
- [ ] Test passes
- [ ] Write failing test: "worked with X" gives medium confidence (0.7)
- [ ] Test fails (returns low confidence)
- [ ] Implement medium confidence patterns
- [ ] Test passes
- [ ] Write failing test: Simple presence gives low confidence (0.5)
- [ ] Test fails (no default)
- [ ] Add default confidence return
- [ ] Test passes
- [ ] Write failing test: Multiple patterns use highest confidence
- [ ] Test fails (using first match)
- [ ] Check all patterns, return highest
- [ ] Test passes
- [ ] Commit: `test(service): add confidence scoring tests`
- [ ] Commit: `feat(service): implement confidence scoring`

### Phase 5: Implement Skill Creation with Event Linking

**Goal**: Persist accepted suggestions and link to events

#### Implementation

```go
// internal/service/career/skill_inference_impl.go

func (svc *DefaultSkillInferenceService) CreateSkillsFromSuggestions(
    ctx context.Context,
    suggestions []SkillSuggestion,
) ([]*domain.Skill, error) {
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }

    createdSkills := []*domain.Skill{}

    for _, suggestion := range suggestions {
        // Check if skill already exists (case-insensitive)
        existing, err := svc.skillRepo.GetByName(ctx, suggestion.Name)

        var skill *domain.Skill

        if err == nil && existing != nil {
            // Skill exists, use it
            skill = existing
        } else {
            // Create new skill
            skill = &domain.Skill{
                ID:        uuid.New().String(),
                Name:      suggestion.Name,
                Category:  suggestion.Category,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }

            if err := svc.skillRepo.Create(ctx, skill); err != nil {
                return nil, fmt.Errorf("failed to create skill %s: %w", suggestion.Name, err)
            }
        }

        // Link skill to events via junction table
        // This is handled by the repository when we update events
        for _, eventID := range suggestion.EventIDs {
            event, err := svc.eventRepo.GetByID(ctx, eventID)
            if err != nil {
                continue // Skip if event not found
            }

            // Add skill ID if not already linked
            if !contains(event.Skills, skill.ID) {
                event.Skills = append(event.Skills, skill.ID)

                // Update event (will save skill associations)
                if err := svc.eventRepo.Update(ctx, event); err != nil {
                    // Log error but continue
                    continue
                }
            }
        }

        // Update LastUsed based on event dates
        if err := svc.updateLastUsed(ctx, skill, suggestion.EventIDs); err != nil {
            // Log warning but don't fail
        }

        createdSkills = append(createdSkills, skill)
    }

    return createdSkills, nil
}

func (svc *DefaultSkillInferenceService) updateLastUsed(
    ctx context.Context,
    skill *domain.Skill,
    eventIDs []string,
) error {
    var mostRecent time.Time

    for _, eventID := range eventIDs {
        event, err := svc.eventRepo.GetByID(ctx, eventID)
        if err != nil {
            continue
        }

        if event.Date.After(mostRecent) {
            mostRecent = event.Date
        }
    }

    if !mostRecent.IsZero() {
        skill.LastUsed = &mostRecent
        skill.UpdatedAt = time.Now()
        return svc.skillRepo.Update(ctx, skill)
    }

    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

**TDD Checklist - Phase 5:**
- [ ] Write failing test: CreateSkillsFromSuggestions creates new skill
- [ ] Test fails (method doesn't exist)
- [ ] Implement CreateSkillsFromSuggestions
- [ ] Test passes
- [ ] Write failing test: Doesn't duplicate existing skills
- [ ] Test fails (duplicate created)
- [ ] Add GetByName check
- [ ] Test passes
- [ ] Write failing test: Links skills to events via junction table
- [ ] Test fails (no linking)
- [ ] Add event update with Skills field
- [ ] Test passes
- [ ] Write failing test: Updates LastUsed to most recent event date
- [ ] Test fails (LastUsed not set)
- [ ] Implement updateLastUsed
- [ ] Test passes
- [ ] Write failing test: Returns list of created/existing skills
- [ ] Test fails (wrong return)
- [ ] Return skill list
- [ ] Test passes
- [ ] Commit: `test(service): add skill creation tests`
- [ ] Commit: `feat(service): implement skill creation and linking`

### Phase 6: Integrate into Burst Confirmation Flow

**Goal**: Add skill inference after fact extraction in burst confirmation

#### Add State to burst_management.go

```go
// internal/cli/intents/burst_management.go

const (
    // ... existing states
    BurstStateExtractingFacts = "extracting_facts"
    BurstStateSkillSuggestions = "skill_suggestions"        // NEW
    BurstStateSkillSuggestionReview = "skill_suggestion_review"  // NEW
)

type SkillSuggestionsLoadedMsg struct {
    Suggestions []career.SkillSuggestion
    Error       error
}

type SkillSuggestionAcceptedMsg struct {
    Suggestion career.SkillSuggestion
    Skill      *domain.Skill
    Error      error
}
```

#### Update Burst Confirmation Flow

```go
// internal/cli/intents/burst_management_intent.go

// After fact extraction completes
func (i *BurstManagementIntent) updateExtractingFactsView(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case FactExtractionCompleteMsg:
        i.state.extractingFacts = false

        if msg.Error != nil {
            i.state.confirmError = msg.Error
            i.state.currentState = BurstStateConfirm
            return nil
        }

        i.state.extractedFactsCount = len(msg.Facts)
        i.state.extractionComplete = true

        // NEW: After facts, infer skills
        return tea.Sequence(
            i.confirmBurstOnly(),
            i.inferSkillsFromBurst(),  // NEW
        )
    }

    return nil
}

func (i *BurstManagementIntent) inferSkillsFromBurst() tea.Cmd {
    return func() tea.Msg {
        if i.state.selectedBurst == nil {
            return SkillSuggestionsLoadedMsg{
                Suggestions: []career.SkillSuggestion{},
            }
        }

        // Load events for this burst
        events := make([]*domain.CareerEvent, 0, len(i.state.selectedBurst.EventIDs))
        for _, eventID := range i.state.selectedBurst.EventIDs {
            event, err := i.context.Service.GetEventByID(i.context.Context, eventID)
            if err == nil {
                events = append(events, event)
            }
        }

        // Infer skills
        inferenceService := career.NewSkillInferenceService(i.context.Service.SkillRepository())
        suggestions, err := inferenceService.InferSkillsFromBurst(
            i.context.Context,
            i.state.selectedBurst,
            events,
        )

        return SkillSuggestionsLoadedMsg{
            Suggestions: suggestions,
            Error:       err,
        }
    }
}
```

**TDD Checklist - Phase 6:**
- [ ] Write failing test: After fact extraction, skill inference triggers
- [ ] Test fails (inference not called)
- [ ] Add inferSkillsFromBurst call after confirmBurstOnly
- [ ] Test passes
- [ ] Write failing test: SkillSuggestionsLoadedMsg is sent
- [ ] Test fails (message not sent)
- [ ] Implement async command
- [ ] Test passes
- [ ] Write failing test: Empty suggestions allowed (no error)
- [ ] Test fails (error on empty)
- [ ] Return empty slice without error
- [ ] Test passes
- [ ] Write failing test: Can accept/reject suggestions
- [ ] Test fails (handlers not implemented)
- [ ] Implement accept/reject (similar to burst suggestions)
- [ ] Test passes
- [ ] Commit: `test(intents): add burst skill inference tests`
- [ ] Commit: `feat(intents): integrate skill inference into burst flow`

### Phase 7: Add Standalone Inference in ManageSkills

**Goal**: Allow users to trigger inference with "i" key

#### Add States and Handler

```go
// internal/cli/intents/manage_skills.go

const (
    // ... existing states
    SkillsStateInferring = "inferring"              // NEW
    SkillsStateSuggestionReview = "suggestion_review"  // NEW
)

// internal/cli/intents/manage_skills_intent.go

func (i *ManageSkillsIntent) updateListView(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        // ... existing handlers

        case "i":  // NEW - Infer skills
            i.state.currentState = SkillsStateInferring
            return i.inferSkillsFromAllEvents()
        }
    }

    return nil
}

func (i *ManageSkillsIntent) inferSkillsFromAllEvents() tea.Cmd {
    return func() tea.Msg {
        // Load all events
        events, err := i.context.Service.ListEvents(i.context.Context, nil)
        if err != nil {
            return SkillSuggestionsLoadedMsg{
                Error: err,
            }
        }

        // Infer skills
        inferenceService := career.NewSkillInferenceService(i.context.Service.SkillRepository())
        suggestions, err := inferenceService.InferSkillsFromEvents(
            i.context.Context,
            events,
        )

        return SkillSuggestionsLoadedMsg{
            Suggestions: suggestions,
            Error:       err,
        }
    }
}
```

**TDD Checklist - Phase 7:**
- [ ] Write failing test: "i" key triggers inference from all events
- [ ] Test fails (key not handled)
- [ ] Add "i" case in updateListView
- [ ] Test passes
- [ ] Write failing test: Shows loading state during inference
- [ ] Test fails (no loading view)
- [ ] Add SkillsStateInferring view
- [ ] Test passes
- [ ] Write failing test: Shows suggestions with confidence
- [ ] Test fails (no review view)
- [ ] Add suggestion review view (reuse pattern)
- [ ] Test passes
- [ ] Write failing test: Accepted suggestions create skills
- [ ] Test fails (not creating)
- [ ] Call CreateSkillsFromSuggestions
- [ ] Test passes
- [ ] Commit: `test(intents): add standalone skill inference tests`
- [ ] Commit: `feat(intents): add standalone skill inference to ManageSkills`

## TDD Checklist (MUST COMPLETE IN ORDER)

### RED Phase
- [ ] Phase 1: Dictionary tests (count, lowercase, duplicates, categories)
- [ ] Phase 2: Interface tests (types, methods, mock)
- [ ] Phase 3: Detection tests (matching, boundaries, context)
- [ ] Phase 4: Confidence tests (high, medium, low patterns)
- [ ] Phase 5: Creation tests (new skills, linking, LastUsed)
- [ ] Phase 6: Burst integration tests (trigger, suggestions, accept)
- [ ] Phase 7: ManageSkills tests ("i" key, loading, review)

### GREEN Phase
- [ ] Phase 1-7: Implement each feature to make tests pass

### REFACTOR Phase
- [ ] Extract common suggestion review UI if duplicated
- [ ] Consider caching keyword map for performance
- [ ] Simplify confidence calculation if complex

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="type(scope): description"`
- [ ] Commit is atomic
- [ ] All tests pass

## Post-Task Checklist
- [ ] `make check-compliance` passes
- [ ] All tests pass (including race detector)
- [ ] Code coverage maintained ≥ 80%
- [ ] Technology dictionary has ~100 entries
- [ ] All checkboxes completed

## Acceptance Criteria
- [ ] Technology dictionary has ~100 comprehensive entries organized by category
- [ ] Skill inference uses word boundary regex (prevents partial matches)
- [ ] Confidence scoring differentiates high/medium/low usage patterns
- [ ] Deduplication merges same skill from multiple events
- [ ] Created skills are linked to source events via junction table
- [ ] LastUsed is set to most recent event date
- [ ] Integration into burst confirmation works (after facts)
- [ ] Standalone "i" key in ManageSkills works
- [ ] Context snippets show actual technology usage (~80 chars)
- [ ] Empty suggestions handled gracefully (no error)

## Expected UX Examples

### After burst confirmation:
```
┌─────────────────────────────────────────────────┐
│ Skill Suggestions (from 5 burst events)         │
│                                                 │
│ Suggestion 1 of 3                               │
│                                                 │
│ Go (backend) - 95% confidence                   │
│ ████████████████████░░░░                        │
│                                                 │
│ Found in 5 events:                              │
│   • "...built API using Go and gRPC for..."    │
│   • "...migrated service to Go for better..."  │
│   • "...wrote Go microservices that handle..." │
│                                                 │
│ [a] Accept  [r] Reject  [A] Accept all  [Esc]  │
└─────────────────────────────────────────────────┘
```

### In ManageSkills after pressing "i":
```
┌─────────────────────────────────────────────────┐
│ Inferring Skills from All Events                │
│                                                 │
│ ⏳ Analyzing 42 events for technology mentions...│
│                                                 │
│ This may take a few moments...                  │
└─────────────────────────────────────────────────┘
```

## Rollback Plan
1. Remove service files (technology_keywords.go, skill_inference*.go)
2. Remove skill suggestion states from burst_management.go
3. Remove inferSkillsFromBurst from burst_management_intent.go
4. Remove "i" key handler from manage_skills_intent.go
5. Run `make check-compliance`
6. All tests should pass

## Dependencies
- Task 45 (burst detection UI) - for similar pattern understanding
- Skill domain model and repository (already exists)

## Next Steps After Completion
- Task 47: Enhanced CV Skills Section (will use skills created by this task)
- Task 44: Wire Enhanced Bullet Generator (independent, can be done in parallel)
