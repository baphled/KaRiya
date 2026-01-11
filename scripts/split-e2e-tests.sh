#!/bin/bash

set -e

# ============================================================================
# E2E Test Splitting Script
# ============================================================================
# Splits MIXED workflow test files into:
# 1. E2E tests (SQLite) - stays in internal/testutil/e2e/
# 2. Navigation tests (Memory) - moves to internal/cli/intents/
#
# Usage: ./scripts/split-e2e-tests.sh <workflow_test_file>
# Example: ./scripts/split-e2e-tests.sh browse_workflow_test.go
# ============================================================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

WORKFLOW_FILE="$1"

if [ -z "$WORKFLOW_FILE" ]; then
    echo -e "${RED}❌ ERROR: Workflow test file required${NC}"
    echo ""
    echo "Usage: $0 <workflow_test_file>"
    echo ""
    echo "Example: $0 browse_workflow_test.go"
    echo ""
    echo "Remaining files to split:"
    ls -1 internal/testutil/e2e/*_workflow_test.go 2>/dev/null | sed 's|internal/testutil/e2e/||'
    exit 1
fi

# Paths
E2E_DIR="internal/testutil/e2e"
INTENTS_DIR="internal/cli/intents"
SOURCE_FILE="$E2E_DIR/$WORKFLOW_FILE"

# Validate source file exists
if [ ! -f "$SOURCE_FILE" ]; then
    echo -e "${RED}❌ ERROR: File not found: $SOURCE_FILE${NC}"
    exit 1
fi

# Extract base name (e.g., "browse" from "browse_workflow_test.go")
BASE_NAME=$(echo "$WORKFLOW_FILE" | sed 's/_workflow_test\.go$//')

# Output files
E2E_OUTPUT="$E2E_DIR/${BASE_NAME}_e2e_test.go"
NAV_OUTPUT="$INTENTS_DIR/${BASE_NAME}_navigation_test.go"

echo ""
echo -e "${BLUE}🔍 Splitting: $WORKFLOW_FILE${NC}"
echo "  E2E tests   → ${BASE_NAME}_e2e_test.go"
echo "  Navigation  → ${BASE_NAME}_navigation_test.go"
echo ""

# ============================================================================
# Step 1: Extract E2E tests (use Setup(GinkgoT()))
# ============================================================================

echo -e "${BLUE}📝 Extracting E2E tests...${NC}"

# Create E2E file with proper header
cat > "$E2E_OUTPUT" << 'EOF'
package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

EOF

# Extract E2E Describe blocks
# Strategy: Find Describe blocks that contain Setup(GinkgoT()) calls
awk '
BEGIN { in_e2e_block = 0; block_depth = 0; buffer = ""; found_setup = 0 }

# Start of Describe block
/^[[:space:]]*Describe\(/ {
    if (block_depth == 0) {
        buffer = $0 "\n"
        block_depth = 1
        found_setup = 0
    } else {
        buffer = buffer $0 "\n"
        block_depth++
    }
    next
}

# Check for Setup(GinkgoT()) - marks as E2E block
/e2e\.Setup\(GinkgoT\(\)\)/ {
    if (block_depth > 0) {
        found_setup = 1
    }
}

# Track block depth
/\{/ {
    if (block_depth > 0) {
        buffer = buffer $0 "\n"
        block_depth++
        next
    }
}

/\}/ {
    if (block_depth > 0) {
        block_depth--
        buffer = buffer $0 "\n"
        
        # End of top-level Describe block
        if (block_depth == 0) {
            if (found_setup) {
                print buffer
            }
            buffer = ""
            found_setup = 0
        }
        next
    }
}

# Accumulate lines within blocks
{
    if (block_depth > 0) {
        buffer = buffer $0 "\n"
    }
}
' "$SOURCE_FILE" >> "$E2E_OUTPUT"

# Add closing brace if needed
echo "})
" >> "$E2E_OUTPUT"

echo -e "${GREEN}✅ E2E tests extracted: $E2E_OUTPUT${NC}"

# ============================================================================
# Step 2: Extract Navigation tests (use SetupWithMemory(GinkgoT()))
# ============================================================================

echo -e "${BLUE}📝 Extracting Navigation tests...${NC}"

# Get the intent name from base_name (e.g., "browse" → "BrowseTimeline")
INTENT_NAME=$(echo "$BASE_NAME" | sed 's/_/ /g' | awk '{for(i=1;i<=NF;i++){$i=toupper(substr($i,1,1)) substr($i,2)}}1' | sed 's/ //g')

# Create Navigation file with proper header
cat > "$NAV_OUTPUT" << EOF
package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("${INTENT_NAME} Navigation", func() {
	var env *e2e.TestEnv

EOF

# Extract Navigation Describe blocks
awk '
BEGIN { in_nav_block = 0; block_depth = 0; buffer = ""; found_setup_memory = 0; skip_outer_describe = 1 }

# Skip the outer "var _ = Describe" line
/^var _ = Describe\(/ {
    skip_outer_describe = 1
    next
}

# Start of inner Describe block (skip outer)
/^[[:space:]]*Describe\(/ {
    if (block_depth == 0 && !skip_outer_describe) {
        buffer = $0 "\n"
        block_depth = 1
        found_setup_memory = 0
    } else if (block_depth == 0) {
        skip_outer_describe = 0
    } else {
        buffer = buffer $0 "\n"
        block_depth++
    }
    next
}

# Check for SetupWithMemory - marks as Navigation block
/e2e\.SetupWithMemory\(GinkgoT\(\)\)/ {
    if (block_depth > 0) {
        found_setup_memory = 1
    }
}

# Track block depth
/\{/ {
    if (block_depth > 0) {
        buffer = buffer $0 "\n"
        block_depth++
        next
    }
}

/\}/ {
    if (block_depth > 0) {
        block_depth--
        buffer = buffer $0 "\n"
        
        # End of top-level Describe block
        if (block_depth == 0) {
            if (found_setup_memory) {
                print buffer
            }
            buffer = ""
            found_setup_memory = 0
        }
        next
    }
}

# Accumulate lines within blocks
{
    if (block_depth > 0) {
        buffer = buffer $0 "\n"
    }
}
' "$SOURCE_FILE" >> "$NAV_OUTPUT"

# Add closing braces
echo "})
" >> "$NAV_OUTPUT"

echo -e "${GREEN}✅ Navigation tests extracted: $NAV_OUTPUT${NC}"

# ============================================================================
# Step 3: Verify files were created
# ============================================================================

echo ""
echo -e "${BLUE}🔍 Verifying output files...${NC}"

if [ ! -s "$E2E_OUTPUT" ]; then
    echo -e "${RED}❌ ERROR: E2E output file is empty${NC}"
    exit 1
fi

if [ ! -s "$NAV_OUTPUT" ]; then
    echo -e "${RED}❌ ERROR: Navigation output file is empty${NC}"
    exit 1
fi

E2E_LINES=$(wc -l < "$E2E_OUTPUT")
NAV_LINES=$(wc -l < "$NAV_OUTPUT")

echo -e "${GREEN}✅ E2E file: $E2E_LINES lines${NC}"
echo -e "${GREEN}✅ Navigation file: $NAV_LINES lines${NC}"

# ============================================================================
# Step 4: Summary
# ============================================================================

echo ""
echo -e "${GREEN}✅ Split complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Review the generated files:"
echo "     - $E2E_OUTPUT"
echo "     - $NAV_OUTPUT"
echo ""
echo "  2. Run tests to verify:"
echo "     ginkgo internal/testutil/e2e/${BASE_NAME}_e2e_test.go"
echo "     ginkgo internal/cli/intents/${BASE_NAME}_navigation_test.go"
echo ""
echo "  3. Delete original file:"
echo "     rm $SOURCE_FILE"
echo ""
echo "  4. Commit changes:"
echo "     git add -A"
echo "     make ai-commit MSG=\"refactor(tests): split ${BASE_NAME}_workflow into e2e and navigation tests\""
echo ""
