---
name: prove-correctness
description: Write tests and provide evidence to prove or disprove claims about code
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Provide concrete evidence for code behaviour through tests, benchmarks, and demonstrations. I help settle disagreements with facts rather than opinions.

## When to use me

Use this skill when:
- Challenging reviewer feedback with evidence
- Proving current behaviour is correct
- Demonstrating a claimed bug doesn't exist
- Comparing performance of alternatives

## Core Principle

**Claims require evidence. Write code that proves your point.**

Instead of arguing, demonstrate:
- Write a test that passes with current code
- Write a benchmark that shows performance
- Create a minimal reproduction that proves behaviour

## Evidence Types

### 1. Unit Tests (Behaviour Proof)

Write tests that demonstrate correct behaviour:

```go
func TestBurstDetection_HandlesEmptyInput(t *testing.T) {
    // Reviewer claimed: "This will panic on empty input"
    // Proof: It handles empty input gracefully
    
    detector := NewBurstDetector()
    result, err := detector.Detect([]Event{})
    
    assert.NoError(t, err)
    assert.Empty(t, result.Bursts)
}
```

### 2. Table Tests (Edge Case Coverage)

Prove multiple scenarios:

```go
func TestCalculateScore_EdgeCases(t *testing.T) {
    // Reviewer claimed: "Edge cases aren't handled"
    // Proof: All edge cases are covered
    
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"zero", 0, 0},
        {"negative", -1, 0},
        {"max int", math.MaxInt, 100},
        {"typical", 50, 50},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateScore(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3. Benchmarks (Performance Proof)

Prove performance claims:

```go
func BenchmarkLookup_Slice(b *testing.B) {
    // Reviewer suggested: "Use a map for O(1) lookup"
    // Counter: With 10 items, slice is actually faster
    
    items := make([]string, 10)
    for i := range items {
        items[i] = fmt.Sprintf("item%d", i)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, item := range items {
            if item == "item5" {
                break
            }
        }
    }
}

func BenchmarkLookup_Map(b *testing.B) {
    items := make(map[string]bool, 10)
    for i := 0; i < 10; i++ {
        items[fmt.Sprintf("item%d", i)] = true
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = items["item5"]
    }
}
```

Run and report:
```bash
go test -bench=BenchmarkLookup -benchmem ./...

# Example output to share:
# BenchmarkLookup_Slice-8    50000000    25.3 ns/op    0 B/op    0 allocs/op
# BenchmarkLookup_Map-8      30000000    42.1 ns/op    0 B/op    0 allocs/op
```

### 4. Integration Tests (System Behaviour)

Prove end-to-end behaviour:

```go
func TestWorkflow_CompletesSuccessfully(t *testing.T) {
    // Reviewer claimed: "This workflow will fail if X happens"
    // Proof: Workflow handles X correctly
    
    env := testenv.New(t)
    env.SetupScenario("X happens")
    
    err := env.RunWorkflow()
    
    assert.NoError(t, err)
    assert.Equal(t, "completed", env.WorkflowState())
}
```

### 5. Reproduction Tests (Bug Verification)

Prove a claimed bug exists (or doesn't):

```go
func TestClaimed_BugDoesNotExist(t *testing.T) {
    // Reviewer reported: "Clicking X causes Y"
    // Investigation: Cannot reproduce
    
    env := testenv.New(t)
    
    // Exact steps from bug report
    env.NavigateTo("X")
    env.Click()
    
    // Expected incorrect behaviour (per report)
    // Actual: Correct behaviour
    assert.NotEqual(t, "Y", env.CurrentState())
    assert.Equal(t, "expected_state", env.CurrentState())
}
```

## Proof Workflow

### When Challenging Feedback

1. **Understand the claim**
   - What exactly is the reviewer claiming?
   - What would "correct" look like to them?

2. **Write a test for their claim**
   - If the test passes, they may be right
   - If the test fails, you have evidence

3. **Share the evidence**
   ```markdown
   I wrote a test to verify this concern:
   
   ```go
   func TestConcern(t *testing.T) {
       // [test code]
   }
   ```
   
   The test passes, demonstrating that [behaviour].
   Let me know if I've misunderstood the concern.
   ```

### When Defending Performance

1. **Write benchmarks for both approaches**
2. **Run multiple times for consistency**
3. **Include memory allocation data**
4. **Share raw numbers**

```markdown
I benchmarked both approaches:

| Approach | ns/op | B/op | allocs/op |
|----------|-------|------|-----------|
| Current (slice) | 25.3 | 0 | 0 |
| Suggested (map) | 42.1 | 0 | 0 |

With our typical data size (~10 items), the slice is ~40% faster.
The O(1) vs O(n) difference only matters at larger scales.
```

### When Proving Architecture

Reference documentation and demonstrate:

```markdown
Our architecture guide specifies this pattern:

> "Domain logic should not depend on infrastructure" 
> - docs/ARCHITECTURE.md#layer-rules

The current implementation follows this:
- Domain layer: `internal/domain/` (no infrastructure imports)
- Service layer: `internal/service/` (orchestrates domain + infra)

Here's the import graph proving no violations:
```bash
go list -f '{{.ImportPath}}: {{.Imports}}' ./internal/domain/...
```
```

## Test Organisation

Place proof tests appropriately:

| Proof Type | Location | Naming |
|------------|----------|--------|
| Bug regression | Same package as code | `TestBugXXX_Description` |
| Behaviour proof | Same package as code | `TestFeature_ProvesBehaviour` |
| Benchmark | Same package as code | `BenchmarkFeature_Comparison` |
| Integration proof | `tests/integration/` | `TestIntegration_Scenario` |

## Evidence Quality Checklist

Before sharing evidence:

- [ ] Test is deterministic (no flaky results)
- [ ] Test name clearly states what's being proved
- [ ] Test covers the exact claim being made
- [ ] Benchmarks run multiple iterations
- [ ] Results are reproducible by reviewer
- [ ] Evidence directly addresses the feedback

## Common Proof Patterns

### "This will panic"
```go
func TestNoPanic_OnEdgeCase(t *testing.T) {
    assert.NotPanics(t, func() {
        FunctionUnderTest(edgeCaseInput)
    })
}
```

### "This is a race condition"
```bash
go test -race ./...
```

### "This leaks goroutines"
```go
func TestNoGoroutineLeak(t *testing.T) {
    before := runtime.NumGoroutine()
    
    // Run the code
    DoSomething()
    
    // Allow cleanup
    time.Sleep(100 * time.Millisecond)
    
    after := runtime.NumGoroutine()
    assert.LessOrEqual(t, after, before+1)
}
```

### "This allocates too much"
```go
func TestAllocation(t *testing.T) {
    result := testing.Benchmark(func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            DoSomething()
        }
    })
    
    // Should allocate 0 bytes
    assert.Equal(t, int64(0), result.AllocedBytesPerOp())
}
```

## Integration with TDD

When writing proof tests:

1. **RED**: Write test that would fail if reviewer is correct
2. **GREEN**: If test passes, you have proof. If fails, investigate.
3. **REFACTOR**: Clean up test for clarity

## Related Skills

- `evaluate-change-request` - Determines if proof is needed
- `tdd-workflow` - Test writing methodology
- `ginkgo-gomega` - BDD test framework
- `critical-thinking` - Analyse claims rigorously
- `respond-to-review` - Share evidence effectively
