---
description: Run security audit on code
agent: plan
subtask: true
---

Perform a security audit.

Load the `security` skill for security checklist.

## Target
$ARGUMENTS

If no specific target, scan the entire codebase.

## Checks

1. **Automated Security Scan**
   ```bash
   make gosec
   ```

2. **Manual Review**
   Check for:
   - SQL injection (string concatenation in queries)
   - Path traversal (unsanitized file paths)
   - Hardcoded secrets (API keys, passwords)
   - Missing input validation
   - Error messages leaking internals
   - Missing timeouts on external calls

3. **Sensitive Files**
   Check that `.gitignore` covers:
   - `.env` files
   - `*.pem`, `*.key` files
   - `credentials.json`
   - `config.local.*`

4. **Dependencies**
   ```bash
   go list -m -json all | jq -r '.Path'
   ```
   Check for known vulnerabilities.

## Report

For each finding:
- Severity: CRITICAL / HIGH / MEDIUM / LOW
- Location: file:line
- Issue: description
- Fix: recommended solution
