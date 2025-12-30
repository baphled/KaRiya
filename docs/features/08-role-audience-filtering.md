# Feature: Role and Audience Filtering

## Purpose
Dynamically generate CV views tailored to specific roles and audiences.

## Role Filtering

### Roles and Inclusion Criteria

#### Principal
- Strategic ownership
- Cross-team leadership
- End-to-end impact
- **Bullet Cap:** 3–4

#### Staff
- Technical leadership
- High-complexity implementation
- **Bullet Cap:** 4–5

#### Engineering Manager (EM)
- Team leadership
- Mentorship
- Delivery accountability
- **Bullet Cap:** 4–5

#### Senior Individual Contributor (Senior IC)
- Deep technical contribution
- System design
- Execution
- **Bullet Cap:** 4–5

## Ranking Rules
1. Ownership > Contribution
2. Strategy > Execution
3. Outcome > Activity
4. Older events compress first if bullet cap exceeded

## Audience Filtering

### Audiences and Emphasis

#### Hiring Manager
- Outcomes
- Ownership
- Business impact
- Leadership

#### Recruiter
- Skills
- Competencies
- High-level achievements

#### Peer
- Technical depth
- Collaboration
- Problem-solving details

## Implementation Rules
- All bullets remain factually traceable
- No aspirational language
- Strength signals and competencies rank relevance

## Acceptance Criteria
- CV bullets respect role-specific caps
- Audience-specific emphasis works correctly
- Maintain factual integrity
- Transparent filtering process

## Non-Functional Requirements
- Fast filtering (≤2s)
- Support multiple simultaneous views
- Consistent ranking across roles
- Scalable to large event collections

