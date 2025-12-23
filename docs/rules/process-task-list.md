## Task Planning + Execution Guide

This guide defines a deterministic, tool-driven process for planning and executing development tasks via Avante.nvim. It is designed to minimize token waste, prevent execution stalls, and ensure reliable task completion while preserving senior-engineer quality standards.

---

## 1. Scope and Assumptions

* Each job is executed as a **fresh prompt**.
* Avante.nvim provides access to tools (filesystem, tests, git, etc.).
* Tools are authoritative when available.
* This guide governs **task planning and execution only**. PRD generation and task generation are separate workflows.

---

## 2. High-Level Workflow

1. **One-time setup (outside execution loop)**

   * Load language-specific guidelines (e.g. `go-guidelines`).
   * Load `senior-engineer-context-guidelines`.
   * Load this document (`process-task-list`).
   * Ensure Avante.nvim tools (filesystem, tests, git MCP) are available.

2. **Planning phase (single job)**

   * Generate a locked task checklist from a PRD or feature summary.
   * Write the checklist to a file (e.g. `task-plan.md`).
   * Confirm correctness once. After confirmation, the checklist becomes immutable.

3. **Execution phase (repeated jobs)**

   * Execute exactly **one task at a time** from the locked checklist.
   * Use tools to inspect state and verify results.
   * Mark the task as complete and stop.

---

## 3. Authority Order

When instructions conflict, the following precedence applies (highest → lowest):

1. Avante.nvim **tool results** (tests, filesystem, git)
2. The **locked task checklist**
3. Language-specific guidelines
4. Senior-engineer best practices

If a higher authority is satisfied, lower authorities must not block completion.

---

## 4. Planning Rules (Checklist Generation Only)

* Output a checklist under a `## Tasks` heading only.
* Each parent task must be independently completable and committable.
* Tasks must be strictly sequential.
* Tasks may include **inline guidance**, but inline guidance is advisory.
* Inline guidance must not override execution completion criteria.
* Once confirmed, the checklist is **read-only**.

### Checklist Immutability Rule

The checklist file may only be modified to change checkbox state (`[ ]` → `[x]`).
Any other modification is forbidden.

---

## 5. Execution Rules (Task Processing Only)

During execution, the assistant must:

* Reference the locked checklist file via tools.
* Select exactly one incomplete task.
* Perform only the work required for that task.
* Use tools to create, modify, and inspect files as needed.
* Use tools to run tests or checks when required.

### Refactoring Constraint

* Refactoring is **permitted only** on code directly modified by the current task.
* Architectural changes or cross-cutting refactors are forbidden unless the task explicitly requires them.

---

## 6. Task Completion Criteria

A task is complete when **all** of the following are true:

* Required files are created or modified.
* Associated tests exist if required by the task.
* Avante.nvim tool reports tests passing **or** tests were not required.
* No Avante.nvim tool reports errors.

### Completion Action

When the criteria above are met:

1. Mark the task checkbox as `[x]`.
2. Stop execution immediately.
3. Do not re-run tools.
4. Do not refactor further.
5. Do not modify any other tasks.

This is a hard stop condition.

---

## 7. Tool Usage Rules

* Tool output is authoritative.
* If tools report success, execution must proceed to completion.
* If tools report failure, fix only the minimal code required to resolve the failure.
* Do not speculate about tool results that are not observed.

---

## 8. Memory and State Rules

* Memory is read-only during execution.
* Do not create, delete, or mutate memory entries as part of task completion.
* All authoritative state must be visible via tools or files.

---

## 9. Cost and Rate-Limit Controls

* Use Claude Sonnet for planning only.
* Use Claude 3.5 Haiku for task execution.
* Execute one task per job.
* Provide only the task index and file references in execution prompts.
* Avoid resending full specifications or checklists; rely on tool access instead.
* Avante.nvim should throttle requests (≥ 2s between calls).

---

## 10. Execution Prompt Template

```
WORKFLOW: PROCESS_TASK_LIST
MODE: EXECUTION

Using the locked checklist file, execute Task <N> only.
Follow the execution rules and completion criteria.
Do not replan, summarize, or modify other tasks.
Stop immediately after marking the task complete.
```

---

## 11. Expected Outcomes

Following this process guarantees:

* Deterministic task completion
* Reliable checkbox updates
* Bounded scope per task
* Minimal token usage
* Clear stopping conditions

This document is the sole authority for task execution behavior.
