---
name: monitoring
description: Post-deployment health checks, observability, and system monitoring
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide post-deployment monitoring, health checks, and observability practices for KaRiya applications.

## When to use me

Use this skill when:
- Deploying new features
- Checking system health post-merge
- Setting up monitoring for new components
- Investigating performance issues
- Validating releases

## Post-Deployment Health Checks

### Immediate Checks (0-5 minutes)

```bash
# 1. Application starts
kariya --version
kariya health

# 2. Database accessible
sqlite3 ~/.kariya/data.db "SELECT 1;"

# 3. No startup errors
kariya 2>&1 | head -20

# 4. Core operations work
kariya events list --limit 5
```

### Short-Term Checks (5-30 minutes)

```bash
# 1. Memory usage stable
ps aux | grep kariya

# 2. No error log growth
tail -f ~/.kariya/logs/error.log

# 3. Response times normal
time kariya events list
```

## Application Health Endpoint

### Implementation Pattern

```go
// internal/health/health.go
package health

type Status struct {
    Status    string            `json:"status"`
    Version   string            `json:"version"`
    Uptime    time.Duration     `json:"uptime"`
    Checks    map[string]Check  `json:"checks"`
    Timestamp time.Time         `json:"timestamp"`
}

type Check struct {
    Status  string        `json:"status"`
    Latency time.Duration `json:"latency,omitempty"`
    Error   string        `json:"error,omitempty"`
}

func (h *Health) Check() *Status {
    status := &Status{
        Status:    "healthy",
        Version:   version.Version,
        Uptime:    time.Since(h.startTime),
        Checks:    make(map[string]Check),
        Timestamp: time.Now(),
    }
    
    // Database check
    status.Checks["database"] = h.checkDatabase()
    
    // Disk space check
    status.Checks["disk"] = h.checkDiskSpace()
    
    // Memory check
    status.Checks["memory"] = h.checkMemory()
    
    // Aggregate status
    for _, check := range status.Checks {
        if check.Status != "healthy" {
            status.Status = "degraded"
        }
    }
    
    return status
}

func (h *Health) checkDatabase() Check {
    start := time.Now()
    err := h.db.Exec("SELECT 1").Error
    latency := time.Since(start)
    
    if err != nil {
        return Check{Status: "unhealthy", Latency: latency, Error: err.Error()}
    }
    if latency > 100*time.Millisecond {
        return Check{Status: "degraded", Latency: latency}
    }
    return Check{Status: "healthy", Latency: latency}
}
```

### CLI Health Command

```go
// cmd/kariya/health.go
var healthCmd = &cobra.Command{
    Use:   "health",
    Short: "Check application health",
    Run: func(cmd *cobra.Command, args []string) {
        status := health.Check()
        
        if outputJSON {
            json.NewEncoder(os.Stdout).Encode(status)
            return
        }
        
        fmt.Printf("Status: %s\n", status.Status)
        fmt.Printf("Version: %s\n", status.Version)
        fmt.Printf("Uptime: %s\n", status.Uptime)
        
        for name, check := range status.Checks {
            fmt.Printf("  %s: %s", name, check.Status)
            if check.Latency > 0 {
                fmt.Printf(" (%s)", check.Latency)
            }
            fmt.Println()
        }
    },
}
```

## Metrics Collection

### Key Metrics for KaRiya

| Metric | Type | Description |
|--------|------|-------------|
| `kariya_events_total` | Counter | Total events processed |
| `kariya_db_query_duration_seconds` | Histogram | Database query latency |
| `kariya_memory_bytes` | Gauge | Memory usage |
| `kariya_goroutines` | Gauge | Active goroutines |
| `kariya_errors_total` | Counter | Error count by type |

### Implementation

```go
// internal/metrics/metrics.go
package metrics

import (
    "runtime"
    "time"
)

type Metrics struct {
    EventsProcessed int64
    ErrorCount      int64
    StartTime       time.Time
    
    // Timing
    DBQueryCount    int64
    DBQueryDuration time.Duration
}

var global = &Metrics{StartTime: time.Now()}

func RecordEvent() {
    atomic.AddInt64(&global.EventsProcessed, 1)
}

func RecordError() {
    atomic.AddInt64(&global.ErrorCount, 1)
}

func RecordDBQuery(duration time.Duration) {
    atomic.AddInt64(&global.DBQueryCount, 1)
    // Note: This is simplified; use proper histogram in production
}

func GetStats() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return map[string]interface{}{
        "events_processed": atomic.LoadInt64(&global.EventsProcessed),
        "error_count":      atomic.LoadInt64(&global.ErrorCount),
        "uptime_seconds":   time.Since(global.StartTime).Seconds(),
        "goroutines":       runtime.NumGoroutine(),
        "memory_alloc_mb":  float64(m.Alloc) / 1024 / 1024,
        "memory_sys_mb":    float64(m.Sys) / 1024 / 1024,
        "gc_runs":          m.NumGC,
    }
}
```

## Logging Best Practices

### Structured Logging

```go
// Use structured logging for queryability
import "log/slog"

func ProcessEvent(ctx context.Context, event *Event) error {
    logger := slog.With(
        "event_id", event.ID,
        "event_type", event.Type,
        "trace_id", ctx.Value("trace_id"),
    )
    
    logger.Info("processing event")
    
    if err := validate(event); err != nil {
        logger.Error("validation failed",
            "error", err,
            "field", err.Field,
        )
        return err
    }
    
    logger.Info("event processed",
        "duration_ms", time.Since(start).Milliseconds(),
    )
    
    return nil
}
```

### Log Levels

| Level | Use For |
|-------|---------|
| `DEBUG` | Development, detailed tracing |
| `INFO` | Normal operations, milestones |
| `WARN` | Recoverable issues, deprecations |
| `ERROR` | Failures requiring attention |

### Log Rotation

```go
// Configure log rotation
import "gopkg.in/natefinch/lumberjack.v2"

func setupLogging() {
    logger := &lumberjack.Logger{
        Filename:   "~/.kariya/logs/app.log",
        MaxSize:    10, // MB
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true,
    }
    slog.SetDefault(slog.New(slog.NewJSONHandler(logger, nil)))
}
```

## Alerting Patterns

### Alert on Errors

```go
// Alert when error rate exceeds threshold
func checkErrorRate() Alert {
    rate := float64(metrics.ErrorsLastHour) / float64(metrics.RequestsLastHour)
    
    if rate > 0.05 { // 5% error rate
        return Alert{
            Severity: "critical",
            Message:  fmt.Sprintf("Error rate %.2f%% exceeds threshold", rate*100),
        }
    }
    return nil
}
```

### Alert on Latency

```go
// Alert when p99 latency exceeds threshold
func checkLatency() Alert {
    p99 := metrics.GetP99Latency()
    
    if p99 > 500*time.Millisecond {
        return Alert{
            Severity: "warning",
            Message:  fmt.Sprintf("P99 latency %s exceeds 500ms", p99),
        }
    }
    return nil
}
```

## Monitoring Checklist

### Pre-Deployment

- [ ] Health endpoint implemented
- [ ] Key metrics identified and instrumented
- [ ] Structured logging in place
- [ ] Log rotation configured

### Post-Deployment

```markdown
## Post-Deployment Monitoring

### Immediate (0-5 min)
- [ ] Application starts without errors
- [ ] Health check returns healthy
- [ ] No error log spikes
- [ ] Database connections established

### Short-Term (5-30 min)
- [ ] Memory usage stable
- [ ] No goroutine leaks
- [ ] Response times normal
- [ ] No unusual errors in logs

### Long-Term (1-24 hours)
- [ ] Memory doesn't grow unbounded
- [ ] No slow query warnings
- [ ] Error rate within normal range
- [ ] All features functioning
```

## Troubleshooting Commands

```bash
# Check running processes
ps aux | grep kariya

# Monitor memory usage
watch -n 5 'ps -o pid,rss,vsz,comm -p $(pgrep kariya)'

# Check open files
lsof -p $(pgrep kariya) | wc -l

# Database size
du -h ~/.kariya/data.db

# Log errors in last hour
grep -c ERROR ~/.kariya/logs/app.log

# Recent error messages
tail -100 ~/.kariya/logs/app.log | grep ERROR

# Disk usage
df -h ~/.kariya
```

## Related Skills

- `incident-response` - When monitoring detects issues
- `logging-observability` - Detailed logging setup
- `rollback-recovery` - When monitoring shows problems
- `fuzz-testing` - Finding issues before they hit production
