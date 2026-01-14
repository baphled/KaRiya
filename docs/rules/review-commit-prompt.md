---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Atomic Commit Review Prompt

---

## Purpose

This prompt helps you review staged changes and ensures you're creating atomic commits that follow project guidelines. Use this **before** finalizing any commit to catch common issues and maintain high-quality Git history.

**IMPORTANT**: Run `make check-compliance` **BEFORE** using this review process. Code quality must pass before commit review.

---

## Required Pre-Commit Workflow

```bash
# 1. REQUIRED: Run compliance check first
make check-compliance

# 2. THEN: Review commit (if needed for manual commits)
bash scripts/review-commit.sh

# 3. RECOMMENDED: Use ai-commit for AI-generated code
make ai-commit MSG="type(scope): description"
```

Or manually follow the steps below.

---

## Manual Review Process

### Step 1: Review Staged Changes

```bash
# Show files that will be committed
git diff --cached --name-status

# Show detailed diff of staged changes
git diff --cached

# Show summary statistics
git diff --cached --stat
```

**Questions to ask:**

1. ✅ Are all these files related to ONE logical change?
2. ✅ Are there any files that should be committed separately?
3. ✅ Are there any generated/temporary files included? (coverage.out, *.exe, etc.)
4. ✅ Are there any sensitive files included? (secrets, credentials, private keys)

---

### Step 2: Verify Build and Tests

```bash
# Verify code builds
go build ./...

# Run all tests
go test ./...

# Check for race conditions
go test -race ./...

# Run linting
go fmt ./...
go vet ./...
```

**Questions to ask:**

1. ✅ Does the code build successfully?
2. ✅ Do all tests pass?
3. ✅ Are there any race conditions?
4. ✅ Is the code properly formatted?

---

### Step 3: Assess Commit Scope

```bash
# Count files changed
git diff --cached --name-only | wc -l

# Count lines changed
git diff --cached --stat | tail -1

# Check which layers are affected
git diff --cached --name-only | grep -E 'domain|service|repository|cli' | cut -d'/' -f2-3 | sort -u
```

**Questions to ask:**

1. ✅ Are you changing more than 10 files? (Consider splitting)
2. ✅ Are you changing more than 500 lines? (Consider splitting, unless initial setup)
3. ✅ Are multiple architectural layers affected? (Might need separate commits)
4. ✅ Does the scope match your intended change?

---

### Step 4: Check Commit Message

Write your commit message first (in a text editor or using a template):

```bash
# Use commit message template
git config commit.template .gitmessage
git commit
```

**Commit Message Checklist:**

- [ ] Type specified (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`)?
- [ ] Scope specified (if applicable)?
- [ ] Subject line under 50 characters?
- [ ] Subject uses imperative mood ("Add" not "Added")?
- [ ] Body explains WHY the change was made?
- [ ] Body describes any side effects or implications?
- [ ] References related issues (e.g., "Closes #123")?
- [ ] No typos or grammatical errors?

**Example Review:**

❌ **Bad:**
```
update service
```

✅ **Good:**
```
feat(service): add event filtering by date range

Implement date range filtering to support timeline views and
historical analysis. Users can now query events between specific
dates, improving performance for large event collections.

The filter is applied at the repository layer using indexed
database queries for optimal performance.

Closes #56
```

---

### Step 5: Logical Coherence Check

Ask yourself:

1. **Can I describe this commit in one sentence?**
   - ✅ Yes → Good sign it's atomic
   - ❌ No → Probably needs splitting

2. **If I revert this commit, will it remove ONE logical feature/fix?**
   - ✅ Yes → Good atomic commit
   - ❌ No → Too much in one commit

3. **Can a code reviewer understand this commit without context?**
   - ✅ Yes → Clear and self-contained
   - ❌ No → Needs better message or splitting

4. **Does this commit leave the codebase in a working state?**
   - ✅ Yes → Properly complete
   - ❌ No → Incomplete, needs more files or should wait

5. **Are tests included for behavior changes?**
   - ✅ Yes → Following TDD principles
   - ❌ No → Add tests before committing

---

### Step 6: Pre-Commit Self-Review

```bash
# Review your changes one more time
git diff --cached | less

# Check for common issues
git diff --cached | grep -E 'TODO|FIXME|console\.log|debugger|XXX'

# Check for secrets/credentials
git diff --cached | grep -iE 'password|secret|api_key|token|credential'
```

**Questions to ask:**

1. ✅ No debug statements left in code?
2. ✅ No TODO comments that should be addressed?
3. ✅ No commented-out code blocks?
4. ✅ No sensitive information exposed?
5. ✅ All imports necessary and used?
6. ✅ No trailing whitespace or formatting issues?

---

## Automated Review Script

Create `scripts/review-commit.sh`:

```bash
#!/bin/bash

set -e

echo "================================================"
echo "🔍 ATOMIC COMMIT REVIEW"
echo "================================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if there are staged changes
if ! git diff --cached --quiet; then
    echo "✅ Staged changes detected"
else
    echo "${RED}❌ No staged changes${NC}"
    exit 1
fi

echo ""
echo "------------------------------------------------"
echo "📋 STAGED FILES"
echo "------------------------------------------------"
git diff --cached --name-status

echo ""
echo "------------------------------------------------"
echo "📊 CHANGE STATISTICS"
echo "------------------------------------------------"
git diff --cached --stat

# Count files and lines
FILE_COUNT=$(git diff --cached --name-only | wc -l | tr -d ' ')
LINE_COUNT=$(git diff --cached --numstat | awk '{add+=$1; del+=$2} END {print add+del}')

echo ""
echo "Files changed: $FILE_COUNT"
echo "Lines changed: $LINE_COUNT"

# Warnings
if [ "$FILE_COUNT" -gt 10 ]; then
    echo "${YELLOW}⚠️  Warning: More than 10 files changed. Consider splitting.${NC}"
fi

if [ "$LINE_COUNT" -gt 500 ]; then
    echo "${YELLOW}⚠️  Warning: More than 500 lines changed. Consider splitting.${NC}"
fi

# Check for generated files
echo ""
echo "------------------------------------------------"
echo "🔍 CHECKING FOR GENERATED FILES"
echo "------------------------------------------------"
if git diff --cached --name-only | grep -E '\.(out|exe|dll|so|dylib|test)$'; then
    echo "${RED}❌ Generated files detected! Remove from staging.${NC}"
    exit 1
else
    echo "${GREEN}✅ No generated files detected${NC}"
fi

# Check for debug statements
echo ""
echo "------------------------------------------------"
echo "🐛 CHECKING FOR DEBUG CODE"
echo "------------------------------------------------"
if git diff --cached | grep -E 'TODO|FIXME|XXX|console\.log|debugger|fmt\.Println\(\"DEBUG'; then
    echo "${YELLOW}⚠️  Warning: Debug statements found${NC}"
else
    echo "${GREEN}✅ No debug statements found${NC}"
fi

# Check for secrets
echo ""
echo "------------------------------------------------"
echo "🔐 CHECKING FOR SECRETS"
echo "------------------------------------------------"
if git diff --cached | grep -iE 'password|secret|api_key|token|credential' | grep -v 'password_hash'; then
    echo "${RED}❌ Potential secrets detected!${NC}"
    exit 1
else
    echo "${GREEN}✅ No secrets detected${NC}"
fi

# Build check
echo ""
echo "------------------------------------------------"
echo "🏗️  BUILD CHECK"
echo "------------------------------------------------"
if go build ./... 2>&1; then
    echo "${GREEN}✅ Build successful${NC}"
else
    echo "${RED}❌ Build failed${NC}"
    exit 1
fi

# Test check
echo ""
echo "------------------------------------------------"
echo "🧪 TEST CHECK"
echo "------------------------------------------------"
if go test ./... 2>&1; then
    echo "${GREEN}✅ All tests pass${NC}"
else
    echo "${RED}❌ Tests failed${NC}"
    exit 1
fi

# Race check
echo ""
echo "------------------------------------------------"
echo "🏁 RACE CONDITION CHECK"
echo "------------------------------------------------"
if go test -race ./... 2>&1; then
    echo "${GREEN}✅ No race conditions detected${NC}"
else
    echo "${RED}❌ Race conditions detected${NC}"
    exit 1
fi

# Formatting check
echo ""
echo "------------------------------------------------"
echo "✨ FORMATTING CHECK"
echo "------------------------------------------------"
UNFORMATTED=$(gofmt -l . 2>&1)
if [ -z "$UNFORMATTED" ]; then
    echo "${GREEN}✅ All files properly formatted${NC}"
else
    echo "${YELLOW}⚠️  Unformatted files:${NC}"
    echo "$UNFORMATTED"
    echo ""
    echo "Run: go fmt ./..."
fi

# Final checklist
echo ""
echo "================================================"
echo "✅ COMMIT READINESS CHECKLIST"
echo "================================================"
echo ""
echo "Before committing, verify:"
echo ""
echo "  [ ] This commit represents ONE logical change"
echo "  [ ] Commit message follows conventional format"
echo "  [ ] Commit message explains WHY, not just WHAT"
echo "  [ ] All tests pass (verified above)"
echo "  [ ] Code is properly formatted (verified above)"
echo "  [ ] No generated files included (verified above)"
echo "  [ ] No secrets or sensitive data included (verified above)"
echo "  [ ] Documentation updated (if needed)"
echo ""
echo "------------------------------------------------"
echo "📝 SUGGESTED COMMIT MESSAGE FORMAT:"
echo "------------------------------------------------"
echo ""
echo "<type>(<scope>): <subject>"
echo ""
echo "<body explaining why this change is needed>"
echo ""
echo "<footer with issue references>"
echo ""
echo "Example:"
echo "feat(service): add event filtering by date range"
echo ""
echo "Implement date range filtering for timeline views."
echo "Improves performance for large event collections."
echo ""
echo "Closes #56"
echo ""
echo "================================================"
echo ""

# Ask for confirmation
read -p "Ready to commit? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "${GREEN}✅ Proceeding with commit${NC}"
    git commit
else
    echo "${YELLOW}⏸️  Commit cancelled${NC}"
    exit 1
fi
```

Make it executable:

```bash
chmod +x scripts/review-commit.sh
```

---

## Common Issues and Solutions

### Issue 1: Too Many Files Changed

**Symptom:** More than 10 files in staging area

**Solution:**
```bash
# Unstage everything
git reset

# Commit by logical groups
git add internal/domain/career/*
git commit -m "feat(domain): add validation logic"

git add internal/service/career/*
git commit -m "feat(service): implement service layer"
```

---

### Issue 2: Mixed Concerns in One File

**Symptom:** Single file has multiple unrelated changes

**Solution:**
```bash
# Use interactive staging
git add -p filename.go

# Stage only related hunks
# 'y' to stage, 'n' to skip, 's' to split, 'q' to quit

# Commit staged hunks
git commit -m "feat: add feature X"

# Stage and commit remaining changes
git add -p filename.go
git commit -m "fix: resolve bug Y"
```

---

### Issue 3: Forgot to Add Tests

**Symptom:** Production code changed but no test changes

**Solution:**
```bash
# Don't commit yet!
# Write tests first

# Add tests
git add internal/service/career/service_test.go

# Amend the commit to include tests
git commit --amend
```

**Better approach:**
```bash
# If not yet committed, unstage everything
git reset

# Commit test first (Red)
git add internal/service/career/service_test.go
git commit -m "test(service): add failing test for feature X"

# Then commit implementation (Green)
git add internal/service/career/service.go
git commit -m "feat(service): implement feature X"
```

---

### Issue 4: Unclear What Changed

**Symptom:** Can't explain commit in one sentence

**Solution:**
```bash
# Review what's staged
git diff --cached --stat

# If multiple logical changes, unstage and split
git reset

# Group related files
git add <files-for-change-1>
git commit -m "Clear message for change 1"

git add <files-for-change-2>
git commit -m "Clear message for change 2"
```

---

### Issue 5: Generated Files Staged

**Symptom:** `coverage.out`, binaries, or build artifacts in staging

**Solution:**
```bash
# Unstage generated files
git reset coverage.out
git reset *.exe

# Add to .gitignore
echo "coverage.out" >> .gitignore
echo "*.exe" >> .gitignore

# Commit .gitignore update
git add .gitignore
git commit -m "chore: ignore generated files"
```

---

## Integration with Development Workflow

### As a Pre-Commit Hook

Install as git hook:

```bash
# Copy review script to hooks
cp scripts/review-commit.sh .git/hooks/pre-commit

# Make executable
chmod +x .git/hooks/pre-commit
```

### As a Make Target

Add to `Makefile`:

```makefile
.PHONY: review-commit

review-commit:
	@bash scripts/review-commit.sh
```

Usage:
```bash
# Before committing
make review-commit
```

### In Your Editor/IDE

#### Neovim Configuration

Add to your `init.lua`:

```lua
-- Keybinding for commit review
vim.keymap.set('n', '<leader>gc', function()
  vim.cmd('terminal bash scripts/review-commit.sh')
end, { desc = 'Review commit before finalizing' })
```

---

## Commit Review Workflow

### Recommended Process

```
1. Stage changes
   └─→ git add <files>

2. Run review script
   └─→ bash scripts/review-commit.sh

3. Script checks:
   ├─→ Staged files analysis
   ├─→ Build verification
   ├─→ Test execution
   ├─→ Race detection
   ├─→ Format checking
   ├─→ Security scanning
   └─→ Readiness checklist

4. If all checks pass:
   └─→ Write commit message
   └─→ Commit

5. If checks fail:
   └─→ Fix issues
   └─→ Return to step 2
```

---

## Quick Reference

### Before Every Commit

```bash
# 1. Review staged changes
git diff --cached --stat

# 2. Run review script
bash scripts/review-commit.sh

# 3. Verify commit message quality
#    - Type: feat/fix/docs/refactor/test/chore
#    - Scope: (domain)/(service)/(repo)/(cli)
#    - Subject: <50 chars, imperative mood
#    - Body: Explains WHY
#    - Footer: Issue references

# 4. Commit
git commit
```

### Emergency Fixes

If you realize after committing:

```bash
# Add forgotten files
git add forgotten-file.go
git commit --amend --no-edit

# Fix commit message
git commit --amend

# Split large commit
git reset --soft HEAD~1
# Then commit in smaller chunks
```

---

## Summary

**Use this prompt/script to:**
- ✅ Catch generated files before committing
- ✅ Verify builds and tests pass
- ✅ Ensure commits are atomic
- ✅ Maintain clear commit messages
- ✅ Follow project conventions
- ✅ Avoid common mistakes

**Remember:**
> "Take 5 minutes to review your commit now, or spend 5 hours debugging history later."

---

**Related Documentation:**
- [Atomic Commits Guidelines](./atomic-commits.md)
- [Senior Engineer Guidelines](./senior-engineer-guidelines.md)
- [Go Guidelines](./go-guidelines.md)

---

*Last Updated: 2025-12-23*
*Version: 1.0*

