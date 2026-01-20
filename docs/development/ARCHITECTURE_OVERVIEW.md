# Architecture Overview

KaRiya project architecture and structure.

---

## Architecture Overview

### Core Principles

The KaRiya TUI is built on a **type-safe, intent-driven architecture**:

1. **Type-Safe Intent Communication**: All intents communicate via `IntentResult[T]`
2. **Clear Intent Boundaries**: Each intent owns only its local state
3. **Predictable State Machines**: Explicit state transitions
4. **Back Navigation with Context**: Full state preservation via metadata
5. **Minimal Global State**: All mutations are local to intents
6. **Compile-Time Safety**: No runtime type assertions

**Visual Architecture**: See [TUI_INTENT_DIAGRAM.md](docs/TUI_INTENT_DIAGRAM.md) for complete architectural diagrams and intent workflows.

**State Matrix**: See [STATE_MATRIX.md](docs/STATE_MATRIX.md) for complete state documentation (10 intents, 64 states, escape behavior).

### The 5 Core Intents

| Intent | Purpose | States | Tests |
|--------|---------|--------|-------|
| CaptureEvent | Capture new career events | Choose Strategy → Form → Review → Confirm | 30+ |
| BrowseTimeline | View career timeline | Timeline → Event Detail | 37 |
| GenerateCV | Generate CVs | Profile → Audience → Role Emphasis → Length → Preview → Review → Confirm | 82+ |
| ExportArtifact | Export artifacts | Select → Configure → Preview → Export | 400+ |
| ConfigureSystem | System configuration | Domain → Settings → Staged Changes → Confirm | 400+ |

### Project Structure

```
internal/
├── cli/
│   ├── app/              # Root Bubble Tea model
│   ├── intents/          # Intent implementations
│   ├── models/           # Legacy UI components
│   ├── components/       # Reusable UI components
│   ├── uikit/            # UIKit Foundation (Task 41)
│   │   ├── theme/        # Theme infrastructure
│   │   ├── primitives/   # Text, Button, ButtonGroup, Input, Badge
│   │   └── containers/   # Box, Overlay
│   ├── behaviors/        # Embeddable Behaviors (Task 42)
│   │   ├── types.go      # Shared types
│   │   ├── table.go      # TableBehavior[T]
│   │   ├── crud.go       # CRUDBehavior[T]
│   │   ├── filter_menu.go # FilterMenuBehavior[T]
│   │   └── sort_menu.go  # SortMenuBehavior[T]
│   ├── context/          # GlobalContext
│   └── styles/           # Lipgloss styling
├── domain/career/        # Domain models
├── repository/career/    # Data access (SQLite)
└── service/career/       # Business logic
```

### Testing Strategy

#### Test Coverage
- **Overall**: 83.5% code coverage
- **Intent Framework**: 88.1%
- **GlobalContext**: 100%
- **Domain Models**: >95%
- **Repository**: >90%
- **Service**: >85%
- **CV Service**: 100% (203 tests, all passing)
- **UIKit Components**: 85%+ (138 tests, all passing)
- **Behaviors Package**: 85%+ (218 tests, all passing)

#### Running Tests
```bash
# All tests
go test -v ./...

# With race detector
go test -race ./...

# Specific package
go test -v ./internal/cli/intents/...

# Ginkgo with focus
ginkgo -r --focus="CaptureEvent" ./internal/cli/intents/

# Benchmarks
go test -bench=. ./internal/cli/intents/
```

#### Test Organization
- **Ginkgo + Gomega** framework for all tests
- **2,430+ total test specs** (all packages)
- **100% pass rate**
- **0 race conditions**
- **Execution time: <5s with race detector**

---
