# Common Development Tasks

Frequently used commands and troubleshooting.

---

## Common Development Tasks

### Adding a New Intent

1. **Create data structures** (`intent_name.go`):
   ```go
   type YourIntentContext struct { ... }
   type YourIntentResult struct { ... }
   type YourIntentModel struct {
       state YourState
       data  *YourIntentContext
       result *IntentResult[*YourIntentResult]
   }
   ```

2. **Implement intent** (`your_intent_intent.go`):
   ```go
   func (y *YourIntentModel) Init(ctx context.Context) tea.Cmd { ... }
   func (y *YourIntentModel) Update(msg tea.Msg) tea.Cmd { ... }
   func (y *YourIntentModel) View() string { ... }
   func (y *YourIntentModel) Result() *IntentResult[interface{}] { ... }
   ```

3. **Write tests** (`your_intent_test.go`):
   - Use Ginkgo/Gomega framework
   - Test state transitions
   - Test view rendering
   - Test result handling

4. **Register with router** in `internal/cli/app/app.go`:
   ```go
   router.RegisterIntent("your_intent", func() intents.Intent {
       return NewYourIntent(context)
   })
   router.RegisterResultHandler("your_intent", func(result) tea.Cmd {
       // Handle result
   })
   ```

### Creating a New Form

**Time**: 30-45 minutes  
**Prerequisites**: Domain object exists, validators identified  
**Reference**: [`docs/rules/FORMS_WORKFLOW_GUIDE.md`](docs/rules/FORMS_WORKFLOW_GUIDE.md) - Complete step-by-step guide

> **NEW (Recommended): `base.FormScreen[T]`**
>
> For new forms, use `base.FormScreen[T]` instead of creating `models/` wrappers:
> - Automatic window resize handling (rebuilds form)
> - Built-in StandardView integration (breadcrumbs, footer)
> - Returns `screens.ScreenResult` (no custom message types needed)
> - No `huh` import needed -- uses `forms.Form` alias
> - ~20 lines per form screen vs ~100-200 lines per wrapper
>
> **Example**: `internal/cli/screens/skills/skill_form.go` (93 lines)
>
> **See**: [`docs/FORMS_GUIDE.md#formscreen-pattern-target`](docs/FORMS_GUIDE.md#formscreen-pattern-target)

> **⚠️ CRITICAL: Form Alignment Rule**
> 
> Forms used in **intents** MUST NOT use `*huh.Form` directly. Either:
> - **New code**: Use `base.FormScreen[T]` in `screens/{feature}/` (recommended)
> - **Legacy code**: Use a wrapper model (like `CaptureForm`, `SkillForm` in `models/`)
>
> Direct use of `*huh.Form` in intents causes **left-alignment issues** because
> form dimensions are captured at creation time and become stale on terminal resize.
>
> **See**: [`docs/FORMS_GUIDE.md#formscreen-pattern-target`](docs/FORMS_GUIDE.md#formscreen-pattern-target)

#### Quick Workflow

1. **Define FormData structure** (`internal/cli/forms/your_form.go`):
   ```go
   type YourFormData struct {
       Field1          string
       Field2          []string // For MultiSelect
       SubmitConfirmed bool     // Required for confirm button
   }
   ```

2. **Create form builder functions**:
   ```go
   // 3 variants for flexibility
   func NewYourForm(obj *career.YourDomain) *huh.Form { ... }
   func NewYourFormWithData(data *YourFormData) *huh.Form { ... }
   func NewYourFormWithDataAndDimensions(data *YourFormData, width, height int) *huh.Form {
       data.SubmitConfirmed = false
       
       fieldsGroup := huh.NewGroup(
           forms.NewInput(forms.FieldConfig{
               Key:         "field1",
               Title:       "Field 1",
               Validate:    forms.Required,
           }).Value(&data.Field1),
       )
       
       return forms.NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
   }
   ```

3. **Create domain conversion functions**:
   ```go
   func GetYourFormData(obj *career.YourDomain) *YourFormData { ... }
   func ApplyYourFormData(obj *career.YourDomain, data *YourFormData) error { ... }
   ```

4. **Write tests** (`your_form_test.go`):
   - Test form creation
   - Test data extraction (GetYourFormData)
   - Test data application (ApplyYourFormData)
   - Test roundtrip conversion
   - Test nil slice handling

5. **Document in FORMS_GUIDE.md**:
   - Add section to "Form Configurations"
   - Add to forms package table
   - Add to tests section

#### Integration Patterns

**Screen Integration (Recommended for new code)**:

```go
// 1. Create screen in internal/cli/screens/{feature}/form_screen.go
type MyFormScreen struct {
    *base.FormScreen[*forms.MyFormData]
}

func NewMyFormScreen() *MyFormScreen {
    formData := &forms.MyFormData{}
    return &MyFormScreen{
        FormScreen: base.NewBaseFormScreen(
            []string{"Main Menu", "My Feature"},
            forms.NewMyFormWithDataAndDimensions,
            formData,
        ),
    }
}

// 2. Use in intent
i.formScreen = myfeature.NewMyFormScreen()
i.formScreen.SetTerminalInfo(i.Width(), i.Height())
i.activeScreen = i.formScreen
```

**Existing screen example**: `internal/cli/screens/skills/skill_form.go`

**Modal Integration** (inline editing - wrapper NOT required):
```go
type EditYourModal struct {
    original  *career.YourDomain  // Preserved (never mutated)
    modified  *career.YourDomain  // Working copy
    result    *ModalEditResult[*career.YourDomain]
    form      *huh.Form           // Direct use OK in modals
    formData  *forms.YourFormData
}
```

**See Complete Guide**: [`docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern`](docs/FORMS_GUIDE.md#form-alignment-and-the-wrapper-pattern)

### Running Tests

```bash
# Create test file
touch internal/cli/intents/your_test.go

# Add Ginkgo test structure
cat > internal/cli/intents/your_test.go << 'EOF'
package intents_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("YourIntent", func() {
    It("should do something", func() {
        Expect(true).To(BeTrue())
    })
})
EOF

# Run tests
go test -v ./internal/cli/intents/...
```

---

## Deployment Guide

### Branching Strategy

KaRiya uses a **dual-branch workflow** for controlled releases:

- **`next` branch**: Integration branch for feature development
  - All feature branches merge here via PR
  - Full CI runs on every merge
  - Staging environment for testing features together
  
- **`main` branch**: Production releases only
  - Only accepts merges from `next` branch
  - Automatic semantic releases on merge
  - Protected branch with strict checks

**See**: [`docs/BRANCHING_STRATEGY.md`](docs/BRANCHING_STRATEGY.md) for complete workflow documentation.

**Quick Workflow**:
```bash
# 1. Create feature branch from next
git checkout next && git pull
git checkout -b feature/my-feature

# 2. Develop and commit (conventional commits)
git commit -m "feat: add new feature"

# 3. Push and create PR to next (NOT main)
git push -u origin feature/my-feature
gh pr create --base next --fill

# 4. After merge to next, release when ready:
#    Create PR: next → main (triggers automatic release)
```

### Pre-Deployment Checklist

**RECOMMENDED**: Run all CI checks locally before pushing:

```bash
# Run ALL CI checks locally (mirrors GitHub Actions exactly)
make ci-local
```

This single command runs:
- ✅ Commitlint validation
- ✅ AI attribution check
- ✅ Code formatting (go fmt)
- ✅ Static analysis (go vet, staticcheck)
- ✅ Tests with race detector and coverage
- ✅ Multi-platform builds (Linux, macOS, Windows)
- ✅ Security scanning (gosec)

**Individual checks** (if needed):

```bash
# 1. Install all CI tools
make ci-install-tools

# 2. Run full test suite
make test

# 3. Check code quality
make fmt
make vet
make staticcheck

# 4. Security scan
make gosec

# 5. Generate coverage report
make coverage

# 6. Build for target platforms
make build  # Current platform
# OR multi-platform:
GOOS=linux GOARCH=amd64 go build -o kariya-linux-amd64 ./cmd/cli
GOOS=darwin GOARCH=amd64 go build -o kariya-darwin-amd64 ./cmd/cli
GOOS=darwin GOARCH=arm64 go build -o kariya-darwin-arm64 ./cmd/cli
GOOS=windows GOARCH=amd64 go build -o kariya-windows-amd64.exe ./cmd/cli
```

### CI/CD Pipeline

**Main CI Workflow** (`.github/workflows/ci.yml`):
- **commitlint** (PR only): Validates commit messages
- **lint**: Code quality (gofmt, vet, staticcheck)
- **test**: Multi-platform testing (Linux, macOS, Windows)
- **build**: Multi-platform builds with artifact upload
- **security**: Gosec security scanning

**PR Validation Workflow** (`.github/workflows/pr-validation.yml`):
- **validate-pr-title**: PR title follows conventional commits
- **check-ai-attribution**: Commits have AI attribution (if applicable)
- **conventional-commits**: All commits follow conventions
- **breaking-changes**: Detects breaking changes
- **size-label**: Auto-labels PR by size

**See Also**:
- [CI Checks Summary](docs/CI_CHECKS_SUMMARY.md) - Complete CI/local command mapping
- [CI Local Guide](docs/CI_LOCAL_GUIDE.md) - Detailed guide for running CI locally

---

## Performance Benchmarks

All benchmarks passing with excellent performance:

| Operation | Target | Actual | Status |
|-----------|--------|--------|--------|
| Intent Init | 50ms | 0.4ms | ✅ |
| View Render | 100ms | 46ms | ✅ |
| State Transition | 10ms | 0.03ms | ✅ |
| Router Operations | 1ms | 0.1ms | ✅ |
| Test Suite | 5s | 1.3s | ✅ |

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./internal/cli/intents/

# Run specific benchmark
go test -bench=BenchmarkCaptureEventInit ./internal/cli/intents/

# With memory profiling
go test -bench=. -benchmem ./internal/cli/intents/
```

---

## Workflow Patterns

### Intent State Machine Pattern

```go
type YourState string
const (
    StateInitial YourState = "initial"
    StateWorking YourState = "working"
    StateFinal   YourState = "final"
)

func (y *YourIntentModel) Update(msg tea.Msg) tea.Cmd {
    switch y.state {
    case StateInitial:
        return y.handleInitial(msg)
    case StateWorking:
        return y.handleWorking(msg)
    case StateFinal:
        return y.handleFinal(msg)
    }
    return nil
}

func (y *YourIntentModel) View() string {
    switch y.state {
    case StateInitial:
        return y.viewInitial()
    case StateWorking:
        return y.viewWorking()
    case StateFinal:
        return y.viewFinal()
    }
    return ""
}
```

### Modal Sub-Flow Pattern

```go
type ModalEditResult[T any] struct {
    Original T
    Modified T
    Accepted bool
    Changes  map[string]interface{}
}

// Usage in intent
if editModal.WasAccepted() {
    y.data.Field = editModal.Modified
} else {
    y.data.Field = editModal.Original
}
```

### Back Navigation with Context Preservation

```go
// When returning result
result := &IntentResult[*YourIntentResult]{
    Status: StatusCompleted,
    Data:   y.result.Data,
}
result.WithMetadata("scroll_position", 42)
result.WithMetadata("selection", selectedID)
return result

// When restoring context
if pos, ok := prevResult.GetMetadata("scroll_position"); ok {
    y.scrollPos = pos.(int)
}
```

---

## Troubleshooting

### Issue: Tests Failing with "Multiple Ginkgo Entry Points"

**Solution**: Ensure only one `_test.go` file per package uses Ginkgo.

### Issue: Race Conditions Detected

**Solution**:
1. Run `go test -race ./...` to identify
2. Add proper synchronization (mutex, channels)
3. Check GlobalContext usage in `internal/cli/context/global.go`

### Issue: High Memory Usage

**Solution**:
1. Check for goroutine leaks with pprof
2. Verify proper cleanup in intent `Init()` and result handling
3. Check repository queries for N+1 issues

### Issue: Slow Tests

**Solution**:
1. Run benchmarks: `go test -bench=. ./internal/cli/intents/`
2. Profile with pprof: `go test -cpuprofile=cpu.prof ./...`
3. Check database queries in repository layer

### Additional Troubleshooting Resources

- **General**: [`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md)
- **CV Issues**: [`docs/guides/CV_TROUBLESHOOTING.md`](docs/guides/CV_TROUBLESHOOTING.md)
- **Error Handling**: [`docs/guides/ERROR_HANDLING_GUIDE.md`](docs/guides/ERROR_HANDLING_GUIDE.md)

---
