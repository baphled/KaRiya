---
description: Run comprehensive compliance and quality checks
agent: build
---

Run all quality and compliance checks.

Load the `check-compliance` skill for check details.

## Checks to Run

1. **Full Compliance**
   ```bash
   make check-compliance
   ```

2. **Architecture Validation**
   ```bash
   make check-intent-architecture
   ```

3. **Pattern Enforcement**
   ```bash
   make check-patterns
   ```

4. **Security Scan**
   ```bash
   make gosec
   ```

5. **Test Suite**
   ```bash
   make test
   ```

6. **Coverage** (for modified packages)
   ```bash
   go test -cover ./path/to/modified/...
   ```

## Report

Summarize:
- Total checks run
- Passed checks
- Failed checks with details
- Suggested fixes for any failures

## Focus Area
$ARGUMENTS
