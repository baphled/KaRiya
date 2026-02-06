---
description: Create and run benchmarks to measure code performance
agent: benchmark
---

# /benchmark - Measure Performance

Create and run benchmarks to measure code performance.

## Skills to Load

- `benchmarking` - Go benchmark patterns
- `performance` - Performance optimization strategies

## Usage

```
/benchmark <what to benchmark>
```

## Examples

```
/benchmark event validation
/benchmark table rendering with 1000 items
/benchmark repository FindAll query
/benchmark concurrent event creation
```

## Process

### 1. Create Benchmark File

```go
// internal/domain/career/event_benchmark_test.go
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
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        event.Validate()
    }
}
```

### 2. Run Baseline

```bash
go test -bench=BenchmarkEventValidation -benchmem ./internal/domain/career/...
```

Output:
```
BenchmarkEventValidation-8    5000000    234 ns/op    48 B/op    2 allocs/op
```

### 3. Save Baseline (If Optimizing)

```bash
go test -bench=. -benchmem -count=5 ./... > baseline.txt
```

### 4. Make Changes

Implement optimization...

### 5. Compare

```bash
go test -bench=. -benchmem -count=5 ./... > new.txt
benchstat baseline.txt new.txt
```

## Common Benchmark Patterns

### Table-Driven

```go
func BenchmarkValidation(b *testing.B) {
    cases := []struct{
        name string
        size int
    }{
        {"small", 10},
        {"medium", 100},
        {"large", 1000},
    }
    
    for _, tc := range cases {
        b.Run(tc.name, func(b *testing.B) {
            data := generateData(tc.size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                validate(data)
            }
        })
    }
}
```

### With Memory

```go
func BenchmarkAllocs(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = NewEvent("test", time.Now())
    }
}
```

### Parallel

```go
func BenchmarkConcurrent(b *testing.B) {
    repo := setupRepo()
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            repo.FindByID(ctx, "test-id")
        }
    })
}
```

## Useful Commands

```bash
# Run all benchmarks
go test -bench=. ./...

# With memory stats
go test -bench=. -benchmem ./...

# Run 5 times for stability
go test -bench=. -count=5 ./...

# Longer run time
go test -bench=. -benchtime=5s ./...

# CPU profile
go test -bench=. -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profile
go test -bench=. -memprofile=mem.prof ./...
go tool pprof mem.prof
```

## Install benchstat

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```
