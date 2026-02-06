---
description: Evaluate and respond to PR review feedback
agent: build
---

Evaluate review feedback critically and craft appropriate responses.

**IMPORTANT: Do not blindly accept all feedback. Evaluate each comment.**

Load skills:
- `evaluate-change-request` - Assess validity
- `prove-correctness` - Write tests as evidence
- `justify-decision` - Explain architectural choices
- `respond-to-review` - Craft responses
- `trade-off-analysis` - Compare alternatives

## Process

### 1. Collect Feedback

```bash
PR_NUM=$(gh pr view --json number -q '.number')
gh pr view $PR_NUM --comments
gh api repos/:owner/:repo/pulls/$PR_NUM/comments
```

### 2. For Each Comment, Evaluate

Use `evaluate-change-request` skill:

| Question | Answer |
|----------|--------|
| What is being requested? | [specific change] |
| What problem does this solve? | [motivation] |
| Is there evidence provided? | [yes/no] |
| Does it align with project standards? | [yes/no/no standard] |
| What are the trade-offs? | [pros/cons] |

### 3. Categorise

| Category | Response Strategy |
|----------|-------------------|
| **Bug/Security** | Accept immediately |
| **Architecture** | Evaluate against docs |
| **Performance** | Require benchmarks |
| **Style** | Follow project standard |
| **Unclear** | Ask clarifying questions |

### 4. Decide Response

For each comment, determine:
- **Accept**: Implement change, respond "Fixed in [SHA]"
- **Challenge**: Provide evidence (tests, benchmarks, docs)
- **Clarify**: Ask specific questions
- **Defer**: Create issue for out-of-scope work

### 5. When Challenging

Use `prove-correctness` skill:
- Write tests demonstrating correct behaviour
- Run benchmarks if performance-related
- Reference documentation

Use `justify-decision` skill:
- Explain the problem being solved
- List alternatives considered
- State trade-offs accepted

### 6. Craft Response

Use `respond-to-review` skill:
- Be professional and specific
- Acknowledge valid points
- Provide evidence when disagreeing
- Offer path forward

### 7. Implement Accepted Changes

For accepted feedback:
```bash
# Make the fix
# Write regression test if bug
git add .
make ai-commit FILE=/tmp/fix.txt
git push
```

### 8. Post Response

```bash
# Reply to specific comment
gh api repos/:owner/:repo/pulls/$PR_NUM/comments/<comment_id>/replies \
  -X POST -f body="Your response"

# Or general comment
gh pr comment $PR_NUM --body "Response text"
```

### 9. Request Re-review

After addressing all feedback:
```bash
gh pr comment $PR_NUM --body "Addressed all feedback. Ready for re-review."
```

## Skills Used
- `evaluate-change-request` - Critical assessment
- `prove-correctness` - Evidence through tests
- `justify-decision` - Explain reasoning
- `respond-to-review` - Professional communication
- `trade-off-analysis` - Compare alternatives
- `critical-thinking` - Rigorous analysis

## Specific Comment to Address (optional)
$ARGUMENTS
