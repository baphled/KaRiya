---
description: Run housekeeping and maintenance tasks on the codebase
agent: build
---

Run maintenance and housekeeping on the codebase.

Load the `housekeeping` skill for maintenance tasks.
Load the `clean-code` skill for cleanup standards.

## Focus Area
$ARGUMENTS

If no specific area, run general maintenance.

## Process

1. **Quick Hygiene**
   ```bash
   make fmt
   make pre-commit
   ```

2. **Find Issues**
   ```bash
   # Dead code
   staticcheck -unused ./...
   
   # TODOs
   grep -rn "TODO\|FIXME" internal/
   
   # Debug code
   grep -rn "fmt.Print" internal/cli/ | grep -v "_test.go"
   ```

3. **Dependency Check**
   ```bash
   go mod tidy
   go list -m -u all | head -20
   ```

4. **Test Health**
   ```bash
   make test
   make coverage
   ```

5. **Documentation**
   ```bash
   make check-docblocks
   ```

6. **Git Cleanup**
   ```bash
   git branch --merged | grep -v "main\|next"
   ```

## Output

Report:
- Issues found
- Actions taken
- Remaining items (if any)
- Recommendations
