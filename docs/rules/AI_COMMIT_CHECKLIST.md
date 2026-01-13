# AI Commit Attribution - Quick Checklist

**⚠️ MANDATORY - ZERO TOLERANCE POLICY**

## ⚠️ CRITICAL: Use `make ai-commit` (REQUIRED)

**Using `git commit` directly for AI code is PROHIBITED.**

```bash
# ✅ REQUIRED METHOD (automatic attribution)
git add <files>
make check-compliance          # MUST pass first
make ai-commit MSG="type(scope): description"

# ❌ NEVER DO THIS for AI code
git commit -m "..."            # REJECTED - no attribution
```

---

## Before Every AI-Generated Commit (STRICT ORDER)

### ✅ Step 1: Code Quality (MANDATORY)

- [ ] AI-generated code has been **thoroughly** reviewed by a human
- [ ] **ALL** tests pass (`make test`)
- [ ] Code follows **ALL** project conventions (SOLID, Go idioms)
- [ ] **NO** security vulnerabilities introduced
- [ ] Performance is acceptable
- [ ] Code is maintainable and understandable

### ✅ Step 2: Compliance Check (MANDATORY)

- [ ] **Run `make check-compliance`** - MUST pass before commit
- [ ] Fix any violations before proceeding
- [ ] Verify commit is **atomic** (ONE logical change only)

### ✅ Step 3: Commit with Attribution (REQUIRED METHOD)

- [ ] **Use `make ai-commit MSG="..."`** (automatic attribution)
- [ ] **NEVER use `git commit` directly** for AI code
- [ ] Message follows conventional commit format
- [ ] Message explains **WHY** the change was made
- [ ] Commit is atomic (single logical change)

### ✅ Format Examples

**Correct** ✅:
```
AI-Generated-By: Avante (Claude 3.5 Sonnet)
AI-Generated-By: Claude (Claude 3.7 Sonnet)
AI-Generated-By: GitHub Copilot (GPT-4)
AI-Generated-By: ChatGPT (GPT-4 Turbo)
AI-Generated-By: Cursor (Claude 3.5 Sonnet)
```

**Incorrect** ❌:
```
AI-Generated-By: Avante Claude 3.5 Sonnet         # Missing parentheses
AI-Generated-By Avante (Claude 3.5 Sonnet)        # Missing colon
AI-Generated-By: Avante                           # Missing model
# AI-Generated-By: Avante (Claude 3.5 Sonnet)    # Still commented
```

### ✅ After Commit

- [ ] Verify attribution: `make check-ai-attribution`
- [ ] Push changes to remote

---

## Quick Commands

```bash
# Setup (one-time - REQUIRED)
make install-git-hooks

# ✅ REQUIRED WORKFLOW (ONLY acceptable method)
git add -p <files>                              # Stage changes
make check-compliance                           # MUST pass before commit
make ai-commit MSG="feat(scope): description"   # Commit with attribution

# ❌ DEPRECATED - DO NOT USE (Manual workflow)
# git commit                  # REJECTED - Use make ai-commit instead

# Verification
make check-ai-attribution   # Check latest commit
make list-ai-commits        # List AI commits
make audit-ai-commits       # Full audit
```

## Enforcement

**Git hooks will REJECT:**
- ❌ Commits without AI attribution for AI code
- ❌ Incorrect attribution format
- ❌ Missing `Reviewed-By` field

**CI will REJECT PRs with:**
- ❌ AI commits without attribution
- ❌ Inconsistent attribution format
- ❌ Non-atomic commits (multiple changes)

**Status:**
| Method | Status | Description |
|--------|--------|-------------|
| `make ai-commit` | ✅ **REQUIRED** | Only acceptable method |
| `git commit` for AI code | ❌ **PROHIBITED** | Will be rejected |
| `git commit` for human code | ✅ **Allowed** | Hook will confirm |

---

## Template Reference

Use this in your commit message for AI-generated code:

```
<type>(<scope>): <subject>

<body explaining WHY>

Human review confirmed:
- All tests pass (X/X passing)
- No security issues
- Follows project conventions

AI-Generated-By: <Assistant Name> (<Model Version>)
Reviewed-By: <Your Name>
Closes #<issue>
```

---

## Common Scenarios

### Scenario 1: Full AI Generation

```
feat(service): add event filtering by date range

Implement date range filtering for timeline views.

Human review confirmed:
- All tests pass (66/66 passing)
- No security issues
- Logic verified correct

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: John Doe
Closes #56
```

### Scenario 2: AI with Human Modifications

```
refactor(repo): optimize query builder

Refactor query building logic for better performance.

Development process:
- Initial implementation by AI
- Manual optimization of allocations
- Added benchmark tests

AI-Generated-By: Avante (Claude 3.5 Sonnet)
AI-Modifications: Optimized hot paths, added benchmarks
Reviewed-By: Jane Smith
```

### Scenario 3: Human-Written Code

```
fix(domain): prevent duplicate tags

Tag deduplication was not working correctly.
Now normalize to lowercase before checking.

Fixes #87
```

(No AI attribution needed - hook will confirm)

---

## Red Flags (Don't Commit)

❌ **Reject AI code if**:
- Tests don't pass
- Logic is incorrect
- Security vulnerabilities present
- Doesn't follow project patterns
- You don't understand the code
- Performance is unacceptable

**When in doubt, regenerate or write manually.**

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Hook not running | `make install-git-hooks` |
| Format validation fails | Check parentheses around model |
| Template not showing | `git config commit.template .gitmessage` |
| Need to bypass for docs | `git commit --no-verify` (docs/config only!) |

---

## Resources

- **Full Guide**: `docs/rules/AI_COMMIT_ATTRIBUTION.md`
- **Setup**: `AI_COMMIT_SETUP.md`
- **Summary**: `AI_COMMIT_ATTRIBUTION_SETUP.md`
- **Help**: `make help`

---

## Remember

🤖 **AI-generated code = MUST have attribution**
👤 **Human review = ALWAYS required**
✅ **Format = AI-Generated-By: Name (Model)**
🔍 **Verify = make check-ai-attribution**

---

**Print this and keep it visible during development!**

