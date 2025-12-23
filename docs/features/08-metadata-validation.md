# Feature: Metadata and Validation

## Purpose
Ensure data integrity, consistency, and quality throughout the career journal system.

## Metadata Management

### Allowed Metadata
- Tags
- Dates
- Companies
- Projects

### Global Metadata Editing
- Adjust metadata without rewriting text
- Changes propagate across related entities
- Maintain event traceability

## Validation Rules

### CareerEvent Validation
- Text not empty
- Text ≤ 2000 characters
- Date ≤ today
- Tags from predefined list

### Allowed Tags
- project
- achievement
- leadership
- technical
- consulting
- research
- product
- mentoring

### Burst Validation
- Requires ≥2 related events
- Meaningful grouping
- Inferred facts must be valid

### Fact Validation
- Competencies from predefined list
- Role fit: [Principal, Staff, EM, Senior IC]
- Audience: [Hiring Manager, Recruiter, Peer]
- Strength: [High, Medium, Low]

### CV View Validation
- Bullets trace to ≥1 Fact
- No aspirational language
- No ungrounded metrics
- Respect role-specific bullet caps

## Competency Management
### Allowed Competencies
- architecture
- delivery
- strategy
- automation
- migration
- system design
- performance
- mentoring
- cross-functional collaboration

## Safety and Conservative Defaults
- Omit if unsure
- No auto-publishing
- Manual review for critical changes

## Acceptance Criteria
- Robust input validation
- Meaningful error messages
- Prevent invalid data entry
- Support global metadata edits
- Maintain data integrity

## Non-Functional Requirements
- Fast validation (≤100ms)
- Scalable to large event collections
- Secure metadata handling
- Audit trail for metadata changes

