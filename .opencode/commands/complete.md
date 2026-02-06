---
description: Verify a task is truly complete with no loose ends
agent: build
---

Verify the task is truly complete.

Load the `task-completer` skill for completion checklist.

## Task to Verify
$ARGUMENTS

## Process

1. **Check Acceptance Criteria**
   - Review original requirements
   - Verify each criterion explicitly

2. **Run All Checks**
   ```bash
   make check-compliance
   make check-intent-architecture
   go test -cover ./path/to/modified/...
   ```

3. **Self-Review**
   ```bash
   git diff
   git status
   ```
   - Would I approve this PR?
   - Is anything unclear?
   - Did I take shortcuts?

4. **Find Loose Ends**
   - Debug code left?
   - TODOs remaining?
   - Commented code?
   - Missing docs?

5. **Clean Up**
   - Remove any loose ends
   - Apply Boy Scout Rule
   - Final verification

6. **Commit**
   - Create descriptive message
   - Use `make ai-commit`

## Output

Report:
- Criteria status (met/not met)
- Checks passed/failed
- Loose ends found and fixed
- Ready to commit (yes/no)
