---
name: time-management
description: Manage time effectively - timeboxing, focus techniques, knowing when to stop, asking for help
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide effective time management during development - staying focused, timeboxing exploration, knowing when to stop, and asking for help at the right time.

## When to use me

- Starting a new task
- Getting stuck on a problem
- Feeling overwhelmed
- Deciding whether to continue or pivot
- Managing multiple priorities

## Core Principles

1. **Time is finite** - Choose where to spend it wisely
2. **Timeboxes create clarity** - Constraints force decisions
3. **Stuck is a signal** - Stop, assess, get help
4. **Focus is fragile** - Protect it actively
5. **Done beats perfect** - Ship and iterate

## Timeboxing

### What is Timeboxing?

```
TIMEBOX: Fixed time limit for an activity.
When time expires: STOP and DECIDE.

Benefits:
- Prevents rabbit holes
- Forces prioritisation
- Creates decision points
- Reduces perfectionism
```

### Standard Timeboxes

| Activity | Timebox | At Expiry |
|----------|---------|-----------|
| Quick investigation | 15-30 min | Have answer or escalate |
| Spike/prototype | 2-4 hours | Decide: proceed or pivot |
| Debugging session | 1-2 hours | Document findings, get help |
| Refactoring | 30-60 min | Commit or revert |
| Code review | 30-60 min | Finish or schedule follow-up |
| Meeting | As scheduled | End on time |

### Timebox Process

```markdown
## Starting a Timebox

1. SET the time limit
2. DEFINE what success looks like
3. START a timer
4. WORK focused until timer
5. STOP when timer ends
6. DECIDE: continue, pivot, or stop

## Example

TIMEBOX: 30 minutes
GOAL: Understand why test is flaky
SUCCESS: Root cause identified OR clear next step

[Work for 30 minutes]

RESULT: Found race condition in test setup
DECISION: Fix it (estimated 15 min)
```

### When Timebox Expires

```markdown
## Decision Tree

Timer ended. What now?

├── Goal achieved?
│   └── YES → Done! Move on.
│   
├── Clear path forward?
│   └── YES → Set new timebox, continue
│   
├── Stuck but making progress?
│   └── YES → One more timebox (max)
│   
├── Stuck, no progress?
│   └── NO → Stop. Get help.
│   
└── Wrong approach entirely?
    └── YES → Stop. Reassess.
```

## Focus Management

### Protecting Focus Time

```markdown
## Focus Blocks

1. SCHEDULE focus time (2-4 hour blocks)
2. ELIMINATE distractions:
   - Notifications off
   - Slack on DND
   - Email closed
   - Phone away
3. SIGNAL unavailability (status, headphones)
4. SINGLE-TASK (one thing only)
```

### The Focus Funnel

```
START OF DAY:
┌─────────────────────────────┐
│ Many possible tasks         │
└─────────────────────────────┘
            │
            ▼ Filter: What's most important?
┌─────────────────────────────┐
│ Important tasks             │
└─────────────────────────────┘
            │
            ▼ Filter: What can I do now?
┌─────────────────────────────┐
│ Actionable tasks            │
└─────────────────────────────┘
            │
            ▼ Pick ONE
┌─────────────────────────────┐
│ THE task                    │
└─────────────────────────────┘
            │
            ▼ FOCUS until done
```

### Handling Interruptions

```markdown
## When Interrupted

1. NOTE where you are (comment, todo)
2. ASSESS interruption urgency
   - Emergency? → Handle now
   - Important? → Schedule time
   - Can wait? → "I'll get back to you"
3. RETURN to task with minimal context loss
```

### Context Switching Cost

```
CONTEXT SWITCH = 15-30 minutes lost

Every switch:
- Lose flow state
- Need to rebuild mental model
- Risk forgetting where you were
- Increased error rate

MINIMISE switches:
- Batch similar tasks
- Protect focus blocks
- Say no to interruptions
- Finish before switching
```

## Knowing When to Stop

### Signs You Should Stop

```markdown
## Red Flags

- Same error for 30+ minutes
- Trying random things hoping they work
- Feeling frustrated or tired
- Making things worse
- Can't explain what you're doing
- "Just one more thing" repeatedly
- Working in circles
```

### The Stop Decision

```markdown
## When to Stop

STOP NOW:
- Tired and making mistakes
- Deadline for something else
- Stuck with no ideas
- Wrong approach confirmed

STOP SOON (finish current step):
- End of timebox
- Natural breaking point
- Need input from others

CONTINUE:
- Making steady progress
- Clear next step
- Energy and focus good
```

### Productive Stopping

```markdown
## Before Stopping

1. DOCUMENT current state
   - What works
   - What doesn't
   - What you tried
   - What to try next

2. COMMIT work in progress (WIP)
   - "WIP: investigating filter bug"
   - Even if incomplete

3. NOTE next step
   - First thing to do when returning

4. CLEAN UP
   - Close tabs
   - Revert experiments
   - Update task status
```

## Asking for Help

### When to Ask

```markdown
## Help Decision Matrix

                    |  Low Urgency  |  High Urgency
─────────────────────────────────────────────────────
Low Confidence     |  Research     |  Ask now
(don't know how)   |  first        |  
─────────────────────────────────────────────────────
High Confidence    |  Do it        |  Do it +
(know how)         |               |  sanity check
```

### The 15-Minute Rule

```
STUCK? Try for 15 minutes.
- Research
- Experiment
- Read docs/code

STILL STUCK after 15 minutes?
- Document what you tried
- Ask for help
- Don't waste hours on what someone could answer in minutes
```

### How to Ask Effectively

```markdown
## Good Help Request

CONTEXT: What are you trying to do?
"I'm adding date filter to the timeline view."

TRIED: What have you attempted?
"I tried adding filter state to the intent, but it resets on navigation."

SPECIFIC: What's the actual question?
"Where should filter state live to persist across screen changes?"

SHOW: Evidence/code/error
"Here's what I have: [code snippet]"
```

### Who to Ask

```markdown
## Escalation Path

1. RUBBER DUCK: Explain to yourself/AI
2. DOCUMENTATION: Official docs, README
3. TEAM CHAT: Quick questions
4. COLLEAGUE: Domain expert
5. MENTOR/LEAD: Architectural questions
6. EXTERNAL: Stack Overflow (with caution)
```

## Daily Rhythm

### Effective Day Structure

```markdown
## Sample Developer Day

MORNING (high energy)
- 30 min: Planning, email, messages
- 3 hours: Deep work (most important task)

MIDDAY (variable energy)
- Lunch
- 1 hour: Meetings, collaboration
- 30 min: Code review

AFTERNOON (medium energy)
- 2 hours: Secondary tasks
- 30 min: Admin, communication

END OF DAY
- 15 min: Document progress
- 15 min: Plan tomorrow
```

### Task Prioritisation

```markdown
## Priority Matrix

            URGENT          NOT URGENT
         ┌─────────────┬─────────────┐
IMPORTANT│ DO NOW      │ SCHEDULE    │
         │ (Crises)    │ (Planning)  │
         ├─────────────┼─────────────┤
NOT      │ DELEGATE    │ ELIMINATE   │
IMPORTANT│ (Interrupts)│ (Time waste)│
         └─────────────┴─────────────┘

Focus on: Important, regardless of urgency
```

## Common Time Traps

| Trap | Solution |
|------|----------|
| Perfectionism | Define "good enough", ship it |
| Rabbit holes | Timebox exploration |
| Over-researching | Time limit, then decide |
| Waiting for perfect info | Decide with 70% info |
| Saying yes to everything | Protect focus time |
| Not asking for help | 15-minute rule |
| Context switching | Batch similar tasks |

## Related Skills

- `scope-management` - Keeping work bounded
- `estimation` - Realistic time estimates
- `pragmatic-problem-solving` - Practical progress
- `checklist-discipline` - Tracking progress
