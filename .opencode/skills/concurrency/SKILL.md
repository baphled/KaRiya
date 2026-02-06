---
name: concurrency
description: Write safe, efficient concurrent Go code - goroutines, channels, sync primitives, and common patterns
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide writing safe concurrent Go code. Concurrency is powerful but dangerous - this skill helps you avoid common pitfalls.

## When to use me

- Writing code with goroutines
- Using channels for communication
- Synchronizing access to shared state
- Debugging race conditions
- Optimizing parallel operations

## Golden Rules

1. **Don't communicate by sharing memory; share memory by communicating**
2. **If in doubt, don't use concurrency**
3. **Always know how goroutines will end**
4. **Race detector is your friend: `go test -race`**

## Goroutine Basics

### Starting Goroutines

```go
// Basic goroutine
go func() {
    // work
}()

// With parameters (capture by value)
go func(id int) {
    fmt.Println(id)
}(i)  // Pass i, don't capture

// WRONG - captures loop variable
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // Race! All print same value
    }()
}

// RIGHT - pass as parameter
for i := 0; i < 10; i++ {
    go func(id int) {
        fmt.Println(id)  // Correct
    }(i)
}
```

### Goroutine Lifecycle

**Always ensure goroutines can exit:**

```go
// WRONG - goroutine leak
func startWorker() {
    go func() {
        for {
            // work forever - can never stop!
        }
    }()
}

// RIGHT - use context for cancellation
func startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return  // Clean exit
            default:
                // work
            }
        }
    }()
}
```

## Channels

### Channel Types

```go
ch := make(chan int)      // Unbuffered - blocks until receive
ch := make(chan int, 10)  // Buffered - blocks when full
ch := make(chan int, 0)   // Same as unbuffered
```

### Channel Operations

```go
// Send (blocks if unbuffered or buffer full)
ch <- value

// Receive (blocks until value available)
value := <-ch

// Receive with ok (check if closed)
value, ok := <-ch
if !ok {
    // channel closed
}

// Close (only sender should close)
close(ch)

// Range over channel (exits when closed)
for value := range ch {
    // process value
}
```

### Channel Direction

```go
// Send-only
func producer(ch chan<- int) {
    ch <- 42
}

// Receive-only
func consumer(ch <-chan int) {
    value := <-ch
}
```

## Common Patterns

### Worker Pool

```go
func workerPool(jobs <-chan Job, results chan<- Result, workers int) {
    var wg sync.WaitGroup
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    
    // Wait for all workers, then close results
    go func() {
        wg.Wait()
        close(results)
    }()
}

// Usage
jobs := make(chan Job, 100)
results := make(chan Result, 100)

workerPool(jobs, results, 5)

// Send jobs
for _, job := range allJobs {
    jobs <- job
}
close(jobs)  // Signal no more jobs

// Collect results
for result := range results {
    // process result
}
```

### Fan-Out, Fan-In

```go
// Fan-out: multiple goroutines read from same channel
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

// Fan-in: merge multiple channels into one
func fanIn(inputs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    output := make(chan int)
    
    for _, ch := range inputs {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                output <- v
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

### Pipeline

```go
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func filter(in <-chan int, pred func(int) bool) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if pred(n) {
                out <- n
            }
        }
    }()
    return out
}

// Usage
nums := generate(1, 2, 3, 4, 5)
squared := square(nums)
evens := filter(squared, func(n int) bool { return n%2 == 0 })

for n := range evens {
    fmt.Println(n)  // 4, 16
}
```

### Context for Cancellation

```go
func doWork(ctx context.Context) error {
    resultCh := make(chan Result)
    errCh := make(chan error)
    
    go func() {
        result, err := expensiveOperation()
        if err != nil {
            errCh <- err
            return
        }
        resultCh <- result
    }()
    
    select {
    case <-ctx.Done():
        return ctx.Err()  // Cancelled or deadline exceeded
    case err := <-errCh:
        return err
    case result := <-resultCh:
        return processResult(result)
    }
}

// Usage with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := doWork(ctx); err != nil {
    if err == context.DeadlineExceeded {
        // Timeout
    }
}
```

### Select Statement

```go
select {
case v := <-ch1:
    // received from ch1
case ch2 <- value:
    // sent to ch2
case <-time.After(time.Second):
    // timeout
case <-ctx.Done():
    // cancelled
default:
    // non-blocking - runs if no other case ready
}
```

## Sync Primitives

### Mutex

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

### RWMutex

```go
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]string
}

func (m *SafeMap) Get(key string) string {
    m.mu.RLock()  // Multiple readers OK
    defer m.mu.RUnlock()
    return m.data[key]
}

func (m *SafeMap) Set(key, value string) {
    m.mu.Lock()  // Exclusive write
    defer m.mu.Unlock()
    m.data[key] = value
}
```

### WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // work
    }(i)
}

wg.Wait()  // Block until all done
```

### Once

```go
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
        instance.init()
    })
    return instance
}
```

### Cond

```go
type Queue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []int
}

func NewQueue() *Queue {
    q := &Queue{}
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *Queue) Enqueue(item int) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
    q.cond.Signal()  // Wake one waiter
}

func (q *Queue) Dequeue() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) == 0 {
        q.cond.Wait()  // Release lock and wait
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}
```

## Race Condition Detection

### Always Test with Race Detector

```bash
# Run tests with race detection
go test -race ./...

# Build with race detection (slower, for debugging)
go build -race ./cmd/myapp
```

### Common Race Patterns

```go
// RACE - shared variable without sync
var counter int
go func() { counter++ }()
go func() { counter++ }()

// RACE - map concurrent access
m := make(map[string]int)
go func() { m["a"] = 1 }()
go func() { _ = m["a"] }()

// RACE - slice append
var slice []int
go func() { slice = append(slice, 1) }()
go func() { slice = append(slice, 2) }()
```

### Fixes

```go
// Fix 1: Mutex
var mu sync.Mutex
var counter int
go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()

// Fix 2: Atomic
var counter int64
go func() { atomic.AddInt64(&counter, 1) }()

// Fix 3: Channel
counterCh := make(chan int, 1)
counterCh <- 0
go func() {
    v := <-counterCh
    counterCh <- v + 1
}()

// Fix 4: sync.Map for maps
var m sync.Map
go func() { m.Store("a", 1) }()
go func() { m.Load("a") }()
```

## Bubble Tea Concurrency

### Async Commands

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case LoadDataMsg:
        // Start async operation
        return m, m.loadDataCmd()
    case DataLoadedMsg:
        // Handle result
        m.data = msg.data
        return m, nil
    }
    return m, nil
}

func (m Model) loadDataCmd() tea.Cmd {
    return func() tea.Msg {
        // This runs in goroutine managed by Bubble Tea
        data, err := fetchData()
        if err != nil {
            return DataErrorMsg{err: err}
        }
        return DataLoadedMsg{data: data}
    }
}
```

### Multiple Concurrent Operations

```go
func (m Model) loadAllDataCmd() tea.Cmd {
    return tea.Batch(
        m.loadUsersCmd(),
        m.loadEventsCmd(),
        m.loadSettingsCmd(),
    )
}
```

## Common Mistakes

| Mistake | Problem | Fix |
|---------|---------|-----|
| Closing channel twice | Panic | Only sender closes, use sync.Once |
| Writing to closed channel | Panic | Check before write or use select |
| Goroutine leak | Memory leak | Use context, ensure exit path |
| Loop variable capture | Race condition | Pass as parameter |
| Forgetting wg.Done() | Deadlock | Use defer wg.Done() |
| Mutex copy | Undefined behavior | Use pointer receiver |

## Performance Considerations

### When to Use Concurrency

| Scenario | Use Concurrency? |
|----------|------------------|
| I/O bound (network, disk) | Yes - goroutines cheap |
| CPU bound, single core | No - overhead not worth it |
| CPU bound, multi-core | Maybe - benchmark first |
| Simple sequential logic | No - keep it simple |

### Benchmarking Concurrent Code

```go
func BenchmarkSequential(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sequentialProcess(data)
    }
}

func BenchmarkConcurrent(b *testing.B) {
    for i := 0; i < b.N; i++ {
        concurrentProcess(data)
    }
}
```

## Related Skills

- `go-expert` - Go language patterns
- `design-patterns` - Design patterns
- `performance` - Performance optimization
- `benchmarking` - Measuring performance
- `debug-test` - Debugging concurrent issues
