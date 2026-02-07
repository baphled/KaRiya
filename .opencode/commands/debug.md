---
description: Debug failing tests or unexpected behavior
agent: build
---

Debug the issue described.

## Skills to Load (MANDATORY)

- `pre-action` - Decision framework before ANY action (always active)
- `debug-test` - Test debugging techniques
- `clean-code` - Understanding expected behavior

## Issue
$ARGUMENTS

## Process

1. **Pre-Action Assessment**
   - STOP: What is failing?
   - THINK: Is this a bug in the code, or is the test wrong?
   - INVESTIGATE: When did this start? What changed?
   - CONFIDENCE: Do I understand the root cause?
   - **ASK if uncertain whether test or code is wrong**

2. **Reproduce the Issue**
   ```bash
   make test
   # or specific test
   make individual-test TEST="TestName"
   ```

3. **Gather Information**
   - What's the error message?
   - When did it start failing? (`git log`)
   - What changed recently? (`git diff`, `git show`)
   - Was the change intentional?

4. **Diagnose**
   - Check for race conditions: `go test -race ./...`
   - Check for nil pointers
   - Check for off-by-one errors
   - Check fixture usage

5. **Common Issues**
   - Race conditions → add synchronization
   - Multiple Ginkgo suites → one suite per package
   - Flaky tests → check for time/order dependencies
   - Coverage issues → find uncovered paths
   - **New validation breaking old tests → verify if validation is intentional**

6. **Fix and Verify**
   - Confirm you're fixing the right thing (test vs code)
   - Write a test that captures the bug (if one doesn't exist)
   - Fix the issue
   - Run full test suite

7. **Compliance Check**
   ```bash
   make check-compliance
   ```
