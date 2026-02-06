---
name: pair-programming
description: Collaborate effectively through pairing - driver/navigator, mob programming, remote pairing
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide effective pair and mob programming - sharing knowledge, catching bugs early, and producing better code through collaboration.

## When to use me

- Onboarding new team members
- Tackling complex problems
- Knowledge sharing
- Code review alternative
- Unfamiliar territory

## Core Principles

1. **Two heads are better than one** - Different perspectives catch more issues
2. **Continuous review** - Problems caught as they're created
3. **Knowledge sharing** - Skills transfer naturally
4. **Engagement matters** - Both people must be actively involved
5. **It's tiring** - Take breaks, don't pair all day

## Pair Programming Styles

### Driver/Navigator

```
DRIVER (keyboard):
- Writes the code
- Focuses on syntax and implementation
- Thinks tactically
- Asks for direction when stuck

NAVIGATOR (thinking):
- Reviews code as it's written
- Thinks about the big picture
- Suggests improvements
- Catches bugs and typos
- Keeps track of next steps

Switch roles every 15-30 minutes.
```

### Ping Pong (TDD)

```
Person A: Write failing test
Person B: Make test pass, write next failing test
Person A: Make test pass, write next failing test
...continue...

Great for TDD discipline and engagement.
```

### Strong Style

```
"For an idea to go from your head into the computer,
it MUST go through someone else's hands."

Navigator: Has the idea, directs
Driver: Implements Navigator's instructions

Forces communication, great for teaching.
```

## Remote Pairing

### Tools

| Tool | Use For |
|------|---------|
| VS Code Live Share | Real-time code editing |
| Tuple | Purpose-built pairing app |
| Screen share + voice | Quick and simple |
| tmux + SSH | Terminal-based |

### Remote Pairing Tips

```markdown
## Setup Checklist

- [ ] Stable internet connection
- [ ] Good microphone (headset recommended)
- [ ] Quiet environment
- [ ] Second monitor helpful
- [ ] Shared access to tools (repo, tickets, docs)

## During Session

- Camera on (optional but helps engagement)
- Verbalise your thinking
- Point with mouse when referencing code
- Take breaks every 45-60 minutes
- Agree on switch signal
```

## Mob Programming

### Structure

```
One computer, whole team:

DRIVER (rotates every 5-15 min):
- Types what the mob decides
- Asks clarifying questions
- Doesn't make unilateral decisions

MOB:
- Discusses approach
- Reviews as code is written
- Suggests improvements
- Keeps driver on track

FACILITATOR (optional):
- Manages rotation timer
- Ensures everyone participates
- Keeps session focused
```

### When to Mob

- Kickoff of new feature (shared understanding)
- Complex problem (multiple perspectives)
- Learning new technology (everyone learns together)
- Critical code (many eyes)

## Effective Pairing Behaviours

### DO

```markdown
## Communication
- Think aloud - share your reasoning
- Ask questions - "What do you think about...?"
- Offer suggestions gently - "What if we tried...?"
- Acknowledge good ideas - "That's clever because..."

## Engagement
- Stay focused - no email/slack checking
- Take breaks - pairing is intense
- Switch roles regularly
- Both people should be learning

## Respect
- Be patient with different skill levels
- Explain your thinking when asked
- Accept that there are multiple valid approaches
- Give your pair time to think
```

### DON'T

```markdown
## Avoid
- Grabbing the keyboard without asking
- Dismissing ideas without consideration
- Checking phone/email while pairing
- Going too fast for your pair
- Being silent for long periods
- Pair programming on trivial tasks

## Red Flags
- One person doing all the work
- Navigator zoning out
- Arguing about style preferences
- Fatigue (pair too long)
```

## Pairing Scenarios

### Expert + Novice

```
Goal: Knowledge transfer

Expert as Navigator:
- Guide the novice driver
- Explain reasoning
- Let them make mistakes (learning opportunities)
- Resist urge to grab keyboard

Novice as Navigator:
- Ask questions constantly
- Don't pretend to understand
- Request explanations
```

### Expert + Expert

```
Goal: Better solution through collaboration

- Challenge each other's assumptions
- Consider alternatives before implementing
- Debate is healthy, but timebox
- Trust each other's expertise
```

### Novice + Novice

```
Goal: Learn together

- Research together
- Don't be afraid to try things
- Ask for help when truly stuck
- Celebrate small wins
```

## Pairing Anti-Patterns

| Anti-Pattern | Problem | Fix |
|--------------|---------|-----|
| **Backseat driver** | Navigator grabs keyboard | Agree on switch signal |
| **Disengaged navigator** | Navigator zones out | Shorter rotations, active role |
| **Steamroller** | One person dominates | Strong style pairing |
| **Code tour** | Driver explains existing code | Both read code silently first |
| **Marathon session** | Pairing for 8 hours | Limit to 4-6 hours max |
| **Forced pairing** | Pairing on everything | Pair on valuable tasks only |

## When NOT to Pair

- Simple, repetitive tasks
- When either person is exhausted
- Research/exploration (parallel then share)
- When one person just needs focus time
- Trivial bug fixes

## AI as Pair Partner

When pairing with AI (like this):

```
AI as Navigator:
- AI suggests approaches
- Human implements
- AI reviews and suggests improvements
- Human makes final decisions

AI as Reference:
- Human drives implementation
- AI answers questions
- AI provides code examples
- Human adapts and integrates
```

## Measuring Pairing Effectiveness

### Signs It's Working

- Both people can explain the code
- Fewer bugs reaching review/production
- Knowledge spreading across team
- People enjoy pairing sessions
- Onboarding is faster

### Signs It's Not Working

- One person always disengaged
- Code quality same or worse
- People avoid pairing
- Constant conflict
- Tasks take much longer

## Related Skills

- `code-reviewer` - Review skills transfer
- `tdd-workflow` - TDD works great with ping-pong
- `clean-code` - Shared standards
