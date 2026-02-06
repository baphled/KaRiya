# /e2e - Write E2E Tests

Write end-to-end tests for TUI workflows using KaRiya's TestEnv harness.

## Skills to Load

- `e2e-testing` - TestEnv patterns and navigation helpers
- `ginkgo-gomega` - BDD testing with Ginkgo
- `test-fixtures` - Factory patterns for test data
- `bubble-tea-expert` - Understanding what's being tested

## Usage

```
/e2e <feature or workflow to test>
```

## Examples

```
/e2e capture workflow
/e2e timeline navigation
/e2e escape behavior for burst management
/e2e filter modal interactions
```

## Process

1. **Identify test file location**
   - `internal/testutil/e2e/{feature}_e2e_test.go` for workflows
   - `internal/testutil/e2e/{feature}_navigation_test.go` for navigation
   - `internal/testutil/e2e/{feature}_escape_test.go` for escape behavior

2. **Check for existing suite**
   - Verify `e2e_suite_test.go` exists
   - Don't create duplicate `TestXxx` entry points

3. **Set up test structure**
   ```go
   var _ = Describe("Feature Workflow", func() {
       var env *e2e.TestEnv
       
       BeforeEach(func() {
           env = e2e.GetSharedEnv()
           env.Reset()
       })
       
       // Tests...
   })
   ```

4. **Create test data with fixtures**
   ```go
   BeforeEach(func() {
       env.Reset()
       events := fixtures.Events(5)
       for _, e := range events {
           env.AddEvent(e)
       }
   })
   ```

5. **Write tests using navigation helpers**
   ```go
   It("completes the workflow", func() {
       env.SelectIntentByName("Feature Name")
       env.NavigateDown()
       env.Confirm()
       
       env.AssertViewContains("Expected content")
   })
   ```

6. **Run and verify**
   ```bash
   go test ./internal/testutil/e2e/... -v -run "Feature"
   ```

## Test Types

| Type | File Pattern | Purpose |
|------|--------------|---------|
| Workflow | `*_e2e_test.go` | Full user workflows |
| Navigation | `*_navigation_test.go` | List navigation, selection |
| Escape | `*_escape_test.go` | Cancel/back behavior |
| State Machine | `*_state_machine_test.go` | State transitions |

## Inline Comments Exception

E2E tests are the ONLY place where inline comments are allowed:

```go
env.SelectIntentByName("Capture") // Navigate to Capture
env.NavigateDown()                 // Select Manual mode
env.Confirm()                      // Confirm selection
```
