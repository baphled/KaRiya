---
name: parallel-execution
description: Maximise efficiency by running independent tasks in parallel - tool calls, tests, builds, investigations
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Identify opportunities to run tasks in parallel and execute them simultaneously. This dramatically reduces total time for multi-step operations.

## When to use me

**Always.** Look for parallelisation opportunities in every task.

## Core Principle

```
SEQUENTIAL (slow):
Task A (2s) → Task B (2s) → Task C (2s) = 6 seconds

PARALLEL (fast):
Task A (2s) ─┐
Task B (2s) ─┼─ = 2 seconds
Task C (2s) ─┘
```

**If tasks are independent, run them in parallel.**

## Identifying Parallel Opportunities

### Independent Tasks (CAN parallelise)

```markdown
Tasks are independent when:
- No task needs output from another
- No shared mutable state
- Order doesn't matter
- No resource conflicts

EXAMPLES:
- Reading multiple files
- Running tests in different packages
- Searching in different directories
- Fetching from multiple APIs
- Checking multiple conditions
```

### Dependent Tasks (MUST sequence)

```markdown
Tasks are dependent when:
- One needs the result of another
- Order matters for correctness
- Shared state must be consistent

EXAMPLES:
- Write file → Read file
- Create branch → Commit to branch
- Build → Test (needs build output)
- Migrate DB → Seed data
```

## Tool Call Parallelisation

### Multiple Reads - PARALLEL

```markdown
NEED: Read 5 files to understand a feature

WRONG (sequential):
1. Read file A
2. Read file B
3. Read file C
4. Read file D
5. Read file E

RIGHT (parallel):
1. Read files A, B, C, D, E simultaneously
```

### Multiple Searches - PARALLEL

```markdown
NEED: Find usages of a function and its tests

WRONG (sequential):
1. Grep for function usage
2. Grep for test files

RIGHT (parallel):
1. Grep for function usage AND grep for tests simultaneously
```

### Multiple Independent Commands - PARALLEL

```markdown
NEED: Check git status, list files, check branch

WRONG (sequential):
1. git status
2. ls -la
3. git branch

RIGHT (parallel):
1. Run all three simultaneously
```

### Dependent Operations - SEQUENTIAL

```markdown
NEED: Create file and then read it back

MUST BE SEQUENTIAL:
1. Write the file
2. THEN read it back (needs file to exist)

Use && for dependencies:
git add . && git commit -m "message"
```

## Test Parallelisation

### Package-Level Parallelism

```bash
# Go runs packages in parallel by default
go test ./...

# Control parallelism
go test -p 4 ./...  # 4 packages at a time

# Run specific packages in parallel
go test ./pkg1/... ./pkg2/... ./pkg3/...
```

### Test-Level Parallelism

```go
func TestParallel(t *testing.T) {
    t.Parallel()  // Mark test as safe for parallel execution
    // ...
}

// In Ginkgo
var _ = Describe("Feature", func() {
    It("test 1", func() {
        // Runs in parallel with other Its
    })
    
    It("test 2", func() {
        // Runs in parallel with other Its
    })
})
```

### What NOT to Parallelise in Tests

```markdown
DON'T parallelise tests that:
- Share database state
- Use same ports
- Modify global variables
- Depend on execution order
- Use shared file system resources
```

## Investigation Parallelisation

### Exploring a Bug

```markdown
WRONG (sequential investigation):
1. Check logs
2. Then check recent commits  
3. Then check related tests
4. Then check configuration

RIGHT (parallel investigation):
1. Simultaneously:
   - Check logs
   - Check recent commits
   - Check related tests
   - Check configuration
2. Synthesise findings
```

### Understanding a Feature

```markdown
WRONG (sequential):
1. Read the intent file
2. Read the screen file
3. Read the tests
4. Read the service

RIGHT (parallel):
1. Read all four files simultaneously
2. Build mental model from combined information
```

## Build Parallelisation

### Go Build

```bash
# Go compiles packages in parallel by default
go build ./...

# Explicit parallelism
go build -p 8 ./...  # 8 parallel compilations
```

### Make Targets

```makefile
# Sequential
all: build test lint

# Parallel (if independent)
all:
	$(MAKE) build &
	$(MAKE) lint &
	wait
	$(MAKE) test  # Test after build completes
```

## Parallel Patterns

### Fan-Out Pattern

```markdown
ONE input → MANY parallel operations → COMBINE results

Example: Search across codebase
1. Start search in /internal
2. Start search in /cmd
3. Start search in /pkg
4. Combine all results
```

### Scatter-Gather Pattern

```markdown
1. SCATTER: Dispatch multiple parallel tasks
2. GATHER: Wait for all to complete
3. PROCESS: Use combined results

Example: Health check
1. Check database (parallel)
2. Check cache (parallel)
3. Check external API (parallel)
4. Return combined status
```

### Pipeline with Parallel Stages

```markdown
Stage 1: [A, B, C] in parallel
    ↓
Stage 2: [D, E] in parallel (needs stage 1)
    ↓
Stage 3: [F] (needs stage 2)

Example: CI Pipeline
1. [Lint, Type Check, Security Scan] - parallel
2. [Unit Tests, Integration Tests] - parallel, after stage 1
3. [Deploy] - after stage 2
```

## Common Parallel Opportunities

| Scenario | Parallel Tasks |
|----------|----------------|
| Starting investigation | Read files, grep, git log, check tests |
| Code review | Read changed files, check tests, run linter |
| Understanding feature | Read intent, screens, services, tests |
| Pre-commit checks | Lint, format check, type check |
| Debugging | Check logs, recent commits, config, state |
| Research | Search docs, search code, search issues |

## Anti-Patterns

### False Parallelism

```markdown
WRONG: Running dependent tasks "in parallel"
- Write file in one call
- Read same file in parallel call
- Result: Race condition, read may fail

RIGHT: Sequence dependent operations
- Write file
- THEN read file
```

### Over-Parallelisation

```markdown
WRONG: 50 parallel API calls
- Overwhelms the server
- May cause rate limiting
- Harder to debug

RIGHT: Controlled parallelism
- Batch into groups of 5-10
- Or use worker pool pattern
```

## Decision Framework

```
For each set of tasks, ask:

1. Are they independent?
   YES → Can parallelise
   NO  → Must sequence

2. Any shared resources?
   YES → Careful coordination needed
   NO  → Safe to parallelise

3. Does order matter?
   YES → Sequence in correct order
   NO  → Parallelise freely

4. Resource constraints?
   YES → Limit parallelism
   NO  → Maximise parallelism
```

## Integration with Workflow

### At Task Start

```markdown
When beginning any task:
1. List all information needed
2. Identify independent reads/searches
3. Execute all independent operations in parallel
4. Process combined results
```

### During Implementation

```markdown
When making changes:
1. Identify independent verifications
2. Run checks in parallel where possible
3. Sequence only what must be sequenced
```

### At Task End

```markdown
When completing:
1. Parallel: git status + git diff + check tests
2. Sequential: git add → git commit (dependent)
3. Parallel: Push + update task tracker (independent)
```

## Related Skills

- `time-management` - Efficient use of time
- `tdd-workflow` - Test execution patterns  
- `concurrency` - Go concurrency patterns
- `devops` - CI/CD parallelisation
