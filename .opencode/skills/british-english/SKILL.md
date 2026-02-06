---
name: british-english
description: Enforce British English spelling, grammar, and conventions in all written content
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Ensure all written content uses British English spelling, grammar, and conventions. This applies to documentation, comments, commit messages, variable names, and user-facing text.

## When to use me

**Always active.** British English is the standard for this project.

## Spelling Differences

### Common Word Endings

| American | British | Example |
|----------|---------|---------|
| -ize | -ise | organise, recognise, initialise |
| -ization | -isation | organisation, initialisation |
| -yze | -yse | analyse, paralyse |
| -or | -our | colour, behaviour, favour, honour |
| -er | -re | centre, metre, theatre |
| -og | -ogue | catalogue, dialogue, analogue |
| -ense | -ence | licence (noun), defence, offence |
| -ed | -t | learnt, spelt, dreamt (both accepted) |

### Common Words

| American | British |
|----------|---------|
| analyze | analyse |
| behavior | behaviour |
| canceled | cancelled |
| center | centre |
| check (money) | cheque |
| color | colour |
| defense | defence |
| dialog | dialogue |
| favor | favour |
| fulfill | fulfil |
| gray | grey |
| honor | honour |
| initialize | initialise |
| jewelry | jewellery |
| labor | labour |
| license (verb) | license |
| license (noun) | licence |
| meter | metre |
| modeling | modelling |
| neighbor | neighbour |
| offense | offence |
| optimize | optimise |
| organize | organise |
| practice (verb) | practise |
| practice (noun) | practice |
| program | programme (but: computer program) |
| realize | realise |
| recognize | recognise |
| serialize | serialise |
| signaling | signalling |
| skeptic | sceptic |
| traveled | travelled |
| vapor | vapour |

### Technical Terms

Some technical terms retain American spelling by convention:

```
KEEP AMERICAN (industry standard):
- program (computer program, not programme)
- disk (hard disk, not disc - but: compact disc)
- analog (in computing contexts)

USE BRITISH:
- colour (in variable names for user-facing features)
- behaviour (BehaviourTree, not BehaviorTree)
- initialise, serialise, etc.
```

## Application in Code

### Variable and Function Names

```go
// CORRECT - British spelling
func (s *Service) InitialiseConnection() error
func (c *Config) GetColour() string
type TableBehaviour struct{}
func (e *Event) Serialise() ([]byte, error)
func Analyse(data []byte) (*Report, error)

// INCORRECT - American spelling
func (s *Service) InitializeConnection() error  // ✗
func (c *Config) GetColor() string              // ✗
type TableBehavior struct{}                     // ✗
```

### Comments and Documentation

```go
// CORRECT
// Initialise sets up the connection and prepares the service for use.
// It honours the configuration settings and organises resources accordingly.
//
// Expected: config with valid centre coordinates
// Returns: initialised service or error
// Side effects: Creates colour-coded log entries

// INCORRECT
// Initialize sets up the connection and prepares the service for use.
// It honors the configuration settings and organizes resources accordingly.
```

### User-Facing Text

```go
// CORRECT
const (
    ErrInvalidColour = "invalid colour specified"
    MsgInitialising  = "Initialising application..."
    LabelFavourites  = "Favourites"
)

// INCORRECT
const (
    ErrInvalidColor = "invalid color specified"
    MsgInitializing = "Initializing application..."
    LabelFavorites  = "Favorites"
)
```

### Commit Messages

```bash
# CORRECT
feat(ui): add colour picker to preferences
fix(service): initialise connection before use
docs: standardise British English throughout

# INCORRECT
feat(ui): add color picker to preferences
fix(service): initialize connection before use
docs: standardize British English throughout
```

## Grammar Conventions

### Collective Nouns

British English treats collective nouns as plural when referring to members:

```
# British
The team are working on the feature.
The company have released their update.

# American (also acceptable in formal writing)
The team is working on the feature.
```

### Quotation Marks

British English traditionally uses single quotes, with double inside:

```
# British
The error message said 'invalid input'.
She said, 'The log shows "connection refused".'

# American
The error message said "invalid input".
```

### Dates

```
# British (day/month/year)
15/03/2024
15 March 2024

# American (month/day/year)
03/15/2024
March 15, 2024

# ISO (preferred in code)
2024-03-15
```

## Exceptions

### Third-Party APIs

When interfacing with external APIs that use American spelling, match their conventions:

```go
// External API uses American spelling - match it
response.Color = api.GetColor()

// Internal code uses British
internalColour := convertToOurFormat(response.Color)
```

### Industry Standard Terms

Some terms are universally American in computing:

```
KEEP AS-IS:
- program (not programme for software)
- disk (hard disk, SSD)
- check (as in health check - but: cheque for payment)
- color (in CSS, as it's the property name)
```

### Existing Codebase Consistency

When modifying existing code that uses American spelling:
1. **Don't mix** - Keep the file consistent
2. **Note for refactoring** - Create tech debt item if widespread
3. **New files** - Always use British spelling

## Enforcement

### Pre-commit Check

Common misspellings to catch:

```bash
# Words that should be British
grep -rn "initialize\|behavior\|color\|center\|analyze" --include="*.go" --include="*.md"
```

### IDE Configuration

VS Code settings:
```json
{
  "cSpell.language": "en-GB",
  "cSpell.words": [
    "behaviour",
    "colour",
    "initialise",
    "serialise"
  ]
}
```

### Review Checklist

When reviewing:
- [ ] Variable names use British spelling
- [ ] Comments use British spelling
- [ ] User-facing strings use British spelling
- [ ] Documentation uses British spelling
- [ ] Commit messages use British spelling

## Quick Reference Card

| Check For | Use Instead |
|-----------|-------------|
| -ize | -ise |
| -ization | -isation |
| -or | -our |
| -er (centre) | -re |
| -ed (traveled) | -lled |
| gray | grey |
| check (money) | cheque |
| defense | defence |
| license (n) | licence |
| practice (v) | practise |
| program (non-tech) | programme |

## Common Mistakes

| Mistake | Correction |
|---------|------------|
| `InitializeService` | `InitialiseService` |
| `colorScheme` | `colourScheme` |
| `TableBehavior` | `TableBehaviour` |
| `centerAlign` | `centreAlign` |
| `favoriteItems` | `favouriteItems` |
| `Analyzing data` | `Analysing data` |
| `optimization` | `optimisation` |

## Related Skills

- `clean-code` - Code style and naming
- `code-reviewer` - Review for British English
- `ai-commit` - Commit message conventions
