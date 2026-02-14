# Canonical Behaviours: Memory vs SQL Repository Gaps

**Date**: 2026-02-14
**Purpose**: Document the 6 identified behavioural differences between memory and SQL repositories, establish canonical (correct) behaviour, and define required fixes.

---

## Gap 1: Validation on Create

### Current Behaviour

**Memory Repository**:
- File: `internal/repository/career/memory/event_repository.go:61-63`
- Calls `event.Validate()` before storing
- Rejects invalid events with validation error

**SQL Repository**:
- File: `internal/repository/career/sql/event_repository.go:49-69`
- Does NOT call `event.Validate()`
- Stores any event, relies on DB constraints

### Canonical Behaviour

**Decision**: Memory behaviour is CORRECT (stricter validation is safer)
- Validation should happen at repository layer
- Prevents invalid data from entering system
- SQL should be updated to match memory

### Change Required

- [ ] Update SQL `EventRepository.Create()` to call `event.Validate()` before storing
- [ ] Ensure validation errors are returned consistently

---

## Gap 2: Pagination Limit=0

### Current Behaviour

**Memory Repository**:
- File: `internal/repository/career/memory/helpers.go:48-53`
- When `limit <= 0`, returns ALL records (no limit)
- Unbounded result set

**SQL Repository**:
- File: `internal/repository/career/sql/event_repository.go:231-237`
- When `limit == 0`, caps at 100 records
- Prevents unbounded queries

### Canonical Behaviour

**Decision**: SQL behaviour is CORRECT (bounded queries are safer)
- Limit=0 should mean "use default limit" (100)
- Prevents accidental full-table scans
- Memory should be updated to match SQL

### Change Required

- [ ] Update memory `ListEvents()` to cap limit at 100 when limit <= 0
- [ ] Apply same pattern to all memory repository List methods

---

## Gap 3: Tag Filtering

### Current Behaviour

**Memory Repository**:
- File: `internal/repository/career/memory/event_repository.go:159-171`
- Exact string match: `event.Tags contains tag`
- Case-sensitive, whole-word match only

**SQL Repository**:
- File: `internal/repository/career/sql/event_repository.go:193-215`
- Uses `LIKE %tag%` pattern
- Substring matching, case-insensitive (SQLite default)

### Canonical Behaviour

**Decision**: Memory behaviour is CORRECT (exact match is safer)
- Substring matching with LIKE is a bug risk (false positives)
- Example: searching for "Go" would match "Going", "Golang", etc.
- SQL LIKE pattern is overly permissive

### Change Required

- [ ] Update SQL `ListEventsByTag()` to use exact match instead of LIKE
- [ ] Ensure both implementations produce identical results

---

## Gap 4: Skill Cross-Repo Sync

### Current Behaviour

**Memory Repository**:
- File: `internal/repository/career/memory/event_repository.go:226-256`
- `LinkSkill()` manually syncs skill back-reference
- Calls `skillRepo.Update()` to add event ID to skill's event list

**SQL Repository**:
- File: `internal/repository/career/sql/event_repository.go:279-305`
- Uses junction table (`event_skills`)
- No manual sync needed; query joins tables

### Canonical Behaviour

**Decision**: Both must produce identical results for the roundtrip:
1. `eventRepo.LinkSkill(eventID, skillID)`
2. `skillRepo.GetSkill(skillID)` → should include eventID in linked events

**Critical**: The junction table approach (SQL) is the source of truth. Memory must replicate this behaviour exactly.

### Change Required

- [ ] Verify memory `LinkSkill()` correctly updates skill's event list
- [ ] Verify memory `GetSkill()` returns all linked events
- [ ] Contract test must verify roundtrip: LinkSkill → GetSkill → verify linked

---

## Gap 5: Delete Cascade

### Current Behaviour

**Memory Repository**:
- File: `internal/repository/career/memory/event_repository.go` (DeleteEvent)
- No cascade: deleting event does NOT remove linked skills
- Orphaned skill records remain

**SQL Repository**:
- File: `internal/repository/career/sql/event_repository.go` (DeleteEvent)
- DB constraints enforce cascade: deleting event removes event_skills junction records
- No orphans

### Canonical Behaviour

**Decision**: SQL behaviour is CORRECT (cascade prevents orphans)
- Deleting an event should remove all its skill links
- Orphaned records are a data integrity risk
- Memory should implement cascade logic

### Change Required

- [ ] Update memory `DeleteEvent()` to also delete all linked skills
- [ ] Implement cascade logic: `for each skill linked to event: skillRepo.UnlinkEvent(skillID, eventID)`

---

## Gap 6: List Ordering

### Current Behaviour

**Memory Repository**:
- File: `internal/repository/career/memory/event_repository.go` (ListEvents)
- Returns events in insertion order (map iteration order is unstable in Go)
- Order is non-deterministic across runs

**SQL Repository**:
- File: `internal/repository/career/sql/event_repository.go:231-237`
- Uses `ORDER BY created_at DESC` (deterministic)
- Consistent ordering across runs

### Canonical Behaviour

**Decision**: SQL behaviour is CORRECT (deterministic ordering is essential)
- Tests depend on predictable ordering
- Non-deterministic results cause flaky tests
- Memory must sort results

### Change Required

- [ ] Update memory `ListEvents()` to sort by `created_at DESC` (matching SQL)
- [ ] Apply same pattern to all memory repository List methods

---

## Summary Table

| Gap | Memory | SQL | Canonical | Fix |
|-----|--------|-----|-----------|-----|
| 1. Validation | Validates | No validation | Memory ✓ | Update SQL |
| 2. Pagination | Unbounded | Capped at 100 | SQL ✓ | Update Memory |
| 3. Tag Filter | Exact match | LIKE substring | Memory ✓ | Update SQL |
| 4. Skill Sync | Manual sync | Junction table | Both must match | Verify Memory |
| 5. Delete Cascade | No cascade | Cascade | SQL ✓ | Update Memory |
| 6. List Ordering | Unstable | Deterministic | SQL ✓ | Update Memory |

---

## Contract Test Design Implications

These canonical behaviours define what the contract tests must verify:

1. **Validation Contract**: Both repos must reject invalid events identically
2. **Pagination Contract**: Both repos must cap limit at 100 when limit <= 0
3. **Tag Filter Contract**: Both repos must return exact matches only
4. **Skill Sync Contract**: LinkSkill + GetSkill roundtrip must work identically
5. **Cascade Contract**: DeleteEvent must remove all linked skills in both repos
6. **Ordering Contract**: ListEvents must return results sorted by created_at DESC

---

## Next Steps

1. Task 5: Create contract tests based on these canonical behaviours
2. Task 5: Fix memory repositories to match canonical behaviours
3. Task 5: Verify all contract tests pass for both implementations
