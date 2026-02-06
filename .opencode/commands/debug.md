---
description: Debug failing tests or unexpected behavior
agent: build
---

Debug the issue described.

Load these skills:
- `debug-test` - Test debugging techniques
- `clean-code` - Understanding expected behavior

## Issue
$ARGUMENTS

## Process

1. **Reproduce the Issue**
   ```bash
   make test
   # or specific test
   make individual-test TEST="TestName"
   ```

2. **Gather Information**
   - What's the error message?
   - When did it start failing?
   - What changed recently?

3. **Diagnose**
   - Check for race conditions: `go test -race ./...`
   - Check for nil pointers
   - Check for off-by-one errors
   - Check fixture usage

4. **Common Issues**
   - Race conditions → add synchronization
   - Multiple Ginkgo suites → one suite per package
   - Flaky tests → check for time/order dependencies
   - Coverage issues → find uncovered paths

5. **Fix and Verify**
   - Write a test that captures the bug (if one doesn't exist)
   - Fix the issue
   - Run full test suite

6. **Compliance Check**
   ```bash
   make check-compliance
   ```
