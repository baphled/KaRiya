---
description: Comprehensive code review of staged or specified changes
agent: plan
subtask: true
---

Perform a comprehensive code review.

Load the `code-reviewer` skill for the review framework.

## What to Review

$ARGUMENTS

If no specific files mentioned, review staged changes:
```bash
git diff --cached
```

## Review Checklist

### 1. Correctness
- Does the code do what it's supposed to?
- Are edge cases handled?
- Are errors handled properly?

### 2. Clean Code (load `clean-code` skill)
- Meaningful names?
- Functions do one thing?
- No duplication?
- No magic numbers?
- No comments inside functions?
- Boy Scout Rule applied?

### 3. Architecture (load `architecture` skill)
- Correct layer placement?
- Dependencies flow downward?
- Follows intent/screen/behavior patterns?

### 4. Security (load `security` skill)
- Input validated?
- No SQL injection risks?
- No hardcoded secrets?

### 5. Testing
- Tests exist and pass?
- Coverage adequate?
- Edge cases covered?

### 6. Documentation
- Exports have godoc?
- Expected/Returns/Side effects sections?

## Output Format

For each issue found:
- [BLOCKING] - Must fix before merge
- [SUGGESTION] - Consider improving
- [QUESTION] - Need clarification
- [NICE] - Good practice worth noting

Also note any Boy Scout Rule opportunities to clean up while we're here.
