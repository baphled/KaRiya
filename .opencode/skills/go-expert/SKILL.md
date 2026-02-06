---
name: go-expert
description: Go language expertise including idioms, patterns, performance, concurrency, and best practices
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Provide Go language expertise including idioms, patterns, performance optimization, concurrency, and best practices.

## When to use me

Use this skill when:
- Writing idiomatic Go code
- Designing interfaces and types
- Working with concurrency
- Optimizing performance
- Debugging Go-specific issues
- Reviewing Go code quality

## Go Philosophy

1. **Simplicity** - Clear is better than clever
2. **Readability** - Code is read more than written
3. **Composition** - Prefer composition over inheritance
4. **Explicit** - No hidden magic
5. **Orthogonality** - Small, focused packages

## Go Idioms

### Error Handling

```go
// GOOD - Handle errors explicitly
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// GOOD - Sentinel errors for expected cases
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // Handle specifically
}

// GOOD - Error types for rich errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// BAD - Ignoring errors
result, _ := doSomething()

// BAD - Panic for normal errors
if err != nil {
    panic(err)
}
```

### Interface Design

```go
// GOOD - Small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// GOOD - Accept interfaces, return structs
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

// GOOD - Define interfaces where they're used
// In consumer package:
type EventStore interface {
    Save(ctx context.Context, event *Event) error
}

// BAD - Large interfaces
type Repository interface {
    Create(...)
    Read(...)
    Update(...)
    Delete(...)
    List(...)
    Count(...)
    // 20 more methods...
}

// BAD - Interface pollution
type Stringer interface {
    String() string
}
// Don't create interfaces for single implementations
```

### Struct Design

```go
// GOOD - Zero value is useful
type Buffer struct {
    data []byte
    // Zero value: empty buffer, ready to use
}

// GOOD - Options pattern for configuration
type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = d
    }
}

func NewServer(opts ...ServerOption) *Server {
    s := &Server{timeout: 30 * time.Second} // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// GOOD - Embedding for composition
type Server struct {
    *http.Server
    logger Logger
}

// BAD - Exported fields that shouldn't be modified
type Config struct {
    APIKey string // Should this be settable?
}
```

### Naming Conventions

```go
// Package names: short, lowercase, no underscores
package user     // GOOD
package userService // BAD
package user_service // BAD

// Getters: no "Get" prefix
func (u *User) Name() string     // GOOD
func (u *User) GetName() string  // BAD

// Interfaces: -er suffix for single method
type Reader interface { Read(...) }
type Stringer interface { String() string }

// Acronyms: consistent case
userID  // GOOD
userId  // BAD
httpClient // GOOD
HTTPClient // Also GOOD (exported)

// Unexported: start lowercase
type service struct{}  // internal
type Service struct{}  // exported
```

## Concurrency Patterns

### Goroutines and Channels

```go
// GOOD - Use context for cancellation
func worker(ctx context.Context, jobs <-chan Job) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case job, ok := <-jobs:
            if !ok {
                return nil
            }
            process(job)
        }
    }
}

// GOOD - Bounded concurrency
func processAll(items []Item, workers int) {
    sem := make(chan struct{}, workers)
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        sem <- struct{}{} // Acquire
        go func(item Item) {
            defer wg.Done()
            defer func() { <-sem }() // Release
            process(item)
        }(item)
    }
    wg.Wait()
}

// GOOD - Fan-out, fan-in
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

// BAD - Goroutine leak
func leaky() {
    ch := make(chan int)
    go func() {
        val := <-ch // Blocks forever if ch never receives
        fmt.Println(val)
    }()
    // Function returns, goroutine stuck
}
```

### Mutexes

```go
// GOOD - Mutex protects data, not code
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

// GOOD - RWMutex for read-heavy workloads
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

// BAD - Lock held during I/O
func (c *Cache) LoadFromDB() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data = db.LoadAll() // Holds lock during slow operation
}
```

### sync.Once

```go
// GOOD - Lazy initialization
type Client struct {
    once sync.Once
    conn *Connection
}

func (c *Client) getConn() *Connection {
    c.once.Do(func() {
        c.conn = connect()
    })
    return c.conn
}
```

## Performance

### Memory Allocation

```go
// GOOD - Pre-allocate slices when size known
items := make([]Item, 0, expectedSize)

// GOOD - Reuse buffers
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process(data []byte) {
    buf := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(buf)
    buf.Reset()
    // Use buf...
}

// GOOD - Avoid string concatenation in loops
var builder strings.Builder
for _, s := range items {
    builder.WriteString(s)
}
result := builder.String()

// BAD - String concatenation in loop
result := ""
for _, s := range items {
    result += s // Allocates new string each time
}
```

### Benchmarking

```go
func BenchmarkProcess(b *testing.B) {
    data := setupTestData()
    b.ResetTimer() // Don't count setup
    
    for i := 0; i < b.N; i++ {
        process(data)
    }
}

func BenchmarkProcessParallel(b *testing.B) {
    data := setupTestData()
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            process(data)
        }
    })
}
```

### Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace
go test -trace=trace.out -bench=.
go tool trace trace.out
```

## Testing Patterns

### Table-Driven Tests

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int
        wantErr bool
    }{
        {"valid", "42", 42, false},
        {"negative", "-1", -1, false},
        {"invalid", "abc", 0, true},
        {"empty", "", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Parse() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Helpers

```go
func setupTestDB(t *testing.T) *DB {
    t.Helper()
    db, err := NewDB(":memory:")
    if err != nil {
        t.Fatalf("failed to create test db: %v", err)
    }
    t.Cleanup(func() {
        db.Close()
    })
    return db
}
```

### Mocking with Interfaces

```go
// Define interface in consumer
type EmailSender interface {
    Send(to, subject, body string) error
}

// Production implementation
type SMTPSender struct { ... }

// Test mock
type MockSender struct {
    SendFunc func(to, subject, body string) error
}

func (m *MockSender) Send(to, subject, body string) error {
    return m.SendFunc(to, subject, body)
}
```

## Common Gotchas

### Loop Variable Capture

```go
// BAD - All goroutines share same variable
for _, item := range items {
    go func() {
        process(item) // Captures loop variable
    }()
}

// GOOD - Pass as parameter
for _, item := range items {
    go func(item Item) {
        process(item)
    }(item)
}

// GOOD (Go 1.22+) - Loop variables are per-iteration
for _, item := range items {
    go func() {
        process(item) // Safe in Go 1.22+
    }()
}
```

### Nil Slices vs Empty Slices

```go
var nilSlice []int      // nil, len=0, cap=0
emptySlice := []int{}   // not nil, len=0, cap=0
madeSlice := make([]int, 0) // not nil, len=0, cap=0

// JSON encoding differs
json.Marshal(nilSlice)   // "null"
json.Marshal(emptySlice) // "[]"
```

### Interface Nil Check

```go
type MyError struct{}
func (e *MyError) Error() string { return "error" }

func returnsError() error {
    var e *MyError = nil
    return e // Returns non-nil interface containing nil pointer
}

err := returnsError()
if err != nil { // TRUE! Interface is not nil
    fmt.Println(err.Error()) // Panic: nil pointer
}

// GOOD - Return nil explicitly
func returnsError() error {
    var e *MyError = nil
    if e == nil {
        return nil
    }
    return e
}
```

### Defer in Loops

```go
// BAD - Defers accumulate until function returns
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close() // Not closed until function ends
}

// GOOD - Use closure or explicit close
for _, file := range files {
    func() {
        f, _ := os.Open(file)
        defer f.Close()
        // Process...
    }()
}
```

## Go Modules

```bash
# Initialize module
go mod init github.com/user/project

# Add dependency
go get github.com/pkg/errors

# Update dependencies
go get -u ./...

# Tidy (remove unused)
go mod tidy

# Vendor dependencies
go mod vendor
```

## Related skills

- `clean-code` - Clean code principles
- `architecture` - System design
- `security` - Secure Go coding
- `debug-test` - Testing and debugging
- `refactor` - Refactoring patterns
