# KaRiya Project Documentation

**Last Updated**: 2026-01-13
**Project Status**: ✅ **PRODUCTION READY - ALL PHASES COMPLETE (100%)**
**Test Coverage**: 240+ tests, 100% pass rate, 0 race conditions
**Code Quality**: All linting checks passing, no technical debt
**Forms**: Huh library integration (Phase 4/5 complete)
**Form Wrappers**: 2 (CaptureForm, SkillForm) - Required for intent forms

---

## ⚠️ CRITICAL: Mandatory Requirements

**Before reading this documentation, read these FIRST:**

1. **[STRICT REQUIREMENTS SUMMARY](docs/rules/STRICT_REQUIREMENTS_SUMMARY.md)** ⭐
   - **ONE-PAGE** summary of ALL mandatory requirements
   - AI assistant identity (senior Go engineer)
   - AI commit attribution (`make ai-commit` ONLY)
   - TDD protocol (Red-Green-Refactor)
   - Code quality standards (SOLID + Go idioms)
   - Refusal criteria

2. **[Senior Engineer Guidelines](docs/rules/senior-engineer-guidelines.md)**
   - Complete engineering standards
   - SOLID principles (enforced)
   - Clean code principles (DRY, KISS, YAGNI)
   - When to refuse to write code

3. **[AI Commit Attribution Rules](docs/rules/AI_COMMIT_ATTRIBUTION.md)**
   - **MANDATORY**: Use `make ai-commit` for ALL AI code
   - **PROHIBITED**: Using `git commit` directly for AI code
   - Zero tolerance policy

**These are NON-NEGOTIABLE. Read them before proceeding.**

---

## Table of Contents

1. [Session Contract](#session-contract)
2. [AI Mandatory Protocol](#ai-mandatory-protocol)
3. [Task Template (Required Format)](#task-template-required-format)
4. [Project Overview](#project-overview)
5. [Quick Start](#quick-start)
6. [Architecture Overview](#architecture-overview)
7. [Development Guidelines & Rules](#development-guidelines--rules)
8. [TUI Development](#tui-development)
9. [User Guides & Features](#user-guides--features)
10. [Implementation Resources](#implementation-resources)
11. [Task Documentation](#task-documentation)
12. [Key Files and Purposes](#key-files-and-purposes)
13. [Recent Fixes](#recent-fixes)
14. [Common Development Tasks](#common-development-tasks)
15. [Deployment Guide](#deployment-guide)
16. [Performance Benchmarks](#performance-benchmarks)
17. [Workflow Patterns](#workflow-patterns)
18. [Troubleshooting](#troubleshooting)
19. [Project Metadata](#project-metadata)

---

## Session Contract

**This contract is displayed when running `make session-start`.**

By proceeding with this work session, you acknowledge and commit to:

1. **TDD Protocol**: Tests are written **BEFORE** implementation code (Red-Green-Refactor)
2. **Compliance First**: `make check-compliance` runs **before AND after** every task
3. **Atomic Commits**: One logical change per commit, with AI attribution if AI-generated
4. **Sequential Tasks**: One task at a time, in checklist order
5. **Token Efficiency**: Tools over text, concise communication, batch operations

**Violation of these rules requires stopping work and correcting before proceeding.**

---

## AI Mandatory Protocol

**⚠️ CRITICAL: These are NON-NEGOTIABLE requirements. Failure to follow these EXACTLY will result in immediate work stoppage.**

### Session Start Requirements (MANDATORY)

The AI assistant **MUST**:

1. **IMMEDIATELY** run `make session-start`
2. **WAIT** for explicit confirmation that it passed
3. If it fails, **REFUSE to proceed** until ALL violations are fixed
4. Display: "Session contract acknowledged. Ready to proceed."

**NO exceptions.** If the user tries to skip this, the AI assistant **MUST REFUSE ALL WORK**.

### Before ANY Code Changes (MANDATORY)

The AI assistant **MUST**:

1. State the specific task being worked on (from task file)
2. Confirm it is **ONE** atomic change (reject if multiple changes)
3. State which test file will be created/modified **FIRST**
4. **WAIT** for explicit user confirmation before proceeding
5. **Acknowledge** that you are a **senior Go engineer** following SOLID principles

**NO code generation** until ALL confirmations are received.

---

# Task Handover Completed

*Merged with existing AGENTS.md content successfully as part of redundancy-free comprehensive documentation.*
