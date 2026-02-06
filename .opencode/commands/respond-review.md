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
- `respond-to-review` - Craft responses (includes progress tracking)
- `trade-off-analysis` - Compare alternatives
- `pr-monitor` - For thread resolution and merge summary

## Process

### 1. Initial Status Report

First, report current PR status to the user:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR STATUS: #[NUM] - [Title]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CI Status: [Passing/Failing/Pending]
Comments: [Total] requiring response

| # | Reviewer | Category | Preview |
|---|----------|----------|---------|
| 1 | @xxx | [type] | "[brief]..." |

Plan: Address in severity order (Security > Bug > Arch > Style)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 2. Collect Feedback

```bash
PR_NUM=$(gh pr view --json number -q '.number')
gh pr view $PR_NUM --comments
gh api repos/:owner/:repo/pulls/$PR_NUM/comments
```

### 3. For Each Comment

#### 3a. Report Starting

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PROCESSING COMMENT #[N]/[Total]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Reviewer: @[username]
Feedback: "[Quote the comment]"
Category: [Bug/Security/Architecture/Style/Unclear]

Evaluating...
```

#### 3b. Evaluate (evaluate-change-request skill)

| Question | Answer |
|----------|--------|
| What is being requested? | [specific change] |
| What problem does this solve? | [motivation] |
| Is there evidence provided? | [yes/no] |
| Does it align with project standards? | [yes/no/no standard] |
| What are the trade-offs? | [pros/cons] |

#### 3c. Decide and Act

| Category | Response Strategy |
|----------|-------------------|
| **Bug/Security** | Accept immediately |
| **Architecture** | Evaluate against docs |
| **Performance** | Require benchmarks |
| **Style** | Follow project standard |
| **Unclear** | Ask clarifying questions |

- **Accept**: Implement change, commit, respond "Fixed in [SHA]"
- **Challenge**: Use `prove-correctness` + `justify-decision`, respond with evidence
- **Clarify**: Ask specific questions
- **Defer**: Create issue, respond with issue link

#### 3d. Post Response and Resolve Thread

```bash
# Reply to comment
gh api repos/:owner/:repo/pulls/$PR_NUM/comments/<comment_id>/replies \
  -X POST -f body="Your response"

# Resolve the thread (if accepted or clarified with answer)
gh api graphql -f query='
  mutation($threadId: ID!) {
    resolveReviewThread(input: {threadId: $threadId}) {
      thread { isResolved }
    }
  }
' -f threadId="$THREAD_ID"
```

#### 3e. Report Completion

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
COMMENT #[N] COMPLETE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Decision: [Accepted/Challenged/Clarified/Deferred]
Reasoning: [Brief explanation]
Action: [What was done]
Commit: [SHA or N/A]
Thread: [Resolved/Pending reviewer response]

Progress: [N]/[Total] | Remaining: [X] comments
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 4. When All Comments Addressed

Provide merge readiness summary using `/pr-ready` command format.

### 5. Request Re-review

```bash
gh pr comment $PR_NUM --body "Addressed all feedback. Ready for re-review."
gh pr edit $PR_NUM --add-reviewer <reviewer>
```

## Progress Summary Table

Maintain throughout:

| # | Reviewer | Category | Decision | Status |
|---|----------|----------|----------|--------|
| 1 | @alice | Bug | Accept | Resolved |
| 2 | @bob | Style | Challenge | Pending |
| 3 | @copilot | Security | Accept | Resolved |

## Skills Used
- `evaluate-change-request` - Critical assessment
- `prove-correctness` - Evidence through tests
- `justify-decision` - Explain reasoning
- `respond-to-review` - Professional communication + progress tracking
- `trade-off-analysis` - Compare alternatives
- `pr-monitor` - Thread resolution, merge summary

## Specific Comment to Address (optional)
$ARGUMENTS
