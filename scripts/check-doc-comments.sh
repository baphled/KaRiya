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

# Exclude false positives for Go generics.
# revive does not understand that "// TypeName[T] ..." is valid godoc for generic types.
# For each violation about comment form, check if the source comment uses generic syntax.
if [ -n "$VIOLATIONS" ]; then
    FILTERED=""
    while IFS= read -r line; do
        if echo "$line" | grep -q 'should be of the form'; then
            FILE=$(echo "$line" | cut -d: -f1)
            LINE_NUM=$(echo "$line" | cut -d: -f2)
            SYMBOL=$(echo "$line" | sed -n 's/.*on exported type \(\w\+\).*/\1/p')
            if [ -z "$SYMBOL" ]; then
                SYMBOL=$(echo "$line" | sed -n 's/.*on exported function \(\w\+\).*/\1/p')
            fi
            if [ -n "$SYMBOL" ]; then
                COMMENT=$(sed -n "${LINE_NUM}p" "$FILE" 2>/dev/null || true)
                if echo "$COMMENT" | grep -q "// ${SYMBOL}\["; then
                    continue
                fi
            fi
        fi
        FILTERED="${FILTERED}${line}
"
    done <<< "$VIOLATIONS"
    VIOLATIONS=$(echo "$FILTERED" | sed '/^$/d')
fi

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
