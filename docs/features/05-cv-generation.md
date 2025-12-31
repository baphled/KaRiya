# Feature: CV Generation

## Purpose
Transform raw career events into credible, audience- and role-specific CV views.

## Core Principles
- Read-only generation
- Explainable and reversible
- No text rewriting
- Conservative defaults

## Generation Rules
### Inclusion Criteria
- Bullet must trace to ≥1 journal entry
- Single-claim bullets only
- Prefer repeated signals

### Exclusion Criteria
- No inferred metrics
- No aspirational language
- No role inflation

### Ranking Priority
1. Ownership → Contribution
2. Strategy → Execution
3. Outcome → Activity

### Compression
- Hard bullet caps per role:
  - Principal: 3–4 bullets
  - Senior IC: 4–5 bullets
- Older roles compress first

## Role-Specific Generation

- Configurable roles with distinct bullet limits and priorities

### Roles
- Principal: Strategic ownership, cross-team leadership
- Staff: Technical leadership, high-complexity implementation
- EM: Team leadership, mentorship, delivery accountability
- Senior IC: Deep technical contribution, system design

## Audience-Specific Filtering
- Hiring Manager: Outcomes, ownership, business impact
- Recruiter: Skills, competencies, high-level achievements
- Peer: Technical depth, collaboration, problem-solving

## Acceptance Criteria
- CV bullets must be traceable to source events
- Not stored in DB, generated on-the-fly
- Respect bullet caps for each role
- Maintain factual accuracy
- Provide source event visibility
- Support multiple role and audience views

## Non-Functional Requirements
- CV generation ≤ 2s for ≤500 events
- Maintain data integrity
- Support global metadata edits
- Preserve event traceability

