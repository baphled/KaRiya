---
description: Create a development task with acceptance criteria
agent: plan
---

Create a development task for the described work.

Load the `create-task` skill for task structure.

## Task
$ARGUMENTS

## Process

1. **Analyze the Work**
   - What needs to be done?
   - Which files/components are affected?
   - What are the dependencies?

2. **Create Task**
   ```bash
   make new-feature TASK="brief description"
   ```

3. **Define Acceptance Criteria**
   - Specific, measurable outcomes
   - Edge cases to handle
   - Error scenarios

4. **Technical Guidance**
   - Files to modify
   - Patterns to use (run `make what-to-use NEED="..."`)
   - Testing requirements

5. **Estimate Size**
   - XS: < 1 hour
   - S: 1-4 hours
   - M: 4-8 hours
   - L: 1-2 days
   - XL: Break down further

## Output

Return the task file path and a summary of:
- Acceptance criteria
- Key technical decisions
- Estimated size
