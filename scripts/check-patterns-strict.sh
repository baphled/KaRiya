#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

VIOLATIONS=0

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 STRICT PATTERN ENFORCEMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

NEW_FILES=$(git diff --cached --name-only --diff-filter=A | grep '\.go$' | grep -v '_test.go' || true)
MODIFIED_FILES=$(git diff --cached --name-only --diff-filter=M | grep '\.go$' | grep -v '_test.go' || true)
ALL_STAGED=$(git diff --cached --name-only | grep '\.go$' || true)

echo "📁 New files: $(echo "$NEW_FILES" | grep -c '.' || echo 0)"
echo "📝 Modified files: $(echo "$MODIFIED_FILES" | grep -c '.' || echo 0)"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1. DEPRECATED PATTERN CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

declare -A DEPRECATED_PATTERNS
DEPRECATED_PATTERNS=(
    ["components\.KeyBadge"]="primitives.HelpKeyBadge()"
    ["components\.StandardView"]="layout.NewScreenLayout()"
    ["components\.NewErrorModal"]="feedback.NewErrorModal()"
    ["components\.NewLoadingModal"]="feedback.NewLoadingModal()"
    ["components\.NewProgressModal"]="feedback.NewProgressModal()"
    ["components\.NewSuccessModal"]="feedback.NewSuccessModal()"
    ["components\.NewWarningModal"]="feedback.NewWarningModal()"
    ["components\.RenderOverlay"]="behaviors.RenderModalOverlay()"
    ["components\.RenderHelpFooter"]="primitives.RenderHelpFooter()"
    ["components\.ModalContainer"]="feedback.ModalContainer"
    ["components\.HelpModal"]="feedback.HelpModal"
    ["components\.ASCIILogo"]="DELETED - remove logo usage"
    ["CreateStandardView"]="layout.NewScreenLayout()"
    ["ThemedNavigationFooter"]="primitives.RenderHelpFooter()"
    ["ThemedCustomFooter"]="primitives.RenderHelpFooter()"
)

for file in $NEW_FILES; do
    for pattern in "${!DEPRECATED_PATTERNS[@]}"; do
        if grep -q "$pattern" "$file" 2>/dev/null; then
            echo -e "${RED}❌ NEW FILE uses deprecated pattern${NC}"
            echo "   File: $file"
            echo "   Pattern: $pattern"
            echo "   Use: ${DEPRECATED_PATTERNS[$pattern]}"
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    done
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2. HARDCODED COLORS CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $NEW_FILES; do
    HARDCODED=$(grep -n 'lipgloss\.Color("#' "$file" 2>/dev/null || true)
    if [ -n "$HARDCODED" ]; then
        echo -e "${RED}❌ NEW FILE has hardcoded colors${NC}"
        echo "   File: $file"
        echo "   Use theme.Primary(), theme.Secondary(), etc."
        echo "$HARDCODED" | head -3 | sed 's/^/   /'
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3. TODO/FIXME CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $ALL_STAGED; do
    TODOS=$(grep -n 'TODO\|FIXME\|XXX\|HACK' "$file" 2>/dev/null | grep -v '_test.go' || true)
    if [ -n "$TODOS" ]; then
        echo -e "${RED}❌ TODO/FIXME comments found${NC}"
        echo "   File: $file"
        echo "   Use task tracking instead of TODO comments"
        echo "$TODOS" | head -3 | sed 's/^/   /'
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4. SKIPPED/PENDING TESTS CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

TEST_FILES=$(git diff --cached --name-only | grep '_test.go$' || true)

for file in $TEST_FILES; do
    SKIPPED=$(grep -n 'XIt(\|XDescribe(\|XContext(\|Skip(' "$file" 2>/dev/null || true)
    if [ -n "$SKIPPED" ]; then
        echo -e "${RED}❌ Skipped tests found${NC}"
        echo "   File: $file"
        echo "   Skipped tests are not allowed"
        echo "$SKIPPED" | head -3 | sed 's/^/   /'
        VIOLATIONS=$((VIOLATIONS+1))
    fi

    PENDING=$(grep -n 'PIt(\|PDescribe(\|PContext(' "$file" 2>/dev/null || true)
    if [ -n "$PENDING" ]; then
        echo -e "${RED}❌ Pending tests found${NC}"
        echo "   File: $file"
        echo "   Pending tests are not allowed"
        echo "$PENDING" | head -3 | sed 's/^/   /'
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5. DEBUG STATEMENTS CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

DEBUG_PATTERNS='fmt\.Println\s*\(|fmt\.Printf\s*\(\s*"DEBUG|spew\.Dump|spew\.Printf|pp\.Print|pp\.Println|pretty\.Print'

for file in $ALL_STAGED; do
    if [[ ! "$file" =~ _test\.go$ ]]; then
        DEBUG=$(grep -nE "$DEBUG_PATTERNS" "$file" 2>/dev/null || true)
        if [ -n "$DEBUG" ]; then
            echo -e "${RED}❌ Debug statements found${NC}"
            echo "   File: $file"
            echo "   Remove debug statements before committing"
            echo "$DEBUG" | head -3 | sed 's/^/   /'
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6. INLINE COMMENTS CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $NEW_FILES; do
    INLINE_COUNT=0
    while IFS= read -r line; do
        if [[ "$line" =~ ^[[:space:]]*//[^/] ]]; then
            NEXT_LINE_NUM=$(($(echo "$line" | cut -d: -f1) + 1))
            NEXT_LINE=$(sed -n "${NEXT_LINE_NUM}p" "$file" 2>/dev/null || true)
            if [[ ! "$NEXT_LINE" =~ ^[[:space:]]*(func|type|var|const|package)[[:space:]] ]]; then
                if [[ ! "$line" =~ //go: ]]; then
                    INLINE_COUNT=$((INLINE_COUNT+1))
                fi
            fi
        fi
    done < <(grep -n '^[[:space:]]*//[^/]' "$file" 2>/dev/null || true)

    if [ "$INLINE_COUNT" -gt 0 ]; then
        echo -e "${RED}❌ Inline comments found in new file${NC}"
        echo "   File: $file"
        echo "   Count: $INLINE_COUNT"
        echo "   Code should be self-documenting. Refactor instead of commenting."
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "7. E2E TEST REQUIREMENTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

NEW_INTENTS=$(echo "$NEW_FILES" | grep '_intent\.go$' || true)
for intent in $NEW_INTENTS; do
    E2E_FILE="${intent%_intent.go}_e2e_test.go"
    if [ ! -f "$E2E_FILE" ]; then
        echo -e "${RED}❌ New intent missing E2E test${NC}"
        echo "   Intent: $intent"
        echo "   Required: $E2E_FILE"
        VIOLATIONS=$((VIOLATIONS+1))
    else
        if ! grep -q "Happy Paths" "$E2E_FILE"; then
            echo -e "${RED}❌ E2E missing Happy Paths section${NC}"
            echo "   File: $E2E_FILE"
            VIOLATIONS=$((VIOLATIONS+1))
        fi
        if ! grep -q "Sad Paths" "$E2E_FILE"; then
            echo -e "${RED}❌ E2E missing Sad Paths section${NC}"
            echo "   File: $E2E_FILE"
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

BUG_FILES=$(git diff --cached --name-only | grep 'bugs/BUG-' || true)
for bug in $BUG_FILES; do
    BUG_NUM=$(echo "$bug" | grep -oP 'BUG-\d+' || true)
    if [ -n "$BUG_NUM" ]; then
        REGRESSION_TEST=$(git diff --cached | grep -l "$BUG_NUM" 2>/dev/null | grep '_test.go' || true)
        if [ -z "$REGRESSION_TEST" ]; then
            echo -e "${RED}❌ Bug fix missing regression test${NC}"
            echo "   Bug: $BUG_NUM"
            echo "   Add E2E test with 'Bug Regressions' section mentioning $BUG_NUM"
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "8. DOCUMENTATION REQUIREMENTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$NEW_INTENTS" ]; then
    WORKFLOW_DOCS=$(git diff --cached --name-only | grep 'docs/workflows/' || true)
    if [ -z "$WORKFLOW_DOCS" ]; then
        echo -e "${RED}❌ New intent requires workflow documentation${NC}"
        echo "   Update docs/workflows/ with new intent documentation"
        VIOLATIONS=$((VIOLATIONS+1))
    fi
fi

NEW_UIKIT=$(echo "$NEW_FILES" | grep 'internal/cli/uikit/' || true)
if [ -n "$NEW_UIKIT" ]; then
    UIKIT_DOC=$(git diff --cached --name-only | grep 'docs/UIKIT_GUIDE.md' || true)
    if [ -z "$UIKIT_DOC" ]; then
        echo -e "${RED}❌ New UIKit component requires UIKIT_GUIDE.md update${NC}"
        VIOLATIONS=$((VIOLATIONS+1))
    fi
fi

NEW_BEHAVIORS=$(echo "$NEW_FILES" | grep 'internal/cli/behaviors/' || true)
if [ -n "$NEW_BEHAVIORS" ]; then
    DEV_DOCS=$(git diff --cached --name-only | grep 'docs/development/' || true)
    if [ -z "$DEV_DOCS" ]; then
        echo -e "${RED}❌ New behavior requires development documentation${NC}"
        VIOLATIONS=$((VIOLATIONS+1))
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "9. COVERAGE CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$ALL_STAGED" ]; then
    PACKAGES=$(echo "$ALL_STAGED" | xargs -I{} dirname {} | sort -u)
    for pkg in $PACKAGES; do
        if [ -d "$pkg" ]; then
            COVERAGE_OUTPUT=$(go test -cover "./$pkg" 2>/dev/null || true)
            COVERAGE=$(echo "$COVERAGE_OUTPUT" | grep -oP 'coverage: \K[0-9.]+' || echo "0")
            if [ -n "$COVERAGE" ]; then
                COVERAGE_INT=${COVERAGE%.*}
                if [ "$COVERAGE_INT" -lt 95 ]; then
                    echo -e "${RED}❌ Package coverage below 95%${NC}"
                    echo "   Package: $pkg"
                    echo "   Coverage: $COVERAGE%"
                    VIOLATIONS=$((VIOLATIONS+1))
                else
                    echo -e "${GREEN}✅ $pkg: ${COVERAGE}%${NC}"
                fi
            fi
        fi
    done
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ All strict pattern checks passed${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
else
    echo -e "${RED}❌ $VIOLATIONS violation(s) found${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 1
fi
