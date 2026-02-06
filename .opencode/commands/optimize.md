---
description: Optimize code performance using profiling and benchmarking
agent: performance
---

# /optimize - Performance Optimization

Optimize code performance using profiling and benchmarking.

## Skills to Load

- `performance` - Optimization strategies
- `benchmarking` - Measuring performance
- `go-expert` - Go-specific optimizations

## Usage

```
/optimize <what to optimize>
```

## Examples

```
/optimize timeline rendering is slow
/optimize event filtering takes too long with large datasets
/optimize memory usage in burst detection
/optimize concurrent repository access
```

## Process

### 1. MEASURE - Profile First

```bash
# CPU Profile
go test -bench=. -cpuprofile=cpu.prof ./path/to/package
go tool pprof cpu.prof
> top10
> list FunctionName

# Memory Profile
go test -bench=. -memprofile=mem.prof ./path/to/package
go tool pprof mem.prof
> top10 -cum
```

### 2. BENCHMARK - Establish Baseline

```bash
go test -bench=BenchmarkTarget -benchmem -count=5 ./... > baseline.txt
```

### 3. ANALYZE - Understand the Bottleneck

Common bottlenecks:
- Excessive allocations (check `allocs/op`)
- String concatenation in loops
- Unnecessary copies (large structs by value)
- N+1 queries
- Lock contention
- Reflection usage

### 4. OPTIMIZE - Targeted Fix

Make ONE change targeting the specific bottleneck.

### 5. VERIFY - Measure Again

```bash
go test -bench=BenchmarkTarget -benchmem -count=5 ./... > new.txt
benchstat baseline.txt new.txt
```

Expect to see improvement:
```
name          old time/op    new time/op    delta
Target-8        234ns ± 2%     156ns ± 1%   -33.33%
```

### 6. CHECK - No Regressions

```bash
go test ./...  # All tests still pass
```

## Common Optimizations

### Pre-allocate Slices

```go
// Before
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// After
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

### Use strings.Builder

```go
// Before
result := ""
for _, s := range strings {
    result += s
}

// After
var b strings.Builder
b.Grow(estimatedSize)
for _, s := range strings {
    b.WriteString(s)
}
result := b.String()
```

### Avoid Reflection

```go
// Before (slow)
reflect.ValueOf(obj).FieldByName("Field")

// After (fast)
obj.Field
```

### Database: Eager Load

```go
// Before (N+1 queries)
events, _ := repo.FindAll(ctx)
for _, e := range events {
    facts, _ := factRepo.FindByEventID(ctx, e.ID)
}

// After (1 query)
repo.db.Preload("Facts").Find(&events)
```

## Warning

**Never optimize without measuring first!**

Premature optimization is the root of all evil. Profile, benchmark, then optimize.
