# KaRiya Bug Tracking

**Purpose**: Track and resolve bugs in a structured, systematic way

**Status**: Active bug tracking system

---

## Bug Lifecycle

```
Reported → Investigation → Root Cause → Fix → Testing → Verification → Closed
```

### States

- **Reported**: Bug identified, initial description created
- **Investigation**: Gathering evidence, reproducing issue
- **Root Cause**: Identified the underlying cause
- **Fix**: Implementation of solution
- **Testing**: Verification via automated and manual tests
- **Verification**: Final acceptance testing
- **Closed**: Bug resolved and verified

---

## Bug Naming Convention

`bug-<number>-<short-description>.md`

**Examples**:
- `bug-001-escape-key-navigation.md`
- `bug-002-form-validation-error.md`
- `bug-003-cv-generation-timeout.md`

---

## Bug Template Structure

Each bug report should include:

1. **Bug Summary**: One-line description
2. **Status**: Current state in lifecycle
3. **Severity**: Critical / High / Medium / Low
4. **Affected Components**: Which parts of codebase
5. **Reproduction Steps**: How to reproduce consistently
6. **Expected Behavior**: What should happen
7. **Actual Behavior**: What actually happens
8. **Investigation Log**: Timestamped notes during investigation
9. **Root Cause**: Identified cause (once known)
10. **Fix Strategy**: Approach to resolution
11. **Testing Plan**: How to verify fix
12. **Verification**: Final acceptance criteria
13. **Related Files**: Code files involved

---

## Bug Severity Levels

### 🔴 Critical
- Application crash or data loss
- Complete feature failure
- Security vulnerability
- Blocks all users

### 🟠 High
- Major feature broken
- Workaround exists but difficult
- Affects most users
- Poor user experience

### 🟡 Medium
- Minor feature issue
- Easy workaround available
- Affects some users
- Degraded user experience

### 🟢 Low
- Cosmetic issue
- Very minor functionality issue
- Affects few users
- Minimal impact

---

## Debugging Workflow

### Phase 1: Reproduction (30-60 min)
- [ ] Document exact steps to reproduce
- [ ] Test in clean environment
- [ ] Verify issue occurs consistently
- [ ] Capture screenshots/logs if applicable
- [ ] Identify affected workflows

### Phase 2: Investigation (1-3 hours)
- [ ] Review relevant code
- [ ] Check related tests (do they pass? should they fail?)
- [ ] Search for similar issues in codebase
- [ ] Identify potential root causes
- [ ] Document investigation findings

### Phase 3: Root Cause Analysis (1-2 hours)
- [ ] Trace execution path
- [ ] Identify exact failure point
- [ ] Determine why it fails
- [ ] Document root cause
- [ ] Assess impact and scope

### Phase 4: Fix Implementation (2-6 hours)
- [ ] Design fix strategy
- [ ] Implement minimal fix
- [ ] Add/update tests
- [ ] Verify fix locally
- [ ] Run full test suite
- [ ] Check for regressions

### Phase 5: Testing & Verification (1-2 hours)
- [ ] Manual testing of fix
- [ ] Automated test coverage
- [ ] Edge case testing
- [ ] Regression testing
- [ ] Performance impact check

### Phase 6: Documentation (30 min)
- [ ] Update bug report with resolution
- [ ] Update relevant documentation
- [ ] Add comments to code if needed
- [ ] Create follow-up tasks if applicable

---

## Active Bugs

| Bug ID | Summary | Severity | Status | Assigned |
|--------|---------|----------|--------|----------|
| 001 | Escape key navigation not working | 🟠 High | Investigation | - |

---

## Closed Bugs

*No closed bugs yet*

---

## Related Documentation

- **Tasks**: `/tasks` - Feature implementation tracking
- **Development Rules**: `/docs/rules` - Development guidelines
- **Testing Guides**: `/docs/development/NAVIGATION_TESTING_GUIDE.md`
- **Troubleshooting**: `/docs/TROUBLESHOOTING.md`

---

**Last Updated**: 2026-01-12
