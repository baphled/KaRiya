# Feature: Burst and Fact Extraction

## Purpose
Automatically group and enrich career events to provide deeper insights and context.

## Core Concepts

### Burst
- Automatically suggested grouping of related career events
- Captures larger initiatives, projects, or phases
- Supports future portfolio/case study creation

### Fact
- Inferred from events
- Contains:
  - Competencies
  - Role fit (Principal, EM, Staff, Senior IC)
  - Audience relevance (Hiring Manager, Recruiter, Peer)
  - Strength signal

## Extraction Rules
- Burst requires ≥2 related CareerEvents
- Facts must link to source events or bursts
- Inference based on text analysis and metadata

## Allowed Competencies
- architecture
- delivery
- strategy
- automation
- migration
- system design
- performance
- mentoring
- cross-functional collaboration

## Validation Rules
- Burst must reference ≥2 CareerEvents
- Inferred facts must have valid references
- No aspirational language
- No ungrounded metrics

## Acceptance Criteria
- System suggests meaningful event groupings
- Facts are accurately inferred from events
- Users can confirm or edit suggested metadata
- Maintains traceability of original events

## Non-Functional Requirements
- Fast inference (≤2s for ≤500 events)
- Support for large event collections (≥10,000 events/user)
- Transparent and explainable inference process

