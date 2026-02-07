# Memory Keeper Usage Guide

This document explains how to use the `memory-keeper` skill with the MCP memory server.

## Prerequisites

1. MCP memory server must be installed and configured
2. Memory server should be available in your MCP server configuration

## Typical Workflows

### Workflow 1: After Fixing a Bug

```
Scenario: You just fixed a race condition in the event list screen.

1. Invoke memory-keeper skill
2. Structure the memory using the Bug Fix template
3. Use MCP memory tool to store:
   - Entity: "kariya-race-condition-event-list"
   - Content: Structured bug fix memory
   - Tags: ["race-condition", "bubble-tea", "event-list", "concurrency"]
   - Scope: "project"
4. Confirm storage successful
```

**Expected Agent Behaviour:**
```
I've fixed the race condition in the event list screen. Let me record this 
in memory so we can reference it if we encounter similar issues.

[Uses MCP memory storage tool]

Memory stored with ID: mem_kariya_race_20260206_001
Tags: race-condition, bubble-tea, event-list, concurrency
Scope: project

This memory can now be recalled when we encounter similar concurrency issues.
```

### Workflow 2: Before Starting Debug Work

```
Scenario: Tests are failing with "Found more than one test suite file"

1. Invoke memory-keeper skill  
2. Query MCP memory for similar issues:
   - Search terms: ["ginkgo", "multiple suite", "test suite file"]
   - Scope: "tool"
3. Review retrieved memories
4. Apply solution from memory if found
5. If no memory exists, proceed with debugging and record afterwards
```

**Expected Agent Behaviour:**
```
I see a Ginkgo error about multiple test suite files. Let me check if 
we've encountered this before.

[Uses MCP memory query tool]

Found 2 related memories:
1. "Ginkgo Suite File Conflict in Repositories" (2025-11-03)
   - Solution: Each package needs exactly ONE *_suite_test.go
   
2. "Duplicate RunSpecs Call" (2025-10-15)
   - Solution: Check for multiple TestXxx functions

Let me apply the first solution...
```

### Workflow 3: During Retrospective

```
Scenario: Post-sprint retrospective identifies recurring test flakiness

1. Conduct retrospective (using retrospective skill)
2. Identify key learnings
3. Invoke memory-keeper to store learnings
4. Create memory entries for each pattern:
   - Pattern/Gotcha template
   - Tag with "test-flakiness", "ginkgo", relevant components
   - Scope: "project" or "tool" depending on specificity
5. Link memories to retrospective document
```

**Expected Agent Behaviour:**
```
The retrospective identified 3 recurring patterns causing test flakiness.
Let me store these in memory:

1. Storing "Eventually() with Default Timeout Too Short"
   [Uses MCP memory storage tool]
   Memory ID: mem_kariya_eventually_timeout_001

2. Storing "Shared Test State Not Reset in BeforeEach"  
   [Uses MCP memory storage tool]
   Memory ID: mem_kariya_shared_state_002

3. Storing "Mock Expectations Not Cleared Between Tests"
   [Uses MCP memory storage tool]  
   Memory ID: mem_kariya_mock_cleanup_003

All memories tagged with: test-flakiness, ginkgo, best-practices
Linked to retrospective: docs/retrospectives/2026-02-06-sprint-24.md
```

## MCP Memory Tool Interface

The memory-keeper skill expects the MCP memory server to provide these capabilities:

### Storage Tool
```
Tool: memory_store
Parameters:
  - entity_id: string (unique identifier)
  - content: string (the memory content)
  - tags: array<string> (searchable tags)
  - metadata: object (scope, type, timestamp, etc.)
```

### Retrieval Tool
```
Tool: memory_query
Parameters:
  - query: string (search terms)
  - tags: array<string> (filter by tags)
  - scope: string (filter by scope)
  - limit: number (max results)
```

### Update Tool
```
Tool: memory_update
Parameters:
  - entity_id: string
  - content: string (updated content)
  - tags: array<string> (updated tags)
```

### Delete Tool
```
Tool: memory_delete
Parameters:
  - entity_id: string
```

## Agent Instructions

When using the memory-keeper skill, the agent should:

1. **Auto-detect storage opportunities:**
   - After successful bug fixes
   - After resolving test failures
   - After debugging sessions
   - During retrospectives

2. **Proactively query before debugging:**
   - When encountering errors
   - When starting work on known problematic areas
   - At session start for context

3. **Structure memories consistently:**
   - Always use appropriate template
   - Include all required fields
   - Tag generously for discoverability
   - Set correct scope

4. **Inform the user:**
   - Tell user when storing a memory
   - Show what was stored and how to recall it
   - Mention when using recalled memory to solve issue

## Example Memory Entries

### Example: Race Condition Bug

```yaml
entity_id: mem_kariya_race_event_list_001
type: bug_fix
scope: project
component: career/ui/event_list
created: 2026-02-06T10:30:00Z

content: |
  **Problem:**
  Event list screen crashes with concurrent map read/write panic when 
  events update during rendering.
  
  **Symptoms:**
  - Panic: concurrent map read and write
  - Only under load, not in normal operation
  - go test -race catches it
  
  **Root Cause:**
  EventListModel reads from shared events map without mutex protection.
  Multiple goroutines update model state concurrently.
  
  **Solution:**
  Added sync.RWMutex to EventListModel:
  - File: internal/cli/intents/career/ui/event_list_model.go:23
  - Reads: RLock()/RUnlock()
  - Writes: Lock()/Unlock()
  
  **Prevention:**
  - Always use mutex for shared state in Bubble Tea models
  - Run `make test-race` before committing
  - Added test: event_list_concurrent_test.go
  
tags:
  - race-condition
  - bubble-tea
  - concurrent-access
  - event-list
  - mutex
  - threading

metadata:
  pr: "#234"
  files_changed:
    - internal/cli/intents/career/ui/event_list_model.go
    - internal/cli/intents/career/ui/event_list_concurrent_test.go
  time_to_fix: "2 hours"
```

### Example: GORM Pattern

```yaml
entity_id: mem_gorm_preload_not_loaded_001
type: pattern
scope: tool
tool: gorm
created: 2025-12-10T14:20:00Z

content: |
  **Pattern Name:**
  GORM Preload Not Loading Related Data
  
  **Context:**
  When querying entities with relationships using GORM repository pattern.
  
  **Problem:**
  Related entities are nil even though they exist in database.
  
  **Incorrect Approach:**
  ```go
  // This won't load events
  career, err := repo.FindByID(ctx, id)
  // career.Events is nil/empty
  ```
  
  **Correct Approach:**
  ```go
  // Use Preload to load relationships
  career, err := repo.FindByIDWithEvents(ctx, id)
  
  // In repository:
  func (r *CareerRepository) FindByIDWithEvents(ctx context.Context, id uuid.UUID) (*Career, error) {
      var career Career
      err := r.db.WithContext(ctx).
          Preload("Events").
          Where("id = ?", id).
          First(&career).Error
      return &career, err
  }
  ```
  
  **Rationale:**
  GORM uses lazy loading by default. Related entities must be explicitly 
  preloaded or queries will return them as nil/empty.
  
tags:
  - gorm
  - preload
  - relationships
  - repository-pattern
  - lazy-loading

metadata:
  related_memories:
    - mem_gorm_n_plus_one_001
  documentation: "https://gorm.io/docs/preload.html"
```

## Memory Search Examples

### Search 1: "How do we handle GORM relationships?"

```
Query: gorm relationships preload
Tags: gorm, repository-pattern
Scope: tool

Results:
1. mem_gorm_preload_not_loaded_001 (relevance: 0.95)
2. mem_gorm_nested_preload_002 (relevance: 0.87)
3. mem_gorm_n_plus_one_001 (relevance: 0.76)
```

### Search 2: "Have we seen this race condition before?"

```
Query: race condition concurrent map
Tags: race-condition, concurrency
Scope: project

Results:
1. mem_kariya_race_event_list_001 (relevance: 0.92)
2. mem_kariya_race_screen_state_003 (relevance: 0.84)
```

### Search 3: "Ginkgo test suite errors"

```
Query: ginkgo test suite file error
Tags: ginkgo, test-suite
Scope: tool

Results:
1. mem_ginkgo_multiple_suite_files_001 (relevance: 0.98)
2. mem_ginkgo_suite_naming_002 (relevance: 0.71)
```

## Integration with Development Workflow

### At Session Start
```bash
make session-start
# Agent automatically:
# 1. Queries memories for current branch/feature area
# 2. Surfaces recent memories from last session
# 3. Shows any unresolved patterns/issues
```

### During Development
```
# When encountering an error:
Agent: "I see a GORM preload error. Let me check memory..."
[Queries memory, finds solution]
Agent: "We've encountered this before in mem_gorm_preload_not_loaded_001"
[Applies stored solution]
```

### After Fixing Issue
```
Agent: "Bug fixed. Let me record this for future reference..."
[Structures memory using template]
[Stores via MCP memory tool]
Agent: "Memory stored as mem_kariya_validation_error_005"
```

### During Code Review
```
Reviewer: "This looks like it might have a race condition"
Agent: "Let me check our memory for similar patterns..."
[Queries race condition memories]
Agent: "Yes, similar to mem_kariya_race_event_list_001. 
       We should add mutex protection here."
```

## Future Enhancements

Potential additions to memory-keeper:

1. **Auto-categorisation:** ML-based tagging and scope detection
2. **Memory linking:** Automatically link related memories
3. **Pattern detection:** Identify when multiple memories indicate systemic issue
4. **Memory analytics:** Show which memories are most useful
5. **Proactive suggestions:** "Based on your current work, you might encounter..."
6. **Memory expiry:** Auto-archive or flag outdated memories
7. **Cross-project search:** Find solutions from other projects
8. **Memory export:** Generate documentation from frequently-accessed memories
