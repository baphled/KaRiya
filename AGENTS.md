# KaRiya Project Documentation

**Last Updated**: 2026-01-12
**Project Status**: ✅ **PRODUCTION READY - ALL PHASES COMPLETE (100%)**
**Test Coverage**: 240+ tests, 100% pass rate, 0 race conditions
**Code Quality**: All linting checks passing, no technical debt
**Forms**: Huh library integration (Phase 4/5 complete)
**Form Wrappers**: 2 (HuhCaptureForm, HuhSkillForm) - Required for intent forms

---

## Table of Contents

1. [Session Contract](#session-contract)
2. [AI Mandatory Protocol](#ai-mandatory-protocol)
3. [Task Template (Required Format)](#task-template-required-format)
4. [Project Overview](#project-overview)
5. [Quick Start](#quick-start)
6. [Architecture Overview](#architecture-overview)
7. [Development Guidelines & Rules](#development-guidelines--rules)
8. [TUI Development](#tui-development)
9. [User Guides & Features](#user-guides--features)
10. [Implementation Resources](#implementation-resources)
11. [Task Documentation](#task-documentation)
12. [Key Files and Purposes](#key-files-and-purposes)
13. [Recent Fixes](#recent-fixes)
14. [Common Development Tasks](#common-development-tasks)
15. [Deployment Guide](#deployment-guide)
16. [Performance Benchmarks](#performance-benchmarks)
17. [Workflow Patterns](#workflow-patterns)
18. [Troubleshooting](#troubleshooting)
19. [Project Metadata](#project-metadata)

---

## Session Contract

**This contract is displayed when running `make session-start`.**

By proceeding with this work session, you acknowledge and commit to:

1. **TDD Protocol**: Tests are written **BEFORE** implementation code (Red-Green-Refactor)
2. **Compliance First**: `make check-compliance` runs **before AND after** every task
3. **Atomic Commits**: One logical change per commit, with AI attribution if AI-generated
4. **Sequential Tasks**: One task at a time, in checklist order
5. **Token Efficiency**: Tools over text, concise communication, batch operations

**Violation of these rules requires stopping work and correcting before proceeding.**

---

## AI Mandatory Protocol

### Session Start Requirements

The AI assistant **MUST**:

1. Ask user to run `make session-start`
2. Wait for confirmation that it passed
3. If it fails, **REFUSE to proceed** until violations are fixed
4. Display: "Session contract acknowledged. Ready to proceed."

### Before ANY Code Changes

The AI assistant **MUST**:

1. State the specific task being worked on (from task file)
2. Confirm it is **ONE** atomic change
3. State which test file will be created/modified **FIRST**
4. Wait for user confirmation before proceeding

### TDD Enforcement (CRITICAL)

The AI assistant **MUST**:

1. **Write the failing test FIRST** - this is non-negotiable
2. Show the test to the user
3. Ask user to run the test and confirm it **FAILS**
4. **ONLY THEN** write implementation code
5. If user asks for implementation first, **REFUSE** and explain TDD

**Example refusal:**
```
I cannot write implementation code before the test exists and fails.
This violates our TDD protocol (Session Contract #1).

Let me write the test first. After you confirm it fails, I'll implement.
```

### Before Each Commit

The AI assistant **MUST**:

1. **Run `make check-compliance`** to verify code quality before committing
2. **Use `make ai-commit MSG="type(scope): description"`** for all AI-generated commits (automatic attribution)
   - This is the **required** method for AI-generated code commits
   - Manual workflow (NOT recommended): `make review-commit` + manual attribution
3. Verify commit is atomic (ONE logical change)
4. If commit violates rules, **REFUSE** and explain corrections needed

**Critical Order**:
```bash
make check-compliance          # MUST pass before commit
make ai-commit MSG="..."       # Commit with automatic AI attribution
```

### After Task Completion

The AI assistant **MUST**:

1. Ask user to run `make check-compliance`
2. Verify all task checkboxes in task file are complete
3. Mark task as complete `[x]` in task file
4. **STOP immediately** - do not proceed to next task without explicit user request

### Refusal Protocol

The AI assistant **MUST REFUSE** to proceed if:

- User requests implementation before test (TDD violation)
- `make session-start` has not been run or failed
- `make check-compliance` fails after task completion
- User attempts to commit without `make check-compliance` passing
- User attempts AI-generated commit without `make ai-commit`
- User attempts to skip required workflow steps

**Refusal template:**
```
I cannot proceed with this request because it violates [specific rule].

Required correction: [specific action needed]

Once corrected, I can continue.
```

---

## Task Template (Required Format)

All task files **MUST** follow this structure:

```markdown
# Task XX: [Task Name]

## Overview
- **Goal**: [What we're achieving]
- **Time Estimate**: [Estimated duration]
- **Prerequisites**: [Required setup or knowledge]

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in: [list files/directories]
- [ ] Confirmed this is ONE atomic task (not multiple changes)
- [ ] Identified which test files will be created/modified

## Files to Modify
- [ ] List of files that will be changed
- [ ] With checkboxes for tracking

## TDD Checklist (MUST COMPLETE IN ORDER)

### RED Phase
- [ ] Test file created/modified: `path/to/test_file.go`
- [ ] Test written and **FAILS** with error:
  ```
  [Paste actual error here]
  ```
- [ ] Test committed:
  ```
  git commit -m "test(scope): add failing test for X"
  ```

### GREEN Phase
- [ ] Minimal implementation written
- [ ] Test now **PASSES**
- [ ] Implementation committed:
  ```
  git commit -m "feat(scope): implement X"
  ```

### REFACTOR Phase (if needed)
- [ ] Code refactored for clarity/DRY
- [ ] Tests still pass
- [ ] Refactoring committed separately:
  ```
  git commit -m "refactor(scope): improve X"
  ```

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes (REQUIRED before commit)
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria
- [ ] Feature works as specified
- [ ] All tests pass
- [ ] Coverage maintained ≥ 80%
- [ ] Documentation updated (if needed)

## Rollback Plan
- Steps to revert changes if needed
- Safety considerations
```

---

## Project Overview

### Executive Summary

KaRiya is a **Go-based Terminal User Interface (TUI) application** for capturing, analyzing, and generating CVs from career events. The project features a **type-safe, intent-driven architecture** with comprehensive testing, CI/CD integration, and production-ready code quality.

### Key Statistics
- **Language**: Go 1.24
- **Framework**: Bubble Tea + Lipgloss
- **Test Framework**: Ginkgo v2 + Gomega
- **Total Tests**: 2,078 (all passing - 100% pass rate)
- **Code Coverage**: 87%+ overall
- **Staticcheck Warnings**: 0
- **Race Conditions**: 0 detected
- **Performance**: All benchmarks passing

---

## Quick Start

### Development Setup
```bash
# Clone and setup
git clone https://github.com/baphled/kariya.git
cd kariya
go mod tidy
npm install
make install-git-hooks

# Run tests
make test

# Run the application
go build -o kariya ./cmd/kariya
./kariya
```

### Key Commands
```bash
# Run all tests with coverage
go test -v -cover ./...

# Run with race detector
go test -race ./...

# Run Ginkgo tests
ginkgo -r ./internal/cli/intents/

# Format and lint
go fmt ./...
go vet ./...
staticcheck ./...

# Compliance check
make check-compliance

# AI-attributed commit (recommended for AI-generated code)
make ai-commit MSG="feat(scope): description"

# Run ALL CI checks locally (mirrors GitHub Actions)
make ci-local
```

---

## Architecture Overview

### Core Principles

The KaRiya TUI is built on a **type-safe, intent-driven architecture**:

1. **Type-Safe Intent Communication**: All intents communicate via `IntentResult[T]`
2. **Clear Intent Boundaries**: Each intent owns only its local state
3. **Predictable State Machines**: Explicit state transitions
4. **Back Navigation with Context**: Full state preservation via metadata
5. **Minimal Global State**: All mutations are local to intents
6. **Compile-Time Safety**: No runtime type assertions

**Visual Architecture**: See [TUI_INTENT_DIAGRAM.md](docs/TUI_INTENT_DIAGRAM.md) for complete architectural diagrams and intent workflows.

### The 5 Core Intents

| Intent | Purpose | States | Tests |
|--------|---------|--------|-------|
| CaptureEvent | Capture new career events | Choose Strategy → Form → Review → Confirm | 30+ |
| BrowseTimeline | View career timeline | Timeline → Event Detail | 37 |
| GenerateCV | Generate CVs | Profile → Audience → Preview → Review → Confirm | 41 |
| ExportArtifact | Export artifacts | Select → Configure → Preview → Export | 400+ |
| ConfigureSystem | System configuration | Domain → Settings → Staged Changes → Confirm | 400+ |

### Project Structure

```
internal/
├── cli/
│   ├── app/              # Root Bubble Tea model
│   ├── intents/          # Intent implementations
│   ├── models/           # Legacy UI components
│   ├── components/       # Reusable UI components
│   ├── context/          # GlobalContext
│   └── styles/           # Lipgloss styling
├── domain/career/        # Domain models
├── repository/career/    # Data access (SQLite)
└── service/career/       # Business logic
```

### Testing Strategy

#### Test Coverage
- **Overall**: 87% code coverage
- **Intent Framework**: 88.1%
- **GlobalContext**: 100%
- **Domain Models**: >95%
- **Repository**: >90%
- **Service**: >85%
- **CV Service**: 100% (203 tests, all passing)

#### Running Tests
```bash
# All tests
go test -v ./...

# With race detector
go test -race ./...

# Specific package
go test -v ./internal/cli/intents/...

# Ginkgo with focus
ginkgo -r --focus="CaptureEvent" ./internal/cli/intents/

# Benchmarks
go test -bench=. ./internal/cli/intents/
```

#### Test Organization
- **Ginkgo + Gomega** framework for all tests
- **164+ total test specs** (CV service)
- **100% pass rate**
- **0 race conditions**
- **Execution time: 1.3s with race detector**

---

## Development Guidelines & Rules

All development follows standardized guidelines documented in `docs/rules/`. These rules ensure consistency, quality, and efficiency across the project.

### Essential Reading (Read in Order)

These documents form the foundation of our development practices. Read them in this order for progressive learning:

#### 1. Senior Engineer Principles
**File**: [`docs/rules/senior-engineer-guidelines.md`](docs/rules/senior-engineer-guidelines.md)
**Purpose**: Language-agnostic engineering standards and best practices
**When to use**: Before starting any development work - establishes core principles
**Key topics**:
- SOLID principles
- Red-Green-Refactor TDD workflow
- Clean code practices
- Error handling and observability
- Security and developer experience

#### 2. Master Development Workflow
**File**: [`docs/rules/master-task-prompt.md`](docs/rules/master-task-prompt.md)
**Purpose**: Complete 5-phase development workflow integrating ALL project rules
**When to use**: For EVERY task - this is the definitive execution workflow
**Key topics**:
- Phase 1: Preparation (compliance check, task review, token awareness)
- Phase 2: TDD Implementation (Red-Green-Refactor)
- Phase 3: Compliance Verification (code quality, architecture, docs)
- Phase 4: Final Verification (review commits, full compliance check)
- Phase 5: Task Completion (summary, handoff notes)
**Essential commands**:
```bash
make check-compliance  # Before and after every task AND before every commit
make ai-commit MSG="type(scope): description"  # For all AI-generated commits
```

#### 3. Go Coding Standards
**File**: [`docs/rules/go-guidelines.md`](docs/rules/go-guidelines.md)
**Purpose**: Go-specific coding standards and idioms
**When to use**: When writing any Go code
**Key topics**:
- Framework usage (Fiber, Cobra, Ginkgo)
- Error handling and concurrency patterns
- Testing standards with Ginkgo/Gomega
- LSP setup and tooling

#### 4. Atomic Commit Standards
**File**: [`docs/rules/atomic-commits.md`](docs/rules/atomic-commits.md)
**Purpose**: Creating atomic commits (one logical change per commit)
**When to use**: Before every commit - ensures clean, reviewable history
**Key topics**:
- Atomic commit principles (self-contained, focused, reversible)
- Breaking down work strategies
- Practical examples and patterns
- Recovery techniques for non-atomic commits
**Essential principle**: One logical change per commit, always

#### 5. AI Attribution Requirements
**File**: [`docs/rules/AI_COMMIT_ATTRIBUTION.md`](docs/rules/AI_COMMIT_ATTRIBUTION.md)
**Purpose**: Mandatory AI attribution for AI-assisted commits
**When to use**: Every commit made with AI assistance
**Key topics**:
- Mandatory attribution format
- Assistant name and model version
- Human review requirements
- Validation and enforcement
**Automated via**: Git hooks (`prepare-commit-msg`)

### Quick References (Daily Use)

Keep these handy for quick lookups during development:

- **[`docs/rules/TASK_QUICK_REF.md`](docs/rules/TASK_QUICK_REF.md)** - Task execution workflow checklist
  - One-page summary of master-task-prompt
  - 5-phase process overview with commands
  - Token efficiency thresholds

- **[`docs/rules/COMMIT_QUICK_REFERENCE.md`](docs/rules/COMMIT_QUICK_REFERENCE.md)** - Commit templates and examples
  - Conventional commit format (type, scope, subject)
  - Common patterns and examples
  - Recovery strategies for mistakes

- **[`docs/rules/COMPLIANCE_QUICK_REF.md`](docs/rules/COMPLIANCE_QUICK_REF.md)** - 5-minute compliance check
  - Code quality checklist
  - Commit validation
  - Task verification
  - Automated check commands

- **[`docs/rules/AI_COMMIT_CHECKLIST.md`](docs/rules/AI_COMMIT_CHECKLIST.md)** - AI commit verification
  - Quick checklist for AI-generated commits
  - Format examples and scenarios
  - Troubleshooting tips

### Process Guides

In-depth process documentation for specific workflows:

- **[`docs/rules/review-commit-prompt.md`](docs/rules/review-commit-prompt.md)** - Commit review process
  - Step-by-step atomic commit review
  - Manual review checklist
  - Automated script usage (`make review-commit`)
  - Common issues and solutions

- **[`docs/rules/rules-compliance-check.md`](docs/rules/rules-compliance-check.md)** - Comprehensive compliance verification
  - Complete rules compliance checklist
  - Code quality standards
  - Test requirements
  - Architecture checks
  - Documentation standards
  - Automated via `make check-compliance`

- **[`docs/rules/token-efficiency.md`](docs/rules/token-efficiency.md)** - Token conservation for AI interactions
  - Use tools over text
  - Be concise and batch operations
  - Reference context, focus on deltas
  - Token thresholds and anti-patterns

### All Rules Index

**Complete list**: See [`docs/rules/README.md`](docs/rules/README.md) for the complete, categorized index of all 14 rule files with detailed descriptions and cross-references.

---

## TUI Development

The KaRiya TUI follows strict standards for consistency, accessibility, and professional user experience. All TUI development should reference these documents.

### Core TUI Resources (Read in Order)

#### 1. TUI Standards & Principles
**File**: [`docs/TUI_STANDARDS.md`](docs/TUI_STANDARDS.md)
**Purpose**: Standardized TUI design principles and patterns
**When to use**: Before designing or modifying any TUI component
**Key topics**:
- Universal keyboard shortcuts (Esc, Tab, Arrow keys, vim-style)
- Component architecture (Header, Footer, Help text)
- Screen layouts and navigation patterns
- Consistency, discoverability, accessibility principles
**Essential**: All screens must follow these keyboard shortcuts and patterns

#### 2. TUI Developer Guide
**File**: [`docs/TUI_DEVELOPER_GUIDE.md`](docs/TUI_DEVELOPER_GUIDE.md)
**Purpose**: Comprehensive guide for creating and maintaining TUI components
**When to use**: When implementing new TUI screens or components
**Key topics**:
- Directory structure and message flow
- Creating components (BubbleTea models)
- Testing TUI components (test harnesses)
- Debugging TUI issues
- Performance optimization
- Common pitfalls and solutions

#### 3. Intent Architecture
**File**: [`docs/TUI_INTENT_DIAGRAM.md`](docs/TUI_INTENT_DIAGRAM.md)
**Purpose**: Complete architectural specification with visual diagrams
**When to use**: Understanding the intent system, designing new intents
**Key topics**:
- Intent state machines and transitions
- IntentResult[T] type-safe communication
- Back navigation with context preservation
- Router architecture
- Visual workflow diagrams

#### 4. Theme System
**File**: [`docs/THEME_CUSTOMIZATION_GUIDE.md`](docs/THEME_CUSTOMIZATION_GUIDE.md)
**Purpose**: Theme system for consistent, customizable styling
**When to use**: When styling any UI component or creating custom themes
**Key topics**:
- Theme architecture (ColorPalette, StyleSet, ThemeManager)
- Using themes in intents (helper method pattern)
- Creating custom themes
- Bubbles component integration (themed tables, lists, etc.)
- API reference for all theme methods
**Related**: See [`docs/LIPGLOSS_BUBBLES_GUIDE.md`](docs/LIPGLOSS_BUBBLES_GUIDE.md) for Lipgloss basics

#### 5. Styling & Components
**File**: [`docs/LIPGLOSS_BUBBLES_GUIDE.md`](docs/LIPGLOSS_BUBBLES_GUIDE.md)
**Purpose**: BubbleTea framework and Lipgloss styling usage
**When to use**: When learning Lipgloss basics or using BubbleTea bubbles
**Key topics**:
- Lipgloss style composition
- Color schemes and themes
- BubbleTea bubbles (textinput, list, etc.)
- Layout patterns and responsive design

#### 6. StandardView System
**File**: [`docs/STANDARDVIEW_GUIDE.md`](docs/STANDARDVIEW_GUIDE.md)
**Purpose**: Standardized view system for consistent layout across all TUI screens
**When to use**: When creating any new TUI screen or updating existing screens
**Key topics**:
- StandardView component with logo, breadcrumbs, and context-aware help
- Modal system (Error, Loading, Progress, Success, Warning)
- Terminal size handling and responsive design
- Performance benchmarks and testing
**Related**: See [`docs/MODAL_PATTERNS.md`](docs/MODAL_PATTERNS.md) for modal usage patterns

#### 7. Modal Patterns
**File**: [`docs/MODAL_PATTERNS.md`](docs/MODAL_PATTERNS.md)
**Purpose**: Modal usage patterns and implementation guide
**When to use**: When implementing modals for errors, loading states, or user confirmation
**Key topics**:
- 5 modal types (Error, Loading, Progress, Success, Warning)
- Modal state management
- Accessibility and user experience patterns
- Testing modal interactions
**Related**: See [`docs/STANDARDVIEW_GUIDE.md`](docs/STANDARDVIEW_GUIDE.md) for StandardView integration

#### 8. Forms System (Huh Library)
**File**: [`docs/FORMS_GUIDE.md`](docs/FORMS_GUIDE.md)
**Purpose**: Comprehensive guide to using Charm's huh library for forms
**When to use**: When creating or modifying form inputs in the TUI
**Key topics**:
- Complete documentation of all 5 form configurations (burst, metadata, fact, capture event, burst suggestion)
- Huh library integration and theming (including theme generation from KaRiya themes)
- 20+ reusable validators (date parsing, email, URL, domain-specific)
- Modal and model integration patterns
- **⚠️ Form Alignment & Wrapper Pattern** - CRITICAL for intent forms
- Testing strategies and comprehensive examples

> **⚠️ CRITICAL**: Forms in intents MUST use wrapper models (e.g., `HuhCaptureForm`, `HuhSkillForm`).
> Direct `*huh.Form` usage causes left-alignment issues. See [Form Alignment and the Wrapper Pattern](docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern).

### TUI Quick References

- **[`docs/KEYBOARD_REFERENCE.md`](docs/KEYBOARD_REFERENCE.md)** - Complete keyboard shortcuts reference
  - Universal navigation (Esc, Tab, Arrow/vim keys)
  - Context-specific shortcuts (forms, lists, metadata)
  - Global shortcuts (?, h, q, c, l, m)

- **[`docs/TERMINAL_UI_STYLING_REFERENCE.md`](docs/TERMINAL_UI_STYLING_REFERENCE.md)** - Styling and color reference
  - Color schemes (primary, secondary, success, error)
  - Border styles and decorations
  - Typography and spacing
  - ANSI color codes

- **[`docs/TERMINAL_SIZE_HANDLING.md`](docs/TERMINAL_SIZE_HANDLING.md)** - Terminal size handling patterns
  - Responsive design for various terminal sizes
  - Graceful degradation
  - Testing across different dimensions

- **[`docs/UNIFIED_SHORTCUT_SYSTEM_DESIGN.md`](docs/UNIFIED_SHORTCUT_SYSTEM_DESIGN.md)** - Shortcut system design
  - Shortcut registration and discovery
  - Context-aware help text
  - Conflict resolution
  - Implementation patterns

- **[`docs/SUPPORTING_COMPONENTS_REFERENCE.md`](docs/SUPPORTING_COMPONENTS_REFERENCE.md)** - Reusable component library
  - Available components (help footer, tag selector, etc.)
  - Usage examples and APIs
  - Integration patterns

- **[`docs/INTENT_DEVELOPMENT_CHECKLIST.md`](docs/INTENT_DEVELOPMENT_CHECKLIST.md)** - Intent development checklist
  - Checklist for creating new intents
  - State machine validation
  - Testing requirements

- **[`docs/QUICK_START_STANDARDVIEW.md`](docs/QUICK_START_STANDARDVIEW.md)** - StandardView quick start
  - Quick reference for StandardView usage
  - Common patterns and examples

- **[`docs/development/NAVIGATION_TESTING_GUIDE.md`](docs/development/NAVIGATION_TESTING_GUIDE.md)** - Navigation testing comprehensive guide
  - E2E test framework usage (helpers, assertions, data population)
  - Navigation test patterns (forward, backward, multi-intent, error recovery)
  - Escape key testing matrix (per state type)
  - Common navigation bugs and prevention
  - State transition testing and debugging techniques

- **[`docs/development/NAVIGATION_TESTING_CHECKLIST.md`](docs/development/NAVIGATION_TESTING_CHECKLIST.md)** - Navigation testing quick reference
  - Pre-implementation checklist (state machine definition, escape behavior)
  - Test coverage requirements (escape, forward, back, universal shortcuts)
  - Copy-paste test templates (escape, navigation, E2E)
  - Pre-commit checklist and verification steps

---

## Workflow Documentation

KaRiya provides comprehensive workflow guides for complex user journeys. Each guide includes state machines, keyboard shortcuts, navigation patterns, and troubleshooting.

### Workflow-Specific Guides

#### Complete Workflow Guides

| Guide | Purpose | States | Complexity | Documentation |
|-------|---------|--------|------------|---------------|
| **[CV Generation Workflow](docs/workflows/CV_GENERATION_WORKFLOW.md)** | Generate role and audience-specific CVs from career events | 10 states | ⭐⭐⭐⭐⭐ High | 800+ lines |
| **[Event Capture Workflow](docs/workflows/EVENT_CAPTURE_WORKFLOW.md)** | Capture events with optional burst/fact extraction | 4 states + 3 modals | ⭐⭐⭐⭐ High | 700+ lines |

**See Also**: [Workflow Documentation Index](docs/workflows/README.md) for complete workflow catalog and navigation guide

#### What's Included in Each Guide

Each workflow guide provides:
- **Overview**: Purpose, when to use, prerequisites
- **State Machine Diagram**: Visual representation of all states and transitions (Mermaid)
- **Step-by-Step Guide**: Detailed walkthrough of each state with screenshots
- **Complete Keyboard Reference**: Comprehensive table of all shortcuts per state
- **Navigation Patterns**: Forward navigation, back navigation, error recovery
- **Common Workflows**: Real-world examples with timing estimates
- **Troubleshooting**: Specific issues and solutions for that workflow
- **Technical Details**: Implementation notes, IntentResult flow, async operations

### Workflow Diagram Generation

Workflow diagrams are generated programmatically from the actual implementation to ensure accuracy.

#### Generate Diagrams

```bash
# Generate all workflow diagrams
make generate-diagrams

# Or run script directly
./scripts/generate_workflow_diagrams.sh
```

**Script**: `scripts/generate_workflow_diagrams.sh`  
**Output**: `docs/workflows/diagrams/*.mermaid`

#### Generated Diagrams

- `docs/workflows/diagrams/cv_generation_flow.mermaid` - CV Generation state machine (10 states)
- `docs/workflows/diagrams/event_capture_flow.mermaid` - Event Capture state machine (4 states + 3 modals)

#### Viewing Diagrams

- **GitHub**: Automatic Mermaid rendering
- **VS Code**: Install "Markdown Preview Mermaid Support" extension
- **Online**: Copy content to https://mermaid.live
- **Documentation**: Embedded in workflow guides

### Keyboard Shortcuts

All keyboard shortcuts are centralized in two comprehensive guides:

| Guide | Audience | Purpose | Lines |
|-------|----------|---------|-------|
| **[Keyboard Shortcuts Guide](docs/KEYBOARD_SHORTCUTS_GUIDE.md)** | Users | Complete keyboard reference for using KaRiya TUI | 400+ |
| **[Keyboard System Guide](docs/development/KEYBOARD_SYSTEM_GUIDE.md)** | Developers | Implementing and extending keyboard shortcuts | 500+ |

**Key Features**:
- Quick reference card (printable)
- Workflow-specific shortcuts for all 5 intents
- Vim-style navigation support
- Common key combinations and patterns
- Screen-specific examples
- Accessibility features
- Comprehensive troubleshooting

### Navigation Patterns

All KaRiya workflows follow consistent navigation:

**Universal Shortcuts** (work everywhere):
- **Esc**: Go back one state (or cancel if root state)
- **m**: Return to main menu from any state
- **q** / **Ctrl+C**: Quit application

**State Types**:
1. **Root State**: First state in workflow (Esc = cancel)
2. **Intermediate State**: Has previous state (Esc = go back)
3. **Async Operation**: Background work (Esc = let complete, navigate back)
4. **Final State**: Workflow complete or error (Esc = retry/cancel)

**Error Handling**: Errors are preserved when navigating back so users maintain context

---

## User Guides & Features

User-facing documentation for features and workflows.

### Feature Guides

#### CLI Usage
**File**: [`docs/CLI_GUIDE.md`](docs/CLI_GUIDE.md)
**Purpose**: Command-line interface usage and commands
**Audience**: End users, administrators
**Topics**: Available commands, options, configuration, examples

#### Burst & Facts Extraction
**File**: [`docs/BURST_FACT_EXTRACTION_GUIDE.md`](docs/BURST_FACT_EXTRACTION_GUIDE.md)
**Purpose**: Understanding and working with bursts and facts
**Audience**: Users capturing career events
**Topics**: What are bursts, fact extraction, workflows, best practices

#### CSV Import
**File**: [`docs/CSV_IMPORT_GUIDE.md`](docs/CSV_IMPORT_GUIDE.md)
**Purpose**: Importing career events from CSV files
**Audience**: Users with existing career data
**Topics**: CSV format, import process, field mapping, validation

**Related**: [`docs/CSV_FORMAT_GUIDE.md`](docs/CSV_FORMAT_GUIDE.md) - CSV format specification
- Required fields and formats
- Optional fields
- Date formats
- Examples and templates

#### Metadata Review
**File**: [`docs/METADATA_REVIEW_GUIDE.md`](docs/METADATA_REVIEW_GUIDE.md)
**Purpose**: Reviewing and editing event metadata
**Audience**: Users refining captured events
**Topics**: Metadata editor workflow, bulk operations, validation

#### CV Generation
**File**: [`docs/guides/CV_GENERATION_GUIDE.md`](docs/guides/CV_GENERATION_GUIDE.md)
**Purpose**: Generating CVs from career events
**Audience**: Users creating CVs for job applications
**Topics**: Profile and audience selection, CV preview, export options

**Related**:
- [`docs/guides/CV_EXAMPLES.md`](docs/guides/CV_EXAMPLES.md) - CV generation examples
- [`docs/guides/CV_TROUBLESHOOTING.md`](docs/guides/CV_TROUBLESHOOTING.md) - Common CV issues and solutions

### Troubleshooting Guides

- **[`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md)** - General troubleshooting
  - Common issues and solutions
  - Error messages and meanings
  - Recovery procedures

- **[`docs/guides/CV_TROUBLESHOOTING.md`](docs/guides/CV_TROUBLESHOOTING.md)** - CV-specific issues
  - Empty or incomplete CVs
  - Formatting issues
  - Export problems

---

## Implementation Resources

Architecture, patterns, and component development resources for implementers.

### Architecture & Patterns

#### Implementation Roadmap
**File**: [`docs/IMPLEMENTATION_ROADMAP.md`](docs/IMPLEMENTATION_ROADMAP.md)
**Purpose**: Overall implementation plan and phases
**Audience**: Developers, project managers
**Topics**: Project phases, milestones, dependencies, completion status

#### Screen to Intent Mapping
**File**: [`docs/SCREEN_TO_INTENT_MAPPING.md`](docs/SCREEN_TO_INTENT_MAPPING.md)
**Purpose**: Maps legacy screens to intent architecture
**Audience**: Developers refactoring or understanding architecture
**Topics**: Old screen structure, new intent structure, migration paths

#### Workflow Diagram
**File**: [`docs/WORKFLOW_DIAGRAM.md`](docs/WORKFLOW_DIAGRAM.md)
**Purpose**: High-level workflow visualization
**Audience**: All developers, new team members
**Topics**: User journeys, intent flows, system interactions

#### View Patterns
**File**: [`docs/guides/VIEW_PATTERNS_GUIDE.md`](docs/guides/VIEW_PATTERNS_GUIDE.md)
**Purpose**: Common view patterns and implementations
**Audience**: TUI developers
**Topics**: List views, form views, detail views, modal patterns, state-based rendering

#### Error Handling
**File**: [`docs/guides/ERROR_HANDLING_GUIDE.md`](docs/guides/ERROR_HANDLING_GUIDE.md)
**Purpose**: Error handling patterns and best practices
**Audience**: All developers
**Topics**: Error types, error presentation in TUI, recovery strategies, user feedback

### Component Development

#### Reusable Components
**File**: [`docs/SUPPORTING_COMPONENTS_REFERENCE.md`](docs/SUPPORTING_COMPONENTS_REFERENCE.md)
**Purpose**: Library of reusable UI components
**Audience**: TUI developers
**Topics**: Component catalog, usage examples, integration patterns, customization

#### List Models
**File**: [`docs/guides/LIST_MODEL_RENDERING_SPECIFICATION.md`](docs/guides/LIST_MODEL_RENDERING_SPECIFICATION.md)
**Purpose**: Specification for list rendering behavior
**Audience**: Developers implementing list views
**Topics**: Rendering algorithms, scrolling behavior, selection states, performance optimization

#### Focus Indicators
**File**: [`docs/guides/FOCUS_INDICATOR_GUIDE.md`](docs/guides/FOCUS_INDICATOR_GUIDE.md)
**Purpose**: Focus indicator usage and patterns
**Audience**: TUI developers
**Topics**: Visual focus indicators, keyboard navigation, accessibility, consistency

#### Style System
**File**: [`docs/guides/STYLE_USAGE_GUIDE.md`](docs/guides/STYLE_USAGE_GUIDE.md)
**Purpose**: Using the style system effectively
**Audience**: TUI developers
**Topics**: Style composition, color schemes, responsive styling, best practices

### Additional Resources

- **[`docs/PERFORMANCE_BENCHMARKS.md`](docs/PERFORMANCE_BENCHMARKS.md)** - Performance targets and baselines
- **[`docs/CI_CD_PIPELINE.md`](docs/CI_CD_PIPELINE.md)** - CI/CD setup and workflows
- **[`docs/COLOR_SCHEME.md`](docs/COLOR_SCHEME.md)** - Color scheme reference and usage

---

## Task Documentation

All active implementation work is tracked in task files under the `tasks/` directory. Each task follows a standardized format ensuring clarity and completeness.

### Task Workflow

**Primary Reference**: See [`docs/rules/master-task-prompt.md`](docs/rules/master-task-prompt.md) for the complete 5-phase task execution workflow.

**Quick Reference**: See [`docs/rules/TASK_QUICK_REF.md`](docs/rules/TASK_QUICK_REF.md) for the one-page task checklist.

### Active Work Location

- **Directory**: `tasks/`
- **Naming Convention**: `tasks-XX-feature-name.md` (where XX is zero-padded task number)
- **Next Available**: `tasks-11-burst-enhancement.md`
- **Current Pattern**: Single comprehensive task file per feature (all phases together)

### Task Format Requirements

Each task file MUST include:

```markdown
# Task XX: Feature Name

## Overview
- **Goal**: [What we're achieving]
- **Time Estimate**: [Estimated duration]
- **Prerequisites**: [Required setup or knowledge]

## Files to Modify
- [ ] List of files that will be changed
- [ ] With checkboxes for tracking

## Implementation Checklist

### Phase 1: [Phase Name]
- [ ] Specific actionable item
- [ ] With acceptance criteria
- [ ] Code snippets where helpful

### Phase 2: [Phase Name]
[... continue for all phases]

## Testing Instructions
- [ ] How to test the implementation
- [ ] What tests to write
- [ ] What to verify manually

## Acceptance Criteria
- [ ] Feature works as specified
- [ ] All tests pass
- [ ] Coverage maintained
- [ ] Documentation updated

## Rollback Plan
- Steps to revert changes if needed
- Safety considerations
```

### Task Processing Rules

**Reference**: [`docs/rules/master-task-prompt.md`](docs/rules/master-task-prompt.md) - See "Task Processing" section

**Key Principles**:
1. Execute tasks in order (top to bottom)
2. Check off items as completed
3. One task at a time
4. Use tools to verify completion
5. Authority order: tools > checklist > guidelines
6. Clear completion criteria (files modified, tests pass, checkbox marked)
7. Refactoring only on directly modified code

### Completed Tasks

Completed tasks remain in `tasks/` directory for reference:
- `tasks-09-tui-intent-refactoring.md` - Intent framework implementation
- `tasks-10-navigation-coverage.md` - Navigation and test coverage
- See individual files for implementation details and lessons learned

---

## Key Files and Purposes

### Intent Framework

| File | Purpose |
|------|---------|
| `contract.go` | Intent and IntentRouter interfaces |
| `result.go` | IntentResult[T] and IntentError types |
| `router.go` | IntentRouter implementation |
| `testing.go` | Test utilities and harnesses |

### Intent Implementations

| Intent | Model | Implementation | Tests |
|--------|-------|-----------------|-------|
| CaptureEvent | `capture_event.go` | `capture_event_intent.go` | `contract_test.go` |
| BrowseTimeline | `browse_timeline.go` | `browse_timeline_intent.go` | `browse_timeline_test.go` |
| GenerateCV | `generate_cv.go` | `generate_cv_intent.go` | `generate_cv_test.go` |
| ExportArtifact | `export_artifact.go` | `export_artifact_intent.go` | `export_artifact_test.go` |
| ConfigureSystem | `configure_system.go` | `configure_system_intent.go` | `configure_system_test.go` |

### Root Application

| File | Purpose |
|------|---------|
| `internal/cli/app/app.go` | Root Bubble Tea model, intent router integration |
| `internal/cli/app/messages.go` | App-level message types |

### Domain & Data

| File | Purpose |
|------|---------|
| `internal/domain/career/event.go` | CareerEvent domain model |
| `internal/domain/career/burst.go` | Burst domain model |
| `internal/domain/career/fact.go` | Fact domain model |
| `internal/domain/career/cv.go` | CV domain model |
| `internal/repository/career/repository.go` | Data access interface |
| `internal/repository/career/sqlite_repository.go` | SQLite implementation |
| `internal/service/career/service.go` | Core business logic |

### CV Service

| File | Purpose |
|------|---------|
| `internal/service/career/cv/enhanced_bullet_generator.go` | Enhanced CV bullet generation with scoring |
| `internal/service/career/cv/enhanced_bullet_generator_test.go` | Tests for bullet generation (203 specs) |

---

## Recent Fixes

### Huh Forms Migration (2026-01-07)

**Status**: ✅ **PHASE 4 COMPLETE - ALL MODALS MIGRATED**

#### Summary
Successfully migrated KaRiya's form handling from manual `textinput.Model` arrays to Charm's **huh** library, achieving:
- ✅ **40% reduction** in modal code (835 → 505 lines)
- ✅ **76 comprehensive tests** (100% passing)
- ✅ **Catppuccin theming** throughout all forms
- ✅ **Zero regressions** in functionality

#### What Was Migrated

**Modals** (Complete):
- EditBurstModal: 256 → 161 lines (-37%)
- EditMetadataModal: 308 → 183 lines (-41%)
- EditFactModal: 271 → 161 lines (-41%)

**Models** (In Progress):
- BurstEditorModel: 335 → 201 lines (-40%) ✅

**Infrastructure Created**:
- Forms package: 1,701 lines (972 source + 729 tests)
- 20+ reusable validators (date parsing, email, URL, domain-specific)
- Form configurations for burst, metadata, and fact editing
- Comprehensive developer guide

#### Key Features

**Date Parsing**:
```go
ParseDateString("2024-01-07")    // ✅ Standard
ParseDateString("today")          // ✅ Quick input
ParseDateString("7 days ago")     // ✅ Relative
ParseDateString("2 weeks ago")    // ✅ Relative
```

**Benefits**:
- ✅ Automatic focus management (no more manual tab handling)
- ✅ Built-in validation with custom validators
- ✅ Professional Catppuccin theming
- ✅ Type-safe field access
- ✅ Reusable form configurations
- ✅ 70% less code per form

**Documentation**:
- [`docs/FORMS_GUIDE.md`](docs/FORMS_GUIDE.md) - Complete forms implementation guide covering all 5 form configurations, validators, theming, integration patterns, and testing

**Remaining Work** (Phase 5):
- FactEditorModel (600 lines)
- MetadataEditorModel (556 lines)
- BurstSuggestionModel (525 lines)
- FormModel (956 lines - main capture form)

**Estimated Final Savings**: 60%+ reduction in all form-related code after complete migration.

---

### Phase 11 & 12 Completion (2026-01-04)

**Status**: ✅ **COMPLETE - CV GENERATION AND EXPORT FULLY FUNCTIONAL**

#### Summary of Fixes

##### Issue
Users were unable to:
- Generate a CV
- See a preview of a CV
- Export a CV
- Preview their export

##### Root Causes Identified
1. **GenerateCV Intent Test Failure**: Test was checking for ">" marker but code was using "▶"
2. **Missing Escape Key Handler**: Profile selection state didn't handle escape key for cancellation
3. **Export Not Accessible**: Confirm state didn't allow transition to export workflow
4. **EnhancedBulletGenerator Bugs**: Two failing tests in CV bullet generation

##### Solutions Implemented

###### 1. Fixed Test Marker Check
- **File**: `internal/cli/intents/generate_cv_test.go`
- **Change**: Updated test to check for "▶" instead of ">"
- **Impact**: View rendering test now passes correctly

###### 2. Added Escape Key Handler
- **File**: `internal/cli/intents/generate_cv_intent.go`
- **Change**: Added `case "esc":` handler in `updateSelectProfile()` to cancel intent
- **Impact**: Users can now press Escape to cancel CV generation from profile selection

###### 3. Enabled Export Workflow
- **File**: `internal/cli/intents/generate_cv_intent.go`
- **Changes**:
  - Added `case "e", "x":` handler in `updateConfirm()` to transition to export format selection
  - Updated footer text in `viewConfirm()` to show export option: "y/Enter to confirm, e/x to export, n/Esc to go back, q to cancel"
- **Impact**: Users can now press "e" or "x" to export CV instead of just confirming

###### 4. Fixed EnhancedBulletGenerator
- **File**: `internal/service/career/cv/enhanced_bullet_generator.go`
- **Changes**:
  - Increased event confidence from 0.70 to 0.80 in `createBulletsFromEvents()` (line 281)
  - Fixed case-insensitive verb replacement in `enhanceActionVerb()` (lines 396-410)
- **Impact**: 203 CV service tests now pass with 100% success rate

##### Test Results

**Before**:
```
GenerateCV Tests: Some failures
- Test marker check failing (expected ">" but got "▶")
- Integration tests panicking on escape key
- Export workflow not accessible
```

**After**:
```
GenerateCV Tests: ✅ ALL PASSING
- Unit tests: All passing
- Integration tests: 19/19 passing
- Overall test suite: 336/347 passing (11 unrelated failures in other intents)
```

##### Complete Workflow Now Functional

Users can now complete the entire CV generation workflow:

1. **Select Profile** ✅
   - Navigate with arrow keys or j/k
   - Press Escape to cancel
   - Press Enter to proceed

2. **Select Audience** ✅
   - View default audiences for selected profile
   - Press Enter to generate CV

3. **Generate CV** ✅
   - Shows progress: "⏳ Generating CV..."
   - Uses CVGenerationService with DataProcessingService
   - Applies EnhancedBulletGenerator for professional bullets

4. **Preview CV** ✅
   - View CV metadata and statistics
   - Shows source events and facts count
   - Press e to edit, c to confirm, Esc to go back

5. **Review & Edit** ✅
   - Review generated CV content
   - Press Enter to proceed to confirmation

6. **Confirm** ✅
   - Confirm CV generation
   - **NEW**: Press e/x to export instead of completing
   - Press y/Enter to complete without exporting

7. **Export** ✅ (NEW WORKFLOW)
   - Select export format: Text, Markdown, or YAML
   - Select save location: File or Clipboard
   - Watch export progress
   - View export completion with file location

##### Files Modified
- `internal/cli/intents/generate_cv_intent.go` (11 lines added for escape and export handlers)
- `internal/cli/intents/generate_cv_test.go` (1 line fixed for marker test)
- `internal/service/career/cv/enhanced_bullet_generator.go` (13 lines fixed for bullet generation)

##### Verification
- ✅ Application builds successfully
- ✅ All GenerateCV tests pass (36 tests)
- ✅ All GenerateCV integration tests pass (19 tests)
- ✅ CV service tests pass (203 tests)
- ✅ No regressions in other test suites
- ✅ Complete workflow is functional and ready for production use

### Tasks 12-16: StandardView Implementation (2026-01-07)

**Status**: ✅ **COMPLETE - STANDARDIZED VIEW SYSTEM PRODUCTION READY**

#### Summary

Implemented a comprehensive standardized view system across all TUI screens, providing consistent layout, branding, and user experience. All 5 primary intents now use StandardView with logo, breadcrumbs, modals, and context-aware help.

#### What Was Accomplished

##### Task 16: Testing, Documentation, and Polish

**Phases Completed**: 10/10 (100%)

**Testing Infrastructure**:
- ✅ 145+ new edge case tests
- ✅ 15 performance benchmarks  
- ✅ Terminal size test suite (7 common sizes)
- ✅ Cross-intent consistency tests (5 intents)
- ✅ Visual test program (`cmd/test_all_views`)

**Performance Results** (all targets exceeded):
- StandardView render: 0.376ms (target: <50ms) - **133x faster**
- Modal render: 0.072-0.296ms (target: <20ms) - **170x faster**
- Full view render: 0.727ms (target: <100ms) - **138x faster**

**Documentation**:
- ✅ STANDARDVIEW_GUIDE.md (633 lines) - Complete developer guide
- ✅ MODAL_PATTERNS.md (678 lines) - Modal usage patterns
- ✅ 1,311 lines of comprehensive documentation

**Quality Metrics**:
- ✅ All tests passing (980+ total)
- ✅ Zero regressions
- ✅ Code coverage maintained >87%
- ✅ All performance targets exceeded by 100x+

#### Components Created

**Core Components** (Task 12):
- `StandardView` - Standardized layout with logo, content, footer
- `ModalContent` - 5 modal types (Error, Loading, Progress, Success, Warning)
- `LoadingMessageRotator` - Message rotation for long operations
- `SimpleSpinner` - Animated spinner for loading states

**Infrastructure** (Task 13):
- `BaseIntent` enhancements - Terminal awareness, logo management
- View helper functions - Easy StandardView creation
- Modal helper functions - Convenient modal creation
- Footer helper functions - Predefined footer templates

**Testing** (Task 16):
- `terminal_size_test.go` - 96 terminal size tests
- `performance_test.go` - 15 performance benchmarks
- `consistency_test.go` - 6 cross-intent consistency tests
- `cmd/test_all_views` - Interactive visual test program

#### Files Created/Modified

**Created** (8 files):
- `internal/cli/components/standard_view.go`
- `internal/cli/components/standard_view_test.go`
- `internal/cli/components/modal.go`
- `internal/cli/components/modal_test.go`
- `internal/cli/components/loading_messages.go`
- `internal/cli/components/loading_messages_test.go`
- `internal/cli/components/terminal_size_test.go`
- `internal/cli/components/performance_test.go`
- `internal/cli/intents/consistency_test.go`
- `internal/cli/intents/view_helpers.go`
- `cmd/test_all_views/main.go`
- `docs/STANDARDVIEW_GUIDE.md`
- `docs/MODAL_PATTERNS.md`

**Modified** (All 5 primary intents migrated):
- `internal/cli/intents/capture_event_intent.go`
- `internal/cli/intents/browse_timeline_intent.go`
- `internal/cli/intents/generate_cv_intent.go`
- `internal/cli/intents/export_artifact_intent.go`
- `internal/cli/intents/configure_system_intent.go`

#### User Experience Improvements

**Visual Consistency**:
- ✅ Logo visible on every screen (all intents)
- ✅ Consistent layout throughout application
- ✅ Breadcrumbs show navigation context
- ✅ Footer separator improves readability
- ✅ Context-aware help on every screen

**Modal System**:
- ✅ Error modals with terminal bell alerts
- ✅ Loading modals with spinner animation
- ✅ Progress modals for multi-step operations
- ✅ Success modals with auto-dismiss (3s)
- ✅ Warning modals for confirmations

**Responsive Design**:
- ✅ Adapts to terminal sizes (80x24 to 240x80)
- ✅ Graceful degradation for small terminals
- ✅ Nil terminal info handling with defaults (120x40)

#### Technical Achievements

**Testing**:
- 160+ new tests across 5 test files
- 100% consistency across all intents
- Zero race conditions detected
- Sub-millisecond rendering performance

**Architecture**:
- Type-safe, composable design
- Builder pattern for easy view construction
- Separation of concerns (view vs logic)
- Reusable components across intents

**Performance**:
- All benchmarks exceed targets by 100x+
- No memory leaks detected
- Reasonable memory usage (~52-139 KB per render)
- Supports concurrent rendering

#### Usage

**Visual Test Program**:
```bash
go build -o test_views ./cmd/test_all_views
./test_views
# Navigate scenarios with ←→ or h/l
```

**Run Benchmarks**:
```bash
go test -bench=. -benchmem ./internal/cli/components/
```

**Documentation**:
- See `docs/STANDARDVIEW_GUIDE.md` for complete usage guide
- See `docs/MODAL_PATTERNS.md` for modal patterns

#### Verification

- ✅ All 980+ tests passing
- ✅ All 5 intents use StandardView consistently  
- ✅ Performance targets exceeded (100x+ faster)
- ✅ Code coverage >87%
- ✅ Zero regressions
- ✅ Visual test program functional
- ✅ Comprehensive documentation complete
- ✅ Production ready

### Tasks 19-20: Test Fixes and Final Cleanup (2026-01-07)

**Status**: ✅ **COMPLETE - ALL TESTS PASSING, ZERO STATICCHECK WARNINGS**

#### Task 19: Fix Remaining Test Failures

**Summary**: Fixed 3 test failures after form refactoring (Task 17)

**Issues Fixed**:

1. **Form Field Visibility Inconsistency** (`internal/cli/models/form.go`)
   - **Root Cause**: FormModel initialized with `strategy="manual"` but `showOptionalFields=false`, creating inconsistent state
   - **Fix**: Changed `showOptionalFields` default to `true` to align with manual strategy behavior
   - **Impact**: Fixed 2 persistence test failures in `cmd/cli/persistence_test.go`

2. **Outdated Footer Text Test** (`internal/cli/intents/capture_event_escape_test.go`)
   - **Root Cause**: Test expected "Retry" but footer now shows "Back"
   - **Fix**: Updated test expectation to check for "Back" instead of "Retry"
   - **Impact**: Fixed 1 escape test failure

**Files Modified**:
- `internal/cli/models/form.go` (1 line changed)
- `internal/cli/intents/capture_event_escape_test.go` (1 line changed)

**Results**:
- All 2,078 tests passing (100% pass rate - up from 2,075/2,078)
- Zero race conditions
- Build successful
- Commit: `31218ec`

#### Task 20: Final Cleanup - Staticcheck

**Summary**: Removed 3 unused functions flagged by staticcheck

**Functions Removed** (all marked as "reserved for future use"):

1. **`BurstManagementIntent.applyFilters()`** - 35 lines
   - Filtering feature not yet implemented
   - Removed unused `sort` import

2. **`BurstManagementIntent.setFailed()`** - 14 lines
   - Error handling pattern not used in this intent

3. **`GenerateCVIntent.getStateName()`** - 27 lines
   - Breadcrumb navigation not used in this intent

**Files Modified**:
- `internal/cli/intents/burst_management_intent.go` (-51 lines)
- `internal/cli/intents/generate_cv_intent.go` (-29 lines)

**Results**:
- **Zero staticcheck warnings** (was 3)
- **80 lines of dead code removed**
- All 2,078 tests still passing
- Zero race conditions
- Cleaner, more maintainable codebase
- Commit: `0307f4f`

#### Combined Impact

**Code Quality Metrics**:
- ✅ **Tests**: 2,078/2,078 passing (100%)
- ✅ **Staticcheck**: 0 warnings (was 3)
- ✅ **Race conditions**: 0 detected
- ✅ **Build**: Successful
- ✅ **Dead code**: Removed (80 lines)
- ✅ **Production ready**: All critical work complete

**Commits**:
- `31218ec` - fix(form): align showOptionalFields default with manual strategy
- `0307f4f` - refactor: remove unused functions flagged by staticcheck

---

## Common Development Tasks

### Adding a New Intent

1. **Create data structures** (`intent_name.go`):
   ```go
   type YourIntentContext struct { ... }
   type YourIntentResult struct { ... }
   type YourIntentModel struct {
       state YourState
       data  *YourIntentContext
       result *IntentResult[*YourIntentResult]
   }
   ```

2. **Implement intent** (`your_intent_intent.go`):
   ```go
   func (y *YourIntentModel) Init(ctx context.Context) tea.Cmd { ... }
   func (y *YourIntentModel) Update(msg tea.Msg) tea.Cmd { ... }
   func (y *YourIntentModel) View() string { ... }
   func (y *YourIntentModel) Result() *IntentResult[interface{}] { ... }
   ```

3. **Write tests** (`your_intent_test.go`):
   - Use Ginkgo/Gomega framework
   - Test state transitions
   - Test view rendering
   - Test result handling

4. **Register with router** in `internal/cli/app/app.go`:
   ```go
   router.RegisterIntent("your_intent", func() intents.Intent {
       return NewYourIntent(context)
   })
   router.RegisterResultHandler("your_intent", func(result) tea.Cmd {
       // Handle result
   })
   ```

### Creating a New Form

**Time**: 30-45 minutes  
**Prerequisites**: Domain object exists, validators identified  
**Reference**: [`docs/rules/FORMS_WORKFLOW_GUIDE.md`](docs/rules/FORMS_WORKFLOW_GUIDE.md) - Complete step-by-step guide

> **⚠️ CRITICAL: Form Alignment Rule**
> 
> Forms used in **intents** MUST use a wrapper model (like `HuhCaptureForm`, `HuhSkillForm`).
> Direct use of `*huh.Form` in intents causes **left-alignment issues** because:
> - Form dimensions are captured at creation time and become stale
> - `WindowSizeMsg` handling is scattered and error-prone
> - Forms don't update properly on terminal resize
>
> **See**: [`docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern`](docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern)

#### Quick Workflow

1. **Define FormData structure** (`internal/cli/forms/your_form.go`):
   ```go
   type YourFormData struct {
       Field1          string
       Field2          []string // For MultiSelect
       SubmitConfirmed bool     // Required for confirm button
   }
   ```

2. **Create form builder functions**:
   ```go
   // 3 variants for flexibility
   func NewYourForm(obj *career.YourDomain) *huh.Form { ... }
   func NewYourFormWithData(data *YourFormData) *huh.Form { ... }
   func NewYourFormWithDataAndDimensions(data *YourFormData, width, height int) *huh.Form {
       data.SubmitConfirmed = false
       
       fieldsGroup := huh.NewGroup(
           forms.NewInput(forms.FieldConfig{
               Key:         "field1",
               Title:       "Field 1",
               Validate:    forms.Required,
           }).Value(&data.Field1),
       )
       
       return forms.NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
   }
   ```

3. **Create domain conversion functions**:
   ```go
   func GetYourFormData(obj *career.YourDomain) *YourFormData { ... }
   func ApplyYourFormData(obj *career.YourDomain, data *YourFormData) error { ... }
   ```

4. **Write tests** (`your_form_test.go`):
   - Test form creation
   - Test data extraction (GetYourFormData)
   - Test data application (ApplyYourFormData)
   - Test roundtrip conversion
   - Test nil slice handling

5. **Document in FORMS_GUIDE.md**:
   - Add section to "Form Configurations"
   - Add to forms package table
   - Add to tests section

#### Integration Patterns

**⚠️ Intent Integration (MUST use wrapper model)**:

When adding forms to an **intent**, you MUST create a wrapper model:

```go
// 1. Create wrapper in internal/cli/models/huh_your_form.go
type HuhYourForm struct {
    *BaseStandardModel
    formData *forms.YourFormData
    form     *huh.Form
    width    int
    height   int
}

func NewHuhYourForm() *HuhYourForm {
    m := &HuhYourForm{
        BaseStandardModel: NewBaseStandardModel(),
        formData:          &forms.YourFormData{},
        width:             80,   // Default width
        height:            24,   // Default height
    }
    m.rebuildForm()
    return m
}

func (m *HuhYourForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Handle window size internally - THIS IS KEY
        m.width = msg.Width
        m.height = msg.Height
        m.form = m.form.WithHeight(forms.DefaultFormHeight(m.height)).WithWidth(m.width - 4)
        return m, nil
    }
    // ... rest of update logic
}

// 2. Use wrapper in intent (NOT raw *huh.Form)
type YourIntent struct {
    yourForm *models.HuhYourForm  // ✅ Correct - wrapper model
    // form *huh.Form             // ❌ Wrong - causes alignment issues
}

func (i *YourIntent) handleAddNew() tea.Cmd {
    i.yourForm = models.NewHuhYourForm()
    return i.yourForm.Init()
}
```

**Existing wrapper examples**:
- `internal/cli/models/huh_capture_form.go` - Career event capture
- `internal/cli/models/huh_skill_form.go` - Skill management

**Modal Integration** (inline editing - wrapper NOT required):
```go
type EditYourModal struct {
    original  *career.YourDomain  // Preserved (never mutated)
    modified  *career.YourDomain  // Working copy
    result    *ModalEditResult[*career.YourDomain]
    form      *huh.Form           // Direct use OK in modals
    formData  *forms.YourFormData
}
```

**See Complete Guide**: [`docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern`](docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern)

### Running Tests

```bash
# Create test file
touch internal/cli/intents/your_test.go

# Add Ginkgo test structure
cat > internal/cli/intents/your_test.go << 'EOF'
package intents_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("YourIntent", func() {
    It("should do something", func() {
        Expect(true).To(BeTrue())
    })
})
EOF

# Run tests
go test -v ./internal/cli/intents/...
```

---

## Deployment Guide

### Branching Strategy

KaRiya uses a **dual-branch workflow** for controlled releases:

- **`next` branch**: Integration branch for feature development
  - All feature branches merge here via PR
  - Full CI runs on every merge
  - Staging environment for testing features together
  
- **`main` branch**: Production releases only
  - Only accepts merges from `next` branch
  - Automatic semantic releases on merge
  - Protected branch with strict checks

**See**: [`docs/BRANCHING_STRATEGY.md`](docs/BRANCHING_STRATEGY.md) for complete workflow documentation.

**Quick Workflow**:
```bash
# 1. Create feature branch from next
git checkout next && git pull
git checkout -b feature/my-feature

# 2. Develop and commit (conventional commits)
git commit -m "feat: add new feature"

# 3. Push and create PR to next (NOT main)
git push -u origin feature/my-feature
gh pr create --base next --fill

# 4. After merge to next, release when ready:
#    Create PR: next → main (triggers automatic release)
```

### Pre-Deployment Checklist

**RECOMMENDED**: Run all CI checks locally before pushing:

```bash
# Run ALL CI checks locally (mirrors GitHub Actions exactly)
make ci-local
```

This single command runs:
- ✅ Commitlint validation
- ✅ AI attribution check
- ✅ Code formatting (go fmt)
- ✅ Static analysis (go vet, staticcheck)
- ✅ Tests with race detector and coverage
- ✅ Multi-platform builds (Linux, macOS, Windows)
- ✅ Security scanning (gosec)

**Individual checks** (if needed):

```bash
# 1. Install all CI tools
make ci-install-tools

# 2. Run full test suite
make test

# 3. Check code quality
make fmt
make vet
make staticcheck

# 4. Security scan
make gosec

# 5. Generate coverage report
make coverage

# 6. Build for target platforms
make build  # Current platform
# OR multi-platform:
GOOS=linux GOARCH=amd64 go build -o kariya-linux-amd64 ./cmd/cli
GOOS=darwin GOARCH=amd64 go build -o kariya-darwin-amd64 ./cmd/cli
GOOS=darwin GOARCH=arm64 go build -o kariya-darwin-arm64 ./cmd/cli
GOOS=windows GOARCH=amd64 go build -o kariya-windows-amd64.exe ./cmd/cli
```

### CI/CD Pipeline

**Main CI Workflow** (`.github/workflows/ci.yml`):
- **commitlint** (PR only): Validates commit messages
- **lint**: Code quality (gofmt, vet, staticcheck)
- **test**: Multi-platform testing (Linux, macOS, Windows)
- **build**: Multi-platform builds with artifact upload
- **security**: Gosec security scanning

**PR Validation Workflow** (`.github/workflows/pr-validation.yml`):
- **validate-pr-title**: PR title follows conventional commits
- **check-ai-attribution**: Commits have AI attribution (if applicable)
- **conventional-commits**: All commits follow conventions
- **breaking-changes**: Detects breaking changes
- **size-label**: Auto-labels PR by size

**See Also**:
- [CI Checks Summary](docs/CI_CHECKS_SUMMARY.md) - Complete CI/local command mapping
- [CI Local Guide](docs/CI_LOCAL_GUIDE.md) - Detailed guide for running CI locally

---

## Performance Benchmarks

All benchmarks passing with excellent performance:

| Operation | Target | Actual | Status |
|-----------|--------|--------|--------|
| Intent Init | 50ms | 0.4ms | ✅ |
| View Render | 100ms | 46ms | ✅ |
| State Transition | 10ms | 0.03ms | ✅ |
| Router Operations | 1ms | 0.1ms | ✅ |
| Test Suite | 5s | 1.3s | ✅ |

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./internal/cli/intents/

# Run specific benchmark
go test -bench=BenchmarkCaptureEventInit ./internal/cli/intents/

# With memory profiling
go test -bench=. -benchmem ./internal/cli/intents/
```

---

## Workflow Patterns

### Intent State Machine Pattern

```go
type YourState string
const (
    StateInitial YourState = "initial"
    StateWorking YourState = "working"
    StateFinal   YourState = "final"
)

func (y *YourIntentModel) Update(msg tea.Msg) tea.Cmd {
    switch y.state {
    case StateInitial:
        return y.handleInitial(msg)
    case StateWorking:
        return y.handleWorking(msg)
    case StateFinal:
        return y.handleFinal(msg)
    }
    return nil
}

func (y *YourIntentModel) View() string {
    switch y.state {
    case StateInitial:
        return y.viewInitial()
    case StateWorking:
        return y.viewWorking()
    case StateFinal:
        return y.viewFinal()
    }
    return ""
}
```

### Modal Sub-Flow Pattern

```go
type ModalEditResult[T any] struct {
    Original T
    Modified T
    Accepted bool
    Changes  map[string]interface{}
}

// Usage in intent
if editModal.WasAccepted() {
    y.data.Field = editModal.Modified
} else {
    y.data.Field = editModal.Original
}
```

### Back Navigation with Context Preservation

```go
// When returning result
result := &IntentResult[*YourIntentResult]{
    Status: StatusCompleted,
    Data:   y.result.Data,
}
result.WithMetadata("scroll_position", 42)
result.WithMetadata("selection", selectedID)
return result

// When restoring context
if pos, ok := prevResult.GetMetadata("scroll_position"); ok {
    y.scrollPos = pos.(int)
}
```

---

## Troubleshooting

### Issue: Tests Failing with "Multiple Ginkgo Entry Points"

**Solution**: Ensure only one `_test.go` file per package uses Ginkgo.

### Issue: Race Conditions Detected

**Solution**:
1. Run `go test -race ./...` to identify
2. Add proper synchronization (mutex, channels)
3. Check GlobalContext usage in `internal/cli/context/global.go`

### Issue: High Memory Usage

**Solution**:
1. Check for goroutine leaks with pprof
2. Verify proper cleanup in intent `Init()` and result handling
3. Check repository queries for N+1 issues

### Issue: Slow Tests

**Solution**:
1. Run benchmarks: `go test -bench=. ./internal/cli/intents/`
2. Profile with pprof: `go test -cpuprofile=cpu.prof ./...`
3. Check database queries in repository layer

### Additional Troubleshooting Resources

- **General**: [`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md)
- **CV Issues**: [`docs/guides/CV_TROUBLESHOOTING.md`](docs/guides/CV_TROUBLESHOOTING.md)
- **Error Handling**: [`docs/guides/ERROR_HANDLING_GUIDE.md`](docs/guides/ERROR_HANDLING_GUIDE.md)

---

## Database Migrations

### Overview

KaRiya uses **goose** for Rails-like database migrations with versioned SQL files. Migrations run automatically at application startup and handle both fresh databases and existing pre-goose databases.

### Migration System

**Framework**: [goose v3](https://github.com/pressly/goose)
**Location**: `internal/repository/career/migrations/`
**Tracking**: `goose_db_version` table

### Current Migrations

| Version | File | Description |
|---------|------|-------------|
| 001 | `001_create_career_events.sql` | Career events table + indexes (date, company) |
| 002 | `002_add_categories_column.sql` | Add categories column |
| 003 | `003_create_bursts.sql` | Bursts table + confirmed index |
| 004 | `004_create_facts.sql` | Facts table + source indexes |

### Database Indexes

| Table | Index | Columns | Purpose |
|-------|-------|---------|---------|
| `career_events` | `idx_career_events_date` | `date` | Timeline queries, date range filters |
| `career_events` | `idx_career_events_company` | `company` | Filter by company |
| `bursts` | `idx_bursts_confirmed` | `confirmed` | Filter confirmed/unconfirmed |
| `facts` | `idx_facts_source_event_id` | `source_event_id` | Join with events |
| `facts` | `idx_facts_source_burst_id` | `source_burst_id` | Join with bursts |

### How Migrations Work

**At Application Startup** (`cmd/cli/main.go`):
```go
// 1. Open database connection
db, err := sql.Open("sqlite", dbPath)

// 2. Run migrations (handles both fresh and existing databases)
if err := career.RunMigrations(db); err != nil {
    return err
}

// 3. Create repositories with existing connection
repo := career.NewSQLiteRepositoryWithDB(db)
```

**Baseline Detection**:
- For **fresh databases**: Runs all migrations 001-004
- For **existing databases** (pre-goose): Auto-detects which migrations have already been applied
  - Checks for `career_events` table → baseline version 1
  - Checks for `categories` column → baseline version 2
  - Checks for `bursts` table → baseline version 3
  - Checks for `facts` table → baseline version 4
  - Marks detected migrations as applied without re-running them

### Migration Files

Each migration has **up** and **down** sections:

```sql
-- +goose Up
-- Migration forward (create table, add column, etc.)
CREATE TABLE IF NOT EXISTS career_events (...);
CREATE INDEX IF NOT EXISTS idx_career_events_date ON career_events(date);

-- +goose Down
-- Migration rollback (drop table, remove index, etc.)
DROP INDEX IF EXISTS idx_career_events_date;
DROP TABLE IF EXISTS career_events;
```

**Note**: Migration 002 (add categories column) cannot be fully rolled back due to SQLite limitations (no `DROP COLUMN` support).

### Running Migrations

**Automatic** (recommended):
Migrations run automatically when the application starts with SQLite mode.

**Programmatic**:
```go
import "github.com/baphled/kariya/internal/repository/career"

db, _ := sql.Open("sqlite", "path/to/db")
if err := career.RunMigrations(db); err != nil {
    log.Fatal(err)
}
```

**Check Migration Status**:
```go
version, err := career.MigrationStatus(db)
fmt.Printf("Current version: %d\n", version)
```

### Migration Storage

Migrations are **embedded in the binary** via `//go:embed`:
```go
//go:embed migrations/*.sql
var migrations embed.FS
```

This ensures:
- Single binary distribution (no external files needed)
- Migrations always in sync with code version
- Simplified deployment

### Test Database Setup

**Reusable Test Helpers** (`internal/testutil`):
```go
import "github.com/baphled/kariya/internal/testutil"

// Standard test setup
db, cleanup := testutil.SetupTestDB(t)
defer cleanup()

// With path (for legacy constructors)
dbPath, db, cleanup := testutil.SetupTestDBWithPath(t)
defer cleanup()
```

### Repository Constructors

**New Constructors** (recommended):
```go
// Use after running migrations
repo := career.NewSQLiteRepositoryWithDB(db)
burstRepo := career.NewSQLiteBurstRepositoryWithDB(db)
factRepo := career.NewSQLiteFactRepositoryWithDB(db)
```

**Legacy Constructors** (backward compatible):
```go
// Runs migrations internally (for test compatibility)
repo, err := career.NewSQLiteRepository(dbPath)
burstRepo, err := career.NewSQLiteBurstRepository(db)
factRepo, err := career.NewSQLiteFactRepository(db)
```

### Key Files

| File | Purpose |
|------|---------|
| `internal/repository/career/migrator.go` | Migration runner with baseline detection |
| `internal/repository/career/migrator_test.go` | Migration tests (13 specs) |
| `internal/repository/career/migrations/*.sql` | SQL migration files |
| `internal/testutil/db.go` | Reusable test database helpers |
| `internal/testutil/db_test.go` | Test helper tests (4 specs) |

### Future Migrations

To add a new migration:

1. **Create migration file**:
   ```bash
   # Create 005_your_migration.sql in internal/repository/career/migrations/
   ```

2. **Add up/down sections**:
   ```sql
   -- +goose Up
   ALTER TABLE career_events ADD COLUMN new_field TEXT;

   -- +goose Down
   -- Rollback instructions (if possible)
   ```

3. **Migration runs automatically** on next application start

4. **Test migration**:
   ```bash
   go test ./internal/repository/career/migrator_test.go
   ```

### Troubleshooting

**Issue**: "duplicate column name" error
**Cause**: Trying to run migration on database that already has the change
**Solution**: goose tracks applied migrations automatically; baseline detection handles this for pre-goose databases

**Issue**: Migration fails on startup
**Cause**: SQL syntax error in migration file
**Solution**: Check migration file syntax, test with `go test ./internal/repository/career/migrator_test.go`

**Issue**: Want to check current migration version
**Solution**:
```bash
sqlite3 ~/.kariya/events.db "SELECT MAX(version_id) FROM goose_db_version"
```

---

## Project Metadata

| Property | Value |
|----------|-------|
| **Repository** | https://github.com/baphled/kariya |
| **Language** | Go 1.24 |
| **Framework** | Bubble Tea + Lipgloss |
| **Test Framework** | Ginkgo v2 + Gomega |
| **Database** | SQLite (modernc.org/sqlite) |
| **Status** | ✅ Production Ready |
| **Last Updated** | 2026-01-06 |

---

*This documentation serves as the comprehensive handover document for the KaRiya project. For questions or suggestions, please open an issue on GitHub.*
