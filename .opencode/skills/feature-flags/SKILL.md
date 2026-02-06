---
name: feature-flags
description: Safe feature rollouts using feature flags, gradual releases, and A/B testing
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide the implementation and use of feature flags for safe feature rollouts, gradual releases, and quick rollbacks.

## When to use me

Use this skill when:
- Rolling out risky new features
- Need ability to quickly disable features
- Implementing A/B testing
- Gradual rollout to users
- Managing beta features

## Feature Flag Fundamentals

### When to Use Feature Flags

| Scenario | Use Flag? | Type |
|----------|-----------|------|
| New risky feature | Yes | Release flag |
| A/B test | Yes | Experiment flag |
| Ops kill switch | Yes | Ops flag |
| Bug fix | No | Direct deployment |
| Small refactor | No | Direct deployment |
| Config change | Maybe | Config, not flag |

### Flag Lifecycle

```
Development → Rollout → 100% Enabled → Code Cleanup → Flag Removed
     ↑                                       ↓
     └───────────── Rollback ────────────────┘
```

## Implementation Pattern

### Basic Feature Flag System

```go
// internal/features/flags.go
package features

import (
    "sync"
)

type Flag string

const (
    NewTimelineUI    Flag = "new_timeline_ui"
    BurstDetection   Flag = "burst_detection"
    ExperimentalCV   Flag = "experimental_cv"
)

type Flags struct {
    mu      sync.RWMutex
    enabled map[Flag]bool
}

var global = &Flags{
    enabled: make(map[Flag]bool),
}

// IsEnabled checks if a feature flag is enabled.
//
// Expected:
//   flag - The feature flag to check
//
// Returns:
//   bool indicating whether the flag is enabled
//
// Side effects:
//   None - read-only operation
func IsEnabled(flag Flag) bool {
    global.mu.RLock()
    defer global.mu.RUnlock()
    return global.enabled[flag]
}

// Enable turns on a feature flag.
//
// Expected:
//   flag - The feature flag to enable
//
// Returns:
//   Nothing
//
// Side effects:
//   Modifies global flag state
func Enable(flag Flag) {
    global.mu.Lock()
    defer global.mu.Unlock()
    global.enabled[flag] = true
}

// Disable turns off a feature flag.
//
// Expected:
//   flag - The feature flag to disable
//
// Returns:
//   Nothing
//
// Side effects:
//   Modifies global flag state
func Disable(flag Flag) {
    global.mu.Lock()
    defer global.mu.Unlock()
    global.enabled[flag] = false
}

// LoadFromConfig loads flags from configuration.
//
// Expected:
//   config - Map of flag names to enabled states
//
// Returns:
//   Nothing
//
// Side effects:
//   Replaces all flag states with config values
func LoadFromConfig(config map[string]bool) {
    global.mu.Lock()
    defer global.mu.Unlock()
    
    for name, enabled := range config {
        global.enabled[Flag(name)] = enabled
    }
}
```

### Configuration-Based Flags

```yaml
# config/features.yaml
features:
  new_timeline_ui: false
  burst_detection: true
  experimental_cv: false
```

```go
// Load from config
func loadFeatures() {
    data, _ := os.ReadFile("config/features.yaml")
    
    var config struct {
        Features map[string]bool `yaml:"features"`
    }
    yaml.Unmarshal(data, &config)
    
    features.LoadFromConfig(config.Features)
}
```

### Environment-Based Flags

```go
// Load from environment
func loadFeaturesFromEnv() {
    if os.Getenv("KARIYA_FEATURE_BURST_DETECTION") == "true" {
        features.Enable(features.BurstDetection)
    }
}
```

## Usage Patterns

### Pattern 1: Simple Boolean Check

```go
func (i *TimelineIntent) View() string {
    if features.IsEnabled(features.NewTimelineUI) {
        return i.newView()
    }
    return i.legacyView()
}
```

### Pattern 2: Strategy Pattern

```go
type EventProcessor interface {
    Process(ctx context.Context, event *Event) error
}

func getProcessor() EventProcessor {
    if features.IsEnabled(features.BurstDetection) {
        return &BurstAwareProcessor{}
    }
    return &StandardProcessor{}
}
```

### Pattern 3: Gradual Rollout

```go
// Percentage-based rollout
func shouldEnableForUser(userID string, flag Flag, percentage int) bool {
    if !features.IsEnabled(flag) {
        return false
    }
    
    // Deterministic hash for consistent experience
    hash := fnv.New32a()
    hash.Write([]byte(userID + string(flag)))
    return int(hash.Sum32()%100) < percentage
}

// Usage
if shouldEnableForUser(user.ID, features.NewTimelineUI, 25) {
    // 25% of users get new UI
}
```

### Pattern 4: Kill Switch

```go
func ProcessEvents(ctx context.Context, events []*Event) error {
    // Kill switch for emergency disable
    if features.IsEnabled(features.DisableEventProcessing) {
        return errors.New("event processing temporarily disabled")
    }
    
    return processEventsInternal(ctx, events)
}
```

## Testing with Feature Flags

### Unit Test Pattern

```go
func TestTimelineView_WithNewUI(t *testing.T) {
    // Enable flag for this test
    features.Enable(features.NewTimelineUI)
    defer features.Disable(features.NewTimelineUI)
    
    intent := NewTimelineIntent()
    view := intent.View()
    
    assert.Contains(t, view, "new-timeline-header")
}

func TestTimelineView_WithLegacyUI(t *testing.T) {
    // Ensure flag is disabled
    features.Disable(features.NewTimelineUI)
    
    intent := NewTimelineIntent()
    view := intent.View()
    
    assert.Contains(t, view, "legacy-timeline-header")
}
```

### Test Helper

```go
// test/features.go
func WithFeature(t *testing.T, flag features.Flag, enabled bool, fn func()) {
    t.Helper()
    
    // Save original state
    original := features.IsEnabled(flag)
    
    // Set test state
    if enabled {
        features.Enable(flag)
    } else {
        features.Disable(flag)
    }
    
    // Restore after test
    defer func() {
        if original {
            features.Enable(flag)
        } else {
            features.Disable(flag)
        }
    }()
    
    fn()
}

// Usage
func TestWithFeatureHelper(t *testing.T) {
    WithFeature(t, features.BurstDetection, true, func() {
        // Test with feature enabled
    })
}
```

## Flag Management

### Adding a New Flag

```go
// 1. Define the flag constant
const (
    MyNewFeature Flag = "my_new_feature"
)

// 2. Add to config/features.yaml
// my_new_feature: false

// 3. Use in code with graceful fallback
func doSomething() {
    if features.IsEnabled(features.MyNewFeature) {
        newImplementation()
    } else {
        oldImplementation()
    }
}

// 4. Add tests for both paths
// 5. Deploy with flag disabled
// 6. Enable flag in config
// 7. Monitor
// 8. Clean up old code path when stable
```

### Removing a Flag (Cleanup)

```bash
# 1. Verify flag is enabled in all environments
grep -r "my_new_feature" config/

# 2. Search for all usages
grep -r "MyNewFeature" internal/

# 3. Remove feature checks, keep new implementation
# 4. Remove flag constant
# 5. Remove from config
# 6. Delete tests for old code path
```

### Cleanup Checklist

```markdown
## Feature Flag Cleanup: [FlagName]

**Flag:** `my_new_feature`
**Added:** YYYY-MM-DD
**Cleanup Target:** YYYY-MM-DD

### Pre-Cleanup
- [ ] Flag enabled 100% for 2+ weeks
- [ ] No incidents related to feature
- [ ] Metrics show feature working correctly

### Cleanup Steps
- [ ] Remove all `IsEnabled` checks for this flag
- [ ] Remove old code path
- [ ] Remove flag constant
- [ ] Remove from config files
- [ ] Update/remove tests
- [ ] Update documentation

### Verification
- [ ] Tests pass
- [ ] No references to flag remain
- [ ] Feature still works
```

## Monitoring Flags

### Flag Status Dashboard

```go
// Health endpoint showing flag status
func handleFlagStatus(w http.ResponseWriter, r *http.Request) {
    status := make(map[string]bool)
    
    for _, flag := range []features.Flag{
        features.NewTimelineUI,
        features.BurstDetection,
        features.ExperimentalCV,
    } {
        status[string(flag)] = features.IsEnabled(flag)
    }
    
    json.NewEncoder(w).Encode(status)
}
```

### CLI Command

```go
// kariya features list
var featuresCmd = &cobra.Command{
    Use:   "features",
    Short: "Manage feature flags",
}

var listFeaturesCmd = &cobra.Command{
    Use:   "list",
    Short: "List all feature flags",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Feature Flags:")
        fmt.Printf("  new_timeline_ui:  %v\n", features.IsEnabled(features.NewTimelineUI))
        fmt.Printf("  burst_detection:  %v\n", features.IsEnabled(features.BurstDetection))
        fmt.Printf("  experimental_cv:  %v\n", features.IsEnabled(features.ExperimentalCV))
    },
}
```

## Best Practices

### DO

- Keep flags short-lived (remove after rollout)
- Name flags descriptively
- Test both flag states
- Document flag purpose
- Set expiration dates for cleanup
- Use kill switches for critical features

### DON'T

- Nest flag checks deeply
- Create dependencies between flags
- Leave flags forever (tech debt)
- Use flags for permanent configuration
- Toggle flags rapidly in production

### Flag Naming Convention

```
[category]_[feature]_[detail]

Examples:
- release_new_timeline_ui
- experiment_cv_layout_v2
- ops_disable_notifications
- beta_burst_detection
```

## Emergency Procedures

### Quick Disable

```bash
# Via environment variable
export KARIYA_FEATURE_BURST_DETECTION=false
kariya restart

# Via config change
sed -i 's/burst_detection: true/burst_detection: false/' config/features.yaml
kariya reload-config

# Via CLI (if implemented)
kariya features disable burst_detection
```

### Rollback Procedure

```markdown
1. Disable the flag immediately
2. Verify feature is disabled
3. Investigate root cause
4. Fix the issue
5. Re-enable flag for small percentage
6. Monitor
7. Gradually increase percentage
```

## Related Skills

- `rollback-recovery` - Quick rollbacks via flags
- `monitoring` - Monitoring feature health
- `release-management` - Coordinating flag-based releases
- `configuration-management` - Managing flag configuration
