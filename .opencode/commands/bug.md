---
description: Create a bug report for an issue
agent: plan
---

Create a bug report for the described issue.

Load the `create-bug` skill for bug report structure.

## Issue
$ARGUMENTS

## Process

1. **Gather Information**
   - What is the expected behavior?
   - What is the actual behavior?
   - What are the steps to reproduce?

2. **Create Bug Report**
   ```bash
   make new-bug BUG="brief description"
   ```

3. **Fill in Details**
   - Summary
   - Steps to reproduce
   - Expected vs actual behavior
   - Severity assessment
   - Environment info

4. **Initial Investigation** (if possible)
   - Identify affected component
   - Check for related issues
   - Note suspected cause

## Output

Return the bug file path and a summary of the bug report.
