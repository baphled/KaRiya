---
name: respond-to-review
description: Craft thoughtful, professional responses to code review feedback
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Craft clear, professional, and constructive responses to code review feedback. I ensure communication is respectful while being technically accurate.

## When to use me

Use this skill when:
- Responding to any review comment
- Accepting feedback gracefully
- Challenging feedback respectfully
- Asking clarifying questions

## Core Principles

1. **Assume good intent** - Reviewers want to help improve the code
2. **Be specific** - Vague responses waste everyone's time
3. **Stay professional** - No defensiveness or sarcasm
4. **Focus on the code** - Not on the person
5. **Be concise** - Respect the reviewer's time

## Response Templates

### Accepting Feedback

**Simple acceptance:**
```markdown
Good catch! Fixed in [commit SHA].
```

**Acceptance with context:**
```markdown
You're right, this was an oversight. I've added [specific fix] and a regression test to prevent this in future.

Fixed in [commit SHA].
```

**Acceptance with learning:**
```markdown
Thanks for catching this! I wasn't aware of [pattern/issue]. 
I've fixed it and will watch for this in future PRs.

Fixed in [commit SHA].
```

### Challenging Feedback

**Respectful disagreement:**
```markdown
Thanks for the suggestion. I considered this approach but chose differently because [reason].

[Evidence: tests/benchmarks/documentation]

Happy to discuss further if you see something I've missed.
```

**Counter with evidence:**
```markdown
I investigated this concern and found:

[Test/benchmark/analysis results]

Based on this evidence, I believe the current approach is correct because [reason].

Let me know if I've misunderstood the concern.
```

**Partial agreement:**
```markdown
You raise a valid point about [aspect]. I've addressed that by [change].

However, for [other aspect], I'd prefer to keep the current approach because [reason with evidence].

Does that address your concern?
```

### Asking Clarification

**Unclear feedback:**
```markdown
Could you clarify what you mean by [quote]? 

I want to make sure I understand the concern before making changes.
Are you suggesting [interpretation A] or [interpretation B]?
```

**Missing context:**
```markdown
I'm not familiar with [pattern/approach] you mentioned. 
Could you point me to an example or documentation?

That would help me evaluate whether it's appropriate here.
```

**Scope question:**
```markdown
This is a good suggestion. To clarify scope:
- Should I address this in this PR, or
- Create a follow-up issue/PR?

It would require changes to [areas] which might be better as a separate effort.
```

### Deferring

**Out of scope:**
```markdown
Good idea! That's outside the scope of this PR though.

I've created issue #[X] to track this. Happy to address in a follow-up.
```

**Needs more discussion:**
```markdown
This is a significant change that would benefit from broader discussion.

I've started a discussion in [location] to get more input before implementing.
```

## Response Structure

### For Simple Feedback

```markdown
[Acknowledgment]. [Action taken].
```

Example:
```markdown
Good catch! Fixed the typo in abc123.
```

### For Complex Feedback

```markdown
## Understanding
[Restate the feedback to confirm understanding]

## Analysis
[Your evaluation of the feedback]

## Evidence
[Tests/benchmarks/documentation if applicable]

## Decision
[What you're doing and why]

## Action
[Specific commits/changes made, or next steps]
```

### For Disagreement

```markdown
Thanks for [specific aspect of feedback].

I've investigated and found [evidence]. Based on this, I believe [conclusion] because [reasoning].

[Specific evidence: tests, benchmarks, documentation]

That said, [acknowledgment of valid points]. [Proposed path forward].

Happy to discuss further.
```

## Tone Guidelines

### DO

- Thank reviewers for their time
- Acknowledge valid points even when disagreeing
- Use "I" statements ("I believe", "I found")
- Ask questions when unsure
- Offer to discuss synchronously for complex topics
- Commit to follow-up actions

### DON'T

- Be defensive ("That's not what I meant")
- Be dismissive ("That doesn't matter")
- Be passive-aggressive ("I guess I'll change it then")
- Attack the reviewer ("You don't understand")
- Make excuses ("I was rushed")
- Ignore feedback (always respond)

## Response Timing

| Feedback Type | Response Time |
|---------------|---------------|
| Blocking issue | Same day |
| Significant feedback | Within 24 hours |
| Minor comments | Before requesting re-review |
| Questions | Same day if possible |

## GitHub Response Mechanics

### Inline Response
```bash
# Reply to specific comment thread
gh api graphql -f query='
  mutation {
    addPullRequestReviewComment(input: {
      pullRequestReviewId: "<REVIEW_ID>",
      body: "Your response here",
      inReplyTo: "<COMMENT_ID>"
    }) {
      comment { id }
    }
  }
'

# Or using REST API for simple reply
gh api repos/:owner/:repo/pulls/$PR_NUM/comments/<COMMENT_ID>/replies \
  -X POST -f body="Your response here"
```

### Resolve Conversation Thread

**IMPORTANT: Always resolve threads after addressing feedback.**

```bash
# Get thread ID from comment
THREAD_ID=$(gh api graphql -f query='
  query {
    repository(owner: "OWNER", name: "REPO") {
      pullRequest(number: PR_NUM) {
        reviewThreads(first: 100) {
          nodes {
            id
            isResolved
            comments(first: 1) {
              nodes { body }
            }
          }
        }
      }
    }
  }
' -q '.data.repository.pullRequest.reviewThreads.nodes[] | select(.isResolved == false) | .id')

# Resolve the thread
gh api graphql -f query='
  mutation {
    resolveReviewThread(input: {threadId: "'"$THREAD_ID"'"}) {
      thread { isResolved }
    }
  }
'
```

**Simplified workflow:**
1. Make the fix
2. Push the commit
3. Reply with "Fixed in [SHA]"
4. Resolve the thread in GitHub UI or via API
5. **Report progress to user**

### Progress Feedback Format

After each comment is addressed, report:

```
## Progress Update

### Comment [N/Total]: [Brief description]
- **Status**: [Accepted/Challenged/Clarified/Deferred]
- **Action**: [What was done]
- **Commit**: [SHA if applicable]
- **Thread**: [Resolved/Pending response]

### Remaining
- [X] comments to address
- [Y] threads unresolved
```

### Request Re-review
After addressing all feedback:
```bash
gh pr edit $PR_NUM --add-reviewer username
```

Or comment:
```markdown
@reviewer I've addressed all feedback. Ready for re-review when you have time.

Summary of changes:
- [Change 1]: Fixed in abc123
- [Change 2]: Fixed in def456
- [Change 3]: Discussed above, keeping current approach
```

## Multi-Round Reviews

### After First Round
```markdown
Thanks for the thorough review! I've addressed all comments:

| Comment | Resolution |
|---------|------------|
| [Issue 1] | Fixed in abc123 |
| [Issue 2] | Fixed in def456 |
| [Issue 3] | Discussed - keeping current |

Ready for re-review.
```

### After Subsequent Rounds
```markdown
Addressed the latest round of feedback:

- [New issue]: Fixed in ghi789

All previous comments remain resolved.
```

## Handling Difficult Situations

### Reviewer is Wrong

```markdown
I appreciate the attention to this area. I investigated and found [evidence that contradicts the feedback].

[Detailed evidence]

I may be missing something though - could you clarify [specific aspect] if I've misunderstood?
```

### Reviewer is Nitpicking

```markdown
Thanks for the attention to detail. For consistency, I'd prefer to follow the existing pattern in [similar code].

If you feel strongly about establishing a new standard, I'd suggest we discuss in [appropriate forum] and update the style guide.
```

### Conflicting Feedback from Multiple Reviewers

```markdown
I'm getting different suggestions from reviewers:
- @reviewer1 suggests [A]
- @reviewer2 suggests [B]

Could we align on the preferred approach? Happy to implement whichever you agree on.
```

### Reviewer Not Responding

```markdown
@reviewer Friendly ping on this PR. Let me know if you need any additional context or if the changes address your concerns.
```

## Progress Tracking

### During Review Response

Keep the user informed with structured updates:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
REVIEW RESPONSE PROGRESS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PR #[NUM]: [Title]
Comments: [Addressed]/[Total]

## Current: Comment [N]
Reviewer: @[username]
Category: [Bug/Architecture/Style/etc.]
Decision: [Accept/Challenge/Clarify]

[Action being taken...]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### After Each Comment

```
Comment [N] complete:
  Decision: Accepted
  Action: Fixed null check
  Commit: abc1234
  Thread: Resolved
```

### Summary Table

Maintain running summary:

| # | Reviewer | Category | Decision | Status |
|---|----------|----------|----------|--------|
| 1 | @alice | Bug | Accept | Resolved |
| 2 | @bob | Style | Challenge | Pending |
| 3 | @alice | Arch | Accept | Resolved |

## Merge Readiness Summary

When all comments are addressed, provide final summary:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR READY FOR MERGE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PR #[NUM]: [Title]
Branch: [branch] -> next

## Review Summary
- Total comments: [N]
- Accepted: [X]
- Challenged: [Y]
- Deferred: [Z]

## Changes Made
| Commit | Description |
|--------|-------------|
| abc123 | Fix null check (comment #1) |
| def456 | Add error handling (comment #3) |

## Threads Status
- Resolved: [N]/[Total]
- Pending: [List any waiting for reviewer]

## CI Status
- Build:
- Tests:
- Coverage:
- Lint:

## Checklist
- [x] All comments addressed
- [x] All threads resolved (or awaiting reviewer)
- [x] CI passing
- [x] No unresolved conversations
- [x] Re-review requested

## Next Steps
[Ready for merge / Waiting for reviewer response on X]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Related Skills

- `evaluate-change-request` - Analyse feedback first
- `prove-correctness` - Generate evidence
- `justify-decision` - Explain reasoning
- `british-english` - Proper spelling/grammar
- `pr-monitor` - Track review status
