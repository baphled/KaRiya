---
name: create-pr
description: Create a pull request targeting the next branch following KaRiya branching strategy
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Create a pull request following KaRiya's branching strategy.

## When to use me

Use this skill when your feature branch is ready for review.

## Pre-PR Validation

```bash
make pre-pr
```

This validates:
- You're not on `main` or `next` branch
- All CI checks pass
- Branch is ready for PR

## Branching Rules

| Rule | Details |
|------|---------|
| Feature branches | Always work on feature branches |
| PR target | Always target `next` (NEVER `main`) |
| `main` updates | Only via `next -> main` release PRs |

## Creating the PR

### 1. Ensure Branch is Up to Date

```bash
git fetch origin
git rebase origin/next
```

### 2. Push to Remote

```bash
git push -u origin HEAD
```

### 3. Create PR

```bash
gh pr create --base next --title "type(scope): description" --body "$(cat <<'EOF'
## Summary

- Brief description of changes
- Key implementation details

## Changes

- List of specific changes made

## Testing

- How the changes were tested
- Coverage information

## Checklist

- [ ] Tests pass (`make test`)
- [ ] Coverage >= 95% for modified packages
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Documentation updated if needed
EOF
)"
```

## PR Title Format

Same as commit message:
```
type(scope): short description
```

Examples:
- `feat(intents): add burst management workflow`
- `fix(screens): resolve modal overlay z-index issue`
- `refactor(forms): migrate to forms package`

## PR Body Template

```markdown
## Summary

1-3 bullet points describing the change

## Changes

- Detailed list of changes
- New files added
- Files modified

## Testing

- Test coverage percentage
- Key scenarios tested

## Screenshots (if UI changes)

Before/After screenshots if applicable
```

## Review Process

1. CI checks must pass
2. At least one approval required
3. All comments resolved
4. Squash merge to `next`

## After Merge

```bash
git checkout next
git pull origin next
git branch -d feature/my-feature
```

## Related skills

- `ai-commit` - Create commits before PR
- `check-compliance` - Validate before PR
- `session-start` - Start new feature work
