---
name: pre-action
description: Mandatory decision-making process before any action - stop, think, choose the right approach
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Skill: pre-action

## What I do

Enforce a mandatory pause before taking any significant action. Every decision - whether to edit code, run a command, load a skill, or respond - must pass through this framework first.

## When to use me

**ALWAYS. Before ANY action.** This skill is non-negotiable.

## The Problem This Solves

Without this framework, AI agents:
- React to symptoms instead of understanding problems
- Use the first approach that comes to mind
- Skip investigation and jump to solutions
- Act on assumptions instead of verified understanding
- Miss opportunities to use the right skill or tool

## MANDATORY Pre-Action Framework

Before taking ANY significant action:

```
PRE-ACTION FRAMEWORK
====================

1. STOP - What am I being asked to do?
   Request: [what the user wants]
   My interpretation: [how I understand it]
   
2. THINK - Do I understand the situation?
   What do I KNOW? [verified facts]
   What do I ASSUME? [unverified beliefs]
   What is UNKNOWN? [gaps in understanding]
   
3. INVESTIGATE - What do I need to find out?
   Questions to answer: [list]
   How to answer them: [tools, skills, resources]
   
4. CHOOSE - What's the right approach?
   Skills needed: [which skills apply]
   Tools needed: [which tools to use]
   Alternatives considered: [other approaches]
   Why this approach: [reasoning]

5. CONFIDENCE - Am I ready to act?
   Level: VERIFIED / SUPPORTED / ASSUMED / UNKNOWN
   If < SUPPORTED: investigate more or ASK

6. ACT or ASK
   If confident: proceed with chosen approach
   If uncertain: ask for clarification
```

## Applying to Different Situations

### When Tests Fail

```
1. STOP - Tests are failing
2. THINK - 
   KNOW: Test output shows specific error
   ASSUME: The test is correct and code is wrong (DANGER!)
   UNKNOWN: Why this validation exists, whether it's intentional
3. INVESTIGATE -
   - When was this code added? (git log)
   - Why was it added? (commit message, PR)
   - Is there a BDD scenario for this? (features/)
4. CHOOSE -
   Skills: debug-test, question-resolver
   Approach: Investigate before changing anything
5. CONFIDENCE - UNKNOWN (haven't investigated yet)
6. ACT - Investigate first, then ASK user before changing
```

### When Asked to Implement Something

```
1. STOP - User wants feature X
2. THINK -
   KNOW: User's stated requirement
   ASSUME: I understand what they really want
   UNKNOWN: Edge cases, constraints, existing patterns
3. INVESTIGATE -
   - How do similar features work in this codebase?
   - What patterns should I follow?
   - Are there existing components to reuse?
4. CHOOSE -
   Skills: software-engineer, architecture, tdd-workflow
   Approach: Research first, then TDD
5. CONFIDENCE - Need to research before implementing
6. ACT - Load skills, investigate codebase, then propose approach
```

### When Facing Uncertainty

```
1. STOP - I'm not sure about X
2. THINK -
   KNOW: [what I'm certain of]
   ASSUME: [what I think but haven't verified]
   UNKNOWN: [what I need to find out]
3. INVESTIGATE -
   Load: question-resolver, epistemic-rigor
   Method: [how to verify]
4. CHOOSE -
   Approach: Verify before proceeding
5. CONFIDENCE - ASSUMED (need to verify)
6. ASK - If I can't verify, ask the user
```

## Decision Tree

```
          [Receive request/observe situation]
                        │
                        ▼
              [Do I understand it?]
               /              \
             NO               YES
             │                 │
             ▼                 ▼
        [Investigate]    [Do I know HOW?]
             │            /          \
             │          NO           YES
             │           │             │
             │           ▼             ▼
             │    [Load skills]   [VERIFIED?]
             │           │         /      \
             │           │       NO       YES
             │           │        │         │
             │           │        ▼         ▼
             └───────────┴──→ [ASK user] [ACT]
```

## Skills to Consider

Before acting, consider which skills apply:

| Situation | Skills to Load |
|-----------|----------------|
| Uncertainty | `question-resolver`, `epistemic-rigor` |
| Code changes | `tdd-workflow`, `clean-code`, `architecture` |
| Debugging | `debug-test`, `prove-correctness` |
| Investigation | `research`, `code-reading` |
| Decision needed | `critical-thinking`, `trade-off-analysis` |
| Complex problem | `systems-thinker`, `pragmatic-problem-solving` |

## Tools to Consider

| Need | Tool |
|------|------|
| Understand code | Read, Grep, Glob |
| Understand history | git log, git show, git blame |
| Verify behaviour | Bash (run tests), Read (check output) |
| Research | WebFetch, Task (explore agent) |
| Think through | sequential-thinking |

## Red Flags - STOP Immediately

If you notice yourself:
- Editing code without understanding why it's there
- "Fixing" something without verifying it's broken
- Acting on the first interpretation
- Skipping investigation because "it's obvious"
- Not considering alternative explanations

**STOP. Go back to step 1.**

## Examples

### BAD: Skipping the framework

```
User: Tests are failing
Agent: [Immediately edits tests to make them pass]

VIOLATION: No investigation, assumed tests were wrong
```

### GOOD: Following the framework

```
User: Tests are failing
Agent: 
  STOP - Tests failing with validation error
  THINK - I ASSUME the test data is wrong, but UNKNOWN if 
          the validation is intentional
  INVESTIGATE - Check git history for when validation was added
  [Finds it was added intentionally with BDD scenario]
  CONFIDENCE - Now VERIFIED that validation is intentional
  ASK - "The validation was added intentionally. Should I 
        update the test data to comply, or reconsider the 
        validation threshold?"
```

## Integration

This skill is the foundation. It activates BEFORE other skills:

1. `pre-action` - Decide what to do
2. Then load appropriate skills based on decision
3. Then use appropriate tools based on skills

## Related Skills

- `epistemic-rigor` - Confidence calibration
- `assumption-tracker` - Track verified vs assumed
- `question-resolver` - Systematic investigation
- `critical-thinking` - Evaluate approaches
- `software-engineer` - Orchestrate technical work
