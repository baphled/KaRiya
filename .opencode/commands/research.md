---
description: Research and understand a codebase area, pattern, or technology
agent: plan
subtask: true
---

Conduct systematic research to understand the topic.

Load the `research` skill for investigation framework.
Load the `critical-thinking` skill to evaluate findings.

## Topic
$ARGUMENTS

## Process

1. **Define the Question**
   - What specifically do I need to understand?
   - What will I do with this knowledge?

2. **Identify Sources**
   - Code (primary source)
   - Tests (expected behavior)
   - Git history (context)
   - Similar implementations (patterns)

3. **Systematic Exploration**
   ```bash
   # Find relevant files
   find internal/cli -name "*keyword*" -type f
   
   # Search for patterns
   grep -rn "pattern" internal/cli/
   
   # Check history
   git log --oneline -20 -- path/to/area/
   ```

4. **Synthesize Findings**
   - How does it work?
   - What patterns are used?
   - What are the key components?
   - What are the dependencies?

5. **Evaluate Critically**
   - What assumptions am I making?
   - What evidence supports my conclusions?
   - What's still unclear?

## Output

Provide:
- Summary of findings
- Key components and their roles
- Patterns identified
- Recommendations for the task at hand
- Open questions (if any)
