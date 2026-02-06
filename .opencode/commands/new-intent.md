---
description: Create a new intent with proper architecture
agent: build
---

Create a new intent following KaRiya architecture.

Load these skills:
- `create-intent` - Intent structure and requirements
- `create-screen` - Screen structure (intents need screens)
- `architecture` - Layer boundaries and patterns
- `clean-code` - Clean implementation
- `component-lookup` - Find correct components to use

## Intent
$ARGUMENTS

## Process

1. **Generate Structure**
   ```bash
   make new-intent NAME=$1
   ```

2. **Implement Core Files**
   - `context.go` - Input parameters with Validate()
   - `constants.go` - State enum
   - `messages.go` - Message types
   - `result.go` - Result type
   - `intent.go` - Implementation

3. **Create Screens** in `screens/{feature}/`
   - List screen with TableBehavior
   - Detail screen (if needed)
   - Form screen (if needed)
   - Modals in `screens/{feature}/modals/`

4. **Write Tests** for each component

5. **Validate**
   ```bash
   make check-intent-architecture
   make check-compliance
   ```

## Requirements

- Embed *BaseIntent
- Implement ScreenResultHandler
- Use typed state enum
- Screens return ScreenResult, don't mutate state
- All handle* functions in handlers.go
- intent.go < 400 lines
