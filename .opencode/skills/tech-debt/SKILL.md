---
name: tech-debt
description: Identify, document, and manage technical debt for systematic reduction
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help identify, document, and systematically address technical debt in the codebase.

## When to use me

Use this skill when:
- Discovering code that needs improvement
- Planning refactoring work
- Assessing code quality
- Making strategic debt decisions

## What is Technical Debt?

Technical debt is the implied cost of future work caused by choosing an easy solution now instead of a better approach that would take longer.

**Not all debt is bad** - sometimes it's a strategic choice to ship faster.

## Types of Technical Debt

| Type | Description | Example |
|------|-------------|---------|
| **Deliberate** | Conscious decision to ship faster | "We'll add tests later" |
| **Accidental** | Learned better approach after | "Now I know a better pattern" |
| **Bit Rot** | Code degraded over time | Outdated dependencies |
| **Design Debt** | Architecture doesn't fit anymore | Monolith needs to split |

## Identifying Tech Debt

### Code Smells

Run pattern checks to find issues:
```bash
make check-patterns
make check-intent-architecture
make golangci-lint
```

### Common Debt Indicators

| Indicator | Sign | Action |
|-----------|------|--------|
| Long functions | > 50 lines | Extract methods |
| High complexity | Deep nesting | Simplify logic |
| Duplicate code | Copy-paste | Extract shared function |
| Missing tests | < 80% coverage | Add tests |
| Outdated deps | Security alerts | Update dependencies |
| Deprecated patterns | Old components | Migrate to new |
| TODO comments | Incomplete work | Complete or track |

### Architecture Debt in KaRiya

```bash
# Find deprecated patterns
grep -r "components.KeyBadge" internal/cli/
grep -r "models.*Form" internal/cli/intents/

# Find large files
find internal/cli/intents -name "*.go" -exec wc -l {} \; | sort -rn | head -20

# Find files without tests
for f in internal/cli/intents/*/*.go; do
  test_file="${f%.go}_test.go"
  [ ! -f "$test_file" ] && echo "Missing: $test_file"
done
```

## Documenting Tech Debt

Create a tech debt item:

```markdown
# DEBT-XXX: [Descriptive Title]

## Summary

Brief description of the debt and its impact.

## Current State

What the code looks like now.

## Desired State

What it should look like after addressing.

## Impact

### Pain Points
- Slows down development because...
- Causes bugs when...
- Makes testing difficult because...

### Risk Assessment
- [ ] High - Blocks features, causes bugs
- [ ] Medium - Slows development
- [ ] Low - Cosmetic, minor inconvenience

## Effort Estimate

| Aspect | Estimate |
|--------|----------|
| Time | X hours/days |
| Files affected | X files |
| Risk of regression | Low/Medium/High |

## Proposed Solution

1. Step 1
2. Step 2
3. Step 3

## Related

- TASK-XXX: Feature blocked by this debt
- BUG-YYY: Bug caused by this debt
```

## Tech Debt Quadrant

Use this to prioritize:

```
                    High Value
                        │
     QUICK WINS         │         MAJOR PROJECTS
     Do these now       │         Plan carefully
                        │
  ──────────────────────┼──────────────────────
                        │
     DON'T BOTHER       │         THANKLESS TASKS
     Low priority       │         Do if time permits
                        │
                    Low Value
        ◄─────────────────────────────────────►
        Low Effort                   High Effort
```

## Managing Debt

### The Boy Scout Rule

**Always leave code cleaner than you found it.**

When working on a feature:
1. Fix obvious issues you encounter
2. Don't scope-creep into large refactors
3. Document significant debt for later

### Debt Budget

Allocate time for debt reduction:
- **20% rule**: 1 day per week for debt
- **Debt sprints**: Periodic focused cleanup
- **With features**: Bundle related cleanup

### When to Pay Debt

**Pay now if:**
- Blocking a feature
- Causing recurring bugs
- Slowing every change in the area
- Security vulnerability

**Pay later if:**
- Area rarely changes
- Low impact on development
- Larger refactor needed first

## Tech Debt Audit

Periodically assess debt:

```markdown
## Tech Debt Audit - [Date]

### High Priority (Pay This Sprint)
1. DEBT-XXX: [Title] - [Reason]
2. DEBT-YYY: [Title] - [Reason]

### Medium Priority (Pay This Quarter)
1. DEBT-ZZZ: [Title]

### Low Priority (Backlog)
1. DEBT-AAA: [Title]

### Resolved Since Last Audit
1. DEBT-BBB: [Title] - Resolved in PR #123
```

## Related skills

- `clean-code` - Clean code principles
- `refactor` - Refactoring command
- `code-reviewer` - Spot debt in reviews
- `architecture` - Architecture debt
