# Testing with Mocks

## Overview

KaRiya uses [GoMock](https://github.com/golang/mock) for all test doubles, with centralized
location and auto-generation via `go:generate` directives.

## Directory Structure

```
internal/testutil/mocks/
├── doc.go                          # Package documentation
├── repository/                     # Repository layer mocks (package: mockrepo)
│   ├── event_repository_mock.go     # Generated from EventRepository
│   ├── burst_repository_mock.go     # Generated from BurstRepository
│   ├── fact_repository_mock.go      # Generated from FactRepository
│   └── skill_repository_mock.go    # Generated from SkillRepository
├── service/                        # Service layer mocks (package: mocksvc)
│   ├── cv_generation_service_mock.go   # Generated from CVGenerationService
│   ├── clipboard_writer_mock.go        # Generated from ClipboardWriter
│   └── skill_inference_service_mock.go # Generated from SkillInferenceService
└── intent/                         # Intent layer mocks (package: mockintent)
    ├── burst_service_mock.go                    # Generated from BurstService
    ├── burst_skill_inference_service_mock.go     # Generated from burst_management.SkillInferenceService
    ├── browse_event_service_mock.go             # Generated from browsetimeline.EventService
    ├── capture_event_service_mock.go            # Generated from captureevent.EventService
    └── skills_skill_inference_service_mock.go   # Generated from skillsmanagement.SkillInferenceService
```

## Generating Mocks

### Using Makefile (recommended)

```bash
make generate-mocks
```

### Using go generate directly

```bash
# All packages
go generate ./...

# Specific layer
go generate ./internal/repository/career/...
go generate ./internal/service/career/...
go generate ./internal/cli/intents/...
```

### Checking mocks are up to date (CI)

```bash
make check-mocks-updated
```

## Import Patterns

Each layer has a distinct package name to avoid collisions:

### Repository Mocks

```go
import (
    mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
    "github.com/golang/mock/gomock"
)

ctrl := gomock.NewController(GinkgoT())

mockRepo := mockrepo.NewMockEventRepository(ctrl)
mockRepo.EXPECT().
    GetByID(gomock.Any(), gomock.Eq("event-1")).
    Return(testEvent, nil).
    Times(1)
```

### Service Mocks

```go
import (
    mocksvc "github.com/baphled/kariya/internal/testutil/mocks/service"
    "github.com/golang/mock/gomock"
)

ctrl := gomock.NewController(GinkgoT())

mockCV := mocksvc.NewMockCVGenerationService(ctrl)
mockCV.EXPECT().
    GenerateCV(gomock.Any(), gomock.Eq("default")).
    Return(testCVView, nil)
```

### Intent Mocks

Intent mocks use prefixed names to disambiguate interfaces with the same name
from different packages:

```go
import (
    mockintent "github.com/baphled/kariya/internal/testutil/mocks/intent"
    "github.com/golang/mock/gomock"
)

ctrl := gomock.NewController(GinkgoT())

// burst_management.BurstService
mockBurst := mockintent.NewMockBurstService(ctrl)

// browsetimeline.EventService (prefixed to avoid collision)
mockBrowse := mockintent.NewMockBrowseEventService(ctrl)

// captureevent.EventService (prefixed to avoid collision)
mockCapture := mockintent.NewMockCaptureEventService(ctrl)

// burst_management.SkillInferenceService (prefixed)
mockBurstSkill := mockintent.NewMockBurstSkillInferenceService(ctrl)

// skillsmanagement.SkillInferenceService (prefixed)
mockSkillsSkill := mockintent.NewMockSkillsSkillInferenceService(ctrl)
```

## Usage Examples

### Success Path

```go
mockRepo.EXPECT().
    GetByID(gomock.Any(), gomock.Eq("event-1")).
    Return(testEvent, nil).
    Times(1)
```

### Error Path

```go
mockRepo.EXPECT().
    GetByID(gomock.Any(), gomock.Eq("event-1")).
    Return(nil, errors.New("not found")).
    Times(1)
```

### Any Parameters

```go
mockRepo.EXPECT().
    List(gomock.Any(), gomock.Any()).
    Return(testEvents, nil)
```

### Custom Logic with DoAndReturn

```go
mockRepo.EXPECT().
    Create(gomock.Any(), gomock.Any()).
    DoAndReturn(func(_ context.Context, event *career.Event) error {
        event.ID = "generated-id"
        return nil
    }).
    Times(1)
```

### Ordered Calls

```go
gomock.InOrder(
    mockRepo.EXPECT().
        GetByID(gomock.Any(), "1").
        Return(event1, nil),
    mockRepo.EXPECT().
        GetByID(gomock.Any(), "2").
        Return(event2, nil),
)
```

### Call Count Tracking

GoMock verifies call counts automatically via `Times()`. For custom counting:

```go
var callCount int
mockSvc.EXPECT().
    GenerateCV(gomock.Any(), gomock.Any()).
    DoAndReturn(func(_ context.Context, _ string) (*career.CVView, error) {
        callCount++
        return testCV, nil
    }).
    AnyTimes()

// After test
Expect(callCount).To(Equal(2))
```

## Adding go:generate to New Interfaces

When creating a new interface that needs mocking, add a `go:generate` directive
above the `package` declaration:

```go
//go:generate mockgen -destination=<relative-path-to-mocks-subdir>/<name>_mock.go -package=<mock-package> <full-import-path> <InterfaceName>

package mypackage
```

**Path calculation**: The `-destination` path is relative to the source file's directory.

**Package names**:
- Repository interfaces: `-package=mockrepo`
- Service interfaces: `-package=mocksvc`
- Intent interfaces: `-package=mockintent`

**Name collisions**: If the interface name conflicts with one from another package
(e.g., `EventService` exists in both `captureevent` and `browsetimeline`), use
`-mock_names=InterfaceName=MockPrefixedName`:

```go
//go:generate mockgen -destination=... -package=mockintent -mock_names=EventService=MockMyFeatureEventService ...
```

## Migration from Hand-Written Mocks

The following hand-written mocks are deprecated and will be migrated in follow-up PRs:

| Old Mock | Location | Replacement |
|----------|----------|-------------|
| `mocks.BurstServiceMock` | `testutil/mocks/burst_service_mock.go` | `mockintent.NewMockBurstService` |
| `mocks.BurstRepositoryMock` | `testutil/mocks/burst_service_mock.go` | `mockrepo.NewMockBurstRepository` |
| `mocks.TestMockRepository` | `repository/career/mocks/mock_helper.go` | `mockrepo.NewMockEventRepository` |
| `mocks.MockEventRepository` (gomock) | `repository/career/mocks/repository_mock.go` | `mockrepo.NewMockEventRepository` |
| `e2e.MockCVGenerationService` | `testutil/e2e/cv_mocks.go` | `mocksvc.NewMockCVGenerationService` |
| `e2e.MockClipboardWriter` | `testutil/e2e/cv_mocks.go` | `mocksvc.NewMockClipboardWriter` |

Use `scripts/migrate-mock-imports.sh` for automated import migration.

## Troubleshooting

### "mockgen: command not found"

Install mockgen:

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

### Generated mock has wrong package name

Check the `-package=` flag in the `go:generate` directive matches the `doc.go`
package declaration in the target directory.

### Name collision in mockintent package

Use `-mock_names=InterfaceName=MockPrefixedName` to give the mock a unique name.

### Mock file not generated in expected location

The `-destination` path is relative to the **source file's directory**, not the
project root. Verify the path calculation by counting directory levels.
