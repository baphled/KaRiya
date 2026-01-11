#!/bin/bash

set -e

# ============================================================================
# Split All E2E Workflow Tests
# ============================================================================
# Processes all remaining *_workflow_test.go files in internal/testutil/e2e/
# ============================================================================

echo "================================================"
echo "🔍 E2E TEST SPLITTING - BATCH MODE"
echo "================================================"
echo ""

# Find all workflow test files
WORKFLOW_FILES=$(ls internal/testutil/e2e/*_workflow_test.go 2>/dev/null || true)

if [ -z "$WORKFLOW_FILES" ]; then
    echo "✅ No workflow test files remaining!"
    echo ""
    echo "All files have been split into:"
    echo "  - E2E tests: internal/testutil/e2e/*_e2e_test.go"
    echo "  - Navigation: internal/cli/intents/*_navigation_test.go"
    exit 0
fi

echo "Found workflow test files to split:"
for file in $WORKFLOW_FILES; do
    basename "$file"
done
echo ""

# Process each file
for filepath in $WORKFLOW_FILES; do
    filename=$(basename "$filepath")
    
    echo "================================================"
    echo "Processing: $filename"
    echo "================================================"
    
    # Run the Python splitter
    python3 scripts/split-e2e-tests.py "$filename"
    
    # Delete the original
    echo "Deleting original file..."
    rm "$filepath"
    
    # Stage changes
    git add -A
    
    # Extract base name for commit message
    base_name=$(echo "$filename" | sed 's/_workflow_test\.go$//')
    
    # Commit
    echo ""
    echo "Creating commit..."
    make ai-commit MSG="refactor(tests): split ${base_name}_workflow into e2e and navigation tests"
    
    echo ""
    echo "✅ $filename complete!"
    echo ""
done

echo "================================================"
echo "✅ ALL WORKFLOW FILES SPLIT!"
echo "================================================"
echo ""
echo "Summary:"
echo "  E2E tests:    internal/testutil/e2e/*_e2e_test.go"
echo "  Navigation:   internal/cli/intents/*_navigation_test.go"
echo ""
echo "Next steps:"
echo "  1. Run all tests: make test"
echo "  2. Verify splits: git log --oneline | head -10"
echo "  3. Push changes: git push"
echo ""
