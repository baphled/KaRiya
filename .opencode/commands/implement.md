---
description: Implement a feature following TDD and clean code principles
agent: build
---

Implement the following feature using TDD and clean code principles.

## Skills to Load (MANDATORY)

- `pre-action` - Decision framework before ANY action (always active)
- `tdd-workflow` - For the Red-Green-Refactor cycle
- `clean-code` - For writing clean, maintainable code
- `architecture` - For proper layer placement
- `component-lookup` - Find correct components to use
- `go-expert` - For idiomatic Go patterns

## Task
$ARGUMENTS

## Process

1. **Pre-Action Assessment**
   - STOP: What exactly is being requested?
   - THINK: Do I understand the requirement? What's ASSUMED?
   - INVESTIGATE: How do similar features work in this codebase?
   - CONFIDENCE: Am I ready to implement, or need clarification?

2. **Understand** - Analyze where this code should live (which layer, package)

3. **Plan** - Break into small, testable increments

4. **TDD Cycle** for each increment:
   - Write a failing test (Red)
   - Write minimal code to pass (Green)
   - Refactor for cleanliness

5. **Boy Scout Rule** - Leave any touched code cleaner than you found it

6. **Verify** - Run `make check-compliance`

## Requirements

- **Apply pre-action framework before each change**
- Follow the architecture patterns (intents orchestrate, screens render)
- Use existing components from `make what-to-use NEED="keyword"`
- No comments inside function bodies
- Meaningful names that reveal intent
- Functions should do one thing
- >= 95% test coverage for new code
