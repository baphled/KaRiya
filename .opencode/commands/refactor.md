---
description: Refactor code following clean code and Boy Scout Rule
agent: build
---

Refactor the specified code following clean code principles.

Load these skills:
- `refactor` - Systematic refactoring process
- `clean-code` - Refactoring principles and Boy Scout Rule
- `code-reviewer` - What to look for
- `architecture` - Ensure refactoring respects layers

## Target
$ARGUMENTS

## Process

1. **Ensure Tests Exist**
   Refactoring requires tests as a safety net. If tests are missing, write them first.

2. **Run Tests**
   ```bash
   make test
   ```
   All tests must pass before refactoring.

3. **Identify Improvements**
   Look for:
   - Long functions (extract methods)
   - Duplicate code (extract to shared)
   - Magic numbers (named constants)
   - Poor names (rename for clarity)
   - Deep nesting (early returns)
   - Dead code (remove)
   - Comments in functions (extract to named methods)

4. **Refactor Incrementally**
   - Small changes
   - Run tests after each change
   - Commit working states

5. **Apply Boy Scout Rule**
   Leave the code cleaner than you found it.

6. **Verify**
   ```bash
   make test
   make check-compliance
   ```
   All tests must still pass.

## Constraints

- NO behavior changes (tests must pass unchanged)
- NO new features
- Focus on structure and clarity
- Maintain or improve test coverage
