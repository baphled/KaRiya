# Development Rules & Guidelines

This directory contains all development standards, guidelines, and workflows for the KaRiya project. These documents define how we work, ensuring consistency, quality, and efficiency.

## 📚 Documentation Overview

### Total Documents: 14 files (~7,400+ lines)
### Archived Documents: 3 files (see docs/archive/rules/)
### Categories: Rules, Guidelines, Quick References, Processes, Best Practices

## 🎯 Quick Access

### Essential Reading (Start Here)
1. **[master-task-prompt.md](master-task-prompt.md)** - Complete 5-phase development workflow
2. **[atomic-commits.md](atomic-commits.md)** - Atomic commit standards
3. **[go-guidelines.md](go-guidelines.md)** - Go coding standards
4. **[AI_COMMIT_ATTRIBUTION.md](AI_COMMIT_ATTRIBUTION.md)** - AI attribution rules

### Quick References (Daily Use)
- **[TASK_QUICK_REF.md](TASK_QUICK_REF.md)** - Task execution workflow
- **[COMMIT_QUICK_REFERENCE.md](COMMIT_QUICK_REFERENCE.md)** - Commit templates
- **[COMPLIANCE_QUICK_REF.md](COMPLIANCE_QUICK_REF.md)** - 5-minute compliance check
- **[AI_COMMIT_CHECKLIST.md](AI_COMMIT_CHECKLIST.md)** - AI commit checklist

## 📖 Document Categories

### 1. Rules & Mandatory Guidelines (3 files)

#### [AI_COMMIT_ATTRIBUTION.md](AI_COMMIT_ATTRIBUTION.md) (652 lines)
- **Type**: Mandatory Rules
- **Purpose**: AI attribution requirements for commits
- **Key Topics**:
  - Mandatory attribution format
  - Assistant name and model version
  - Human review requirements
  - Validation and enforcement

#### [atomic-commits.md](atomic-commits.md) (1000 lines)
- **Type**: Guidelines
- **Purpose**: Creating atomic commits (one logical change per commit)
- **Key Topics**:
  - Atomic commit principles
  - Breaking down work strategies
  - Practical examples
  - Recovery techniques

#### [go-guidelines.md](go-guidelines.md) (194 lines)
- **Type**: Language Guidelines
- **Purpose**: Go-specific coding standards
- **Key Topics**:
  - Framework usage (Fiber, Beego, Cobra)
  - Error handling and concurrency
  - Testing with Ginkgo
  - LSP setup

### 2. Quick References (4 files)

#### [AI_COMMIT_CHECKLIST.md](AI_COMMIT_CHECKLIST.md) (182 lines)
- Quick checklist for AI-generated commits
- Format examples and scenarios
- Troubleshooting tips

#### [COMMIT_QUICK_REFERENCE.md](COMMIT_QUICK_REFERENCE.md) (296 lines)
- One-page commit reference
- Templates and examples
- Common patterns
- Recovery strategies

#### [COMPLIANCE_QUICK_REF.md](COMPLIANCE_QUICK_REF.md) (91 lines)
- 5-minute compliance check
- Code quality, commits, tasks
- Automated check commands

#### [TASK_QUICK_REF.md](TASK_QUICK_REF.md) (173 lines)
- Task execution workflow
- 5-phase process overview
- Commands and templates
- Token efficiency thresholds

### 3. Process & Workflows (4 files)

#### [master-task-prompt.md](master-task-prompt.md) (768 lines)
- **Type**: Master Workflow
- **Purpose**: Integrates ALL project rules into 5-phase workflow
- **Phases**:
  1. Preparation - Planning and setup
  2. TDD Implementation - Red-Green-Refactor
  3. Compliance - Code quality checks
  4. Verification - Testing and validation
  5. Completion - Final checks and commit
- **Key Topics**:
  - Rule integration
  - Token efficiency
  - Examples and troubleshooting

#### [process-task-list.md](process-task-list.md) (161 lines)
- **Type**: Process Guidelines
- **Purpose**: Deterministic task planning and execution
- **Key Topics**:
  - Authority order for task sources
  - Checklist immutability rules
  - Completion criteria

#### [review-commit-prompt.md](review-commit-prompt.md) (650 lines)
- **Type**: Process / Checklist
- **Purpose**: Step-by-step atomic commit review
- **Key Topics**:
  - Manual review checklist
  - Automated script usage
  - Common issues and solutions
  - Integration patterns

#### [rules-compliance-check.md](rules-compliance-check.md) (648 lines)
- **Type**: Comprehensive Checklist
- **Purpose**: Complete rules compliance verification
- **Key Topics**:
  - Code quality standards
  - Commit validation
  - Test requirements
  - Architecture checks
  - Documentation standards

### 4. Best Practices (2 files)

#### [senior-engineer-guidelines.md](senior-engineer-guidelines.md) (103 lines)
- **Type**: Guidelines / Best Practices
- **Purpose**: Language-agnostic engineering standards
- **Key Topics**:
  - SOLID principles
  - Red-Green-Refactor workflow
  - Error handling and observability
  - Security and developer experience

#### [token-efficiency.md](token-efficiency.md) (426 lines)
- **Type**: Guidelines / Best Practices
- **Purpose**: Token conservation strategies for AI interactions
- **Key Topics**:
  - Use tools, be concise
  - Batch operations
  - Reference context
  - Focus on deltas
  - Thresholds and anti-patterns

## 🚀 Getting Started

### For New Developers

**Day 1: Essential Reading**
1. [master-task-prompt.md](master-task-prompt.md) - Learn the workflow
2. [atomic-commits.md](atomic-commits.md) - Understand commit standards
3. [go-guidelines.md](go-guidelines.md) - Go coding standards

**Week 1: Deep Dive**
4. [AI_COMMIT_ATTRIBUTION.md](AI_COMMIT_ATTRIBUTION.md) - If using AI
5. [senior-engineer-guidelines.md](senior-engineer-guidelines.md) - Best practices
6. [token-efficiency.md](token-efficiency.md) - AI interaction efficiency

**Ongoing: Quick References**
- Keep quick references handy for daily work
- Review process documents as needed
- Use compliance checklists regularly

### For AI Assistants

**Essential Rules**
1. [master-task-prompt.md](master-task-prompt.md) - Complete workflow
2. [AI_COMMIT_ATTRIBUTION.md](AI_COMMIT_ATTRIBUTION.md) - Attribution rules
3. [token-efficiency.md](token-efficiency.md) - Token conservation
4. [process-task-list.md](process-task-list.md) - Task processing

**Quick References**
- [TASK_QUICK_REF.md](TASK_QUICK_REF.md)
- [AI_COMMIT_CHECKLIST.md](AI_COMMIT_CHECKLIST.md)
- [COMPLIANCE_QUICK_REF.md](COMPLIANCE_QUICK_REF.md)

## 🎯 Using This Documentation

### Finding What You Need

| I need to... | Check this file |
|--------------|----------------|
| Understand the complete workflow | [master-task-prompt.md](master-task-prompt.md) |
| Make a commit | [COMMIT_QUICK_REFERENCE.md](COMMIT_QUICK_REFERENCE.md) |
| Attribute AI work | [AI_COMMIT_ATTRIBUTION.md](AI_COMMIT_ATTRIBUTION.md) |
| Write Go code | [go-guidelines.md](go-guidelines.md) |
| Check compliance | [COMPLIANCE_QUICK_REF.md](COMPLIANCE_QUICK_REF.md) |
| Save tokens | [token-efficiency.md](token-efficiency.md) |
| Review commits | [review-commit-prompt.md](review-commit-prompt.md) |

### Reading Strategy

1. **Start with Quick Refs**: Get the essentials fast
2. **Read Full Guides**: When implementing new practices
3. **Reference During Work**: Keep quick refs open
4. **Review Periodically**: Refresh knowledge regularly

## 📊 Documentation Statistics

| Category | Files | Lines | Purpose |
|----------|-------|-------|---------|
| Rules & Guidelines | 3 | ~1,850 | Mandatory standards |
| Quick References | 4 | ~742 | Fast lookup |
| Process & Workflows | 4 | ~2,227 | Step-by-step guides |
| Best Practices | 2 | ~529 | Recommended practices |
| **Total** | **14** | **~7,400+** | **Complete coverage** |

## 📦 Archived Documentation

Some files have been archived as workflows evolved. See [`../archive/rules/README.md`](../archive/rules/README.md) for:
- `generate-prd.md` - Replaced by direct task file creation
- `generate-tasks.md` - Replaced by single comprehensive task files
- `task-instructions.md` - Redundant with master-task-prompt.md

## 🔗 Integration & Cross-References

### Document Relationships

- **master-task-prompt.md** references:
  - atomic-commits.md
  - token-efficiency.md
  - rules-compliance-check.md

- **review-commit-prompt.md** references:
  - atomic-commits.md
  - AI_COMMIT_ATTRIBUTION.md

- Quick references complement full guides:
  - TASK_QUICK_REF.md ↔ master-task-prompt.md
  - COMMIT_QUICK_REFERENCE.md ↔ atomic-commits.md
  - AI_COMMIT_CHECKLIST.md ↔ AI_COMMIT_ATTRIBUTION.md

### Related Documentation

- **Setup Guides**: [../setup/](../setup/) - Implementation summaries
- **PRDs**: [../PRD_*.md](../) - Product requirements
- **Features**: [../features/](../features/) - Feature specifications

## 💡 Best Practices for Using Rules

### Daily Workflow

1. **Morning**: Review [TASK_QUICK_REF.md](TASK_QUICK_REF.md)
2. **During Work**: Follow [master-task-prompt.md](master-task-prompt.md)
3. **Before Commit**: Check [COMMIT_QUICK_REFERENCE.md](COMMIT_QUICK_REFERENCE.md)
4. **End of Day**: Run compliance check [COMPLIANCE_QUICK_REF.md](COMPLIANCE_QUICK_REF.md)

### Token Efficiency

When using AI assistants:
- Reference rules by filename instead of pasting content
- Use quick references first, full guides only when needed
- Follow [token-efficiency.md](token-efficiency.md) strategies

### Compliance

Regular compliance checks:
```bash
make check-compliance  # Run automated checks
```

Review:
- [COMPLIANCE_QUICK_REF.md](COMPLIANCE_QUICK_REF.md) - 5-minute check
- [rules-compliance-check.md](rules-compliance-check.md) - Comprehensive check

## 🛠️ Automation Support

### Makefile Targets (Referenced in Rules)
```bash
make test              # Run all tests
make coverage          # Generate coverage report
make review-commit     # Atomic commit review
make check-compliance  # Compliance verification
make staticcheck       # Advanced static analysis
```

### Git Hooks (Referenced in Rules)
- **prepare-commit-msg** - AI attribution helper
- **commit-msg** - Commitlint validation

### Scripts (Referenced in Rules)
- `scripts/review-commit.sh` - Commit review automation
- `scripts/check-compliance.sh` - Compliance automation

## ✅ Compliance Checklist

After reading the rules, ensure you understand:

- [ ] 5-phase development workflow
- [ ] Atomic commit principles
- [ ] AI attribution requirements (if using AI)
- [ ] Go coding standards
- [ ] Token efficiency strategies
- [ ] Compliance check process

## 📝 Maintaining Rules

### When to Update

- Process improvements discovered
- New tools or practices adopted
- Common issues need documentation
- Team feedback suggests changes

### Contribution Guidelines

1. Discuss changes with team
2. Update relevant rules document
3. Update related quick references
4. Update cross-references
5. Update this README if structure changes

## ℹ️ Getting Help

### Questions About Rules

- Check quick references first
- Read full guide for deep understanding
- Review examples in documents
- Ask team for clarification

### Suggesting Improvements

- Document the issue or suggestion
- Propose the change with rationale
- Discuss with team
- Update documentation after agreement

---

**Last Updated**: 2026-01-06
**Directory**: `docs/rules/`
**Total Files**: 14
**Total Lines**: ~7,400+
**Status**: Production-ready and comprehensive

