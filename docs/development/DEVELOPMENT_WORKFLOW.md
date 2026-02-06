# Development Workflow

Guidelines and rules for KaRiya development.

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

## UIKit Component Standards

When writing TUI code, use the UIKit component library for consistency:

### Required UIKit Components

| Need | Use | Location |
|------|-----|----------|
| View layout | `layout.ScreenLayout` | `uikit/layout/` |
| Footer badges | `primitives.HelpKeyBadge()` | `uikit/primitives/` |
| Text styling | `primitives.Title()`, `Body()`, `Muted()` | `uikit/primitives/` |
| Modals | `feedback.Modal` + `behaviors.RenderModalOverlay()` | `uikit/feedback/` |
| Boxes/containers | `containers.NewBox()` | `uikit/containers/` |

### Quick Lookup

```bash
make what-to-use NEED="modal"   # Component usage examples
make check-patterns             # Verify pattern compliance
```

### Full Reference

See [`docs/UIKIT_GUIDE.md`](../UIKIT_GUIDE.md) for complete component documentation.

---

## VHS Demo Generation

When creating new features or modifying workflows, generate visual documentation using VHS.

### Why Generate Demos?

1. **Automated documentation** - GIFs always match current behavior
2. **PR evidence** - Visual proof of feature functionality
3. **Visual regression testing** - Detect unintended UI changes
4. **Marketing material** - Professional demos for README/docs

### Quick Start

```bash
# Generate all workflow demos
make vhs-demos

# Generate specific workflow demo
make vhs-capture
make vhs-browse
make vhs-cv

# Generate feature demo for PR
make vhs-feature FEATURE=your-feature
```

### Creating Feature Demos

For new features, create demo tapes covering three scenarios:

| Scenario | Purpose | Template |
|----------|---------|----------|
| **Happy path** | Successful workflow completion | `happy-path.tape` |
| **Sad path** | Error handling and validation | `sad-path.tape` |
| **Edge cases** | Cancel, back navigation, empty states | `edge-cases.tape` |

```bash
# 1. Create feature directory
mkdir -p demos/vhs/features/my-feature

# 2. Copy templates
cp demos/vhs/features/template/*.tape demos/vhs/features/my-feature/

# 3. Edit tapes for your specific feature

# 4. Generate demos
make vhs-feature FEATURE=my-feature

# 5. Include in PR description
```

### Full Reference

See:
- [`docs/workflows/WORKFLOW_DOCUMENTATION_GUIDE.md`](../workflows/WORKFLOW_DOCUMENTATION_GUIDE.md) - Complete guide
- [`demos/vhs/README.md`](../../demos/vhs/README.md) - VHS setup and usage

---
