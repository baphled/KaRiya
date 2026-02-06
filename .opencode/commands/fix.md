---
description: Fix a bug following TDD with regression test
agent: build
---

Fix the following bug using TDD approach.

Load these skills:
- `tdd-workflow` - For writing regression test first
- `debug-test` - If tests are failing unexpectedly
- `clean-code` - For clean fix implementation

## Bug
$ARGUMENTS

## Process

1. **Reproduce** - Understand the bug and when it occurs
2. **Regression Test First** - Write a test that fails due to the bug
3. **Fix** - Implement the minimal fix to make the test pass
4. **Verify** - Ensure all tests pass and no regressions
5. **Cleanup** - Apply Boy Scout Rule to touched code
6. **Compliance** - Run `make check-compliance`

## Requirements

- The regression test must fail before the fix
- The regression test must pass after the fix
- No other tests should break
- Apply clean code principles to the fix
- >= 95% coverage maintained
