---
name: fuzz-testing
description: Advanced testing with Go's built-in fuzzing for finding edge cases and crashes
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide the use of Go's built-in fuzz testing to find edge cases, crashes, and security vulnerabilities that traditional tests miss.

## When to use me

Use this skill when:
- Testing input parsing/validation
- Testing data transformation functions
- Looking for crash bugs
- Security testing input handlers
- Testing serialisation/deserialisation

## Go Fuzzing Basics

### Fuzz Test Structure

```go
// event_fuzz_test.go
package domain

import (
    "testing"
)

func FuzzParseEvent(f *testing.F) {
    // Add seed corpus (known good inputs)
    f.Add("2026-01-15", "work", "Completed feature X")
    f.Add("", "", "")
    f.Add("invalid-date", "unknown", "test")
    
    // Fuzz function
    f.Fuzz(func(t *testing.T, date, eventType, description string) {
        // This will be called with random mutations
        event, err := ParseEvent(date, eventType, description)
        
        if err != nil {
            // Error is acceptable, but shouldn't panic
            return
        }
        
        // If parsing succeeded, result should be valid
        if event.Date.IsZero() {
            t.Error("parsed event has zero date")
        }
    })
}
```

### Running Fuzz Tests

```bash
# Run fuzz test for 30 seconds
go test -fuzz=FuzzParseEvent -fuzztime=30s ./internal/domain/...

# Run until failure found
go test -fuzz=FuzzParseEvent ./internal/domain/...

# Run with more workers
go test -fuzz=FuzzParseEvent -parallel=8 ./internal/domain/...

# Run specific fuzz test
go test -fuzz=FuzzParseEvent -run=FuzzParseEvent ./internal/domain/...
```

## Seed Corpus

### Adding Seeds

Seeds help the fuzzer start with meaningful inputs:

```go
func FuzzParseJSON(f *testing.F) {
    // Add real-world examples
    f.Add([]byte(`{"name":"test","value":123}`))
    f.Add([]byte(`{}`))
    f.Add([]byte(`[]`))
    
    // Add edge cases
    f.Add([]byte(`{"name":""}`))           // Empty string
    f.Add([]byte(`{"value":-1}`))          // Negative
    f.Add([]byte(`{"value":9999999999}`))  // Large number
    
    // Add malformed inputs
    f.Add([]byte(`{`))                     // Incomplete
    f.Add([]byte(`null`))                  // Null
    f.Add([]byte(``))                      // Empty
    
    f.Fuzz(func(t *testing.T, data []byte) {
        var v MyStruct
        _ = json.Unmarshal(data, &v)  // Should not panic
    })
}
```

### Corpus Directory

Fuzzer stores interesting inputs in `testdata/fuzz/<TestName>/`:

```
testdata/fuzz/FuzzParseEvent/
├── corpus/
│   ├── 0a1b2c3d...   # Seed that found edge case
│   └── 4e5f6a7b...   # Another interesting input
└── crashers/
    └── 8c9d0e1f...   # Input that caused crash
```

## Fuzzing Patterns

### Pattern 1: Parser Testing

```go
func FuzzParseDate(f *testing.F) {
    f.Add("2026-01-15")
    f.Add("01/15/2026")
    f.Add("January 15, 2026")
    f.Add("")
    f.Add("not a date")
    
    f.Fuzz(func(t *testing.T, input string) {
        date, err := ParseDate(input)
        
        if err == nil {
            // If parsing succeeded, verify round-trip
            formatted := date.Format("2006-01-02")
            reparsed, err := ParseDate(formatted)
            if err != nil {
                t.Errorf("round-trip failed: %s -> %s -> error", input, formatted)
            }
            if !date.Equal(reparsed) {
                t.Errorf("round-trip mismatch: %s != %s", date, reparsed)
            }
        }
    })
}
```

### Pattern 2: Serialisation Round-Trip

```go
func FuzzEventRoundTrip(f *testing.F) {
    f.Add("test", int64(1234567890), "description")
    
    f.Fuzz(func(t *testing.T, name string, timestamp int64, desc string) {
        original := &Event{
            Name:        name,
            Timestamp:   time.Unix(timestamp, 0),
            Description: desc,
        }
        
        // Serialize
        data, err := json.Marshal(original)
        if err != nil {
            return  // Some inputs may not be serializable
        }
        
        // Deserialize
        var decoded Event
        if err := json.Unmarshal(data, &decoded); err != nil {
            t.Errorf("failed to unmarshal: %v", err)
            return
        }
        
        // Compare
        if original.Name != decoded.Name {
            t.Errorf("name mismatch: %q != %q", original.Name, decoded.Name)
        }
    })
}
```

### Pattern 3: State Machine

```go
func FuzzStateMachine(f *testing.F) {
    f.Add([]byte{0, 1, 2, 3})  // Sequence of actions
    
    f.Fuzz(func(t *testing.T, actions []byte) {
        sm := NewStateMachine()
        
        for _, action := range actions {
            switch action % 4 {
            case 0:
                sm.Start()
            case 1:
                sm.Process("data")
            case 2:
                sm.Pause()
            case 3:
                sm.Stop()
            }
            
            // State machine should never be in invalid state
            if !sm.IsValidState() {
                t.Error("invalid state reached")
            }
        }
    })
}
```

### Pattern 4: SQL Injection Prevention

```go
func FuzzSQLQuery(f *testing.F) {
    f.Add("normal input")
    f.Add("'; DROP TABLE users; --")
    f.Add("1 OR 1=1")
    f.Add("admin'--")
    
    f.Fuzz(func(t *testing.T, input string) {
        // Build query using parameterized queries
        query, args := buildQuery(input)
        
        // Verify no SQL injection possible
        if strings.Contains(query, input) && strings.Contains(input, "'") {
            t.Error("potential SQL injection: raw input in query")
        }
        
        // Verify args are properly parameterized
        if len(args) == 0 && input != "" {
            t.Error("input not parameterized")
        }
    })
}
```

### Pattern 5: Memory Safety

```go
func FuzzSliceOperations(f *testing.F) {
    f.Add([]byte{1, 2, 3}, 0, 3)
    f.Add([]byte{}, 0, 0)
    
    f.Fuzz(func(t *testing.T, data []byte, start, end int) {
        // This should never panic even with invalid indices
        result := safeSlice(data, start, end)
        
        // Result should always be valid
        if result == nil {
            t.Error("result should never be nil")
        }
        
        // Length should be bounded
        if len(result) > len(data) {
            t.Error("result longer than input")
        }
    })
}
```

## Handling Crashes

When fuzzer finds a crash:

```bash
# 1. Run the failing case
go test -run=FuzzParseEvent/crasher_name ./internal/domain/...

# 2. Get the input that caused crash
cat testdata/fuzz/FuzzParseEvent/crashers/crasher_name

# 3. Add regression test
func TestParseEventCrasher(t *testing.T) {
    // This input was found by fuzzing
    input := string(mustReadFile(t, "testdata/fuzz/FuzzParseEvent/crashers/crasher_name"))
    
    // Should not panic
    _, _ = ParseEvent(input)
}
```

## CI Integration

### Running in CI

```yaml
# .github/workflows/fuzz.yml
name: Fuzz Testing

on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM
  workflow_dispatch:      # Manual trigger

jobs:
  fuzz:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run fuzz tests
        run: |
          go test -fuzz=Fuzz -fuzztime=5m ./...
      
      - name: Upload crashers
        if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: fuzz-crashers
          path: testdata/fuzz/**/crashers/
```

### Makefile Target

```makefile
# Run fuzz tests for 1 minute each
.PHONY: fuzz
fuzz:
	@echo "Running fuzz tests..."
	@for test in $$(go test -list 'Fuzz.*' ./... 2>/dev/null | grep '^Fuzz'); do \
		echo "Fuzzing $$test..."; \
		go test -fuzz=$$test -fuzztime=1m ./... || exit 1; \
	done

# Run specific fuzz test
.PHONY: fuzz-one
fuzz-one:
	go test -fuzz=$(FUZZ) -fuzztime=$(TIME) ./...
```

## Best Practices

### What to Fuzz

| Good Targets | Why |
|-------------|-----|
| Parsers | Complex input handling |
| Deserializers | Security-critical |
| Validators | Edge case discovery |
| Encoders/Decoders | Round-trip correctness |
| State machines | Invalid state detection |

### What NOT to Fuzz

| Bad Targets | Why |
|------------|-----|
| Database operations | Slow, side effects |
| Network calls | External dependencies |
| File system | Side effects |
| Random number generators | Non-deterministic |

### Fuzz Test Checklist

- [ ] Add meaningful seed corpus
- [ ] Test for panics (implicit)
- [ ] Test for hangs (use timeout)
- [ ] Verify invariants in fuzz function
- [ ] Add regression tests for found bugs
- [ ] Run regularly in CI

## Debugging Fuzz Failures

```bash
# Reproduce with verbose output
go test -v -run=FuzzParseEvent/specific_input ./...

# Run with race detector
go test -race -fuzz=FuzzParseEvent -fuzztime=30s ./...

# Get coverage during fuzzing
go test -fuzz=FuzzParseEvent -fuzztime=30s -coverprofile=fuzz.out ./...
go tool cover -html=fuzz.out
```

## Related Skills

- `tdd-workflow` - Integrating fuzz tests with TDD
- `security` - Finding security vulnerabilities
- `prove-correctness` - Using fuzz tests as evidence
- `debug-test` - Debugging fuzz failures
