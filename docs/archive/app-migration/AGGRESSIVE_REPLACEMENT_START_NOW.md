---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Aggressive Replacement - Start NOW! 🚀

**Status**: READY TO BEGIN
**Timeline**: 3 weeks
**Approach**: Complete replacement, no transition
**First Task**: Phase 1 - Preparation (9 hours)

---

## What's Happening

You're replacing the entire screen-based architecture (1766 lines of app.go with 31 screens and 37+ model fields) with a clean, intent-driven system (10 intents, ~200 lines of app.go).

**Why**:
- ✅ You're the sole user
- ✅ All 5 core intents already exist and work
- ✅ No backward compatibility needed
- ✅ Much cleaner and simpler
- ✅ Takes only 3 weeks

---

## Week 1: Preparation & Create Missing Intents

### Days 1-2: Phase 1 Preparation (9 hours)

#### Task 1.1: Audit Legacy Models (4 hours)

**What to do**:
1. List all 29 legacy screen models:
   ```bash
   ls -la internal/cli/models/*.go | grep -v test
   ```

2. For each model, document:
   - What screen it handles
   - What features it provides
   - What data it displays
   - What actions it supports

3. Create a mapping spreadsheet:
   ```
   Model Name | Screen | Features | Maps To Intent | Status
   ========================================================
   FormModel | CaptureScreen | Create event | CaptureEvent ✅
   ListModel | ListScreen | List events | BrowseTimeline ✅
   ...
   ```

**Output**: `docs/LEGACY_MODEL_AUDIT.md` (reference document)

#### Task 1.2: Identify Intent Gaps (3 hours)

**What to do**:
1. Review the mapping above
2. Identify screens without intents:
   - BurstListScreen → Needs BurstManagement intent
   - FactListScreen → Needs FactManagement intent
   - ImportReviewScreen → Needs ImportWizard intent
   - MetadataReviewScreen → Needs MetadataEditor intent
   - BulkOperationsScreen → Needs BulkOperations intent

3. Document each missing intent:
   - What screens it replaces
   - What states it needs
   - What data it handles

**Output**: `docs/MISSING_INTENTS_SPEC.md` (specification)

#### Task 1.3: Plan Intent Implementation (2 hours)

**What to do**:
1. For each missing intent, outline:
   - States (e.g., StateListBursts, StateEditBurst, StateConfirm)
   - Context data needed
   - Result data to return
   - Views to render

2. Estimate effort for each:
   - BurstManagement: 16 hours
   - FactManagement: 16 hours
   - ImportWizard: 12 hours
   - MetadataEditor: 10 hours
   - BulkOperations: 10 hours

**Output**: Implementation plan ready for Phase 2

---

### Days 3-7: Phase 2 Create Missing Intents (64 hours)

#### Task 2.1: BurstManagement Intent (16 hours)

**What to do**:
1. Create `internal/cli/intents/burst_management.go`:
   - Define BurstManagementState enum
   - Define BurstManagementContext
   - Define BurstManagementResult
   - Define BurstManagementModel

2. Create `internal/cli/intents/burst_management_intent.go`:
   - Implement Init()
   - Implement Update() with state transitions
   - Implement View() with state-specific rendering
   - Implement Result()

3. Create `internal/cli/intents/burst_management_test.go`:
   - State transition tests
   - View rendering tests
   - Result handling tests
   - Integration tests

4. Register with router in app.go

**Effort**: 16 hours (2 days)
**Tests**: 30+ specs
**Coverage**: >90%

#### Task 2.2: FactManagement Intent (16 hours)

**What to do**:
1. Create `internal/cli/intents/fact_management.go`
2. Create `internal/cli/intents/fact_management_intent.go`
3. Create `internal/cli/intents/fact_management_test.go`
4. Register with router

**Effort**: 16 hours (2 days)
**Tests**: 40+ specs
**Coverage**: >90%

#### Task 2.3: ImportWizard Intent (12 hours)

**What to do**:
1. Create `internal/cli/intents/import_wizard.go`
2. Create `internal/cli/intents/import_wizard_intent.go`
3. Create `internal/cli/intents/import_wizard_test.go`
4. Register with router

**Effort**: 12 hours (1.5 days)
**Tests**: 25+ specs
**Coverage**: >90%

#### Task 2.4: MetadataEditor Intent (10 hours)

**What to do**:
1. Create `internal/cli/intents/metadata_editor.go`
2. Create `internal/cli/intents/metadata_editor_intent.go`
3. Create `internal/cli/intents/metadata_editor_test.go`
4. Register with router

**Effort**: 10 hours (1.25 days)
**Tests**: 20+ specs
**Coverage**: >90%

#### Task 2.5: BulkOperations Intent (10 hours)

**What to do**:
1. Create `internal/cli/intents/bulk_operations.go`
2. Create `internal/cli/intents/bulk_operations_intent.go`
3. Create `internal/cli/intents/bulk_operations_test.go`
4. Register with router

**Effort**: 10 hours (1.25 days)
**Tests**: 20+ specs
**Coverage**: >90%

**End of Phase 2**:
- ✅ All 10 intents created
- ✅ 164+ tests passing
- ✅ Ready for Phase 3

---

## Week 2: Rebuild app.go

### Days 1-2: Phase 3 Rebuild app.go (14 hours)

#### Task 3.1: Write New Minimal app.go (8 hours)

**What to do**:
1. Start fresh with `internal/cli/app/app.go`
2. Keep only what's essential:
   - Minimal constants (HomeScreen, HelpScreen, MainMenuScreen, QuitScreen)
   - Minimal model fields (~7 total)
   - Simple Update() method (~30 lines)
   - Simple View() method (~20 lines)
   - Menu selection handler
   - Intent activation handler

3. Template structure:
   ```go
   // Imports (minimal)
   // Constants (4 screens only)
   // Model struct (5-7 fields)
   // NewModel() (50 lines)
   // Init() (5 lines)
   // Update() (30 lines)
   // View() (20 lines)
   // Helper methods (handleMenuSelection, activateIntent, etc.)
   ```

4. Register all 10 intents with router

**Effort**: 8 hours
**Result**: ~200 lines total (vs 1766 before)

#### Task 3.2: Delete Legacy Code (2 hours)

**What to do**:
1. Remove all old screen constants
2. Remove all old model fields
3. Remove all old Update() logic
4. Remove all old View() logic
5. Remove all old helper methods

**Effort**: 2 hours
**Result**: Clean slate, no legacy pollution

#### Task 3.3: Update Service Layer (4 hours)

**What to do**:
1. Review service interfaces
2. Remove screen-specific methods
3. Consolidate interfaces
4. Ensure all intents have dependencies

**Effort**: 4 hours

**End of Phase 3**:
- ✅ New app.go complete (~200 lines)
- ✅ All legacy code removed
- ✅ Code compiles cleanly

---

### Days 3-5: Phase 4 Testing & Validation (34 hours)

#### Task 4.1: Integration Testing (20 hours)

**What to do**:
1. Test all intents activate from menu
2. Test complete workflows for each intent
3. Test result handling
4. Test navigation (back, main menu, etc.)
5. Test edge cases

**Effort**: 20 hours
**Result**: All workflows functional

#### Task 4.2: Performance & Quality (14 hours)

**What to do**:
1. Run benchmarks (2 hours)
2. Performance profiling (4 hours)
3. Code review & cleanup (4 hours)
4. Documentation update (4 hours)

**Effort**: 14 hours
**Result**: Optimized, documented, production-ready

**End of Phase 4**:
- ✅ All tests passing (164+)
- ✅ All workflows functional
- ✅ Performance optimized
- ✅ Documentation updated

---

## Week 3: Final Polish & Merge

### Days 1-3: Cleanup & Final Testing (30 hours)

#### Task 5.1: Final Integration Testing (10 hours)

**What to do**:
1. Full end-to-end testing of all features
2. Test on clean system
3. Verify no edge cases broken
4. Manual testing of all intents

#### Task 5.2: Code Quality Review (8 hours)

**What to do**:
1. Run linter: `golangci-lint run ./...`
2. Run formatter: `gofmt -l ./...`
3. Run tests: `go test -v ./...`
4. Run race detector: `go test -race ./...`

#### Task 5.3: Documentation (8 hours)

**What to do**:
1. Update README.md with new architecture
2. Update CLI_GUIDE.md with intent descriptions
3. Update ARCHITECTURE.md
4. Create migration summary

#### Task 5.4: Merge & Deploy (4 hours)

**What to do**:
1. Final code review
2. Merge feature branch to main
3. Tag release
4. Deploy

---

## Daily Checklist

### Each Day

```
Morning (30 min):
- [ ] Review today's tasks
- [ ] Check no blockers from yesterday
- [ ] Plan commits

Work (6-8 hours):
- [ ] Implement assigned tasks
- [ ] Write tests as you go
- [ ] Commit frequently

Evening (30 min):
- [ ] Run full test suite
- [ ] Verify no regressions
- [ ] Document blockers
- [ ] Plan next day
```

### Each Week

```
Monday:
- [ ] Review week's goals
- [ ] Plan daily tasks
- [ ] Sync any blockers

Friday:
- [ ] Review progress
- [ ] Run full test suite
- [ ] Document accomplishments
- [ ] Plan next week
```

---

## Key Commands

### Development
```bash
# Run tests
go test -v ./...

# Run with race detector
go test -race ./...

# Run linter
golangci-lint run ./...

# Format code
gofmt -w ./...

# Build
go build -o kariya ./cmd/kariya

# Run
./kariya
```

### Testing Specific Components
```bash
# Test new intents
go test -v ./internal/cli/intents/...

# Test app
go test -v ./internal/cli/app/...

# Test specific intent
go test -v -run BurstManagement ./internal/cli/intents/...
```

### Git Workflow
```bash
# Create feature branch
git checkout -b feat/aggressive-replacement

# Commit frequently
git commit -m "feat(intents): implement burst management intent"

# Before merging, sync with main
git rebase main

# Final merge
git checkout main
git merge feat/aggressive-replacement
git push
```

---

## Success Indicators

### By End of Week 1
- ✅ Phase 1 complete (audit done)
- ✅ Phase 2 started (first intents created)
- ✅ All new code has tests

### By End of Week 2
- ✅ All 5 new intents created
- ✅ All 10 intents have tests
- ✅ New app.go complete
- ✅ Code compiles and tests pass

### By End of Week 3
- ✅ All workflows functional
- ✅ All tests passing
- ✅ Performance optimized
- ✅ Documentation updated
- ✅ Ready to merge

---

## If You Get Stuck

### Common Issues

**Issue**: Intent doesn't activate from menu
**Solution**: Check router is registered, check activateIntent() is called

**Issue**: Tests failing
**Solution**: Check imports, verify test setup, check assertions

**Issue**: Compilation errors
**Solution**: Run `go mod tidy`, check import paths

**Issue**: Performance slow
**Solution**: Profile with pprof, check for N+1 queries, optimize hot paths

### Resources

- `docs/APP_GO_AGGRESSIVE_REPLACEMENT_PLAN.md` - Full plan
- `docs/APP_GO_MIGRATION_CODE_EXAMPLES.md` - Code examples
- `internal/cli/intents/contract.go` - Intent interface
- `internal/cli/intents/capture_event_intent.go` - Reference implementation

---

## Ready to Start?

### Step 1: Create Feature Branch
```bash
git checkout -b feat/aggressive-replacement
```

### Step 2: Start Phase 1
```bash
# Begin auditing legacy models
ls -la internal/cli/models/*.go | grep -v test
```

### Step 3: Document Your Plan
Create `docs/LEGACY_MODEL_AUDIT.md` with your findings

### Step 4: Proceed Through Phases
Follow the timeline above, committing frequently

---

## Timeline at a Glance

```
Week 1 (40 hours):
├─ Days 1-2: Audit & plan (9 hours)
└─ Days 3-7: Create 5 missing intents (64 hours) - overlapping

Week 2 (48 hours):
├─ Days 1-2: Rebuild app.go (14 hours)
└─ Days 3-5: Testing & validation (34 hours)

Week 3 (30 hours):
├─ Days 1-3: Final polish (30 hours)
└─ Ready to merge and deploy
```

**Total**: 121 hours (~3 weeks full-time)

---

## You've Got This! 💪

Everything is planned, all dependencies are ready, and you have clear milestones.

**Start now with Phase 1 (Preparation)** - it's just 9 hours of auditing and planning.

Then move through the phases systematically, testing as you go.

By week 3, you'll have a clean, modern, intent-driven TUI that's easier to maintain and extend.

**Good luck! 🚀**


