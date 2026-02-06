---
name: benchmarking
description: Go benchmarking for measuring and optimising code performance
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Benchmarking Skill

You are an expert in Go benchmarking for measuring and optimizing code performance.

## Overview

Go has built-in benchmarking support via the `testing` package. Use benchmarks to measure performance, identify bottlenecks, and validate optimizations.

---

## Writing Benchmarks

### Basic Benchmark

```go
// event_benchmark_test.go
package career_test

import (
    "testing"
    "github.com/baphled/kariya/internal/domain/career"
)

func BenchmarkEventValidation(b *testing.B) {
    event := &career.Event{
        Text: "Led successful project delivery",
        Date: time.Now(),
    }
    
    b.ResetTimer()  // Exclude setup from timing
    
    for i := 0; i < b.N; i++ {
        event.Validate()
    }
}
```

### Benchmark with Setup

```go
func BenchmarkEventService_Create(b *testing.B) {
    // Setup - not timed
    repo := memory.NewEventRepository()
    service := career.NewService(repo)
    ctx := context.Background()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        b.StopTimer()  // Pause for per-iteration setup
        text := fmt.Sprintf("Event %d", i)
        b.StartTimer()
        
        _, _ = service.Create(ctx, text, time.Now())
    }
}
```

### Table-Driven Benchmarks

```go
func BenchmarkEventValidation(b *testing.B) {
    cases := []struct {
        name  string
        event *career.Event
    }{
        {"small_text", &career.Event{Text: "Small"}},
        {"medium_text", &career.Event{Text: strings.Repeat("x", 500)}},
        {"large_text", &career.Event{Text: strings.Repeat("x", 2000)}},
    }
    
    for _, tc := range cases {
        b.Run(tc.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                tc.event.Validate()
            }
        })
    }
}
```

### Memory Benchmarks

```go
func BenchmarkEventCreation(b *testing.B) {
    b.ReportAllocs()  // Report memory allocations
    
    for i := 0; i < b.N; i++ {
        _ = &career.Event{
            ID:   fmt.Sprintf("event-%d", i),
            Text: "Test event",
            Date: time.Now(),
        }
    }
}
```

---

## Running Benchmarks

### Basic Commands

```bash
# Run all benchmarks in package
go test -bench=. ./internal/domain/career/...

# Run specific benchmark
go test -bench=BenchmarkEventValidation ./...

# Run with memory stats
go test -bench=. -benchmem ./...

# Run multiple times for stability
go test -bench=. -count=5 ./...

# Set minimum run time
go test -bench=. -benchtime=5s ./...

# Run specific iterations
go test -bench=. -benchtime=1000x ./...
```

### Output Format

```
BenchmarkEventValidation-8    5000000    234 ns/op    48 B/op    2 allocs/op
│                         │   │          │            │          │
│                         │   │          │            │          └─ allocations per op
│                         │   │          │            └─ bytes allocated per op
│                         │   │          └─ nanoseconds per operation
│                         │   └─ number of iterations
│                         └─ GOMAXPROCS
└─ benchmark name
```

### Comparing Results

```bash
# Save baseline
go test -bench=. -benchmem ./... > old.txt

# Make changes, then compare
go test -bench=. -benchmem ./... > new.txt

# Compare with benchstat
benchstat old.txt new.txt
```

**Install benchstat:**
```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

**Output:**
```
name                old time/op    new time/op    delta
EventValidation-8     234ns ± 2%     156ns ± 1%   -33.33%  (p=0.008 n=5+5)

name                old alloc/op   new alloc/op   delta
EventValidation-8     48.0B ± 0%     32.0B ± 0%   -33.33%  (p=0.008 n=5+5)
```

---

## Benchmark Patterns

### Sub-benchmarks for Comparison

```go
func BenchmarkFind(b *testing.B) {
    items := generateItems(10000)
    target := items[5000]
    
    b.Run("linear", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            linearSearch(items, target)
        }
    })
    
    b.Run("binary", func(b *testing.B) {
        sorted := sortItems(items)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            binarySearch(sorted, target)
        }
    })
    
    b.Run("map", func(b *testing.B) {
        m := buildMap(items)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = m[target.ID]
        }
    })
}
```

### Scaling Benchmarks

```go
func BenchmarkEventFilter(b *testing.B) {
    for _, size := range []int{10, 100, 1000, 10000} {
        events := generateEvents(size)
        
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                filterByCompany(events, "Acme Corp")
            }
        })
    }
}
```

### Parallel Benchmarks

```go
func BenchmarkConcurrentAccess(b *testing.B) {
    repo := memory.NewEventRepository()
    ctx := context.Background()
    
    // Pre-populate
    for i := 0; i < 1000; i++ {
        repo.Save(ctx, fixtures.Event(i))
    }
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, _ = repo.FindByID(ctx, "event-500")
        }
    })
}
```

---

## Repository Benchmarks

```go
func BenchmarkEventRepository(b *testing.B) {
    b.Run("SQL", func(b *testing.B) {
        db := setupTestDB(b)
        repo := sql.NewEventRepository(db)
        benchmarkRepository(b, repo)
    })
    
    b.Run("Memory", func(b *testing.B) {
        repo := memory.NewEventRepository()
        benchmarkRepository(b, repo)
    })
}

func benchmarkRepository(b *testing.B, repo career.EventRepository) {
    ctx := context.Background()
    
    b.Run("Save", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            event := fixtures.Event(i)
            _ = repo.Save(ctx, event)
        }
    })
    
    b.Run("FindByID", func(b *testing.B) {
        // Setup
        event := fixtures.Event(1)
        repo.Save(ctx, event)
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = repo.FindByID(ctx, event.ID)
        }
    })
    
    b.Run("FindAll", func(b *testing.B) {
        // Setup
        for i := 0; i < 100; i++ {
            repo.Save(ctx, fixtures.Event(i))
        }
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = repo.FindAll(ctx)
        }
    })
}
```

---

## UI/TUI Benchmarks

```go
func BenchmarkTableRender(b *testing.B) {
    theme := themes.NewDefaultTheme()
    items := fixtures.Events(100)
    
    table := behaviors.NewTableBehavior[*career.Event](
        theme,
        columns,
        eventFormatter,
    )
    table.SetItems(items)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = table.View()
    }
}

func BenchmarkScreenUpdate(b *testing.B) {
    screen := NewEventListScreen(fixtures.Events(100))
    keyMsg := tea.KeyMsg{Type: tea.KeyDown}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = screen.Update(keyMsg)
    }
}
```

---

## Best Practices

### DO: Reset Timer After Setup

```go
func BenchmarkWithSetup(b *testing.B) {
    // Expensive setup
    data := loadLargeDataset()
    
    b.ResetTimer()  // Don't count setup
    
    for i := 0; i < b.N; i++ {
        process(data)
    }
}
```

### DO: Use b.ReportAllocs()

```go
func BenchmarkAllocs(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        // Code that allocates
    }
}
```

### DO: Run Multiple Times

```bash
# Run 10 times for statistical significance
go test -bench=. -count=10 ./...
```

### DON'T: Benchmark with I/O

```go
// WRONG - I/O dominates, not useful
func BenchmarkWithIO(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Println("output")  // I/O is slow and variable
    }
}

// RIGHT - Benchmark computation only
func BenchmarkComputation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        compute()
    }
}
```

### DON'T: Let Compiler Optimize Away

```go
// WRONG - Compiler might optimize away unused result
func BenchmarkOptimizedAway(b *testing.B) {
    for i := 0; i < b.N; i++ {
        compute()  // Result unused, might be optimized out
    }
}

// RIGHT - Use result to prevent optimization
var result int

func BenchmarkCorrect(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = compute()
    }
    result = r  // Use result
}
```

---

## Profiling Integration

```bash
# CPU profile
go test -bench=. -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profile
go test -bench=. -memprofile=mem.prof ./...
go tool pprof mem.prof

# Block profile (contention)
go test -bench=. -blockprofile=block.prof ./...
go tool pprof block.prof
```

---

## CI Integration

```yaml
# .github/workflows/benchmark.yml
- name: Run benchmarks
  run: go test -bench=. -benchmem ./... | tee bench.txt

- name: Compare benchmarks
  uses: benchmark-action/github-action-benchmark@v1
  with:
    tool: 'go'
    output-file-path: bench.txt
```

---

## Related Skills

- `performance` - Overall performance optimization strategies
- `ginkgo-gomega` - Testing framework context
- `gorm-repository` - Database benchmarks
- `bubble-tea-expert` - TUI rendering benchmarks
