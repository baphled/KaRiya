---
description: Fix architecture violations detected by check-intent-architecture
agent: build
---

Fix architecture violations in the codebase.

Load these skills:
- `fix-architecture` - Diagnose and fix violations
- `architecture` - Layer patterns and rules
- `check-compliance` - Validate fixes

## Target
$ARGUMENTS

If no specific target, run full architecture check.

## Process

1. **Identify Violations**
   ```bash
   make check-intent-architecture
   ```

2. **Diagnose Each Issue**
   - What's the specific violation?
   - What's the correct pattern?

3. **Fix Incrementally**
   - One violation at a time
   - Run check after each fix
   - Ensure tests pass

4. **Common Fixes**
   - Untyped state → Add typed enum
   - Missing BaseIntent → Embed *BaseIntent
   - Modal in intent → Move to screens/{feature}/modals/
   - huh import → Use forms package
   - Intent too large → Extract to handlers.go, helpers.go

5. **Verify**
   ```bash
   make check-intent-architecture
   make test
   make check-compliance
   ```

## Output

Report:
- Violations found
- Fixes applied
- Verification status
