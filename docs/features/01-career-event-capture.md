# Feature: Career Event Capture

## Purpose
Enable senior engineers and consultants to capture career events in a low-friction, ad-hoc manner.

## Core Concept
A Career Event is the smallest input unit, representing:
- Free-form text describing a project, outcome, or responsibility
- Optional metadata: date, company, project

## Input Modes
1. Timeline journaling (default)
2. CV backfill (import existing CVs)

## Validation Rules
- Text cannot be empty
- Text length ≤ 2000 characters
- Date must be ≤ today
- Tags must be from allowed values

## Allowed Tags
- project
- achievement
- leadership
- technical
- consulting
- research
- product
- mentoring

## Acceptance Criteria
- Users can input career events with optional metadata
- System validates input against defined rules
- Provides clear feedback for invalid input
- Supports mobile and desktop layouts

## Non-Functional Requirements
- Input process should be quick and intuitive
- Minimal friction for event capture
- Support for both structured and unstructured input

