---
description: Start a KaRiya development session with environment validation
agent: build
---

Start a new KaRiya development session.

First, load the `session-start` skill to understand the session protocol.

Then run:
```bash
make session-start
```

If it fails, help me fix the issues before proceeding.

Acknowledge the critical rules:
1. Feature branches only (never commit to next/main)
2. TDD workflow (test first)
3. Use `make ai-commit` for commits
4. Run `make check-compliance` before and after tasks

$ARGUMENTS
