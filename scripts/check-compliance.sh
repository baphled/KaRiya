#!/bin/bash

set -e

echo "================================================"
echo "🔍 RULES COMPLIANCE CHECK"
echo "================================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

VIOLATIONS=0
WARNINGS=0

# Helper function for checks
check_pass() {
    echo -e "${GREEN}✅ Pass${NC}"
}

check_fail() {
    echo -e "${RED}❌ Fail${NC} - $1"
    VIOLATIONS=$((VIOLATIONS+1))
}

check_warn() {
    echo -e "${YELLOW}⚠️  Warning${NC} - $1"
    WARNINGS=$((WARNINGS+1))
}

# ============================================
# 1. CODE QUALITY CHECKS
# ============================================
echo "📋 CODE QUALITY"
echo "------------------------------------------------"

# Formatting
echo -n "Formatting (gofmt): "
UNFORMATTED=$(gofmt -l . 2>&1 | grep -v '^vendor/' | grep '\.go$' || true)
if [ -z "$UNFORMATTED" ]; then
    check_pass
else
    check_fail "Run: go fmt ./..."
    echo "  Unformatted files:"
    echo "$UNFORMATTED" | sed 's/^/    /'
fi

# Build
echo -n "Build: "
if go build ./... > /dev/null 2>&1; then
    check_pass
else
    check_fail "Fix build errors"
fi

# Tests
echo -n "Tests: "
if go test ./... > /dev/null 2>&1; then
    check_pass
else
    check_fail "Fix failing tests"
fi

# Race Conditions
echo -n "Race Detection: "
if go test -race ./... > /dev/null 2>&1; then
    check_pass
else
    check_fail "Fix race conditions"
fi

# Vet
echo -n "Go Vet: "
if go vet ./... > /dev/null 2>&1; then
    check_pass
else
    check_fail "Fix vet warnings"
fi

echo ""

# ============================================
# 2. TEST COVERAGE
# ============================================
echo "📊 TEST COVERAGE"
echo "------------------------------------------------"

if command -v ginkgo &> /dev/null; then
    # Using Ginkgo
    COVERAGE_OUTPUT=$(go test -cover ./... 2>/dev/null || true)
    if [ -n "$COVERAGE_OUTPUT" ]; then
        COVERAGE=$(echo "$COVERAGE_OUTPUT" | grep -oP 'coverage: \K[0-9.]+' | awk '{sum+=$1; count++} END {if(count>0) print sum/count; else print 0}')
        COVERAGE_INT=$(printf "%.0f" "$COVERAGE")

        if [ "$COVERAGE_INT" -ge 80 ]; then
            echo -e "Coverage: ${GREEN}${COVERAGE}% ✅${NC}"
        elif [ "$COVERAGE_INT" -ge 70 ]; then
            echo -e "Coverage: ${YELLOW}${COVERAGE}% ⚠️${NC} (Target: 80%)"
            check_warn "Coverage below 80%"
        else
            echo -e "Coverage: ${RED}${COVERAGE}% ❌${NC} (Target: 80%)"
            check_fail "Coverage significantly below 80%"
        fi
    else
        check_warn "Could not calculate coverage"
    fi
else
    check_warn "Ginkgo not installed, skipping coverage check"
fi

echo ""

# ============================================
# 3. STAGED CHANGES (Commit Compliance)
# ============================================
echo "📝 STAGED CHANGES"
echo "------------------------------------------------"

if git diff --cached --quiet; then
    echo -e "${BLUE}No staged changes${NC}"
else
    FILE_COUNT=$(git diff --cached --name-only | wc -l | tr -d ' ')
    LINE_COUNT=$(git diff --cached --numstat | awk '{add+=$1; del+=$2} END {print add+del}')

    echo "Files staged: $FILE_COUNT"
    echo "Lines changed: $LINE_COUNT"

    # Check file count
    if [ "$FILE_COUNT" -gt 10 ]; then
        check_warn "More than 10 files staged (consider atomic commits)"
    fi

    # Check line count
    if [ "$LINE_COUNT" -gt 500 ]; then
        check_warn "More than 500 lines changed (consider splitting)"
    fi

    # Check for generated files
    GENERATED=$(git diff --cached --name-only | grep -E '\.(out|exe|dll|so|test)$' || true)
    if [ -n "$GENERATED" ]; then
        check_fail "Generated files in staging:"
        echo "$GENERATED" | sed 's/^/    /'
    fi

    # Check for coverage files
    COVERAGE_FILES=$(git diff --cached --name-only | grep -E 'coverage\.(out|html)' || true)
    if [ -n "$COVERAGE_FILES" ]; then
        check_fail "Coverage files in staging (should be in .gitignore):"
        echo "$COVERAGE_FILES" | sed 's/^/    /'
    fi

    # Check for debug statements
    DEBUG_PATTERNS='TODO|FIXME|XXX|fmt\.Println\("DEBUG'
    if git diff --cached | grep -E "$DEBUG_PATTERNS" > /dev/null 2>&1; then
        check_warn "Debug statements found in staged changes"
    fi
fi

echo ""

# ============================================
# 4. ARCHITECTURAL COMPLIANCE (DDD)
# ============================================
echo "🏛️  ARCHITECTURAL COMPLIANCE"
echo "------------------------------------------------"

# Check domain layer purity (no external deps)
echo -n "Domain Layer Purity: "
DOMAIN_IMPORTS=$(find internal/domain -name "*.go" -type f 2>/dev/null | xargs grep -h "^import" 2>/dev/null | grep -E "(service|repository|cli)" || true)
if [ -z "$DOMAIN_IMPORTS" ]; then
    check_pass
else
    check_fail "Domain layer has dependencies on other layers"
fi

# Check proper layer structure exists
echo -n "Layer Structure: "
REQUIRED_LAYERS=("internal/domain" "internal/service" "internal/repository")
MISSING_LAYERS=()
for layer in "${REQUIRED_LAYERS[@]}"; do
    if [ ! -d "$layer" ]; then
        MISSING_LAYERS+=("$layer")
    fi
done

if [ ${#MISSING_LAYERS[@]} -eq 0 ]; then
    check_pass
else
    check_warn "Missing layers: ${MISSING_LAYERS[*]}"
fi

echo ""

# ============================================
# 5. DOCUMENTATION
# ============================================
echo "📚 DOCUMENTATION"
echo "------------------------------------------------"

# README
echo -n "README.md: "
if [ -f "README.md" ]; then
    LINES=$(wc -l < README.md)
    if [ "$LINES" -gt 10 ]; then
        check_pass
    else
        check_warn "README exists but is sparse (${LINES} lines)"
    fi
else
    check_fail "Missing README.md"
fi

# AGENTS.md
echo -n "AGENTS.md: "
if [ -f "AGENTS.md" ]; then
    check_pass
else
    check_warn "Consider creating AGENTS.md for handover documentation"
fi

# .gitignore
echo -n ".gitignore: "
if [ -f ".gitignore" ]; then
    if grep -q "coverage.out" .gitignore && grep -q "*.exe" .gitignore; then
        check_pass
    else
        check_warn "gitignore missing common patterns"
    fi
else
    check_fail "Missing .gitignore"
fi

echo ""

# ============================================
# 6. TESTING STANDARDS (Ginkgo/Gomega)
# ============================================
echo "🧪 TESTING STANDARDS"
echo "------------------------------------------------"

# Check for test files
TEST_COUNT=$(find . -name "*_test.go" -type f | wc -l)
GO_FILES=$(find . -name "*.go" -not -name "*_test.go" -not -path "./vendor/*" -type f | wc -l)

echo -n "Test Files: "
if [ "$TEST_COUNT" -gt 0 ]; then
    RATIO=$(awk "BEGIN {printf \"%.1f\", $TEST_COUNT/$GO_FILES}")
    echo -e "${GREEN}${TEST_COUNT} test files${NC} (ratio: ${RATIO}:1)"
else
    check_fail "No test files found"
fi

# Check for Ginkgo suite files
echo -n "Ginkgo Suites: "
SUITE_COUNT=$(find . -name "suite_test.go" -type f | wc -l)
if [ "$SUITE_COUNT" -gt 0 ]; then
    check_pass
else
    check_warn "No Ginkgo suite files found"
fi

echo ""

# ============================================
# 7. DEPENDENCY HEALTH
# ============================================
echo "📦 DEPENDENCY HEALTH"
echo "------------------------------------------------"

# go.mod exists
echo -n "go.mod: "
if [ -f "go.mod" ]; then
    check_pass
else
    check_fail "Missing go.mod"
fi

# go.sum exists
echo -n "go.sum: "
if [ -f "go.sum" ]; then
    check_pass
else
    check_warn "Missing go.sum (run: go mod tidy)"
fi

# Check for tidiness
echo -n "Module Tidiness: "
if go mod tidy -diff > /dev/null 2>&1; then
    check_pass
else
    check_warn "Modules not tidy (run: go mod tidy)"
fi

echo ""

# ============================================
# 8. FILE ORGANIZATION
# ============================================
echo "📁 FILE ORGANIZATION"
echo "------------------------------------------------"

# Check for proper internal/ structure
echo -n "Internal Structure: "
if [ -d "internal" ]; then
    check_pass
else
    check_warn "No internal/ directory (non-standard)"
fi

# Check for cmd/ structure
echo -n "Binary Structure: "
if [ -d "cmd" ]; then
    check_pass
else
    check_warn "No cmd/ directory (consider for binaries)"
fi

# Check for Makefile
echo -n "Makefile: "
if [ -f "Makefile" ]; then
    check_pass
else
    check_warn "No Makefile (consider for task automation)"
fi

echo ""

# ============================================
# 9. GIT HEALTH
# ============================================
echo "🔀 GIT HEALTH"
echo "------------------------------------------------"

# Check if in git repo
echo -n "Git Repository: "
if git rev-parse --git-dir > /dev/null 2>&1; then
    check_pass
else
    check_fail "Not a git repository"
fi

# Check for uncommitted changes (excluding staged)
echo -n "Working Directory: "
if git diff --quiet; then
    echo -e "${GREEN}Clean${NC}"
else
    echo -e "${YELLOW}Has unstaged changes${NC}"
fi

# Check recent commit message quality (if commits exist)
if git log -1 --pretty=%B > /dev/null 2>&1; then
    echo -n "Last Commit Format: "
    LAST_MSG=$(git log -1 --pretty=%B | head -1)
    if echo "$LAST_MSG" | grep -qE '^(feat|fix|docs|style|refactor|test|chore|perf)(\(.+\))?: .{1,50}$'; then
        check_pass
    else
        check_warn "Last commit doesn't follow conventional format"
    fi
fi

echo ""

# ============================================
# SUMMARY
# ============================================
echo "================================================"
echo "SUMMARY"
echo "================================================"

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL CHECKS PASSED${NC}"
    echo ""
    echo "Your project is compliant with all rules!"
    exit 0
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  ${WARNINGS} WARNING(S) FOUND${NC}"
    echo ""
    echo "Project is functional but has minor issues."
    echo "Review warnings above for improvements."
    exit 0
else
    echo -e "${RED}❌ ${VIOLATIONS} VIOLATION(S) FOUND${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  ${WARNINGS} WARNING(S) FOUND${NC}"
    fi
    echo ""
    echo "Fix violations before proceeding:"
    echo "  1. Review failed checks above"
    echo "  2. Apply suggested fixes"
    echo "  3. Re-run: make check-compliance"
    exit 1
fi

