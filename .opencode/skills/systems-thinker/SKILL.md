---
name: systems-thinker
description: Understand complex systems, interconnections, and emergent behaviors for better architectural decisions
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Understand complex systems, see interconnections, anticipate emergent behaviors, and make better architectural decisions.

## When to use me

Use this skill when:
- Making architectural decisions
- Debugging complex issues
- Understanding unexpected behaviors
- Planning system changes
- Evaluating impacts of changes
- Designing new features

## Systems Thinking Principles

1. **Everything is connected** - Changes ripple through systems
2. **The whole differs from parts** - Emergent properties matter
3. **Cause and effect are not linear** - Feedback loops exist
4. **Today's solutions are tomorrow's problems** - Trade-offs persist
5. **There are no side effects, only effects** - Unintended doesn't mean unimportant

## Systems Concepts

### Feedback Loops

**Reinforcing (Positive) Loops** - Amplify change:
```
More bugs → More firefighting → Less prevention → More bugs
```

**Balancing (Negative) Loops** - Resist change:
```
More features → More complexity → Slower development → Fewer features
```

### Emergence

System behaviors that arise from component interactions:

```
Individual components: Intent, Screen, Behavior, UIKit
Emergent behavior: Consistent user experience, maintainability
```

### Boundaries

Where does the system end? What's inside/outside?

```
┌─────────────────────────────────────────────┐
│ KaRiya System                               │
│  ┌─────────────────────────────────────┐    │
│  │ CLI Application                      │    │
│  │  ┌──────┐ ┌──────┐ ┌──────┐        │    │
│  │  │Intent│→│Screen│→│UIKit │        │    │
│  │  └──────┘ └──────┘ └──────┘        │    │
│  └─────────────────────────────────────┘    │
│                    ↓                         │
│  ┌─────────────────────────────────────┐    │
│  │ Service Layer                        │    │
│  └─────────────────────────────────────┘    │
│                    ↓                         │
│  ┌─────────────────────────────────────┐    │
│  │ Repository Layer                     │    │
│  └─────────────────────────────────────┘    │
│                    ↓                         │
│  ┌─────────────────────────────────────┐    │
│  │ SQLite Database                      │    │ ← Boundary
│  └─────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
                     ↓
        ┌────────────────────┐
        │ External: File System │  ← Outside boundary
        └────────────────────┘
```

### Leverage Points

Where small changes have big effects:

| Leverage | Example | Impact |
|----------|---------|--------|
| High | Architecture patterns | Affects everything |
| High | Interface definitions | All implementations |
| Medium | Component design | Feature area |
| Low | Implementation details | Local only |

## Systems Analysis Tools

### Mapping the System

```markdown
## System Map: Timeline Feature

### Components
- TimelineIntent (orchestration)
- ListScreen (display)
- DetailModal (drill-down)
- EventService (business logic)
- EventRepository (persistence)

### Connections
- Intent → Screen: Update/View cycle
- Screen → Intent: ScreenResult messages
- Intent → Service: Data operations
- Service → Repository: Persistence

### Data Flows
- User input → Intent → Screen → View
- Load events: Intent → Service → Repository → Service → Intent → Screen

### Control Flows
- Navigation: Router → Intent → Screen
- State changes: Intent manages, Screen reflects
```

### Identifying Dependencies

```markdown
## Dependency Analysis: Adding Skill Filter

### Direct Dependencies
- SkillService (need to query skills)
- EventService (need to filter events)
- FilterModal (UI component)

### Indirect Dependencies
- SkillRepository (via SkillService)
- Event-Skill junction (data relationship)

### Affected By
- Skill changes (need to refresh)
- Event changes (filter results change)

### Affects
- Timeline display (filtered results)
- Export functionality (filtered data)
- Statistics (filtered calculations)
```

### Impact Analysis

Before making changes, trace impacts:

```markdown
## Impact Analysis: Change Event Date Field Type

### Primary Impact
- Event struct changes
- Migration required
- Repository methods change

### Secondary Impact
- Service layer adapts
- Serialization changes
- Test fixtures update

### Tertiary Impact
- API responses change
- Import/export format changes
- Reports recalculate

### Unknown Impacts
- External integrations?
- User expectations?
- Performance characteristics?
```

## Systems Patterns in KaRiya

### Layer Architecture as System

```
Higher layers depend on lower layers.
Lower layers are more stable.
Changes flow down more easily than up.

Intent (high, volatile)
   ↓
Screen (medium)
   ↓
UIKit (low, stable)
   ↓
Theme (lowest, most stable)
```

### Message-Based Communication

```
Components communicate via messages, not direct calls.
This creates loose coupling but requires message coordination.

Intent ←→ Screen via ScreenResult
Intent ←→ Service via Commands/Results
```

### State Machine Pattern

```
Intents are state machines.
State determines behavior.
Transitions are explicit.

StateList → StateDetail → StateEdit → StateList
    ↑           ↓
    ←←←←← Cancel ←←←←←
```

## Systems Thinking Questions

### Understanding
- What are the components?
- How do they interact?
- Where are the boundaries?
- What emerges from interactions?

### Causality
- What causes what?
- Are there feedback loops?
- What are the delays?
- What's the root cause vs. symptom?

### Change
- Where are the leverage points?
- What are the ripple effects?
- What might resist change?
- What unintended effects might occur?

### Stability
- What keeps the system stable?
- What could destabilize it?
- Where are the single points of failure?
- What are the recovery mechanisms?

## Avoiding Systems Pitfalls

### Local Optimization
Optimizing one part at expense of whole:
```
❌ Make this screen super fast (adds complexity everywhere else)
✅ Keep consistent patterns (overall system simpler)
```

### Ignoring Delays
Effects aren't immediate:
```
❌ "Tests slow us down" (ignoring future bug cost)
✅ "Tests slow us now, speed us later" (see full loop)
```

### Shifting the Burden
Treating symptoms instead of causes:
```
❌ Add more error handling (treating symptoms)
✅ Fix why errors occur (treating cause)
```

### Eroding Goals
Lowering standards instead of meeting them:
```
❌ "Coverage is hard, let's lower threshold"
✅ "Coverage is hard, let's improve tests"
```

## Related skills

- `architecture` - System structure
- `critical-thinking` - Analyze interactions
- `devils-advocate` - Challenge system assumptions
- `research` - Understand existing system
- `tech-debt` - System health over time
