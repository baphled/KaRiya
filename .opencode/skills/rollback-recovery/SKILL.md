---
name: rollback-recovery
description: Handling failed deployments, reverting changes, and recovery procedures
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide rollback and recovery procedures when deployments fail, bugs are discovered post-merge, or system recovery is needed.

## When to use me

Use this skill when:
- A deployment causes issues
- A merged PR introduces bugs
- Need to revert changes quickly
- Database migration fails
- System needs recovery

## Rollback Decision Framework

```
Issue detected post-merge
    │
    ├─ Severity assessment
    │   ├─ Critical (system down, data loss) → Immediate rollback
    │   ├─ High (major feature broken) → Rollback within 1 hour
    │   ├─ Medium (minor feature broken) → Hotfix or rollback
    │   └─ Low (cosmetic issue) → Fix forward
    │
    └─ Rollback type
        ├─ Git revert → Safe, preserves history
        ├─ Git reset → Dangerous, rewrites history
        └─ Database restore → Last resort
```

## Git Revert (Preferred Method)

### Revert Single Commit

```bash
# Identify the problematic commit
git log --oneline -10

# Revert it (creates new commit)
git revert <commit-sha>

# Or revert without committing (to review first)
git revert --no-commit <commit-sha>
git diff --cached  # Review changes
git commit -m "revert: undo problematic change

Reverts <commit-sha> due to [reason]"
```

### Revert Merge Commit

```bash
# Find the merge commit
git log --oneline --merges -5

# Revert merge (specify parent)
# -m 1 means keep the first parent (usually main/next)
git revert -m 1 <merge-commit-sha>
```

### Revert Multiple Commits

```bash
# Revert a range (oldest to newest)
git revert --no-commit <oldest-sha>^..<newest-sha>
git commit -m "revert: undo commits from <oldest> to <newest>

Reverting due to [reason]"
```

## Emergency Rollback Procedure

### 1. Immediate Response

```bash
# Get current state
git log --oneline -5
gh pr list --state merged --limit 5

# Identify last known good state
LAST_GOOD_TAG=$(git describe --tags --abbrev=0)
echo "Last good: $LAST_GOOD_TAG"
```

### 2. Create Rollback PR

```bash
# Create rollback branch
git checkout -b hotfix/rollback-to-$LAST_GOOD_TAG

# Revert all commits since last good
git revert --no-commit HEAD...$LAST_GOOD_TAG

# Review changes
git diff --cached

# Commit
git commit -m "revert: emergency rollback to $LAST_GOOD_TAG

Issue: [describe the problem]
Impact: [what was affected]
Root cause: [if known, or 'under investigation']"

# Push and create PR
git push -u origin hotfix/rollback-to-$LAST_GOOD_TAG

gh pr create \
  --base main \
  --title "EMERGENCY: Rollback to $LAST_GOOD_TAG" \
  --body "## Emergency Rollback

**Issue:** [description]
**Severity:** Critical
**Impact:** [affected users/features]

This reverts all changes since $LAST_GOOD_TAG.

### Post-Merge Actions
- [ ] Verify system stability
- [ ] Notify stakeholders
- [ ] Create incident report
"
```

### 3. Fast-Track Merge

```bash
# For emergencies, can merge without full review
gh pr merge --admin --squash
```

## Database Rollback

### SQLite Backup Recovery

```bash
# KaRiya uses SQLite - check for backups
ls -la ~/.kariya/*.db*

# Restore from backup
cp ~/.kariya/data.db.backup ~/.kariya/data.db

# Verify integrity
sqlite3 ~/.kariya/data.db "PRAGMA integrity_check;"
```

### Migration Rollback

If a GORM migration failed:

```go
// Create down migration
func (m *Migration) Down(db *gorm.DB) error {
    // Reverse the up migration
    return db.Migrator().DropColumn(&Event{}, "new_column")
}
```

```bash
# Run down migration
go run cmd/migrate/main.go down
```

## Recovery Scenarios

### Scenario 1: Bad Merge to Main

```bash
# 1. Identify the bad merge
BAD_MERGE=$(git log --oneline --merges -1 main)

# 2. Revert it
git checkout main
git pull
git revert -m 1 $BAD_MERGE
git push

# 3. Fix on next branch
git checkout next
# Fix the issue properly
# Create new PR when ready
```

### Scenario 2: Corrupted Database

```bash
# 1. Stop the application
# 2. Check for corruption
sqlite3 ~/.kariya/data.db "PRAGMA integrity_check;"

# 3. If corrupted, try recovery
sqlite3 ~/.kariya/data.db ".recover" | sqlite3 ~/.kariya/data_recovered.db

# 4. Verify recovery
sqlite3 ~/.kariya/data_recovered.db "SELECT COUNT(*) FROM events;"

# 5. Replace if successful
mv ~/.kariya/data.db ~/.kariya/data.db.corrupted
mv ~/.kariya/data_recovered.db ~/.kariya/data.db
```

### Scenario 3: Configuration Regression

```bash
# 1. Check git history for config changes
git log --oneline -- config/ .kariya/

# 2. Restore previous config
git show <good-commit>:config/settings.yaml > config/settings.yaml

# 3. Restart application
```

### Scenario 4: Go Module Issues

```bash
# 1. Clear module cache
go clean -modcache

# 2. Restore go.mod/go.sum
git checkout HEAD~1 -- go.mod go.sum

# 3. Re-download dependencies
go mod download

# 4. Verify
go build ./...
```

## Post-Rollback Checklist

```markdown
## Post-Rollback Verification

### Immediate (within 5 minutes)
- [ ] Application starts successfully
- [ ] Core functionality works
- [ ] No error logs flooding
- [ ] Database accessible

### Short-term (within 1 hour)
- [ ] All critical paths tested
- [ ] Monitoring shows normal metrics
- [ ] No user complaints

### Follow-up
- [ ] Incident report created
- [ ] Root cause identified
- [ ] Fix developed and tested
- [ ] New PR created with fix
- [ ] Post-mortem scheduled
```

## Incident Report Template

```markdown
# Incident Report: [Title]

**Date:** YYYY-MM-DD
**Duration:** X hours Y minutes
**Severity:** Critical/High/Medium
**Status:** Resolved/Investigating

## Summary

[1-2 sentence summary of what happened]

## Timeline

| Time | Event |
|------|-------|
| HH:MM | Issue detected |
| HH:MM | Rollback initiated |
| HH:MM | System restored |
| HH:MM | Root cause identified |

## Impact

- Users affected: [number]
- Features affected: [list]
- Data loss: [none/description]

## Root Cause

[Detailed explanation of what caused the issue]

## Resolution

[What was done to fix it]

## Prevention

[Changes to prevent recurrence]

## Lessons Learned

- [Lesson 1]
- [Lesson 2]

## Action Items

- [ ] [Action] - @owner - Due: YYYY-MM-DD
```

## Prevention Strategies

### Pre-Merge

```bash
# Always run full validation
make check-compliance
make test
go test -race ./...
```

### Feature Flags

```go
// Use feature flags for risky changes
if features.IsEnabled("new_processing") {
    return newProcessing(data)
}
return oldProcessing(data)
```

### Staged Rollout

```bash
# Deploy to subset first
# Monitor
# Expand gradually
```

## Related Skills

- `incident-response` - Handling production incidents
- `monitoring` - Detecting issues early
- `git-advanced` - Complex git operations
- `feature-flags` - Safe feature rollouts
