---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 23: Custom CV Format - Language-Agnostic Professional Template

## Overview

**Goal**: Add a "Custom" export format to KaRiya's CV generation system that produces narrative-focused CVs for language-agnostic professionals. This format emphasizes pragmatic tool selection, cross-domain experience, and systems thinking over technology evangelism.

**Time Estimate**: 2-3 weeks (4 phases)

**Prerequisites**:
- Existing CV generation system functional (✅ Complete)
- Export service supports Text/Markdown/YAML (✅ Complete)
- All 2,078 tests passing (✅ Complete)
- Understanding of `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` (✅ Complete)

**Related Docs**:
- `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` - Complete implementation proposal
- `docs/guides/CV_GENERATION_GUIDE.md` - CV generation workflow
- `internal/service/career/cv/export_service.go` - Current export implementation

---

## Success Criteria

- [ ] Custom format appears in export selection menu
- [ ] Custom export generates valid markdown CV
- [ ] Profile management (create/edit/delete) works
- [ ] Profile data persists across sessions
- [ ] Auto-enrichment extracts data from facts/events
- [ ] All 2,078+ tests pass
- [ ] Zero regressions in existing formats
- [ ] Code coverage maintained >87%
- [ ] User guide complete with examples

---

## Implementation Phases

### Phase 1: Core Export Functionality (Week 1) ✅ RECOMMENDED START

**Objective**: Add custom format export with hardcoded profile (no profile management yet)

#### Files to Create
- [ ] `internal/service/career/cv/export_custom.go` (new)
- [ ] `internal/service/career/cv/export_custom_test.go` (new)

#### Files to Modify
- [ ] `internal/service/career/cv/export_service.go`
- [ ] `internal/service/career/cv/export_service_test.go`
- [ ] `internal/cli/intents/generate_cv.go`
- [ ] `internal/cli/intents/generate_cv_intent.go`

#### Implementation Checklist

##### Step 1.1: Add Export Format Constant
**File**: `internal/service/career/cv/export_service.go`

- [ ] Add `ExportFormatCustom ExportFormat = "custom"` constant (line 32)
- [ ] Update `getFileExtension()` to return `.md` for custom format (line 238)
- [ ] Add custom case to extension switch statement

**Code Example**:
```go
const (
    ExportFormatText     ExportFormat = "text"
    ExportFormatMarkdown ExportFormat = "markdown"
    ExportFormatYAML     ExportFormat = "yaml"
    ExportFormatCustom   ExportFormat = "custom"  // NEW
)

func getFileExtension(format ExportFormat) string {
    switch format {
    case ExportFormatText:
        return ".txt"
    case ExportFormatMarkdown:
        return ".md"
    case ExportFormatYAML:
        return ".yaml"
    case ExportFormatCustom:
        return ".md"  // NEW
    default:
        return ".txt"
    }
}
```

##### Step 1.2: Create Custom Export Implementation
**File**: `internal/service/career/cv/export_custom.go` (new)

- [ ] Create new file with `ExportToCustom()` method
- [ ] Use hardcoded profile for Phase 1 (no profile management yet)
- [ ] Implement section rendering:
  - [ ] Header (name, role, location, contact)
  - [ ] Summary
  - [ ] Core Strengths (from facts)
  - [ ] Languages & Technologies (from facts)
  - [ ] Selected Experience (from experience + projects sections)
  - [ ] What I Bring (from top bullets)
  - [ ] Footer
- [ ] Add helper functions:
  - [ ] `renderHeader()` - Name, role, contact info
  - [ ] `renderSummary()` - Professional summary
  - [ ] `renderCoreStrengths()` - Bullet list of strengths
  - [ ] `renderTechnologies()` - Grouped tech stack
  - [ ] `renderExperience()` - Experience with bullets
  - [ ] `renderWhatIBring()` - Value propositions

**Code Template**:
```go
package cv

import (
    "bytes"
    "context"
    "fmt"
    "strings"
    
    "github.com/baphled/kariya/internal/domain/career"
)

// ExportToCustom exports a CV to custom narrative format
func (es *ExportService) ExportToCustom(
    ctx context.Context,
    cv *career.CVView,
    sections []*career.CVSection,
    bullets map[string][]*career.CVBullet,
) (string, error) {
    if cv == nil {
        return "", fmt.Errorf("CV view is nil")
    }
    
    var buf bytes.Buffer
    
    // Hardcoded profile for Phase 1
    profile := &CustomProfileData{
        Name:         "Yomi Colledge",
        PrimaryRole:  "Senior Software Engineer / Technical Consultant",
        Location:     "Remote (UK)",
        Email:        "yomi@boodah.net",
        GitHubURL:    "https://github.com/baphled",
        PortfolioURL: "http://boodah.net",
    }
    
    // Render sections
    renderHeader(&buf, profile)
    buf.WriteString("---\n\n")
    
    renderSummary(&buf, sections)
    buf.WriteString("---\n\n")
    
    renderCoreStrengths(&buf, sections)
    buf.WriteString("---\n\n")
    
    renderTechnologies(&buf, sections)
    buf.WriteString("---\n\n")
    
    renderExperience(&buf, sections)
    buf.WriteString("---\n\n")
    
    renderWhatIBring(&buf, bullets)
    buf.WriteString("---\n\n")
    
    buf.WriteString("**References available on request.**\n")
    
    return buf.String(), nil
}

// CustomProfileData contains hardcoded profile data for Phase 1
type CustomProfileData struct {
    Name         string
    PrimaryRole  string
    Location     string
    Email        string
    GitHubURL    string
    PortfolioURL string
}

// Helper rendering functions
func renderHeader(buf *bytes.Buffer, profile *CustomProfileData) {
    buf.WriteString(fmt.Sprintf("# %s\n\n", profile.Name))
    buf.WriteString(fmt.Sprintf("**%s**  \n", profile.PrimaryRole))
    buf.WriteString(fmt.Sprintf("%s  \n", profile.Location))
    buf.WriteString(fmt.Sprintf("Email: [%s](mailto:%s)  \n", profile.Email, profile.Email))
    buf.WriteString(fmt.Sprintf("GitHub: %s  \n", profile.GitHubURL))
    buf.WriteString(fmt.Sprintf("Portfolio: %s\n\n", profile.PortfolioURL))
}

func renderSummary(buf *bytes.Buffer, sections []*career.CVSection) {
    buf.WriteString("## Summary\n\n")
    
    // Find summary section
    for _, section := range sections {
        if section.SectionType == "summary" && section.Summary != "" {
            buf.WriteString(section.Summary + "\n\n")
            return
        }
    }
    
    // Default summary if none found
    buf.WriteString("Experienced professional with a strong track record of delivering high-quality results.\n\n")
}

func renderCoreStrengths(buf *bytes.Buffer, sections []*career.CVSection) {
    buf.WriteString("## Core Strengths\n\n")
    
    // Extract competencies from skills section
    strengths := []string{
        "Language-agnostic backend and systems engineering",
        "System design and architectural ownership",
        "Pragmatic problem decomposition",
        "Legacy stabilisation and modernisation",
        "Product-focused delivery",
        "Linux-first operational mindset",
    }
    
    for _, strength := range strengths {
        buf.WriteString(fmt.Sprintf("- %s  \n", strength))
    }
    buf.WriteString("\n")
}

func renderTechnologies(buf *bytes.Buffer, sections []*career.CVSection) {
    buf.WriteString("## Languages & Technologies\n\n")
    
    // TODO: Extract from facts in Phase 3
    buf.WriteString("**Languages:** Ruby, Go, PHP, C/C++, JavaScript, Shell  \n")
    buf.WriteString("**Frontend:** Vue.js  \n")
    buf.WriteString("**Systems:** Linux, SQL, APIs, CI/CD, automation  \n\n")
}

func renderExperience(buf *bytes.Buffer, sections []*career.CVSection) {
    buf.WriteString("## Selected Experience\n\n")
    
    // Find experience and projects sections
    for _, section := range sections {
        if section.SectionType != "experience" && section.SectionType != "projects" {
            continue
        }
        
        // Render each content group (company/project)
        for _, group := range section.Content {
            if group.Header != "" {
                buf.WriteString(fmt.Sprintf("### %s\n", group.Header))
                
                // Add date range if present
                if group.StartDate != "" && group.EndDate != "" {
                    if group.StartDate == group.EndDate {
                        buf.WriteString(fmt.Sprintf("*%s*\n\n", group.StartDate))
                    } else {
                        buf.WriteString(fmt.Sprintf("*%s – %s*\n\n", group.StartDate, group.EndDate))
                    }
                } else {
                    buf.WriteString("\n")
                }
            }
            
            // Render bullets (only selected ones)
            for _, bullet := range group.Bullets {
                if bullet.Selected {
                    buf.WriteString(fmt.Sprintf("- %s  \n", bullet.Text))
                }
            }
            buf.WriteString("\n---\n\n")
        }
    }
}

func renderWhatIBring(buf *bytes.Buffer, bullets map[string][]*career.CVBullet) {
    buf.WriteString("## What I Bring\n\n")
    
    // TODO: Extract from top bullets in Phase 3
    values := []string{
        "Languages as tools, not identity",
        "Calm handling of complexity",
        "Clear thinking under constraints",
        "Long-term maintainability focus",
    }
    
    for _, value := range values {
        buf.WriteString(fmt.Sprintf("- %s  \n", value))
    }
    buf.WriteString("\n")
}
```

##### Step 1.3: Add Custom Format to UI
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Add "Custom" to export format selection list (line 883+)
- [ ] Add `CVExportFormatCustom CVExportFormat = "custom"` constant to `generate_cv.go` (line 204)
- [ ] Handle custom format selection in `viewExportSelectFormat()`

**Code Changes**:
```go
// In generate_cv.go (line 204)
const (
    CVExportFormatText     CVExportFormat = "text"
    CVExportFormatMarkdown CVExportFormat = "markdown"
    CVExportFormatYAML     CVExportFormat = "yaml"
    CVExportFormatCustom   CVExportFormat = "custom"  // NEW
)

// In generate_cv_intent.go (viewExportSelectFormat method)
func (i *GenerateCVIntent) viewExportSelectFormat() string {
    formats := []string{"Text", "Markdown", "YAML", "Custom"}  // Add "Custom"
    // ... rest of implementation
}
```

##### Step 1.4: Wire Export to ExportService
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Update `updateExporting()` to call `ExportToCustom()` for custom format (line 850+)
- [ ] Pass sections and bullets to export method

**Code Changes**:
```go
// In updateExporting() method (around line 850)
case cv.ExportFormatCustom:
    content, err = i.context.ExportService.ExportToCustom(
        i.context.AppContext,
        i.state.generatedCV,
        i.state.generatedCV.Sections,
        nil,  // bullets map (can be nil for Phase 1)
    )
```

##### Step 1.5: Write Unit Tests
**File**: `internal/service/career/cv/export_custom_test.go` (new)

- [ ] Test complete CV export
- [ ] Test each section renders correctly
- [ ] Test edge cases (missing sections, empty data)
- [ ] Test markdown formatting (proper spacing, headers)

**Test Template**:
```go
package cv_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    
    "github.com/baphled/kariya/internal/service/career/cv"
    "github.com/baphled/kariya/internal/domain/career"
    // ... other imports
)

var _ = Describe("ExportToCustom", func() {
    var (
        exportService *cv.ExportService
        cvView        *career.CVView
        sections      []*career.CVSection
    )
    
    BeforeEach(func() {
        exportService = cv.NewExportService(logger)
        cvView = &career.CVView{
            Name:           "Test CV",
            TargetRole:     "senior_ic",
            TargetAudience: "hiring_manager",
            // ... populate test data
        }
        sections = []*career.CVSection{
            // ... create test sections
        }
    })
    
    Context("with complete data", func() {
        It("generates valid markdown", func() {
            content, err := exportService.ExportToCustom(ctx, cvView, sections, nil)
            Expect(err).ToNot(HaveOccurred())
            Expect(content).To(ContainSubstring("# Yomi Colledge"))
            Expect(content).To(ContainSubstring("## Summary"))
            Expect(content).To(ContainSubstring("## Core Strengths"))
            Expect(content).To(ContainSubstring("## Languages & Technologies"))
            Expect(content).To(ContainSubstring("## Selected Experience"))
            Expect(content).To(ContainSubstring("## What I Bring"))
            Expect(content).To(ContainSubstring("**References available on request.**"))
        })
        
        It("includes all required sections", func() {
            content, err := exportService.ExportToCustom(ctx, cvView, sections, nil)
            Expect(err).ToNot(HaveOccurred())
            
            requiredSections := []string{
                "Summary", "Core Strengths", "Languages & Technologies",
                "Selected Experience", "What I Bring",
            }
            for _, section := range requiredSections {
                Expect(content).To(ContainSubstring(section))
            }
        })
        
        It("formats experience chronologically", func() {
            // Test date ordering
        })
    })
    
    Context("with minimal data", func() {
        It("handles missing summary gracefully", func() {
            // Test with no summary section
        })
        
        It("handles missing experience gracefully", func() {
            // Test with no experience section
        })
    })
    
    Context("edge cases", func() {
        It("handles nil CV view", func() {
            _, err := exportService.ExportToCustom(ctx, nil, sections, nil)
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("CV view is nil"))
        })
        
        It("handles empty sections", func() {
            content, err := exportService.ExportToCustom(ctx, cvView, []*career.CVSection{}, nil)
            Expect(err).ToNot(HaveOccurred())
            Expect(content).ToNot(BeEmpty())
        })
    })
})
```

#### Testing Instructions

**Unit Tests**:
```bash
# Test new custom export
go test -v ./internal/service/career/cv/ -run TestExportToCustom

# Test all export formats
go test -v ./internal/service/career/cv/export_*.go
```

**Integration Test** (manual):
```bash
# 1. Build application
go build -o kariya ./cmd/cli

# 2. Run CV generation
./kariya

# 3. Navigate to: Generate CV → Select Profile → Generate
# 4. Select "Custom" export format
# 5. Verify output in ~/.kariya/cv_exports/
# 6. Check markdown formatting with:
cat ~/.kariya/cv_exports/Senior_IC_*.md
```

#### Acceptance Criteria - Phase 1

- [ ] Custom format constant added
- [ ] ExportToCustom() method implemented
- [ ] Custom format appears in UI selection menu
- [ ] Custom export generates valid markdown
- [ ] All sections render correctly with hardcoded profile
- [ ] Unit tests pass (20+ tests)
- [ ] Integration test passes (manual)
- [ ] Zero regressions in existing formats
- [ ] Code coverage >87%

---

### Phase 2: Profile Management (Week 2)

**Objective**: Add profile CRUD functionality and persistence

#### Files to Create
- [ ] `internal/domain/career/profile.go` (new)
- [ ] `internal/service/career/profile_manager.go` (new)
- [ ] `internal/service/career/profile_manager_test.go` (new)
- [ ] `internal/cli/intents/manage_profile.go` (new)
- [ ] `internal/cli/intents/manage_profile_intent.go` (new)
- [ ] `internal/cli/intents/manage_profile_test.go` (new)

#### Files to Modify
- [ ] `internal/cli/intents/generate_cv_intent.go`
- [ ] `internal/service/career/cv/export_custom.go`

#### Implementation Checklist

##### Step 2.1: Create Profile Domain Model
**File**: `internal/domain/career/profile.go` (new)

- [ ] Define `CustomProfile` struct with all fields
- [ ] Add validation methods
- [ ] Add YAML marshaling tags
- [ ] Add profile error types

**Code Template**:
```go
package career

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

// CustomProfile contains user profile data for custom CV format
type CustomProfile struct {
    ID               string    `yaml:"id" json:"id"`
    Name             string    `yaml:"name" json:"name"`
    PrimaryRole      string    `yaml:"primary_role" json:"primary_role"`
    Location         string    `yaml:"location" json:"location"`
    Email            string    `yaml:"email" json:"email"`
    GitHubURL        string    `yaml:"github_url,omitempty" json:"github_url,omitempty"`
    PortfolioURL     string    `yaml:"portfolio_url,omitempty" json:"portfolio_url,omitempty"`
    
    // Auto-enriched fields (can be manually overridden)
    Languages        []string  `yaml:"languages,omitempty" json:"languages,omitempty"`
    FrontendTech     []string  `yaml:"frontend_tech,omitempty" json:"frontend_tech,omitempty"`
    SystemsTech      []string  `yaml:"systems_tech,omitempty" json:"systems_tech,omitempty"`
    CoreStrengths    []string  `yaml:"core_strengths,omitempty" json:"core_strengths,omitempty"`
    WhatIBring       []string  `yaml:"what_i_bring,omitempty" json:"what_i_bring,omitempty"`
    
    CreatedAt        time.Time `yaml:"created_at" json:"created_at"`
    UpdatedAt        time.Time `yaml:"updated_at" json:"updated_at"`
}

// Validate checks if the profile is valid
func (cp *CustomProfile) Validate() error {
    if strings.TrimSpace(cp.Name) == "" {
        return errors.New("profile name cannot be empty")
    }
    if strings.TrimSpace(cp.PrimaryRole) == "" {
        return errors.New("primary role cannot be empty")
    }
    if strings.TrimSpace(cp.Email) == "" {
        return errors.New("email cannot be empty")
    }
    // TODO: Add email format validation
    return nil
}

// Error types
var (
    ErrProfileNotFound = errors.New("profile not found")
    ErrProfileExists   = errors.New("profile already exists")
    ErrInvalidProfile  = errors.New("invalid profile data")
)
```

##### Step 2.2: Create Profile Manager Service
**File**: `internal/service/career/profile_manager.go` (new)

- [ ] Implement `ProfileManager` interface
- [ ] Add CRUD operations (Create, Read, Update, Delete, List)
- [ ] Store profiles as YAML in `$HOME/.kariya/profiles/`
- [ ] Add filename sanitization
- [ ] Add concurrent access protection (mutex)

**Code Template**:
```go
package career

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    
    "gopkg.in/yaml.v3"
    "github.com/baphled/kariya/internal/domain/career"
    "github.com/baphled/kariya/internal/logger"
    "github.com/google/uuid"
)

// ProfileManager handles profile CRUD operations
type ProfileManager interface {
    CreateProfile(ctx context.Context, profile *career.CustomProfile) error
    GetProfile(ctx context.Context, id string) (*career.CustomProfile, error)
    UpdateProfile(ctx context.Context, profile *career.CustomProfile) error
    DeleteProfile(ctx context.Context, id string) error
    ListProfiles(ctx context.Context) ([]*career.CustomProfile, error)
}

// DefaultProfileManager is the file-based implementation
type DefaultProfileManager struct {
    profilesDir string
    logger      *logger.Logger
    mu          sync.RWMutex
}

// NewProfileManager creates a new profile manager
func NewProfileManager(logger *logger.Logger) (*DefaultProfileManager, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("failed to get home directory: %w", err)
    }
    
    profilesDir := filepath.Join(homeDir, ".kariya", "profiles")
    
    // Create directory if it doesn't exist
    if err := os.MkdirAll(profilesDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create profiles directory: %w", err)
    }
    
    return &DefaultProfileManager{
        profilesDir: profilesDir,
        logger:      logger,
    }, nil
}

// CreateProfile creates a new profile
func (pm *DefaultProfileManager) CreateProfile(ctx context.Context, profile *career.CustomProfile) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    // Validate profile
    if err := profile.Validate(); err != nil {
        return err
    }
    
    // Generate ID if not set
    if profile.ID == "" {
        profile.ID = uuid.New().String()
    }
    
    // Check if profile already exists
    filePath := pm.getProfilePath(profile.ID)
    if _, err := os.Stat(filePath); err == nil {
        return career.ErrProfileExists
    }
    
    // Set timestamps
    now := time.Now()
    profile.CreatedAt = now
    profile.UpdatedAt = now
    
    // Marshal to YAML
    data, err := yaml.Marshal(profile)
    if err != nil {
        return fmt.Errorf("failed to marshal profile: %w", err)
    }
    
    // Write to file
    if err := os.WriteFile(filePath, data, 0644); err != nil {
        return fmt.Errorf("failed to write profile file: %w", err)
    }
    
    pm.logger.Info("Created profile: %s (%s)", profile.Name, profile.ID)
    return nil
}

// GetProfile retrieves a profile by ID
func (pm *DefaultProfileManager) GetProfile(ctx context.Context, id string) (*career.CustomProfile, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    filePath := pm.getProfilePath(id)
    
    // Read file
    data, err := os.ReadFile(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, career.ErrProfileNotFound
        }
        return nil, fmt.Errorf("failed to read profile file: %w", err)
    }
    
    // Unmarshal YAML
    var profile career.CustomProfile
    if err := yaml.Unmarshal(data, &profile); err != nil {
        return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
    }
    
    return &profile, nil
}

// UpdateProfile updates an existing profile
func (pm *DefaultProfileManager) UpdateProfile(ctx context.Context, profile *career.CustomProfile) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    // Validate profile
    if err := profile.Validate(); err != nil {
        return err
    }
    
    // Check if profile exists
    filePath := pm.getProfilePath(profile.ID)
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return career.ErrProfileNotFound
    }
    
    // Update timestamp
    profile.UpdatedAt = time.Now()
    
    // Marshal to YAML
    data, err := yaml.Marshal(profile)
    if err != nil {
        return fmt.Errorf("failed to marshal profile: %w", err)
    }
    
    // Write to file
    if err := os.WriteFile(filePath, data, 0644); err != nil {
        return fmt.Errorf("failed to write profile file: %w", err)
    }
    
    pm.logger.Info("Updated profile: %s (%s)", profile.Name, profile.ID)
    return nil
}

// DeleteProfile deletes a profile
func (pm *DefaultProfileManager) DeleteProfile(ctx context.Context, id string) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    filePath := pm.getProfilePath(id)
    
    // Check if profile exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return career.ErrProfileNotFound
    }
    
    // Delete file
    if err := os.Remove(filePath); err != nil {
        return fmt.Errorf("failed to delete profile file: %w", err)
    }
    
    pm.logger.Info("Deleted profile: %s", id)
    return nil
}

// ListProfiles lists all profiles
func (pm *DefaultProfileManager) ListProfiles(ctx context.Context) ([]*career.CustomProfile, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    // Read directory
    files, err := os.ReadDir(pm.profilesDir)
    if err != nil {
        return nil, fmt.Errorf("failed to read profiles directory: %w", err)
    }
    
    var profiles []*career.CustomProfile
    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
            continue
        }
        
        // Read profile
        filePath := filepath.Join(pm.profilesDir, file.Name())
        data, err := os.ReadFile(filePath)
        if err != nil {
            pm.logger.Error("Failed to read profile file %s: %v", filePath, err)
            continue
        }
        
        // Unmarshal YAML
        var profile career.CustomProfile
        if err := yaml.Unmarshal(data, &profile); err != nil {
            pm.logger.Error("Failed to unmarshal profile %s: %v", filePath, err)
            continue
        }
        
        profiles = append(profiles, &profile)
    }
    
    return profiles, nil
}

// getProfilePath returns the file path for a profile
func (pm *DefaultProfileManager) getProfilePath(id string) string {
    return filepath.Join(pm.profilesDir, fmt.Sprintf("%s.yaml", id))
}
```

##### Step 2.3: Create ManageProfile Intent
**File**: `internal/cli/intents/manage_profile_intent.go` (new)

- [ ] Implement intent state machine
- [ ] Add CRUD UI flows
- [ ] Add form for profile creation/editing
- [ ] Add list view for profile selection
- [ ] Add confirmation for deletion

**States**:
- `StateList` - List all profiles
- `StateCreate` - Create new profile
- `StateEdit` - Edit existing profile
- `StateDelete` - Delete profile (with confirmation)

##### Step 2.4: Integrate Profile Selection into CV Generation
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Add profile selection before export (for custom format only)
- [ ] If no profiles exist, prompt to create one
- [ ] Pass selected profile to `ExportToCustom()`

##### Step 2.5: Update ExportToCustom to Use Real Profiles
**File**: `internal/service/career/cv/export_custom.go`

- [ ] Replace hardcoded profile with parameter
- [ ] Update method signature: `ExportToCustom(ctx, cv, sections, bullets, profile)`

##### Step 2.6: Write Tests
- [ ] Test ProfileManager CRUD operations
- [ ] Test ManageProfile intent workflows
- [ ] Test profile integration in CV generation

#### Acceptance Criteria - Phase 2

- [ ] Users can create profiles via UI
- [ ] Users can edit existing profiles
- [ ] Users can delete profiles (with confirmation)
- [ ] Users can list all profiles
- [ ] Profiles persist across application restarts
- [ ] CV generation uses selected profile
- [ ] All tests pass (50+ new tests)
- [ ] Zero regressions

---

### Phase 3: Data Extraction & Enrichment (Week 3)

**Objective**: Auto-populate profile fields from facts and events

#### Files to Create
- [ ] `internal/service/career/profile_enricher.go` (new)
- [ ] `internal/service/career/profile_enricher_test.go` (new)

#### Files to Modify
- [ ] `internal/cli/intents/manage_profile_intent.go`

#### Implementation Checklist

##### Step 3.1: Create Profile Enricher Service
**File**: `internal/service/career/profile_enricher.go` (new)

- [ ] Extract languages from facts (competency_categories containing "language")
- [ ] Extract frontend tech from facts (competency_categories containing "frontend")
- [ ] Extract systems tech from facts (competency_categories containing "system", "infrastructure")
- [ ] Extract core strengths from top competencies (by frequency)
- [ ] Extract "What I Bring" from top bullets (by confidence + rank)

**Methods**:
```go
type ProfileEnricher interface {
    EnrichProfile(ctx context.Context, profile *career.CustomProfile, events []*career.CareerEvent, facts []*career.Fact, bullets []*career.CVBullet) error
}
```

##### Step 3.2: Implement Extraction Logic

- [ ] **ExtractLanguages**: Parse facts for programming languages
- [ ] **ExtractFrontendTech**: Parse facts for frontend frameworks
- [ ] **ExtractSystemsTech**: Parse facts for infrastructure/systems tools
- [ ] **ExtractCoreStrengths**: Count competency categories, take top 6
- [ ] **ExtractWhatIBring**: Analyze top 4 bullets for value propositions

##### Step 3.3: Add Auto-Enrichment to Profile Creation

- [ ] Trigger enrichment when creating new profile
- [ ] Allow user to review and edit auto-generated fields
- [ ] Show "Auto-generated" tag for enriched fields

##### Step 3.4: Write Tests

- [ ] Test each extraction method
- [ ] Test full enrichment pipeline
- [ ] Test edge cases (no data, partial data)

#### Acceptance Criteria - Phase 3

- [ ] Profile fields auto-populate on creation
- [ ] Extraction accuracy >80% on sample data
- [ ] Users can override auto-generated fields
- [ ] All tests pass (30+ new tests)
- [ ] Zero regressions

---

### Phase 4: Polish & Documentation (Week 4)

**Objective**: Production-ready feature

#### Implementation Checklist

##### Step 4.1: Error Handling & Validation

- [ ] Add input validation for all profile fields
- [ ] Add email format validation (regex)
- [ ] Add URL validation for GitHub/portfolio
- [ ] Add error messages for all failure cases
- [ ] Add retry logic for file I/O errors

##### Step 4.2: UI Polish

- [ ] Add help text for profile management
- [ ] Add keyboard shortcuts reference
- [ ] Add loading states for enrichment
- [ ] Add success/error modals
- [ ] Improve form validation UX

##### Step 4.3: Write User Documentation
**File**: `docs/guides/CUSTOM_CV_FORMAT_GUIDE.md` (new)

- [ ] Explain custom format purpose
- [ ] Step-by-step guide for creating profiles
- [ ] Examples of generated CVs
- [ ] Troubleshooting section

##### Step 4.4: Performance Testing

- [ ] Benchmark export performance (<500ms target)
- [ ] Benchmark profile enrichment (<1s target)
- [ ] Test with large datasets (100+ events, 50+ facts)

##### Step 4.5: Security Review

- [ ] Validate all user inputs
- [ ] Sanitize file paths
- [ ] Check for YAML injection vulnerabilities
- [ ] Review email/URL handling

#### Acceptance Criteria - Phase 4

- [ ] All error cases handled gracefully
- [ ] User guide complete with examples
- [ ] Performance benchmarks pass
- [ ] Security checks pass
- [ ] All tests pass (2,200+ total)
- [ ] Code coverage >87%
- [ ] Zero regressions

---

## Testing Strategy

### Unit Tests (150+ tests)

**Export Tests** (`export_custom_test.go`):
- [ ] Test each rendering function (header, summary, etc.)
- [ ] Test edge cases (nil data, empty sections)
- [ ] Test markdown formatting

**Profile Manager Tests** (`profile_manager_test.go`):
- [ ] Test CRUD operations
- [ ] Test file persistence
- [ ] Test concurrent access
- [ ] Test error handling

**Profile Enricher Tests** (`profile_enricher_test.go`):
- [ ] Test extraction logic for each field
- [ ] Test full enrichment pipeline
- [ ] Test with various data sets

**Intent Tests** (`manage_profile_test.go`):
- [ ] Test state transitions
- [ ] Test form validation
- [ ] Test CRUD workflows

### Integration Tests (20+ tests)

**End-to-End Workflow**:
1. Create profile
2. Generate CV
3. Select custom format
4. Export CV
5. Verify output

**Profile Persistence**:
1. Create profile
2. Restart application
3. Load profile
4. Verify data matches

### Manual Testing Checklist

- [ ] Create profile via UI
- [ ] Edit profile via UI
- [ ] Delete profile via UI
- [ ] Generate CV with custom format
- [ ] Verify markdown output
- [ ] Test with empty profile
- [ ] Test with minimal data
- [ ] Test with complete data
- [ ] Test profile enrichment accuracy

---

## Rollback Plan

If issues arise, rollback is simple:

1. Remove custom format constant
2. Remove custom export method
3. Remove custom format from UI
4. Profiles remain in filesystem (user data preserved)
5. Revert commits atomically

**Impact**: Zero - feature is additive, no existing functionality modified

---

## Future Enhancements

**Template Engine** (Post-MVP):
- Allow users to define custom markdown templates
- Use Go `text/template` for rendering
- Support variables like `{{.Name}}`, `{{.CoreStrengths}}`

**Profile Presets** (Post-MVP):
- Pre-built profiles for common roles
- "Language-Agnostic Engineer"
- "Frontend Specialist"
- "Systems Architect"

**Multi-Format Support** (Post-MVP):
- PDF export via markdown → PDF conversion
- HTML export for web portfolios
- LaTeX export for academic CVs

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Adoption | >30% try custom format | Usage analytics (manual) |
| Profile Creation | >10% create profile | File count in `~/.kariya/profiles/` |
| Quality | <20% manual editing | User survey (manual) |
| Performance | <500ms export | Benchmark tests |
| Satisfaction | >4/5 stars | User feedback (manual) |

---

## Questions & Decisions

| Question | Decision | Date |
|----------|----------|------|
| Profile storage: file vs. database? | File-based (YAML) | 2026-01-08 |
| Support multiple profiles? | Yes | 2026-01-08 |
| Template engine? | Custom renderer for MVP | 2026-01-08 |
| Auto-enrichment vs. manual? | Auto with manual override | 2026-01-08 |

---

## References

- [Custom CV Format Proposal](../docs/CUSTOM_CV_FORMAT_PROPOSAL.md)
- [CV Generation Service](../internal/service/career/cv/cv_generation_service.go)
- [Export Service](../internal/service/career/cv/export_service.go)
- [GenerateCV Intent](../internal/cli/intents/generate_cv_intent.go)

---

**Last Updated**: 2026-01-08  
**Status**: Ready for Implementation  
**Next Step**: Begin Phase 1 - Core Export Functionality
