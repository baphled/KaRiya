# KaRiya Documentation

Welcome to the KaRiya Career Journal & CV Generator documentation. This directory contains all project documentation organized by purpose.

## 📚 Documentation Structure

### Product Requirements Documents (PRDs)
- **[PRD_MASTER.md](PRD_MASTER.md)** - Authoritative product requirements document with complete system architecture, data models, and CV generation rules
- **[PRD_USER_STORIES.md](PRD_USER_STORIES.md)** - User stories, functional requirements, and phased rollout approach
- **[PRD_CLI.md](PRD_CLI.md)** - CLI interface specifications using BubbleTea framework

### Setup Guides
- **[setup/](setup/)** - All setup and installation guides
  - AI commit attribution setup
  - Atomic commits setup
  - CI/CD pipeline setup
  - Master task prompt setup
  - Rules compliance setup

### Development Rules & Guidelines
- **[rules/](rules/)** - Development standards and guidelines (16 documents)
  - AI commit attribution rules
  - Atomic commit guidelines
  - Go coding standards
  - Task execution workflows
  - Token efficiency best practices
  - Quick reference guides

### Feature Specifications
- **[features/](features/)** - Detailed feature specifications (16 documents)
  - Career event capture
  - Burst & fact extraction
  - CV generation
  - Data model schema
  - User experience principles
  - Role & audience filtering
  - Metadata validation
  - Export & integration

### TUI Development (MANDATORY)
- **[UIKIT_GUIDE.md](UIKIT_GUIDE.md)** - **UIKit component library** (REQUIRED for all new code)
  - Primitives: Text, Button, Badge, Input
  - Containers: Box, Overlay
  - Nil theme guard patterns
  - Migration guide from lipgloss
- **[TUI_DEVELOPER_GUIDE.md](TUI_DEVELOPER_GUIDE.md)** - Comprehensive TUI development guide
- **[TUI_STANDARDS.md](TUI_STANDARDS.md)** - Design principles and accessibility
- **[THEME_CUSTOMIZATION_GUIDE.md](THEME_CUSTOMIZATION_GUIDE.md)** - Theme system documentation
- **[MODAL_PATTERNS.md](MODAL_PATTERNS.md)** - Modal implementation patterns
- **[STANDARDVIEW_GUIDE.md](STANDARDVIEW_GUIDE.md)** - StandardView system

### User & Developer Guides
- **[guides/](guides/)** - User and developer guides (8 documents)
  - CV generation guide and examples
  - CV troubleshooting
  - View patterns guide
  - Error handling guide
  - List model rendering specification
  - Focus indicator guide
  - Style usage guide

### Audits & Analysis
- **[audits/](audits/)** - Code audits and analysis reports (9 documents)
  - Legacy screen audits
  - Intent framework readiness
  - UI/UX consistency audits
  - Focus indicator audits
  - Error display audits
  - Navigation audits

### Infrastructure & Testing
- **[CI_CD_PIPELINE.md](CI_CD_PIPELINE.md)** - Complete CI/CD pipeline documentation with GitHub Actions, commitlint, and semantic-release
- **[CI_CD_QUICK_REF.md](CI_CD_QUICK_REF.md)** - Quick reference for CI/CD operations
- **[integration-test-strategy.md](integration-test-strategy.md)** - Integration testing strategy and scenarios

### Tools & Editor Setup
- **[tools/](tools/)** - Tool-specific documentation
  - Editor configuration (Neovim/neotest setup)

### Reports & Analysis
- **[reports/](reports/)** - Test reports and analysis documents
  - Test verification reports (dated)

### Archived Documentation
- **[archive/](archive/)** - Historical documentation (52+ documents)
  - **[archive/phases/](archive/phases/)** - Completed phase reports (20 files)
  - **[archive/plans/](archive/plans/)** - Completed implementation plans (6 files)
  - **[archive/app-migration/](archive/app-migration/)** - App.go migration documentation (6 files)
  - **[archive/investigations/](archive/investigations/)** - One-time analysis reports (10 files)
  - **[archive/sessions/](archive/sessions/)** - Session summaries (6 files)
  - **[archive/summaries/](archive/summaries/)** - Implementation summaries (2 files)
  - **[archive/proposals/](archive/proposals/)** - Architecture proposals (2 files)
  - **[archive/rules/](archive/rules/)** - Archived rule documents (5 files)

## 🚀 Quick Start

### For New Developers
1. Start with [README.md](../README.md) in project root
2. Read [AGENTS.md](../AGENTS.md) for comprehensive handover
3. Review [PRD_MASTER.md](PRD_MASTER.md) for product overview
4. Check [setup/](setup/) for environment setup
5. Review [rules/](rules/) for development standards

### For Setting Up Your Environment
1. Follow guides in [setup/](setup/)
2. Configure your editor using [tools/editor-setup/](tools/editor-setup/)
3. Review [CI_CD_PIPELINE.md](CI_CD_PIPELINE.md) for commit standards

### For Understanding Features
1. Review [PRD_MASTER.md](PRD_MASTER.md) for complete system design
2. Dive into [features/](features/) for specific feature details
3. Check [PRD_CLI.md](PRD_CLI.md) for CLI interface specs
4. Read [guides/](guides/) for user and developer guides

### For Development Workflows
1. Read [rules/master-task-prompt.md](rules/master-task-prompt.md) for the complete workflow
2. Use quick references in [rules/](rules/) for daily operations
3. Follow [rules/atomic-commits.md](rules/atomic-commits.md) for commit guidelines

## 📖 Document Categories

### Rules (Mandatory)
These documents define mandatory development practices:
- AI commit attribution
- Atomic commits
- Go guidelines
- Code review process

### Guidelines (Best Practices)
These documents provide recommended best practices:
- Senior engineer guidelines
- Token efficiency
- Task processing
- Compliance checks

### Quick References
One-page cheat sheets for common operations:
- [rules/COMMIT_QUICK_REFERENCE.md](rules/COMMIT_QUICK_REFERENCE.md)
- [rules/COMPLIANCE_QUICK_REF.md](rules/COMPLIANCE_QUICK_REF.md)
- [rules/TASK_QUICK_REF.md](rules/TASK_QUICK_REF.md)
- [CI_CD_QUICK_REF.md](CI_CD_QUICK_REF.md)
- [IMPLEMENTATION_INDEX.md](IMPLEMENTATION_INDEX.md) - Complete implementation documentation index

### Process Documentation
Step-by-step process guides:
- [rules/generate-prd.md](rules/generate-prd.md)
- [rules/generate-tasks.md](rules/generate-tasks.md)
- [rules/master-task-prompt.md](rules/master-task-prompt.md) (includes task processing rules)

## 🔍 Finding What You Need

### "How do I set up...?"
→ Check [setup/](setup/)

### "What are the rules for...?"
→ Check [rules/](rules/)

### "What features does the system have?"
→ Check [PRD_MASTER.md](PRD_MASTER.md) or [features/](features/)

### "How do I implement...?"
→ Check [rules/master-task-prompt.md](rules/master-task-prompt.md)

### "What's the commit format?"
→ Check [CI_CD_QUICK_REF.md](CI_CD_QUICK_REF.md)

### "How do I run tests?"
→ Check [integration-test-strategy.md](integration-test-strategy.md) or [tools/editor-setup/](tools/editor-setup/)

### "Where are the old phase reports?"
→ Check [archive/phases/](archive/phases/)

### "How was feature X implemented?"
→ Check [archive/plans/](archive/plans/) or [archive/investigations/](archive/investigations/)

## 🎯 Documentation Standards

### Document Types
- **PRD**: Product requirements with user stories and acceptance criteria
- **Rules**: Mandatory development standards
- **Guidelines**: Recommended best practices
- **Setup**: Installation and configuration instructions
- **Reference**: Quick lookup cheat sheets
- **Process**: Step-by-step workflows

### File Naming Conventions
- `PRD_*.md` - Product requirements documents
- `*_QUICK_REF.md` - Quick reference guides
- `*_SETUP.md` - Setup and installation guides
- `[number]-*.md` - Sequential feature specifications
- Lowercase with hyphens for most other files

### Cross-References
Documents are heavily cross-referenced. Follow links to related topics for comprehensive understanding.

## 📊 Documentation Statistics

- **Total Documents**: ~160 files (108 active + 52 archived)
- **PRDs**: 3 files
- **Setup Guides**: 7 files
- **Rules & Guidelines**: 13 files
- **Feature Specs**: 16 files
- **User/Developer Guides**: 8 files
- **Audits**: 9 files
- **Infrastructure Docs**: 3 files
- **Reports**: 2 files
- **Archived**: 52 files (in archive/ subdirectories)

## 🛠️ Maintaining Documentation

### When to Update
- Feature implementations: Update PRDs and feature specs
- Process changes: Update rules and guidelines
- New tools: Add to tools/ directory
- Test results: Add dated reports to reports/

### Conventions
- Keep documents focused and single-purpose
- Create quick references for frequently accessed info
- Date reports and time-sensitive documents
- Cross-reference related documents
- Use Markdown formatting consistently

## 📝 Related Files

- **[../README.md](../README.md)** - Project README
- **[../AGENTS.md](../AGENTS.md)** - Comprehensive handover document
- **[../Makefile](../Makefile)** - Build and test commands
- **[../tasks/](../tasks/)** - Task tracking documents

## ℹ️ Help & Support

For questions about:
- **Product features**: See PRDs
- **Development process**: See rules/
- **Setup issues**: See setup/
- **Testing**: See integration-test-strategy.md
- **CI/CD**: See CI_CD_PIPELINE.md

---

**Last Updated**: 2026-01-08
**Status**: Production-ready documentation (recently reorganized)
**Maintainer**: Development Team

