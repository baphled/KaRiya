# Breaking Change Notice: YAML CV Export Bullet Structure

## What Changed
- The YAML CV export now preserves bullet structure for `summary` and `highlights` fields as YAML lists (`[]string`).
- Previously, these fields were exported as block strings (prose). Now, they are exported as proper YAML arrays.

## Impact
- **YAML consumers** expecting block strings for `summary` or `highlights` must update their parsers to handle lists.
- **Text and Markdown exports** remain unchanged and unaffected by this refactor.

## Migration Guidance
- Update any downstream tools, scripts, or integrations that consume the YAML CV export to handle `summary` and `highlights` as lists.
- Example:

```yaml
summary:
  - "Heading text"
  - "Bullet one summary"
  - "Bullet two summary"
highlights:
  - "Highlight one"
  - "Highlight two"
```

## Reason
- This change ensures bullet structure is preserved for clarity and consistency across all CV export formats.
- It aligns YAML output with Text and Markdown exports, which already treat bullets as lists.

## Compatibility
- No changes to bullet generation or scoring logic.
- No changes to Text or Markdown export formats.

---

**If you rely on block string summary/highlights in YAML, update your code to handle lists.**
