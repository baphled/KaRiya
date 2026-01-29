#!/bin/bash
# Check that all exported Go symbols have documentation comments.
# Uses standalone revive because golangci-lint's revive integration
# does not report missing doc comments (only stuttering names).

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Ensure revive is installed.
if ! command -v revive &>/dev/null; then
    echo -e "${RED}revive not found. Install with: go install github.com/mgechev/revive@latest${NC}"
    exit 1
fi

echo "Checking exported symbol documentation..."

# Run revive with the exported rule only.
# Exclude test files and deprecated packages (models/, components/).
OUTPUT=$(revive -config revive.toml \
    -exclude ./internal/cli/models/... \
    -exclude ./internal/cli/components/... \
    ./... 2>&1 || true)

# Filter to only exported rule violations (ignore dot-imports, etc.).
VIOLATIONS=$(echo "$OUTPUT" | grep "exported" | grep -v "_test.go" || true)

if [ -z "$VIOLATIONS" ]; then
    echo -e "${GREEN}All exported symbols have documentation comments.${NC}"
    exit 0
fi

COUNT=$(echo "$VIOLATIONS" | wc -l)
echo -e "${RED}Found $COUNT undocumented exported symbols:${NC}"
echo "$VIOLATIONS"
echo ""
echo "Fix: Add a doc comment above each exported type, function, var, or const."
echo "See: docs/conventions/GO_DOCUMENTATION_RULES.md"
exit 1
