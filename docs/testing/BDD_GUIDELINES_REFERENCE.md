# BDD Guidelines - KaRiya Project Reference

**Status**: KaRiya-Specific Bridge Document  
**Date**: February 12, 2026  

---

## ⚠️ Important: Global Standards Now Available

KaRiya's BDD testing standards have been integrated into a **Global BDD Standards Library** that applies to ALL projects, not just KaRiya.

### Global Standards Location
```
/home/baphled/Projects/Standards/BDD/
```

### Why This Matters
- **Universal**: Applies to web apps, APIs, CLIs, mobile, TUI, streaming, batch processing
- **Framework-Neutral**: Works with Cucumber, Gherkin, Behave, Jest, Ginkgo, etc.
- **Language-Neutral**: Go, Python, JavaScript, Java, Ruby, C#, PHP, etc.
- **Single Source of Truth**: One set of standards for all projects
- **Always Current**: Updates in one place benefit all projects

---

## Quick Navigation

### For KaRiya Developers

1. **Quick Reference** (5 min read)
   - See [Global BDD README](file:///home/baphled/Projects/Standards/BDD/README.md)

2. **Best Practices** (20 min read)
   - See [Global BDD Best Practices](file:///home/baphled/Projects/Standards/BDD/BDD_BEST_PRACTICES.md)
   - Focus on TUI/CLI examples section

3. **Good vs Bad Examples** (15 min read)
   - See [Global Examples](file:///home/baphled/Projects/Standards/BDD/BDD_EXAMPLES_GOOD_VS_BAD.md)
   - Review "CLI Tools" and "TUI Applications" sections

4. **Anti-Patterns** (Reference material)
   - See [Global Anti-Patterns](file:///home/baphled/Projects/Standards/BDD/BDD_ANTI_PATTERNS.md)
   - Use in code reviews to identify issues

5. **Decision Framework** (Decision making)
   - See [Global Decision Framework](file:///home/baphled/Projects/Standards/BDD/BDD_DECISION_FRAMEWORK.md)
   - Use when deciding BDD vs Unit vs Integration tests

### For Code Review in KaRiya

**Check scenario against**:
1. Global BDD Anti-Patterns (identify common mistakes)
2. Global Examples (compare to similar scenario type)
3. Global Decision Framework (verify it's BDD, not Unit/Integration)

---

## KaRiya-Specific Decisions

The following have been decided for KaRiya and align with global standards:

### ✅ What We Test in BDD (KaRiya)

- Career event capture workflows
- Event enrichment and review processes
- Burst suggestion and acceptance
- Skill inference and management
- Multi-feature integration (events + skills + bursts)
- Error handling and edge cases
- Critical user paths (main functionality)

**Global Principle**: Test business outcomes, not implementation

### ❌ What We Don't Test in BDD (KaRiya)

- Keyboard shortcut handling (j/k, arrows, /, f, s, etc.)
- Modal opening/closing mechanics
- Form field tab navigation and focus
- Button styling and visibility
- Animation timing
- Screen transitions
- Focus management
- Pure navigation logic

**Global Principle**: These belong in Unit Tests

### 🎯 Put Instead in Unit Tests (KaRiya)

- Keyboard handlers (individual key handlers)
- Modal display and lifecycle
- Form field mechanics (tab order, focus)
- UI styling and appearance
- Animation timing
- Focus management
- Input validation
- Field clearing

**Global Principle**: Unit tests cover implementation details

---

## Sync with Global Standards

### Current Status
✅ KaRiya BDD cleanup aligned with global standards  
✅ All 277 remaining BDD scenarios follow global principles  
✅ All unit tests follow global allocation (40% of test suite)  
✅ All commits follow global message standards  

### When Global Standards Update
If global BDD standards are updated, KaRiya will:
1. Review the update
2. Assess impact on current scenarios
3. Update any affected KaRiya documentation
4. Train team on changes
5. Refactor scenarios if needed

See [Global Sync Instructions](file:///home/baphled/Projects/Standards/BDD/SYNC_INSTRUCTIONS.md) for details.

---

## Local KaRiya Documentation

While global standards are our primary reference, KaRiya also maintains:

### Original KaRiya Documents
- **BDD_GUIDELINES.md** - Archived KaRiya-specific version
- **BDD_GOOD_VS_BAD_EXAMPLES.md** - Archived KaRiya examples

These are kept for historical reference but **should not be updated** as they're superseded by global standards.

---

## Getting Help

### Questions About BDD?

1. **"When should I use BDD?"**
   → See Global [Decision Framework](file:///home/baphled/Projects/Standards/BDD/BDD_DECISION_FRAMEWORK.md)

2. **"Is this an anti-pattern?"**
   → See Global [Anti-Patterns](file:///home/baphled/Projects/Standards/BDD/BDD_ANTI_PATTERNS.md)

3. **"How do I write a good scenario?"**
   → See Global [Best Practices](file:///home/baphled/Projects/Standards/BDD/BDD_BEST_PRACTICES.md)

4. **"Show me an example similar to mine"**
   → See Global [Examples](file:///home/baphled/Projects/Standards/BDD/BDD_EXAMPLES_GOOD_VS_BAD.md)

5. **"How do I adopt these standards?"**
   → See Global [Adoption Guide](file:///home/baphled/Projects/Standards/BDD/ADOPTING_IN_YOUR_PROJECT.md)

---

## Summary

**KaRiya uses the Global BDD Standards Library for all testing guidance.**

All future BDD scenarios should:
- ✅ Reference the global standards
- ✅ Follow the global anti-patterns checklist
- ✅ Use the global decision framework
- ✅ Align with global principles

**See**: `/home/baphled/Projects/Standards/BDD/README.md` to get started.

---

**Last Updated**: February 12, 2026  
**Maintains**: Alignment with global BDD standards
