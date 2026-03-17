#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Generating GoMock mocks..."
echo "Project root: $PROJECT_ROOT"
echo ""

cd "$PROJECT_ROOT"

echo "=== Repository layer ==="
go generate ./internal/repository/career/...

echo "=== Service layer ==="
go generate ./internal/service/career/cv/...
go generate ./internal/service/career/skillinference/...

echo "=== Intent layer ==="
go generate ./internal/tui/intents/captureevent/...
go generate ./internal/tui/intents/burst_management/...
go generate ./internal/tui/intents/skillsmanagement/...
go generate ./internal/tui/intents/browsetimeline/...

echo ""
echo "Mock generation complete!"
echo "Generated files:"
find "$PROJECT_ROOT/internal/testutil/mocks" -name "*_mock.go" -type f | sort
