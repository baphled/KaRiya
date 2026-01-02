# KaRiya TUI Intent Architecture: Implementation Enhancements

## Overview

This document addresses key recommendations for enhancing the implementation of the KaRiya TUI intent architecture. These enhancements build on the solid foundation defined in [TUI_INTENT_DIAGRAM.md](TUI_INTENT_DIAGRAM.md) and [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md).

---

## 1. Cross-Intent Metadata & Global Context

### Problem
Some data is naturally shared across multiple intents (e.g., user preferences, last selected profile, application settings) without violating intent boundaries. Currently, the architecture passes data through `IntentResult` metadata, which can be cumbersome for frequently accessed data.

### Solution: Lightweight GlobalContext

**Design**:
```go
type GlobalContext struct {
    mu               sync.RWMutex
    preferences      *UserPreferences      // Non-mutable, read-only
    lastSelectedItem map[string]interface{} // Transient UI state
    appConfig        *AppConfig            // Non-mutable, read-only
}

// Read-only access
func (gc *GlobalContext) GetPreferences() *UserPreferences {
    gc.mu.RLock()
    defer gc.mu.RUnlock()
    return gc.preferences
}

// Transient state (no persistence)
func (gc *GlobalContext) SetLastSelectedItem(key string, value interface{}) {
    gc.mu.Lock()
    defer gc.mu.Unlock()
    gc.lastSelectedItem[key] = value
}

func (gc *GlobalContext) GetLastSelectedItem(key string) interface{} {
    gc.mu.RLock()
    defer gc.mu.RUnlock()
    return gc.lastSelectedItem[key]
}
```

**Key Principles**:
- ✅ **Read-only**: Preferences and config are immutable
- ✅ **Transient**: Last selected items are session-only
- ✅ **No Persistence**: GlobalContext never writes to database
- ✅ **Thread-Safe**: Proper mutex protection
- ✅ **Minimal**: Only truly shared, non-mutable data

**Usage in Intents**:
```go
func (c *CaptureEventIntent) Init(ctx context.Context) tea.Cmd {
    globalCtx := ctx.Value("global_context").(*GlobalContext)
    c.form.SetPreferredEventType(globalCtx.GetLastSelectedItem("eventType"))
    return nil
}

func (c *CaptureEventIntent) Result() *IntentResult[interface{}] {
    globalCtx := c.ctx.Value("global_context").(*GlobalContext)
    globalCtx.SetLastSelectedItem("eventType", c.form.EventType)
    return c.result
}
```

**Benefits**:
- Reduces boilerplate in metadata passing
- Maintains intent independence (can still work without GlobalContext)
- Improves UX with remembered preferences
- No violation of intent boundaries

**Implementation Checklist**:
- [ ] Define GlobalContext struct
- [ ] Implement thread-safe accessors
- [ ] Add to root model
- [ ] Pass via context to all intents
- [ ] Update intents to use GlobalContext for preferences
- [ ] Add tests for GlobalContext
- [ ] Document usage in AGENTS.md

---

## 2. Async Feedback in TUI: Progress Indicators

### Problem
Long-running operations (export, enrichment, CV generation) need user feedback. The current architecture supports async operations but lacks clear UX patterns for progress indication.

### Solution: Non-Blocking Progress Feedback

**Design Pattern**:

```go
// In ExportIntent or any async operation
type AsyncOperation struct {
    id       string
    status   AsyncStatus
    progress int // 0-100
    message  string
    err      error
}

type AsyncStatus string
const (
    StatusPending    AsyncStatus = "pending"
    StatusInProgress AsyncStatus = "in_progress"
    StatusComplete   AsyncStatus = "complete"
    StatusError      AsyncStatus = "error"
)

// Message for progress updates
type AsyncProgressMsg struct {
    Op       *AsyncOperation
    Timestamp time.Time
}

// In intent's Update:
case StateExportInProgress:
    return intent, intent.startAsyncExport()

func (e *ExportIntent) startAsyncExport() tea.Cmd {
    return func() tea.Msg {
        // Start goroutine with progress updates
        go e.performExportWithProgress()
        return nil
    }
}

func (e *ExportIntent) performExportWithProgress() {
    for progress := 0; progress <= 100; progress += 10 {
        time.Sleep(500 * time.Millisecond)

        e.operation.progress = progress
        e.operation.message = fmt.Sprintf("Exporting... %d%%", progress)

        // Send progress update via channel
        // This will be picked up as AsyncProgressMsg in Update
    }
}
```

**View with Progress Indicator**:
```go
func (e *ExportIntent) viewExportInProgress() string {
    width := 40
    filled := (e.operation.progress * width) / 100

    progressBar := "[" + strings.Repeat("=", filled) +
                   strings.Repeat(" ", width-filled) + "]"

    return fmt.Sprintf(
        "Exporting artifact...\n%s\n%d%%\n\n%s",
        progressBar,
        e.operation.progress,
        e.operation.message,
    )
}
```

**Spinner Alternative** (for indeterminate operations):
```go
// Use bubbletea/spinner component
type SpinnerIntent struct {
    spinner *spinner.Model
    // ... other fields
}

func (s *SpinnerIntent) viewInProgress() string {
    return fmt.Sprintf(
        "%s Processing...\n\nPress Ctrl+C to cancel",
        s.spinner.View(),
    )
}
```

**Implementation Checklist**:
- [ ] Define AsyncOperation and AsyncStatus types
- [ ] Create progress update message type
- [ ] Implement progress bar rendering
- [ ] Add spinner component to export intent
- [ ] Test progress updates with timing
- [ ] Implement cancellation during async ops
- [ ] Document async patterns in AGENTS.md

**Benefits**:
- Clear user feedback during long operations
- Non-blocking: UI remains responsive
- Cancellable operations
- Consistent across all async workflows

---

## 3. CI/CD Integration for Quality Assurance

### Problem
Without automated checks, code quality and test coverage can degrade. Manual verification is error-prone and time-consuming.

### Solution: Comprehensive CI/CD Pipeline

**GitHub Actions Workflow** (`.github/workflows/intent-checks.yml`):

```yaml
name: Intent Architecture Checks

on:
  pull_request:
    paths:
      - 'internal/cli/intents/**'
      - '.github/workflows/intent-checks.yml'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run Tests
        run: go test -v -race -coverprofile=coverage.out ./internal/cli/intents/...

      - name: Check Coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 90" | bc -l) )); then
            echo "Coverage ${coverage}% is below 90% threshold"
            exit 1
          fi

      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: ./internal/cli/intents/...

  type-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run go vet
        run: go vet ./internal/cli/intents/...

      - name: Run go fmt check
        run: |
          if [ "$(gofmt -s -l ./internal/cli/intents/ | wc -l)" -gt 0 ]; then
            echo "Code is not formatted with gofmt"
            exit 1
          fi

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        run: go build -v ./cmd/kariya

  property-based-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run Property-Based Tests
        run: go test -v -run Property ./internal/cli/intents/...
```

**Pre-Commit Hook** (`.git/hooks/pre-commit`):
```bash
#!/bin/bash

# Format code
gofmt -s -w ./internal/cli/intents/...

# Run linter
golangci-lint run ./internal/cli/intents/...
if [ $? -ne 0 ]; then
    echo "Linting failed"
    exit 1
fi

# Run tests
go test -race ./internal/cli/intents/...
if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi

# Check coverage
coverage=$(go test -coverprofile=coverage.out ./internal/cli/intents/... | \
    go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

if (( $(echo "$coverage < 90" | bc -l) )); then
    echo "Coverage ${coverage}% is below 90% threshold"
    rm coverage.out
    exit 1
fi

rm coverage.out
```

**Implementation Checklist**:
- [ ] Create `.github/workflows/intent-checks.yml`
- [ ] Set up codecov integration
- [ ] Create pre-commit hook
- [ ] Configure branch protection rules
- [ ] Document CI/CD in README
- [ ] Add coverage badge to README
- [ ] Set up coverage reports

**Benefits**:
- Automatic quality checks on every PR
- Coverage reports prevent regressions
- Consistent code formatting
- Early detection of issues

---

## 4. Performance Benchmarking & Early Detection

### Problem
UI rendering performance issues can accumulate unnoticed. Large datasets (e.g., timeline with 1000+ events) can cause lag without explicit benchmarking.

### Solution: Baseline Performance Tests

**Benchmark Tests** (`internal/cli/intents/browse/browse_bench_test.go`):

```go
func BenchmarkBrowseTimelineRendering(b *testing.B) {
    intent := NewBrowseTimelineIntent()

    // Load test data: 1000 events
    events := generateTestEvents(1000)
    intent.events = events

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = intent.View()
    }

    // Ensure rendering completes in <100ms per frame
    b.StopTimer()
    avgTime := b.Elapsed() / time.Duration(b.N)
    if avgTime > 100*time.Millisecond {
        b.Errorf("Rendering too slow: %v per frame (target: <100ms)", avgTime)
    }
}

func BenchmarkCaptureEventFormValidation(b *testing.B) {
    form := NewCaptureForm()
    testData := generateTestFormData()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = form.Validate(testData)
    }

    // Validation should be <1ms
    b.StopTimer()
    avgTime := b.Elapsed() / time.Duration(b.N)
    if avgTime > 1*time.Millisecond {
        b.Errorf("Validation too slow: %v (target: <1ms)", avgTime)
    }
}

func BenchmarkGenerateCVPreview(b *testing.B) {
    intent := NewGenerateCVIntent()
    intent.profile = generateTestProfile()
    intent.audience = generateTestAudience()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = intent.generatePreview()
    }

    // Preview generation should be <500ms
    b.StopTimer()
    avgTime := b.Elapsed() / time.Duration(b.N)
    if avgTime > 500*time.Millisecond {
        b.Errorf("Preview generation too slow: %v (target: <500ms)", avgTime)
    }
}
```

**CI/CD Integration**:
```yaml
  benchmarks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run Benchmarks
        run: go test -bench=. -benchmem ./internal/cli/intents/...

      - name: Store Benchmark Results
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: output.txt
          github-token: ${{ secrets.GITHUB_TOKEN }}
          auto-push: true
```

**Performance Thresholds**:
- Rendering: <100ms per frame (60 FPS)
- Form validation: <1ms
- List filtering: <50ms
- CV preview generation: <500ms
- Export operation: <5s (with progress feedback)

**Implementation Checklist**:
- [ ] Create benchmark tests for each intent
- [ ] Set performance thresholds
- [ ] Add benchmarks to CI/CD
- [ ] Store benchmark history
- [ ] Document performance characteristics
- [ ] Create performance testing guide

**Benefits**:
- Catch performance regressions early
- Establish baseline metrics
- Optimize critical paths
- Better user experience

---

## 5. Implementation Timeline with Enhancements

### Enhanced Phase 1 (1.5 weeks → 2 weeks)
- Foundation & Core Infrastructure
- **+ GlobalContext implementation**
- **+ CI/CD pipeline setup**
- **+ Benchmark baseline tests**

### Enhanced Phase 2 (2 weeks → 2.5 weeks)
- CaptureEvent Intent Implementation
- **+ Progress feedback pattern**
- **+ Performance benchmarks**

### Enhanced Phase 3 (4 weeks → 4.5 weeks)
- Remaining Core Intents
- **+ Async feedback for each async intent**

### Enhanced Phase 4 (2 weeks → 2.5 weeks)
- Integration & Polish
- **+ CI/CD validation**
- **+ Performance optimization based on benchmarks**

**Total Enhanced Timeline**: 10.5 weeks (vs. 9.5 weeks baseline)

---

## Implementation Priority

### Must Have (Phase 1-4)
1. Type-safe intent architecture ✅
2. GlobalContext for shared preferences
3. Progress feedback for async operations
4. CI/CD pipeline with coverage checks

### Should Have (Phase 2-4)
1. Performance benchmarks
2. Pre-commit hooks
3. Coverage reports

### Nice to Have (Phase 5+)
1. Advanced analytics
2. Custom performance dashboards
3. A/B testing framework

---

## Maintenance & Monitoring

### Weekly Checks
- [ ] Review coverage reports
- [ ] Check benchmark trends
- [ ] Review error logs
- [ ] Check performance metrics

### Monthly Reviews
- [ ] Analyze usage patterns
- [ ] Review performance trends
- [ ] Identify optimization opportunities
- [ ] Plan for next phase

### Quarterly Audits
- [ ] Full architecture review
- [ ] Capacity planning
- [ ] Technology updates
- [ ] Security review

---

## References

- [TUI_INTENT_DIAGRAM.md](TUI_INTENT_DIAGRAM.md) - Core architecture
- [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md) - Base roadmap
- [IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md) - Detailed checklist
- [AGENTS.md](../AGENTS.md) - Agent guidelines

---

## Next Steps

1. **Week 1**: Review enhancements with team
2. **Week 1**: Set up CI/CD pipeline
3. **Week 1**: Create GlobalContext types
4. **Week 2**: Implement progress feedback pattern
5. **Week 2**: Create performance benchmarks
6. **Week 3+**: Continue with core implementation

---

*Last Updated: 2026-01-02*
*Status: Enhancement Specification Ready*

