---
description: Identify and document technical debt
agent: plan
subtask: true
---

Analyze and document technical debt in the specified area.

Load the `tech-debt` skill for debt assessment framework.

## Target
$ARGUMENTS

If no specific target, analyze recent changes or a key area.

## Process

1. **Scan for Debt Indicators**
   ```bash
   make check-patterns
   make check-intent-architecture
   ```

2. **Identify Issues**
   - Long functions (> 50 lines)
   - Duplicate code
   - Missing tests
   - Deprecated patterns
   - TODO/FIXME comments
   - Architecture violations

3. **Assess Each Item**
   - Impact (High/Medium/Low)
   - Effort to fix
   - Risk if left unfixed

4. **Prioritize**
   Use the quadrant:
   - Quick Wins (High value, Low effort) → Do now
   - Major Projects (High value, High effort) → Plan
   - Thankless Tasks (Low value, High effort) → Maybe later
   - Don't Bother (Low value, Low effort) → Skip

5. **Document**
   Create debt items for significant findings.

## Output

Report including:
- Summary of debt found
- Prioritized list
- Recommended actions
- Quick wins to tackle immediately
