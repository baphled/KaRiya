---
description: Prepare and create a properly attributed commit
agent: build
---

Prepare and create a commit for the current changes.

Load these skills:
- `ai-commit` - Commit format and attribution
- `code-reviewer` - Pre-commit review

## Process

1. **Review Changes**
   ```bash
   git status
   git diff --cached
   ```
   If nothing staged, help identify what should be committed.

2. **Pre-commit Checks**
   ```bash
   make check-compliance
   ```
   Fix any issues before proceeding.

3. **Generate Commit Message**
   Based on the changes, create a commit message file:
   ```bash
   cat > /tmp/commit.txt << 'EOF'
   type(scope): short description

   Explanation of WHY the change was made.
   EOF
   ```

   Types: feat, fix, docs, refactor, test, chore
   Scopes: domain, service, cli, intents, screens, uikit, forms

4. **Create Commit**
   ```bash
   make ai-commit FILE=/tmp/commit.txt
   ```

5. **Verify**
   ```bash
   git log -1
   ```

## Additional Context
$ARGUMENTS
