---
description: Clean up code applying Boy Scout Rule
agent: build
---

Apply the Boy Scout Rule to clean up the specified area.

Load the `clean-code` skill for cleanup guidelines.

## Target
$ARGUMENTS

## Cleanup Tasks

1. **Find Issues**
   - Dead code (unused functions, variables)
   - Commented-out code
   - TODO/FIXME comments (complete or create issues)
   - Magic numbers without named constants
   - Poor variable/function names
   - Duplicate code
   - Long functions that should be split
   - Comments inside function bodies

2. **Prioritize**
   - Quick wins first (dead code, magic numbers)
   - Then structural improvements
   - Document any tech debt that can't be fixed now

3. **Clean Up**
   - Remove dead code
   - Replace magic numbers with constants
   - Rename unclear identifiers
   - Extract duplicate code
   - Split long functions
   - Replace comments with better names

4. **Verify**
   ```bash
   make test
   make check-compliance
   ```
   All tests must pass.

5. **Commit**
   ```bash
   make ai-commit FILE=/tmp/commit.txt
   ```
   Use type `refactor` or `chore` for cleanup commits.

## Constraints

- NO behavior changes
- Tests must continue to pass
- One logical cleanup per commit
