---
name: logging-observability
description: Implement structured logging, tracing, and metrics for debugging and monitoring
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide implementation of logging, tracing, and metrics for effective debugging and monitoring. Good observability makes problems visible and diagnosable.

## When to use me

- Adding logging to new code
- Debugging issues in development
- Setting up monitoring
- Investigating production problems
- Reviewing logging practices

## Core Principles

1. **Structured over text** - Machine-parseable logs
2. **Context propagation** - Trace requests through system
3. **Appropriate levels** - Right verbosity for situation
4. **Actionable information** - Logs should help solve problems

## Structured Logging

### Use Key-Value Pairs

```go
// GOOD - Structured logging
logger.Info("event created",
    "event_id", event.ID,
    "user_id", userID,
    "event_type", event.Type,
    "duration_ms", duration.Milliseconds(),
)

// BAD - String interpolation
logger.Info(fmt.Sprintf("event %s created by user %s", event.ID, userID))
```

### Consistent Field Names

```go
// Standard fields
"request_id"    // Request correlation
"user_id"       // User identifier
"event_id"      // Entity identifiers
"duration_ms"   // Timing in milliseconds
"error"         // Error message
"stack"         // Stack trace
"component"     // System component
"operation"     // What's being done
```

### Log Levels

| Level | Use For | Example |
|-------|---------|---------|
| **Debug** | Development details, verbose | Variable values, flow tracing |
| **Info** | Normal operations | Request completed, job started |
| **Warn** | Recoverable issues | Retry attempted, deprecated usage |
| **Error** | Failures requiring attention | Request failed, connection lost |

```go
// Debug - Development only, very verbose
logger.Debug("processing item",
    "item_id", item.ID,
    "current_index", i,
    "total_items", len(items),
)

// Info - Normal operations
logger.Info("request completed",
    "request_id", reqID,
    "method", r.Method,
    "path", r.URL.Path,
    "status", status,
    "duration_ms", duration.Milliseconds(),
)

// Warn - Something unexpected but handled
logger.Warn("retry attempted",
    "operation", "database_connect",
    "attempt", attempt,
    "max_attempts", maxAttempts,
    "error", err.Error(),
)

// Error - Something failed
logger.Error("request failed",
    "request_id", reqID,
    "error", err.Error(),
    "stack", string(debug.Stack()),
)
```

## Context Propagation

### Request Context

```go
// Add request ID to context
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Extract and log
func (s *Service) DoSomething(ctx context.Context) error {
    requestID := ctx.Value(requestIDKey).(string)
    
    s.logger.Info("doing something",
        "request_id", requestID,
        "operation", "do_something",
    )
    // ...
}
```

### Logger with Context

```go
// Create logger with base fields
func NewLogger(ctx context.Context) *Logger {
    logger := baseLogger
    
    if requestID := ctx.Value(requestIDKey); requestID != nil {
        logger = logger.With("request_id", requestID)
    }
    if userID := ctx.Value(userIDKey); userID != nil {
        logger = logger.With("user_id", userID)
    }
    
    return logger
}

// Use throughout request
logger := NewLogger(ctx)
logger.Info("starting operation")
// All logs automatically include request_id, user_id
```

## What to Log

### Always Log

```go
// Request boundaries
logger.Info("request started", "method", method, "path", path)
logger.Info("request completed", "status", status, "duration_ms", ms)

// State changes
logger.Info("event created", "event_id", id)
logger.Info("event deleted", "event_id", id)

// Errors
logger.Error("operation failed", "error", err, "context", ctx)

// Security events
logger.Warn("authentication failed", "user", user, "reason", reason)
logger.Info("permission granted", "user", user, "resource", resource)
```

### Never Log

```go
// NEVER log these
password
api_key
token
secret
credit_card
ssn
personal_health_info

// Mask sensitive data
logger.Info("user authenticated",
    "user_id", userID,
    "email", maskEmail(email),  // j***@example.com
)
```

### Conditional Logging

```go
// Expensive to compute - check level first
if logger.IsDebugEnabled() {
    logger.Debug("detailed state",
        "state", expensiveStateSnapshot(),
    )
}
```

## Error Logging

### Include Context

```go
// GOOD - Rich context
logger.Error("failed to save event",
    "event_id", event.ID,
    "user_id", userID,
    "error", err.Error(),
    "event_type", event.Type,
)

// BAD - No context
logger.Error("save failed", "error", err)
```

### Log Once at Boundary

```go
// DON'T log at every level
func (r *Repo) Save(e *Event) error {
    err := r.db.Save(e)
    if err != nil {
        r.logger.Error("save failed", "error", err)  // Logged here
        return err
    }
    return nil
}

func (s *Service) CreateEvent(e *Event) error {
    err := s.repo.Save(e)
    if err != nil {
        s.logger.Error("save failed", "error", err)  // Logged again!
        return err
    }
    return nil
}

// DO log once at appropriate boundary
func (r *Repo) Save(e *Event) error {
    err := r.db.Save(e)
    if err != nil {
        return fmt.Errorf("save event %s: %w", e.ID, err)  // Wrap, don't log
    }
    return nil
}

func (s *Service) CreateEvent(e *Event) error {
    err := s.repo.Save(e)
    if err != nil {
        s.logger.Error("failed to create event",  // Log at service boundary
            "event_id", e.ID,
            "error", err.Error(),
        )
        return err
    }
    return nil
}
```

## Metrics

### Counter (Things that increment)

```go
// Requests, errors, items processed
requestsTotal.Inc()
errorsTotal.WithLabels("type", "validation").Inc()
eventsCreated.Inc()
```

### Gauge (Current value)

```go
// Active connections, queue size, temperature
activeConnections.Set(count)
queueSize.Set(len(queue))
```

### Histogram (Distribution)

```go
// Request duration, response size
requestDuration.Observe(duration.Seconds())
responseSize.Observe(float64(size))
```

### Key Metrics

```go
// RED metrics for services
requests_total           // Rate
request_errors_total     // Errors  
request_duration_seconds // Duration

// USE metrics for resources
utilisation  // % time busy
saturation   // Queue depth
errors       // Error count
```

## Tracing

### Span Creation

```go
func (s *Service) CreateEvent(ctx context.Context, event *Event) error {
    ctx, span := tracer.Start(ctx, "CreateEvent")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("event.id", event.ID),
        attribute.String("event.type", string(event.Type)),
    )
    
    // Pass context to downstream calls
    if err := s.repo.Save(ctx, event); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }
    
    return nil
}
```

## KaRiya Logging Patterns

### Service Layer

```go
func (s *EventService) Create(ctx context.Context, event *Event) (*Event, error) {
    logger := s.logger.With(
        "operation", "create_event",
        "event_type", event.Type,
    )
    
    logger.Debug("creating event")
    
    created, err := s.repo.Save(ctx, event)
    if err != nil {
        logger.Error("failed to create event", "error", err)
        return nil, fmt.Errorf("create event: %w", err)
    }
    
    logger.Info("event created", "event_id", created.ID)
    return created, nil
}
```

### CLI/TUI Layer

```go
// Log state transitions
func (i *Intent) transitionTo(state IntentState) {
    i.logger.Debug("state transition",
        "from", i.state,
        "to", state,
    )
    i.state = state
}

// Log user actions
func (i *Intent) handleKeyPress(key string) {
    i.logger.Debug("key pressed",
        "key", key,
        "state", i.state,
    )
}
```

## Related Skills

- `error-handling` - Error management
- `debug-test` - Debugging techniques
- `incident-response` - Using logs in incidents
- `security` - Secure logging practices
