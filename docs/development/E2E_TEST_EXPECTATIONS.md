# E2E Test Expectations for CaptureEvent

**Created**: 2026-01-14  
**Purpose**: Clarify what the E2E test should validate based on ACTUAL documented workflow  
**Reference**: `docs/workflows/EVENT_CAPTURE_WORKFLOW.md`

---

## TL;DR - What We Got Wrong

**Initial misunderstanding**: We thought the workflow should be:
```
Form → Review → Submit → Save → Enrichment → Post-Save Review → Complete
```

**ACTUAL documented workflow** (line 427-430):
```
Choose Strategy → Form → Review → Submit → Success
```

**Key insight**: Enrichment happens **AFTER save** but it's **OPTIONAL** and **USER-TRIGGERED**, not automatic!

---

## The Actual Implementation (As-Built)

From `docs/workflows/EVENT_CAPTURE_WORKFLOW.md` lines 420-442:

### Three Supported Paths

#### 1. Minimal Path (Quick Submit)
```
Choose Strategy (Quick) → Form → Ctrl+S → Submit → Success
```
- User presses Ctrl+S from form
- Skips pre-save review entirely
- Event saves directly

#### 2. Standard Path (With Review)
```
Choose Strategy → Form → Review → Submit → Success
```
- User fills form and presses Enter
- Reviews event details before saving
- Submits and completes
- **This is the MOST COMMON path**

#### 3. Full Path (With Enrichment Editing)
```
Choose Strategy → Form → Review → Edit Metadata → Review →
    Edit Bursts → Review → Edit Facts → Review → Submit → Success
```
- User manually triggers enrichment with 'b' or 'f' keys
- Edits suggested bursts/facts
- Returns to review after each edit
- Submits when satisfied

### Key Documentation Points

**From line 438-442**:
- Bursts and facts are **OPTIONAL** and **USER-TRIGGERED** during Review
- User can skip Review entirely with `Ctrl+S` from Form
- Metadata editing happens in modal overlays, not during initial capture
- **Enrichment (burst/fact extraction) happens AFTER event save, not before**

**From line 451-460** (explaining the difference from PRD):
- Burst suggestion is **optional** (triggered by 'b' key in Review)
- Fact enrichment is **optional** (triggered by 'f' key in Review)  
- Both happen **during review**, not automatically
- User can bypass both and go straight to Submit

**Why?**
- User control: Users can capture events quickly without waiting for AI processing
- Performance: Enrichment can be done later in bulk
- Flexibility: Not all events need bursts/facts immediately

---

## What "AFTER event save" Actually Means

The statement "Enrichment happens AFTER event save, not before" (line 442) means:

### NOT This (What We Thought):
```
❌ Form → Review → Submit → Save → [Automatic Enrichment] → Post-Save Review → Complete
```

### But This (What It Actually Means):
```
✅ Form → Review → Submit → Save → Complete
   Later: User manually navigates to BrowseTimeline → Views event → Triggers enrichment
```

**OR** (if user triggers enrichment BEFORE submit):
```
✅ Form → Review → [User presses 'b' or 'f'] → Enrichment happens → Edit bursts/facts → Submit → Complete
```

The "AFTER save" refers to the architectural decision that enrichment:
- Operates on **persisted events** (not in-memory draft events)
- Can be triggered **at any time** after the event exists in DB
- Is **not automatic** during the capture workflow

---

## What the E2E Test SHOULD Validate

### Test 1: Standard Path (MOST IMPORTANT)
```go
It("should complete standard capture workflow", func() {
    // Choose Strategy → Form → Review → Submit → Success
    env.SelectIntentByName("capture_event")
    env.Confirm() // Select Quick strategy
    
    // Fill form
    env.FillFormFields(...) // Need to implement this helper
    
    // Go to review
    env.SubmitForm() // Navigate to Submit button and press it
    
    // Should show Review state
    env.AssertViewContains("Review")
    
    // Submit event
    env.Confirm() // Press Enter to submit
    
    // Should show Success state or return to main menu
    env.AssertViewContainsAny("Success", "saved successfully", "Capture Event")
    
    // Verify event was saved
    env.AssertEventCount(1)
})
```

### Test 2: Quick Submit Path
```go
It("should support quick submit with Ctrl+S", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm() // Select Quick strategy
    
    // Fill form
    env.FillFormFields(...)
    
    // Press Ctrl+S to skip review
    env.PressKey(tea.KeyCtrlS)
    
    // Should go directly to Submit (skip Review)
    env.AssertViewContainsAny("Saving", "Success")
    
    // Verify event saved
    env.AssertEventCount(1)
})
```

### Test 3: Enrichment is Optional (No Auto Post-Save Review)
```go
It("should NOT automatically show post-save enrichment review", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm()
    
    // Fill and submit
    env.FillFormFields(...)
    env.SubmitForm()
    env.Confirm() // Submit from review
    
    // After success, should return to MAIN MENU (not post-save review)
    // This is CORRECT behavior according to docs
    env.AssertViewContains("Capture Event") // Back at main menu
    
    // Event saved
    env.AssertEventCount(1)
    
    // Bursts/facts may or may not be inferred (implementation detail)
    // But user is NOT shown a post-save review screen
})
```

### Test 4: Pre-Save Enrichment (User-Triggered)
```go
It("should allow user to trigger enrichment before submitting", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm()
    
    // Fill and go to review
    env.FillFormFields(...)
    env.SubmitForm()
    
    // Should be in Review state
    env.AssertViewContains("Review")
    
    // User presses 'b' to trigger burst suggestion
    env.PressKeyRune('b')
    
    // Should show burst modal or burst editing UI
    env.AssertViewContainsAny("Burst", "Suggest")
    
    // User accepts/edits bursts, returns to review
    env.Confirm()
    
    // Now submit
    env.Confirm()
    
    // Event saved with associated bursts
    env.AssertEventCount(1)
})
```

### Test 5: Cancel Behavior
```go
It("should not persist event when cancelled", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm()
    
    env.FillFormFields(...)
    env.Cancel() // Cancel before submit
    
    // Should return to main menu without saving
    env.AssertEventCount(0)
})
```

### Test 6: Escape Navigation
```go
It("should navigate back through states with Escape", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm() // Choose strategy
    
    // At form, press Esc → go back to Choose Strategy
    env.Cancel()
    env.AssertViewContains("Strategy")
    
    // Press Esc again → go back to Main Menu
    env.Cancel()
    env.AssertViewContains("Capture Event")
})
```

---

## What Our Current E2E Test Got Wrong

### Incorrect Expectation (Lines 94-108)

```go
// Step 9: CRITICAL CHECK - Should return to Review state (post-save)
// Expected: Review state with enriched data (bursts and facts)
// NOT expected: Going back to main menu immediately
```

**This is WRONG**. According to the documentation:
- ✅ Going back to main menu IS the correct behavior
- ❌ Automatic post-save review does NOT exist in the documented workflow
- ✅ Enrichment is optional and user-triggered, not automatic

### Correct Expectation

```go
// Step 9: After success, should return to main menu
view = env.GetView()
Expect(view).To(ContainSubstring("Capture Event"), 
    "After successful save, should return to main menu")

// Event should be persisted
env.AssertEventCount(1)

// Post-save enrichment review does NOT happen automatically
// User can navigate to BrowseTimeline later to trigger enrichment
```

---

## The Confusion: Where "Post-Save Review" Came From

### The Document Said (Line 442):
> "Enrichment (burst/fact extraction) happens AFTER event save, not before"

### We Interpreted This As:
"After saving, the user is automatically shown a review screen with enriched bursts/facts"

### What It Actually Means:
"Architecturally, enrichment operates on saved events (from DB), not on draft events in memory. Users can trigger enrichment at any time after an event is saved, either:
- BEFORE submitting (press 'b' or 'f' in Review state)
- AFTER submitting (navigate to BrowseTimeline, view event, trigger enrichment)"

---

## The Manual Trigger Workflow (Documented on Lines 317-327)

From the Review state keyboard shortcuts:

| Key | Action | Result |
|-----|--------|--------|
| `Ctrl+S` or `Enter` | Confirm | Submit event |
| `e` | Edit metadata | Open metadata editor modal |
| **`b`** | **Edit bursts** | **Open burst suggestion modal** |
| **`f`** | **Edit facts** | **Open fact editor modal** |
| `Esc` | Back | Return to form |

**This is where enrichment happens** - User manually presses 'b' or 'f' in the Review state (BEFORE submitting), NOT automatically after submitting.

---

## What We Should Actually Test

### Core Happy Path (PRIORITY 1)
```
Choose Strategy → Form → Review → Submit → Success → Main Menu
```
**Validates**: Basic capture flow works, event persists, completes cleanly

### Quick Submit (PRIORITY 2)
```
Choose Strategy → Form → Ctrl+S → Submit → Success → Main Menu
```
**Validates**: Skip review shortcut works

### Manual Enrichment Trigger (PRIORITY 3)
```
Choose Strategy → Form → Review → Press 'b' → Edit Bursts → Review → Submit
```
**Validates**: User-triggered enrichment works

### Cancel/Escape (PRIORITY 4)
```
Various cancel/escape scenarios
```
**Validates**: Navigation and non-persistence of cancelled workflows

---

## Recommended E2E Test Rewrite

```go
var _ = Describe("CaptureEvent E2E Workflow", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.Setup(GinkgoT())
    })

    AfterEach(func() {
        env.Cleanup()
    })

    Describe("Standard Capture Workflow", func() {
        It("should complete: Choose → Form → Review → Submit → Main Menu", func() {
            env.SelectIntentByName("capture_event")
            env.Confirm() // Quick strategy

            // Fill form (need helper for huh form submission)
            // TODO: Implement env.FillCaptureForm(text, date)
            
            // Review
            // TODO: Assert review state
            
            // Submit
            env.Confirm()
            
            // Should complete and return to main menu
            view := env.GetView()
            Expect(view).To(ContainSubstring("Capture Event"))
            
            // Event persisted
            env.AssertEventCount(1)
        })
    })

    Describe("Quick Submit Path", func() {
        It("should skip review with Ctrl+S", func() {
            // TODO: Implement
        })
    })

    Describe("Cancel Behavior", func() {
        It("should not persist when cancelled", func() {
            // TODO: Implement
        })
    })

    // Manual enrichment testing can be added later
    // This requires implementing burst/fact modal interactions
})
```

---

## Action Items

1. **Fix Current E2E Test**:
   - Remove expectations for automatic post-save review
   - Expect main menu return after successful save
   - Focus on standard path first

2. **Implement Form Submission Helper**:
   - `env.FillCaptureForm(text, date)` - Fills form and submits
   - Handles huh form interaction properly

3. **Start Simple**:
   - Get standard path working first
   - Add quick submit second
   - Add manual enrichment later (nice-to-have)

4. **Update Documentation**:
   - Clarify that post-save review is NOT automatic
   - Explain enrichment is user-triggered
   - Remove any references to automatic post-save enrichment

---

## Summary

**What We Thought**: Automatic post-save enrichment review is missing  
**Reality**: Post-save enrichment review was never part of the workflow  
**Actual Flow**: Choose → Form → Review → Submit → Main Menu  
**Enrichment**: User-triggered via 'b' or 'f' keys in Review (before submit) or later in BrowseTimeline  

**E2E Test Should Validate**: The actual documented workflow, not an imagined one.
