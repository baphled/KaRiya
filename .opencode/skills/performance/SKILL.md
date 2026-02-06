# Performance Skill

You are an expert in Go performance optimization, profiling, and writing efficient code.

## Overview

Performance optimization in Go requires measurement, understanding bottlenecks, and targeted improvements. Never optimize without profiling first.

---

## The Performance Process

```
1. MEASURE (Profile)
   └─> Identify actual bottlenecks
   
2. ANALYZE (Understand)
   └─> Why is it slow?
   
3. OPTIMIZE (Targeted)
   └─> Fix specific bottleneck
   
4. VERIFY (Benchmark)
   └─> Confirm improvement
   └─> Check for regressions
```

**Golden Rule:** Measure first, optimize second. Premature optimization is the root of all evil.

---

## Profiling

### CPU Profiling

```go
import "runtime/pprof"

// In tests
go test -cpuprofile=cpu.prof -bench=. ./...

// Analyze
go tool pprof cpu.prof
> top10          # Top 10 functions by CPU
> list FuncName  # Line-by-line breakdown
> web            # Visual graph (needs graphviz)
```

### Memory Profiling

```go
// In tests
go test -memprofile=mem.prof -bench=. ./...

// Analyze
go tool pprof mem.prof
> top10 -cum     # By cumulative allocations
> list FuncName  # Line-by-line allocations
```

### Runtime Profiling

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    // Enable pprof endpoint
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // Your app...
}
```

```bash
# Profile running application
go tool pprof http://localhost:6060/debug/pprof/profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

### Trace

```bash
go test -trace=trace.out ./...
go tool trace trace.out
```

---

## Common Bottlenecks

### 1. Memory Allocations

**Problem:** Excessive allocations cause GC pressure

```go
// SLOW - Allocates on every call
func ProcessItems(items []Item) []Result {
    var results []Result
    for _, item := range items {
        results = append(results, process(item))  // Grows, reallocates
    }
    return results
}

// FAST - Pre-allocate
func ProcessItems(items []Item) []Result {
    results := make([]Result, 0, len(items))  // Pre-allocate capacity
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}
```

### 2. String Concatenation

```go
// SLOW - Creates new string each time
func BuildQuery(parts []string) string {
    result := ""
    for _, p := range parts {
        result += p + ", "  // O(n²) allocations
    }
    return result
}

// FAST - Use strings.Builder
func BuildQuery(parts []string) string {
    var b strings.Builder
    b.Grow(estimatedSize)  // Pre-allocate if size known
    for i, p := range parts {
        if i > 0 {
            b.WriteString(", ")
        }
        b.WriteString(p)
    }
    return b.String()
}
```

### 3. Unnecessary Copies

```go
// SLOW - Copies entire struct
func (e Event) Process() {  // Value receiver copies
    // ...
}

// FAST - Pointer avoids copy (for large structs)
func (e *Event) Process() {  // Pointer receiver
    // ...
}

// Note: Small structs (< 64 bytes) may be faster by value
```

### 4. Map Lookups in Loops

```go
// SLOW - Map lookup each iteration
for _, item := range items {
    if categories[item.Category] {  // Map lookup
        process(item)
    }
}

// FAST - Cache lookup result
categorySet := make(map[string]bool)
for _, c := range targetCategories {
    categorySet[c] = true
}
for _, item := range items {
    if categorySet[item.Category] {
        process(item)
    }
}
```

### 5. Reflection

```go
// SLOW - Reflection is expensive
func GetField(obj interface{}, name string) interface{} {
    return reflect.ValueOf(obj).FieldByName(name).Interface()
}

// FAST - Direct access
func GetField(obj *Event, name string) interface{} {
    switch name {
    case "Text":
        return obj.Text
    case "Date":
        return obj.Date
    }
    return nil
}
```

### 6. Interface Conversions

```go
// SLOW - Type assertion in hot path
func Process(items []interface{}) {
    for _, item := range items {
        e := item.(*Event)  // Type assertion
        process(e)
    }
}

// FAST - Use concrete type
func Process(items []*Event) {
    for _, item := range items {
        process(item)
    }
}
```

---

## TUI Performance

### Efficient Rendering

```go
// SLOW - Rebuild entire view
func (s *Screen) View() string {
    var b strings.Builder
    for _, item := range s.items {  // All items every render
        b.WriteString(renderItem(item))
    }
    return b.String()
}

// FAST - Render only visible items
func (s *Screen) View() string {
    var b strings.Builder
    start, end := s.visibleRange()
    for i := start; i < end; i++ {
        b.WriteString(renderItem(s.items[i]))
    }
    return b.String()
}
```

### Cache Computed Values

```go
type TableBehavior struct {
    items       []*Event
    cachedView  string
    viewDirty   bool
}

func (t *TableBehavior) View() string {
    if !t.viewDirty && t.cachedView != "" {
        return t.cachedView
    }
    t.cachedView = t.render()
    t.viewDirty = false
    return t.cachedView
}

func (t *TableBehavior) SetItems(items []*Event) {
    t.items = items
    t.viewDirty = true  // Invalidate cache
}
```

### Debounce Updates

```go
// Debounce rapid updates
type Screen struct {
    updateTimer *time.Timer
    pendingMsg  tea.Msg
}

func (s *Screen) Update(msg tea.Msg) tea.Cmd {
    if s.updateTimer != nil {
        s.updateTimer.Stop()
    }
    
    s.pendingMsg = msg
    s.updateTimer = time.AfterFunc(16*time.Millisecond, func() {
        s.processUpdate(s.pendingMsg)
    })
    
    return nil
}
```

---

## Database Performance

### N+1 Query Problem

```go
// SLOW - N+1 queries
events, _ := repo.FindAll(ctx)
for _, e := range events {
    facts, _ := factRepo.FindByEventID(ctx, e.ID)  // Query per event!
    e.Facts = facts
}

// FAST - Eager loading with GORM
events, _ := repo.db.Preload("Facts").Find(&events)

// Or batch query
eventIDs := extractIDs(events)
facts, _ := factRepo.FindByEventIDs(ctx, eventIDs)
factsByEvent := groupByEventID(facts)
for _, e := range events {
    e.Facts = factsByEvent[e.ID]
}
```

### Index Usage

```go
// Ensure indexes exist for query patterns
// migrations/002_add_indexes.sql
CREATE INDEX idx_events_date ON events(date);
CREATE INDEX idx_events_company ON events(company);
CREATE INDEX idx_events_company_date ON events(company, date);  // Composite
```

### Query Optimization

```go
// SLOW - Fetching all columns
repo.db.Find(&events)

// FAST - Select only needed columns
repo.db.Select("id", "title", "date").Find(&events)

// SLOW - Loading all then filtering
events, _ := repo.FindAll(ctx)
filtered := filter(events, criteria)

// FAST - Filter in database
repo.db.Where("company = ?", company).Find(&events)
```

---

## Concurrency Performance

### Worker Pool

```go
func ProcessConcurrently(items []*Event, workers int) []*Result {
    jobs := make(chan *Event, len(items))
    results := make(chan *Result, len(items))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                results <- process(item)
            }
        }()
    }
    
    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Collect results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var output []*Result
    for r := range results {
        output = append(output, r)
    }
    return output
}
```

### Sync.Pool for Reuse

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func Render() string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    // Use buf...
    return buf.String()
}
```

### Reduce Lock Contention

```go
// SLOW - Single lock for all operations
type Cache struct {
    mu    sync.RWMutex
    items map[string]*Item
}

// FAST - Sharded locks
type ShardedCache struct {
    shards [256]*cacheShard
}

type cacheShard struct {
    mu    sync.RWMutex
    items map[string]*Item
}

func (c *ShardedCache) Get(key string) *Item {
    shard := c.shards[hash(key)%256]
    shard.mu.RLock()
    defer shard.mu.RUnlock()
    return shard.items[key]
}
```

---

## Memory Optimization

### Reduce Struct Size

```go
// INEFFICIENT - Poor alignment (40 bytes)
type Event struct {
    Active    bool      // 1 byte + 7 padding
    ID        int64     // 8 bytes
    Flag      bool      // 1 byte + 7 padding
    Timestamp int64     // 8 bytes
    Count     int32     // 4 bytes + 4 padding
}

// EFFICIENT - Better alignment (24 bytes)
type Event struct {
    ID        int64     // 8 bytes
    Timestamp int64     // 8 bytes
    Count     int32     // 4 bytes
    Active    bool      // 1 byte
    Flag      bool      // 1 byte + 2 padding
}
```

### Slice vs Array

```go
// Use arrays for fixed small sizes (stack allocated)
type Point struct {
    Coords [3]float64  // Fixed size, no pointer
}

// Use slices for variable sizes
type Polygon struct {
    Points []Point  // Variable size
}
```

---

## Performance Checklist

Before optimizing:
- [ ] Have I profiled to identify the actual bottleneck?
- [ ] Is this code in a hot path?
- [ ] Will the optimization provide measurable improvement?

During optimization:
- [ ] Am I measuring with benchmarks?
- [ ] Am I keeping the code readable?
- [ ] Do I have tests to prevent regressions?

After optimization:
- [ ] Did benchmarks confirm improvement?
- [ ] Did I check for regressions elsewhere?
- [ ] Is the code still maintainable?

---

## Anti-Patterns

### DON'T: Optimize Without Measuring

```go
// WRONG - Assuming this is faster
func FasterMaybe() {
    // Complex "optimization" without proof
}

// RIGHT - Measure first
// 1. Profile to find bottleneck
// 2. Benchmark current implementation
// 3. Optimize
// 4. Benchmark again to verify
```

### DON'T: Sacrifice Readability

```go
// WRONG - Unreadable for marginal gain
func Process(d []byte) int {
    return int(d[0])<<24 | int(d[1])<<16 | int(d[2])<<8 | int(d[3])
}

// RIGHT - Readable (use binary package)
func Process(d []byte) int {
    return int(binary.BigEndian.Uint32(d))
}
```

### DON'T: Premature Optimization

```go
// WRONG - Optimizing before it's a problem
func GetName(user *User) string {
    // Using unsafe pointer tricks for "performance"
    // when simple return user.Name works fine
}

// RIGHT - Simple first, optimize if needed
func GetName(user *User) string {
    return user.Name
}
```

---

## Related Skills

- `benchmarking` - Measuring performance
- `gorm-repository` - Database performance
- `bubble-tea-expert` - TUI rendering performance
- `go-expert` - Go idioms for performance
