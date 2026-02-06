---
description: Implement a feature following TDD and clean code principles
agent: build
---

Implement the following feature using TDD and clean code principles.

Load these skills for guidance:
- `tdd-workflow` - For the Red-Green-Refactor cycle
- `clean-code` - For writing clean, maintainable code
- `architecture` - For proper layer placement
- `component-lookup` - Find correct components to use
- `go-expert` - For idiomatic Go patterns

## Task
$ARGUMENTS

## Process

1. **Understand** - Analyze where this code should live (which layer, package)
2. **Plan** - Break into small, testable increments
3. **TDD Cycle** for each increment:
   - Write a failing test (Red)
   - Write minimal code to pass (Green)
   - Refactor for cleanliness
4. **Boy Scout Rule** - Leave any touched code cleaner than you found it
5. **Verify** - Run `make check-compliance`

## Requirements

- Follow the architecture patterns (intents orchestrate, screens render)
- Use existing components from `make what-to-use NEED="keyword"`
- No comments inside function bodies
- Meaningful names that reveal intent
- Functions should do one thing
- >= 95% test coverage for new code
