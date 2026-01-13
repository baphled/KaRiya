# AI Commit Attribution Rules

**⚠️ MANDATORY - ZERO TOLERANCE POLICY**

---

## Overview

**ALL commits created by AI assistants MUST be clearly labeled as AI-generated.**

**This is NON-NEGOTIABLE. Commits without proper AI attribution will be rejected.**

This document establishes **MANDATORY and STRICTLY ENFORCED** rules for attributing commits created by AI assistants (OpenCode, Avante, Claude, GitHub Copilot, ChatGPT, etc.) to ensure transparency, traceability, and accountability in the project's git history.

### Zero Tolerance Policy

- ❌ **REJECTED**: Any AI-generated commit without attribution
- ❌ **REJECTED**: Any commit using `git commit` directly for AI code
- ✅ **REQUIRED**: Use `make ai-commit` for ALL AI-generated commits
- ✅ **REQUIRED**: Human review documented for every AI commit

---

## Table of Contents

1. [Core Requirements](#core-requirements)
2. [Commit Message Format](#commit-message-format)
3. [AI Attribution Markers](#ai-attribution-markers)
4. [Examples](#examples)
5. [Automation](#automation)
6. [Verification](#verification)
7. [Human Review Requirements](#human-review-requirements)

---

## Core Requirements

### Mandatory AI Attribution (ZERO EXCEPTIONS)

**EVERY commit created with AI assistance MUST include:**

1. **AI Assistant Name**: Which AI tool generated the code (e.g., OpenCode, Avante, Claude, Copilot)
2. **Model Version**: The specific model used (e.g., Claude Sonnet 4.5, Claude 3.5 Sonnet, GPT-4)
3. **Attribution Marker**: `AI-Generated-By:` in the commit footer
4. **Human Review**: `Reviewed-By:` in the commit footer

### Non-Negotiable Rules (STRICT ENFORCEMENT)

**MUST DO (MANDATORY):**
- ✅ **ALWAYS** use `make ai-commit MSG="..."` for AI-generated commits
- ✅ **ALWAYS** include AI attribution in commit messages (automatic via make ai-commit)
- ✅ **ALWAYS** specify the exact model used
- ✅ **ALWAYS** review AI-generated code before committing
- ✅ **ALWAYS** run `make check-compliance` before committing
- ✅ **ALWAYS** run all tests before committing
- ✅ **ALWAYS** document human review with `Reviewed-By:`

**NEVER DO (REJECTED):**
- ❌ **NEVER** use `git commit` directly for AI-generated code
- ❌ **NEVER** commit AI-generated code without attribution
- ❌ **NEVER** use vague attribution like "AI-assisted" without details
- ❌ **NEVER** skip human review of AI-generated changes
- ❌ **NEVER** commit without passing `make check-compliance`
- ❌ **NEVER** use `--no-verify` flag (except emergency situations)

**Consequences of Violation:**
- Commit will be rejected by git hooks
- PR will be rejected by CI
- Work must be restarted with proper attribution

---

## Commit Message Format

### Standard Format with AI Attribution

```
<type>(<scope>): <subject line>

<body explaining WHY this change is needed>
<body explaining WHAT was changed>
<body describing human review/validation performed>

AI-Generated-By: <assistant-name> (<model-version>)
Reviewed-By: <human-name>
<footer with issue references>
```

### Required Fields

1. **AI-Generated-By**: Name of AI assistant and model version
2. **Reviewed-By**: Name of human who reviewed and approved the changes

### Optional Fields

- **AI-Prompt**: Brief description of the prompt given to AI (if relevant)
- **AI-Modifications**: Changes made by human after AI generation
- **Co-Authored-By**: If multiple humans collaborated on review

---

## AI Attribution Markers

### Format

```
AI-Generated-By: <Assistant Name> (<Model Version>)
```

### Recognized Assistants

#### Avante
```
AI-Generated-By: Avante (Claude 3.5 Sonnet)
```

#### Claude
```
AI-Generated-By: Claude (Claude 3.5 Sonnet)
AI-Generated-By: Claude (Claude 3.7 Sonnet)
AI-Generated-By: Claude (Claude 3 Opus)
```

#### GitHub Copilot
```
AI-Generated-By: GitHub Copilot (GPT-4)
AI-Generated-By: GitHub Copilot (GPT-3.5 Turbo)
```

#### ChatGPT
```
AI-Generated-By: ChatGPT (GPT-4)
AI-Generated-By: ChatGPT (GPT-4 Turbo)
AI-Generated-By: ChatGPT (o1)
```

#### Cursor
```
AI-Generated-By: Cursor (Claude 3.5 Sonnet)
AI-Generated-By: Cursor (GPT-4)
```

#### OpenCode
```
AI-Generated-By: OpenCode (Claude Sonnet 4.5)
AI-Generated-By: OpenCode (Claude Sonnet 4)
AI-Generated-By: OpenCode (Claude 3.5 Sonnet)
```

### Reviewed-By Format

```
Reviewed-By: Your Name <your.email@example.com>
```

Or simply:

```
Reviewed-By: Your Name
```

---

## Examples

### Example 1: Feature Implementation

```
feat(service): add event filtering by date range

Implement date range filtering for timeline views to support
historical analysis and performance improvements for large
event collections.

Changes include:
- Add StartDate and EndDate fields to ListFilters
- Implement date comparison in filter logic
- Add comprehensive unit tests for edge cases

Human review confirmed:
- All tests pass (100% coverage maintained)
- No breaking changes to existing API
- Performance acceptable for expected dataset sizes

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: John Doe <john@example.com>
Closes #56
```

### Example 2: Bug Fix

```
fix(domain): prevent duplicate tags in event validation

Tag deduplication was not working correctly when tags had
different casing (e.g., "Technical" vs "technical"). This
caused validation to fail incorrectly.

Solution: Normalize all tags to lowercase before deduplication
check, ensuring case-insensitive comparison.

Human modifications:
- Added additional test case for mixed-case scenarios
- Updated error message for clarity

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: Jane Smith
Fixes #87
```

### Example 3: Test Addition

```
test(classification): add coverage for edge cases

Add comprehensive tests for:
- Multi-category classification
- Empty text handling
- Special characters in event text
- Null/undefined event fields

All tests pass. Coverage increased from 84% to 100%.

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: Bob Johnson
```

### Example 4: Documentation

```
docs(rules): add AI commit attribution guidelines

Create comprehensive guidelines for attributing AI-generated
commits to ensure transparency and traceability.

Includes:
- Mandatory attribution format
- Examples for all common AI assistants
- Verification scripts and automation
- Human review requirements

AI-Generated-By: Claude (Claude 3.7 Sonnet)
Reviewed-By: Team Lead
```

### Example 5: Refactoring

```
refactor(repo): extract query building into helper method

Consolidate duplicate query construction code from List()
and Count() methods into a shared buildQueryFilters() helper.
No behavior changes.

Human verification:
- Ran full test suite (31/31 passing)
- Verified no performance regression
- Confirmed query output identical to previous implementation

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: Sarah Chen
```

### Example 6: Multiple AI Sessions

```
feat(cli): add interactive event capture form

Implement BubbleTea-based form with real-time validation,
character counter, and date parsing.

Features:
- Text input with 2000 character limit
- Smart date parsing (ISO, relative, special keywords)
- Visual focus indicators and navigation
- Inline validation with error display

Development process:
- Initial form structure generated by AI
- Validation logic refined by human
- Styling manually adjusted for consistency

AI-Generated-By: Avante (Claude 3.5 Sonnet)
AI-Modifications: Adjusted date parsing logic, refined error messages
Reviewed-By: Alex Rodriguez <alex@example.com>
```

---

## Automation

### MANDATORY: `make ai-commit` Command

**⚠️ THIS IS THE ONLY ACCEPTABLE METHOD FOR AI-GENERATED COMMITS**

**You MUST use `make ai-commit` for ALL AI-generated commits:**

```bash
# Stage your changes
git add -p internal/cli/forms/validators.go

# Create AI-attributed commit (ONLY METHOD ALLOWED)
make ai-commit MSG="feat(forms): add date validation helpers"
```

**What it does (AUTOMATIC):**
1. ✅ Validates commit message follows conventional commit format
2. ✅ Checks that changes are staged
3. ✅ Automatically adds AI attribution (`AI-Generated-By: OpenCode (Claude Sonnet 4.5)`)
4. ✅ Automatically adds reviewer attribution from `git config user.name`
5. ✅ Creates the commit with proper formatting
6. ✅ Enforces all project commit standards

**Environment variables** (optional - for non-OpenCode assistants):
```bash
# Override agent/model if using different AI assistant
AI_AGENT="Cursor" AI_MODEL="Claude 3.5 Sonnet" make ai-commit MSG="feat: ..."
AI_AGENT="Avante" AI_MODEL="Claude 3.5 Sonnet" make ai-commit MSG="fix: ..."
```

**Default values (OpenCode):**
- `AI_AGENT=OpenCode`
- `AI_MODEL=Claude Sonnet 4.5`

**Advantages (WHY THIS IS MANDATORY):**
- ✅ Zero-effort AI attribution (automatic)
- ✅ Impossible to forget attribution
- ✅ Automatic format validation (prevents bad commits)
- ✅ Consistent attribution format (across entire project)
- ✅ Prevents common mistakes (missing fields, wrong format)
- ✅ Enforces project standards (conventional commits)
- ✅ Integrated with git hooks (double validation)

**Comparison: Manual vs. make ai-commit**

| Feature | Manual `git commit` | `make ai-commit` |
|---------|-------------------|------------------|
| AI attribution | ❌ Manual (error-prone) | ✅ Automatic |
| Format validation | ❌ Post-commit via hook | ✅ Pre-commit validation |
| Consistency | ❌ Human error risk | ✅ 100% consistent |
| Speed | ⚠️ Slower (manual typing) | ✅ Fast (one command) |
| Compliance | ❌ Easy to forget | ✅ Always compliant |
| **Status** | ❌ **DEPRECATED - DO NOT USE** | ✅ **REQUIRED - ONLY METHOD** |

### Alternative: Git Commit Template

Create a commit message template that includes AI attribution placeholders:

#### Create `.gitmessage` Template

```bash
# Create commit template file
cat > ~/.gitmessage << 'EOF'
# <type>(<scope>): <subject>
#
# <body>
#
# AI-Generated-By: <Assistant Name> (<Model Version>)
# Reviewed-By: <Your Name>
# Closes #<issue>
#
# Types: feat, fix, docs, style, refactor, test, chore, perf
# Scopes: domain, service, repo, cli, logger
#
# Remember:
# - Explain WHY, not just WHAT
# - Include human review notes for AI-generated code
# - Verify all tests pass before committing
EOF

# Configure git to use the template
git config --global commit.template ~/.gitmessage
```

Or for project-local configuration:

```bash
# Create template in project
cat > .gitmessage << 'EOF'
# <type>(<scope>): <subject>
#
# <body>
#
# AI-Generated-By: <Assistant Name> (<Model Version>)
# Reviewed-By: <Your Name>
# Closes #<issue>
EOF

# Configure for project only
git config commit.template .gitmessage
```

### Prepare-Commit-Msg Hook

Create a git hook to automatically add AI attribution template:

```bash
# Create hook
cat > .git/hooks/prepare-commit-msg << 'EOF'
#!/bin/bash

# Get commit message file
COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# Only add template for regular commits (not merges, etc.)
if [ -z "$COMMIT_SOURCE" ]; then
    # Check if AI attribution is missing
    if ! grep -q "AI-Generated-By:" "$COMMIT_MSG_FILE"; then
        # Add AI attribution footer template
        echo "" >> "$COMMIT_MSG_FILE"
        echo "# AI ATTRIBUTION (remove this section if not AI-generated):" >> "$COMMIT_MSG_FILE"
        echo "# AI-Generated-By: Avante (Claude 3.5 Sonnet)" >> "$COMMIT_MSG_FILE"
        echo "# Reviewed-By: Your Name" >> "$COMMIT_MSG_FILE"
    fi
fi
EOF

# Make executable
chmod +x .git/hooks/prepare-commit-msg
```

### Commit-Msg Hook (Validation)

Create a hook to validate AI attribution is present:

```bash
cat > .git/hooks/commit-msg << 'EOF'
#!/bin/bash

COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Skip if merge commit
if echo "$COMMIT_MSG" | grep -q "^Merge"; then
    exit 0
fi

# Check for AI markers in code changes
STAGED_FILES=$(git diff --cached --name-only)
CODE_FILES=$(echo "$STAGED_FILES" | grep -E '\.(go|js|ts|py|java|c|cpp|rs)$' || true)

# If code files are present, check for AI attribution
if [ -n "$CODE_FILES" ]; then
    # Look for AI attribution marker
    if ! echo "$COMMIT_MSG" | grep -q "AI-Generated-By:"; then
        echo ""
        echo "❌ ERROR: AI-Generated code detected without attribution"
        echo ""
        echo "If this code was AI-generated, add attribution:"
        echo ""
        echo "  AI-Generated-By: <Assistant> (<Model>)"
        echo "  Reviewed-By: <Your Name>"
        echo ""
        echo "If this code was NOT AI-generated, you can proceed."
        echo ""
        echo "To bypass this check (human-written code only):"
        echo "  git commit --no-verify"
        echo ""
        read -p "Was this code AI-generated? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Please add AI attribution and commit again."
            exit 1
        fi
    else
        # Validate format
        if ! echo "$COMMIT_MSG" | grep -q "AI-Generated-By: .* (.*)"; then
            echo ""
            echo "❌ ERROR: AI attribution format incorrect"
            echo ""
            echo "Expected format:"
            echo "  AI-Generated-By: <Assistant Name> (<Model Version>)"
            echo ""
            echo "Examples:"
            echo "  AI-Generated-By: Avante (Claude 3.5 Sonnet)"
            echo "  AI-Generated-By: Claude (Claude 3.7 Sonnet)"
            echo ""
            exit 1
        fi

        # Check for human review
        if ! echo "$COMMIT_MSG" | grep -q "Reviewed-By:"; then
            echo ""
            echo "⚠️  WARNING: AI-generated code should include 'Reviewed-By:'"
            echo ""
            echo "Add: Reviewed-By: <Your Name>"
            echo ""
            read -p "Continue without review attribution? (y/n): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
    fi
fi

exit 0
EOF

chmod +x .git/hooks/commit-msg
```

---

## Verification

### Check AI Attribution in History

```bash
# List all commits with AI attribution
git log --all --grep="AI-Generated-By:"

# Count AI-generated commits
git log --all --grep="AI-Generated-By:" --oneline | wc -l

# List AI-generated commits by assistant
git log --all --grep="AI-Generated-By: Avante"
git log --all --grep="AI-Generated-By: Claude"
git log --all --grep="AI-Generated-By: GitHub Copilot"

# Show commits WITHOUT AI attribution in code files
git log --all --no-merges --oneline --name-only | \
  grep -E '\.(go|js|ts|py)$' -B 1 | \
  grep -v "AI-Generated-By:"
```

### Verify Latest Commit

```bash
# Check if latest commit has AI attribution
if git log -1 --pretty=%B | grep -q "AI-Generated-By:"; then
    echo "✅ AI attribution present"
    git log -1 --pretty=%B | grep "AI-Generated-By:"
else
    echo "⚠️  No AI attribution found"
fi
```

### Project-Wide Audit

```bash
# Create audit script
cat > scripts/audit-ai-commits.sh << 'EOF'
#!/bin/bash

echo "================================================"
echo "🤖 AI COMMIT ATTRIBUTION AUDIT"
echo "================================================"
echo ""

TOTAL_COMMITS=$(git log --all --oneline --no-merges | wc -l)
AI_COMMITS=$(git log --all --oneline --no-merges --grep="AI-Generated-By:" | wc -l)

echo "Total commits: $TOTAL_COMMITS"
echo "AI-attributed commits: $AI_COMMITS"
echo ""

if [ $AI_COMMITS -gt 0 ]; then
    PERCENTAGE=$(echo "scale=2; ($AI_COMMITS / $TOTAL_COMMITS) * 100" | bc)
    echo "Percentage AI-generated: ${PERCENTAGE}%"
    echo ""

    echo "AI Assistants Used:"
    echo "-------------------"
    git log --all --grep="AI-Generated-By:" --pretty=%B | \
      grep "AI-Generated-By:" | \
      sort | uniq -c | sort -rn
    echo ""
fi

# Check for code commits without AI attribution
echo "Checking for code commits without attribution..."
CODE_COMMITS_NO_AI=$(git log --all --oneline --no-merges --name-only | \
  awk '/\.(go|js|ts|py|java|c|cpp|rs)$/{getline prev; if(prev !~ /AI-Generated-By:/) print prev}' | \
  wc -l)

if [ $CODE_COMMITS_NO_AI -gt 0 ]; then
    echo "⚠️  Found $CODE_COMMITS_NO_AI code commits without AI attribution"
    echo "    (These may be human-written, which is fine)"
else
    echo "✅ All code commits properly attributed"
fi

echo ""
echo "================================================"
EOF

chmod +x scripts/audit-ai-commits.sh
```

---

## Human Review Requirements

### Mandatory Review Steps

Before committing AI-generated code, ALWAYS:

1. **Read the Code**: Understand what the AI generated
2. **Run Tests**: Ensure all tests pass
3. **Check Logic**: Verify the logic is correct
4. **Review Edge Cases**: Consider scenarios the AI might have missed
5. **Check Standards**: Ensure code follows project conventions
6. **Verify Security**: Check for security vulnerabilities
7. **Test Manually**: If applicable, manually test the functionality

### Review Checklist

Add this to your commit message body:

```
Human review performed:
- [ ] Code logic verified correct
- [ ] All tests pass (X/X passing)
- [ ] No security vulnerabilities introduced
- [ ] Follows project coding standards
- [ ] Edge cases considered
- [ ] Performance acceptable
- [ ] Documentation updated (if needed)
```

### Red Flags (Reject AI Code)

Reject AI-generated code if:
- ❌ Tests don't pass
- ❌ Logic is incorrect or incomplete
- ❌ Security vulnerabilities present
- ❌ Doesn't follow project patterns
- ❌ Performance is unacceptable
- ❌ You don't understand the code

**When in doubt, regenerate or write manually.**

---

## Integration with Review Script

Update `scripts/review-commit.sh` to check for AI attribution:

```bash
# Add this section to review-commit.sh

echo ""
echo "------------------------------------------------"
echo "🤖 AI ATTRIBUTION CHECK"
echo "------------------------------------------------"

# Check if code files are staged
CODE_FILES=$(git diff --cached --name-only | grep -E '\.(go|js|ts|py|java|c|cpp|rs)$' || true)

if [ -n "$CODE_FILES" ]; then
    echo "Code files detected in commit."
    echo ""
    echo "⚠️  REMINDER: If AI-generated, commit message MUST include:"
    echo ""
    echo "  AI-Generated-By: <Assistant Name> (<Model Version>)"
    echo "  Reviewed-By: <Your Name>"
    echo ""
    echo "Example:"
    echo "  AI-Generated-By: Avante (Claude 3.5 Sonnet)"
    echo "  Reviewed-By: John Doe"
    echo ""
else
    echo "No code files in this commit."
fi
```

---

## Makefile Integration

Add AI attribution helpers to Makefile:

```makefile
# AI Attribution Commands
.PHONY: check-ai-attribution
check-ai-attribution:
	@echo "Checking latest commit for AI attribution..."
	@git log -1 --pretty=%B | grep "AI-Generated-By:" || \
		echo "⚠️  No AI attribution found in latest commit"

.PHONY: audit-ai-commits
audit-ai-commits:
	@./scripts/audit-ai-commits.sh

.PHONY: list-ai-commits
list-ai-commits:
	@echo "AI-Generated Commits:"
	@git log --all --grep="AI-Generated-By:" --oneline

.PHONY: ai-stats
ai-stats:
	@echo "AI Commit Statistics:"
	@echo "Total AI commits: $(shell git log --all --grep='AI-Generated-By:' --oneline | wc -l)"
	@echo ""
	@echo "By Assistant:"
	@git log --all --grep="AI-Generated-By:" --pretty=%B | \
		grep "AI-Generated-By:" | sort | uniq -c
```

---

## Summary

### Key Takeaways

1. **ALL AI-generated commits MUST include attribution**
2. **Format**: `AI-Generated-By: <Assistant> (<Model>)`
3. **Always include**: `Reviewed-By: <Human Name>`
4. **Human review is mandatory** for all AI-generated code
5. **Use automation** (hooks, templates) to enforce rules
6. **Audit regularly** to ensure compliance

### Quick Reference

```bash
# Set up commit template
git config commit.template .gitmessage

# Check latest commit
make check-ai-attribution

# Audit all commits
make audit-ai-commits

# List AI commits
make list-ai-commits
```

---

**Version**: 1.0
**Created**: 2025-12-23
**Status**: Active and Mandatory

**Remember**: Transparency is key. Always attribute AI-generated code clearly.

