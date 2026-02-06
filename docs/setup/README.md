# Setup Guides

This directory contains comprehensive setup and installation guides for the KaRiya project. These guides walk you through setting up your development environment and understanding the project's automation systems.

## 📋 Available Guides

### Core Setup Guides

1. **[AI_COMMIT_ATTRIBUTION_SETUP.md](AI_COMMIT_ATTRIBUTION_SETUP.md)**
   - Comprehensive AI commit attribution system setup
   - Created documentation, tools, hooks, and validation
   - Usage instructions and workflow integration

2. **[AI_COMMIT_SETUP.md](AI_COMMIT_SETUP.md)**
   - Quick 5-minute setup guide for AI commit attribution
   - Step-by-step installation and verification
   - Common workflows and troubleshooting

3. **[ATOMIC_COMMITS_SETUP.md](ATOMIC_COMMITS_SETUP.md)**
   - Atomic commit implementation and workflow
   - Review scripts and automation tools
   - Daily workflow examples

4. **[CI_CD_SETUP.md](CI_CD_SETUP.md)**
   - Complete CI/CD pipeline setup guide
   - GitHub Actions, commitlint, and semantic-release
   - Release automation and commit format enforcement

5. **[MASTER_TASK_PROMPT_SETUP.md](MASTER_TASK_PROMPT_SETUP.md)**
   - Master task execution system setup
   - 5-phase workflow integration
   - Rule integration and compliance checks

6. **[RULES_COMPLIANCE_SETUP.md](RULES_COMPLIANCE_SETUP.md)**
   - Rules compliance and token efficiency system
   - Compliance check scripts and automation
   - Makefile targets and workflow integration

7. **[BRANCH_CLEANUP_SETUP.md](BRANCH_CLEANUP_SETUP.md)**
   - Automatic branch deletion configuration
   - GitHub auto-delete of merged branches
   - Local cleanup script and Git aliases

## 🚀 Getting Started

### For New Developers

**Recommended Setup Order:**

1. Start with **CI_CD_SETUP.md** to understand commit standards
2. Follow **ATOMIC_COMMITS_SETUP.md** for commit best practices
3. Review **AI_COMMIT_SETUP.md** for AI attribution (if using AI assistance)
4. Read **MASTER_TASK_PROMPT_SETUP.md** for the complete development workflow
5. Check **RULES_COMPLIANCE_SETUP.md** for quality standards

### For Quick Setup

If you just want to start coding:

1. **Essential**: [CI_CD_SETUP.md](CI_CD_SETUP.md) - Learn commit format
2. **Recommended**: [AI_COMMIT_SETUP.md](AI_COMMIT_SETUP.md) - If using AI tools
3. **Optional**: Other guides for advanced workflows

## 📖 Guide Categories

### Commit & Version Control
- **AI_COMMIT_ATTRIBUTION_SETUP.md** - AI attribution system
- **AI_COMMIT_SETUP.md** - Quick AI setup
- **ATOMIC_COMMITS_SETUP.md** - Atomic commit practices
- **BRANCH_CLEANUP_SETUP.md** - Branch cleanup automation

### CI/CD & Automation
- **CI_CD_SETUP.md** - Pipeline and automation
- **RULES_COMPLIANCE_SETUP.md** - Quality checks

### Workflows
- **MASTER_TASK_PROMPT_SETUP.md** - Complete task execution

## 🎯 What Each Guide Covers

### AI_COMMIT_ATTRIBUTION_SETUP.md (411 lines)
- **Purpose**: Comprehensive AI attribution setup
- **Contents**:
  - System overview and requirements
  - Documentation created
  - Tools and hooks installed
  - Workflow integration
  - Troubleshooting

### AI_COMMIT_SETUP.md (359 lines)
- **Purpose**: Quick start for AI attribution
- **Contents**:
  - 5-minute setup steps
  - Commit workflows (human/AI/mixed)
  - Verification commands
  - Common issues and fixes

### ATOMIC_COMMITS_SETUP.md (508 lines)
- **Purpose**: Atomic commit implementation
- **Contents**:
  - Atomic commit principles
  - Review scripts and tools
  - .gitignore setup
  - Makefile targets
  - Daily workflow examples

### CI_CD_SETUP.md (566 lines)
- **Purpose**: Complete CI/CD pipeline
- **Contents**:
  - GitHub Actions workflows
  - Commitlint configuration
  - Semantic versioning
  - Release automation
  - Commit format rules

### MASTER_TASK_PROMPT_SETUP.md (478 lines)
- **Purpose**: Master task system
- **Contents**:
  - 5-phase workflow
  - Rule integration
  - Token efficiency
  - Compliance checks
  - Examples and troubleshooting

### RULES_COMPLIANCE_SETUP.md (149+ lines)
- **Purpose**: Compliance system
- **Contents**:
  - Token efficiency strategies
  - Compliance check scripts
  - Makefile targets
  - Workflow integration

## 🔗 Related Documentation

### In ../rules/
These setup guides reference the authoritative rules in `../rules/`:
- **[../rules/AI_COMMIT_ATTRIBUTION.md](../rules/AI_COMMIT_ATTRIBUTION.md)** - Full AI attribution rules
- **[../rules/atomic-commits.md](../rules/atomic-commits.md)** - Complete atomic commit guide
- **[../rules/master-task-prompt.md](../rules/master-task-prompt.md)** - Master workflow prompt

### Quick References
- **[../rules/COMMIT_QUICK_REFERENCE.md](../rules/COMMIT_QUICK_REFERENCE.md)** - Commit quick ref
- **[../rules/COMPLIANCE_QUICK_REF.md](../rules/COMPLIANCE_QUICK_REF.md)** - Compliance checklist
- **[../CI_CD_QUICK_REF.md](../CI_CD_QUICK_REF.md)** - CI/CD quick ref

## 💡 Usage Tips

### When to Use These Guides

- **First Time Setup**: Follow all guides in recommended order
- **Quick Reference**: Use for specific setup questions
- **Troubleshooting**: Check individual guides for issues
- **Onboarding**: Share with new team members

### Reading Strategies

1. **Skim First**: Get overview of what's available
2. **Focus on Essentials**: Start with CI/CD and commits
3. **Deep Dive**: Read thoroughly when implementing
4. **Revisit**: Come back when troubleshooting

## ✅ Verification Checklist

After following the setup guides, verify:

- [ ] Can make conventional commits
- [ ] Commitlint validates commit messages
- [ ] AI attribution system works (if using AI)
- [ ] Atomic commit workflow understood
- [ ] Master task prompt available
- [ ] Compliance checks can run

## 🛠️ Setup Tools Reference

### Makefile Targets (from guides)
```bash
make test                # Run all tests
make coverage           # Generate coverage report
make review-commit      # Review atomic commits
make check-compliance   # Run compliance checks
```

### Git Hooks (if installed)
- **prepare-commit-msg** - AI attribution helper
- **commit-msg** - Commitlint validation

### Scripts (if installed)
- `scripts/review-commit.sh` - Atomic commit review
- `scripts/check-compliance.sh` - Compliance verification

## 📊 Setup Time Estimates

| Guide | Time Required | Complexity |
|-------|--------------|------------|
| CI_CD_SETUP.md | 30-45 min | Medium |
| AI_COMMIT_SETUP.md | 5-10 min | Low |
| ATOMIC_COMMITS_SETUP.md | 15-20 min | Low |
| MASTER_TASK_PROMPT_SETUP.md | 20-30 min | Medium |
| RULES_COMPLIANCE_SETUP.md | 10-15 min | Low |
| **Total** | **80-120 min** | - |

## 🔄 Keeping Setup Current

These guides document the setup process and should be updated when:
- Tools or dependencies change
- New automation is added
- Setup procedures are improved
- Issues are discovered and fixed

## ℹ️ Getting Help

If you encounter issues during setup:

1. Check the troubleshooting section in each guide
2. Review the related rules documentation
3. Check the quick reference guides
4. Verify prerequisites are met
5. Ensure you're in the project root directory

## 📝 Next Steps

After completing setup:

1. Read [../README.md](../../README.md) for project overview
2. Review [../rules/](../rules/) for development standards
3. Check [../PRD_MASTER.md](../PRD_MASTER.md) for product requirements
4. Explore [../features/](../features/) for feature specifications
5. Start developing following the master task workflow!

---

**Last Updated**: 2025-12-23
**Directory**: `docs/setup/`
**Status**: Complete and production-ready

