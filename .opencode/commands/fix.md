---
description: Fix a bug following TDD with regression test
agent: build
---

Fix the following bug using TDD approach.

## Skills to Load (MANDATORY)

- `pre-action` - Decision framework before ANY action (always active)
- `tdd-workflow` - For writing regression test first
- `debug-test` - If tests are failing unexpectedly
- `clean-code` - For clean fix implementation

## Bug
$ARGUMENTS

## Process

1. **Pre-Action Assessment**
   - STOP: What is the reported bug?
   - THINK: What do I KNOW about this? What am I ASSUMING?
   - INVESTIGATE: Is this actually a bug, or intentional behaviour?
   - CONFIDENCE: Am I VERIFIED it's a bug?
   - If ASSUMED: investigate more or ASK before proceeding

2. **Reproduce** - Understand the bug and when it occurs

3. **Regression Test First** - Write a test that fails due to the bug

4. **Fix** - Implement the minimal fix to make the test pass

5. **Verify** - Ensure all tests pass and no regressions

6. **Cleanup** - Apply Boy Scout Rule to touched code

7. **Compliance** - Run `make check-compliance`

## Requirements

- **Verify it's actually a bug before fixing**
- The regression test must fail before the fix
- The regression test must pass after the fix
- No other tests should break
- Apply clean code principles to the fix
- >= 95% coverage maintained
