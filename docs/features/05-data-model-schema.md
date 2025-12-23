# Feature: Data Model and Schema

## Purpose
Define a canonical, validated data model for career events, bursts, facts, and CV views.

## Data Entities

### CareerEvent
```yaml
CareerEvent:
  required:
    - id: string (UUID v4)
    - text: string
    - created_at: datetime
  optional:
    - date: date
    - company: string
    - project: string
    - tags: array[string]
    - bursts: array[string]
```

### Burst
```yaml
Burst:
  required:
    - id: string (UUID v4)
    - event_ids: array[string] (≥2)
  optional:
    - name: string
    - inferred_facts: array[string]
```

### Fact
```yaml
Fact:
  required:
    - id: string (UUID v4)
    - competencies: array[string]
    - role_fit: array[string]
    - audience: array[string]
    - strength: string
  optional:
    - source_event_ids: array[string]
    - burst_ids: array[string]
```

### CVView
```yaml
CVView:
  required:
    - id: string (UUID v4)
    - target_role: string
    - audience: string
    - selected_fact_ids: array[string]
    - bullets: array[string]
```

## Global Validation Rules
1. All IDs must be UUID v4
2. Valid references between entities
3. Text ≤ 2000 characters
4. No aspirational language
5. No ungrounded numeric metrics

## Relationships
- CareerEvent can belong to multiple Bursts
- Burst can infer multiple Facts
- Facts can be selected for multiple CVViews

## Acceptance Criteria
- Entities have clear, enforceable schemas
- References between entities are valid
- No dangling or orphaned IDs
- Consistent metadata across entities

## Non-Functional Requirements
- Support scalability (≥10,000 events/user)
- Data encryption at rest and in transit
- Auditability of changes
- Support for metadata edits and propagation

