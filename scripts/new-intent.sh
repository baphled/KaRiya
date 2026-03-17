#!/bin/bash
# scripts/new-intent.sh
# 
# Generator script for creating new intent subdirectory structure.
# Usage: make new-intent NAME=my_feature
#
# This creates:
#   intents/{feature}/
#     - context.go
#     - result.go
#     - constants.go
#     - messages.go
#     - intent.go
#   screens/{feature}/
#     - list_screen.go
#     - modals/.gitkeep

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get the name parameter
NAME=$1
if [[ -z "$NAME" ]]; then
    echo -e "${RED}Error: NAME parameter required${NC}" >&2
    echo ""
    echo "Usage: make new-intent NAME=feature_name"
    echo ""
    echo "Examples:"
    echo "  make new-intent NAME=skill_management"
    echo "  make new-intent NAME=event_capture"
    exit 1
fi

# Convert to various naming conventions
PACKAGE_NAME=$(echo "$NAME" | tr '[:upper:]' '[:lower:]' | tr '-' '_')
# Convert snake_case to PascalCase for type names
FEATURE_NAME=$(echo "$PACKAGE_NAME" | sed -r 's/(^|_)([a-z])/\U\2/g')
INTENT_NAME="Intent"
CONTEXT_NAME="Context"
RESULT_NAME="Result"
STATE_NAME="State"

# Paths
INTENT_DIR="internal/tui/intents/${PACKAGE_NAME}"
SCREEN_DIR="internal/tui/views/${PACKAGE_NAME}"
TEMPLATE_DIR="examples/intent_subdirectory_template"

echo -e "${BLUE}Creating intent structure for: ${GREEN}${NAME}${NC}"
echo ""

# Check if directories already exist
if [[ -d "$INTENT_DIR" ]]; then
    echo -e "${YELLOW}Warning: Intent directory already exists: ${INTENT_DIR}${NC}"
    read -p "Overwrite? (y/N): " confirm
    if [[ "$confirm" != "y" ]] && [ "$confirm" != "Y" ]]; then
        echo "Aborted."
        exit 1
    fi
fi

# Create directories
echo -e "${BLUE}Creating directories...${NC}"
mkdir -p "$INTENT_DIR"
mkdir -p "$SCREEN_DIR"
mkdir -p "${SCREEN_DIR}/modals"

# Function to process template
process_template() {
    local template="$1"
    local output="$2"
    
    if [[ ! -f "$template" ]]; then
        echo -e "${YELLOW}Warning: Template not found: ${template}${NC}"
        return
    fi
    
    sed -e "s/{feature}/${PACKAGE_NAME}/g" \
        -e "s/{Feature}/${FEATURE_NAME}/g" \
        -e "s/{FEATURE}/${FEATURE_NAME^^}/g" \
        -e "s/{Entity}/${FEATURE_NAME}/g" \
        "$template" > "$output"
    
    echo -e "  ${GREEN}Created${NC}: $output"
}

# Generate intent files
echo -e "${BLUE}Generating intent files...${NC}"
process_template "${TEMPLATE_DIR}/context.go.template" "${INTENT_DIR}/context.go"
process_template "${TEMPLATE_DIR}/result.go.template" "${INTENT_DIR}/result.go"
process_template "${TEMPLATE_DIR}/constants.go.template" "${INTENT_DIR}/constants.go"
process_template "${TEMPLATE_DIR}/messages.go.template" "${INTENT_DIR}/messages.go"
process_template "${TEMPLATE_DIR}/intent.go.template" "${INTENT_DIR}/intent.go"

# Generate screen files
echo -e "${BLUE}Generating screen files...${NC}"
if [[ -f "${TEMPLATE_DIR}/screens/list_screen.go.template" ]]; then
    process_template "${TEMPLATE_DIR}/screens/list_screen.go.template" "${SCREEN_DIR}/list_screen.go"
fi

# Create .gitkeep for modals
touch "${SCREEN_DIR}/modals/.gitkeep"
echo -e "  ${GREEN}Created${NC}: ${SCREEN_DIR}/modals/.gitkeep"

# Create a basic test file
cat > "${INTENT_DIR}/intent_test.go" << EOF
package ${PACKAGE_NAME}_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func Test${FEATURE_NAME}Intent(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "${FEATURE_NAME} Intent Suite")
}

var _ = Describe("${FEATURE_NAME}Intent", func() {
    Describe("NewIntent", func() {
        It("creates a new intent with valid context", func() {
            // TODO: Implement test
            Skip("Not implemented")
        })
    })

    Describe("Init", func() {
        It("initializes the intent and loads data", func() {
            // TODO: Implement test
            Skip("Not implemented")
        })
    })

    Describe("Update", func() {
        Context("when receiving key messages", func() {
            It("handles quit key", func() {
                // TODO: Implement test
                Skip("Not implemented")
            })
        })
    })
})
EOF
echo -e "  ${GREEN}Created${NC}: ${INTENT_DIR}/intent_test.go"

echo ""
echo -e "${GREEN}Successfully created intent structure!${NC}"
echo ""
echo "Files created:"
echo -e "  ${BLUE}intents/${PACKAGE_NAME}/${NC}"
echo "    - context.go     (business logic, input parameters)"
echo "    - result.go      (output type)"
echo "    - constants.go   (state enum)"
echo "    - messages.go    (message types)"
echo "    - intent.go      (broker implementation)"
echo "    - intent_test.go (test suite)"
echo ""
echo -e "  ${BLUE}screens/${PACKAGE_NAME}/${NC}"
echo "    - list_screen.go (list view)"
echo "    - modals/        (feature-specific modals)"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Edit context.go to add input parameters and business logic"
echo "  2. Edit constants.go to define your state machine states"
echo "  3. Edit messages.go to define your message types"
echo "  4. Edit intent.go to implement your workflow"
echo "  5. Add screens in screens/${PACKAGE_NAME}/ as needed"
echo "  6. Run: make check-intent-architecture"
echo ""
echo -e "${BLUE}Reference implementation:${NC} intents/browse_timeline/"
