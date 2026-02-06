---
description: Develop a feature using BDD smallest-change workflow with scenario first then incremental implementation
agent: dev
---

# /bdd - BDD Feature Development

Develop a feature using BDD smallest-change workflow: scenario first, then incremental implementation.

## Skills to Load

- `cucumber` - BDD scenarios and smallest-change workflow
- `ginkgo-gomega` - KaRiya's test framework
- `tdd-workflow` - Red-Green-Refactor cycle
- `clean-code` - Refactoring phase

## Usage

```
/bdd <feature description>
```

## Examples

```
/bdd user can filter events by company
/bdd quick capture saves event to timeline
/bdd escape from detail modal returns to list
/bdd skill linking updates both event and skill
```

## Process

### 1. Write the Scenario First

```gherkin
Scenario: Filter events by company
  Given I have events from "Acme Corp" and "Other Inc"
  When I filter by company "Acme Corp"
  Then I only see events from "Acme Corp"
```

Translated to Ginkgo:

```go
It("filters events by company", func() {
    env.AddEvent(fixtures.EventWith(fixtures.WithCompany("Acme Corp")))
    env.AddEvent(fixtures.EventWith(fixtures.WithCompany("Other Inc")))
    
    env.SelectIntentByName("Browse Timeline")
    env.OpenFilterModal()
    env.SelectCompany("Acme Corp")
    env.Confirm()
    
    env.AssertViewContains("Acme Corp")
    env.AssertViewNotContains("Other Inc")
})
```

### 2. Run - See It Fail (RED)

```bash
go test ./internal/testutil/e2e/... -v -run "filters events"
# FAIL - OpenFilterModal doesn't exist
```

### 3. Smallest Change to Pass ONE Thing

```go
// Add empty method - just enough to compile
func (e *TestEnv) OpenFilterModal() {
    // Empty for now
}
```

### 4. Run Again

```bash
go test ./... -run "filters events"
# FAIL - Filter modal doesn't show
```

### 5. Next Smallest Change

```go
func (e *TestEnv) OpenFilterModal() {
    e.PressKeyRune('f')  // Open filter
}
```

### 6. Repeat Until GREEN

Each cycle:
1. Run test
2. See what fails
3. Make smallest fix
4. Run test again

### 7. Refactor (CLEAN)

Once all steps pass:
- Extract duplication
- Improve naming
- Clean up code
- Keep tests green!

## Commit Cadence

```
feat(timeline): add scenario for company filter [RED]
feat(timeline): add filter modal trigger [GREEN step 1]
feat(timeline): add company selector [GREEN step 2]
feat(timeline): apply company filter [GREEN all steps]
refactor(timeline): extract filter logic [REFACTOR]
```

## Key Principles

1. **Scenario defines "done"** - Don't add features not in scenario
2. **One step at a time** - Resist urge to implement everything
3. **Smallest change** - Just enough to pass current step
4. **Run tests constantly** - After every change
5. **Refactor only when green** - Never during RED phase
