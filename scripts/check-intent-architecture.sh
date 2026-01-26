#!/bin/bash

# Intent Architecture Enforcement
# This script enforces strict architectural rules for intent implementations.
# It catches violations that would otherwise slip through code review.
#
# POLICY:
# - NEW intents: violations BLOCK commits
# - EXISTING (legacy) intents: violations become WARNINGS (non-blocking)
#
# This allows incremental migration while enforcing rules on new code.
#
# REFERENCE IMPLEMENTATIONS: browse_timeline (PR #115), burst_management (PR #117)
#
# ============================================================================
# COMPLETE FILE STRUCTURE FOR NEW INTENTS
# ============================================================================
#
# intents/{feature_name}/
# ├── constants.go              # REQUIRED: State enum type and constants
# ├── constants_test.go         # REQUIRED: Tests for state constants
# ├── context.go                # REQUIRED: IntentContext struct + Validate()
# ├── context_test.go           # REQUIRED: Tests for context
# ├── result.go                 # REQUIRED: Result struct
# ├── result_test.go            # REQUIRED: Tests for result
# ├── messages.go               # REQUIRED: ALL *Msg types
# ├── messages_test.go          # REQUIRED: Tests for messages
# ├── intent.go                 # REQUIRED: NewIntent, Init, Update, View, Result
# ├── intent_test.go            # REQUIRED: Tests for intent
# ├── types.go                  # RECOMMENDED: Intent struct definition
# ├── types_test.go             # Tests for types (if types.go exists)
# ├── handlers.go               # RECOMMENDED: ScreenResultHandler methods
# ├── helpers.go                # RECOMMENDED: Helper methods
# ├── helpers_test.go           # Tests for helpers (if helpers.go exists)
# ├── interfaces.go             # RECOMMENDED: Service interfaces for DI
# ├── filters.go                # OPTIONAL: Domain-specific filter logic
# ├── {feature}_suite_test.go   # REQUIRED: Ginkgo test suite entry point
# └── {feature}_e2e_test.go     # RECOMMENDED: End-to-end tests
#
# screens/{feature_name}/
# ├── {name}_list_screen.go         # List view screen
# ├── {name}_list_screen_test.go    # Tests
# ├── {name}_detail_screen.go       # Detail view screen (if needed)
# ├── {name}_detail_screen_test.go  # Tests
# ├── {feature}_suite_test.go       # REQUIRED: Ginkgo test suite
# └── modals/                       # REQUIRED if intent uses modals
#     ├── modals_suite_test.go      # REQUIRED: Ginkgo test suite for modals
#     ├── {modal}_modal.go          # Each modal
#     ├── {modal}_modal_test.go     # Tests for each modal
#     ├── helpers.go                # OPTIONAL: Shared modal helpers
#     └── helpers_test.go           # Tests for helpers
#
# ============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

VIOLATIONS=0
WARNINGS=0
LEGACY_WARNINGS=0

# ============================================
# LEGACY INTENT DETECTION
# ============================================
# Intents that exist in the base branch are "legacy" - violations become warnings.
# New intents must comply with all rules (violations block commits).

BASE_BRANCH="${BASE_BRANCH:-origin/next}"

# Build list of legacy intent files from base branch
LEGACY_INTENTS=""
if git rev-parse --verify "$BASE_BRANCH" >/dev/null 2>&1; then
    # Get flat structure intents
    LEGACY_FLAT=$(git ls-tree --name-only "$BASE_BRANCH" internal/cli/intents/ 2>/dev/null | grep "_intent.go$" || true)
    # Get subdirectory intents
    LEGACY_SUBDIRS=$(git ls-tree -d --name-only "$BASE_BRANCH" internal/cli/intents/ 2>/dev/null | grep -v "types$" || true)
    LEGACY_INTENTS="$LEGACY_FLAT $LEGACY_SUBDIRS"
fi

# Check if a file is from a legacy intent
# Returns 0 (true) if legacy, 1 (false) if new
is_legacy_intent() {
    local file="$1"
    local intent_name=""
    
    # Extract intent name from file path
    if [[ "$file" == *"_intent.go" ]]; then
        # Flat structure: capture_event_intent.go -> capture_event
        intent_name=$(basename "$file" _intent.go)
    elif [[ "$file" == */intents/*/* ]]; then
        # Subdirectory structure: intents/browse_timeline/intent.go -> browse_timeline
        intent_name=$(echo "$file" | sed -n 's|.*/intents/\([^/]*\)/.*|\1|p')
    else
        # Context file: browse_timeline.go -> browse_timeline
        intent_name=$(basename "$file" .go)
    fi
    
    # Check if this intent exists in base branch
    if [ -n "$LEGACY_INTENTS" ]; then
        if echo "$LEGACY_INTENTS" | grep -q "${intent_name}"; then
            return 0  # Legacy
        fi
    fi
    
    return 1  # New
}

# Report a violation or legacy warning based on intent status
# Usage: report_issue "file" "Check Name" "Description" "Rule" "Example"
report_issue() {
    local file="$1"
    local check_name="$2"
    local description="$3"
    local rule="$4"
    local example="$5"
    
    if is_legacy_intent "$file"; then
        echo -e "${CYAN}📋 LEGACY WARNING: $description${NC}"
        echo "   File: $file"
        echo "   Rule: $rule"
        if [ -n "$example" ]; then
            echo ""
            echo "$example" | sed 's/^/   /'
        fi
        echo ""
        echo -e "   ${CYAN}Note: This is a legacy intent. Fix recommended but not required.${NC}"
        echo ""
        LEGACY_WARNINGS=$((LEGACY_WARNINGS+1))
    else
        echo -e "${RED}❌ VIOLATION: $description${NC}"
        echo "   File: $file"
        echo "   Rule: $rule"
        if [ -n "$example" ]; then
            echo ""
            echo "$example" | sed 's/^/   /'
        fi
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
}

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏛️  INTENT ARCHITECTURE ENFORCEMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "${BLUE}Policy: NEW intents → BLOCK violations | LEGACY intents → WARN only${NC}"
echo ""

# ============================================
# 1. TYPED STATE ENUM CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1. TYPED STATE ENUM"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

INTENT_FILES=$(find internal/cli/intents -name '*_intent.go' -not -name '*_test.go' 2>/dev/null || true)

for file in $INTENT_FILES; do
    INTENT_NAME=$(basename "$file" _intent.go)
    
    # Check if intent has state field of type string (untyped)
    UNTYPED_STATE=$(grep -n "currentState.*string\s*$" "$file" 2>/dev/null || true)
    
    if [ -n "$UNTYPED_STATE" ]; then
        report_issue "$file" "Typed State Enum" "Untyped state field" \
            "State fields must use typed enum, not raw string" \
            "Required pattern:
type ${INTENT_NAME^}State string

const (
    State... ${INTENT_NAME^}State = \"...\"
)"
    fi
    
    # Check if intent has state field but missing typed enum definition
    HAS_STATE=$(grep -q "state.*State" "$file" && echo "yes" || echo "no")
    HAS_TYPED_ENUM=$(grep -q "type.*State string" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_STATE" = "yes" ] && [ "$HAS_TYPED_ENUM" = "no" ]; then
        report_issue "$file" "Typed State Enum" "Missing typed state enum definition" \
            "Has state field but no 'type <Name>State string' declaration" ""
    fi
done

echo ""

# ============================================
# 2. FLATTENED STATE MODEL CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2. FLATTENED STATE MODEL"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

WRAPPED_STATE_FILES=$(grep -l "^\s*state\s\+\*.*Model" internal/cli/intents/*_intent.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$WRAPPED_STATE_FILES" ]; then
    for file in $WRAPPED_STATE_FILES; do
        report_issue "$file" "Flattened State" "Wrapped state model found" \
            "Intent state must be flattened directly into the intent struct" \
            "BAD:  state *BrowseTimelineModel
GOOD: Flatten all model fields directly into intent"
    done
else
    echo -e "${GREEN}✅ All intents have flattened state${NC}"
fi

echo ""

# ============================================
# 3. CONTEXT USAGE CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3. CONTEXT USAGE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

BG_CONTEXT_FILES=$(grep -l "context\.Background()" internal/cli/intents/*_intent.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$BG_CONTEXT_FILES" ]; then
    for file in $BG_CONTEXT_FILES; do
        report_issue "$file" "Context Usage" "context.Background() in intent" \
            "Intents must use i.getContext() for cancellation support" \
            "BAD:  ctx := context.Background()
GOOD: ctx := i.getContext()"
    done
else
    echo -e "${GREEN}✅ All intents use proper context${NC}"
fi

echo ""

# ============================================
# 4. DEAD CODE CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4. DEAD CODE MARKERS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

DEAD_CODE_FILES=$(grep -l "should not be reached\|LEGACY.*not.*reached\|This case is kept for backward compatibility but should not be reached" internal/cli/intents/*.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$DEAD_CODE_FILES" ]; then
    for file in $DEAD_CODE_FILES; do
        report_issue "$file" "Dead Code" "Dead code markers found" \
            "Remove unreachable code or create cleanup task" ""
    done
else
    echo -e "${GREEN}✅ No dead code markers${NC}"
fi

echo ""

# ============================================
# 5. BASEINTENT EMBEDDING CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5. BASEINTENT EMBEDDING"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

MISSING_BASE=$(grep -L "\*BaseIntent" internal/cli/intents/*_intent.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$MISSING_BASE" ]; then
    for file in $MISSING_BASE; do
        report_issue "$file" "BaseIntent Embedding" "Intent missing *BaseIntent" \
            "All intents must embed *BaseIntent" \
            "Required pattern:
type MyIntent struct {
    *intents.BaseIntent  // <- REQUIRED
    ...
}"
    done
else
    echo -e "${GREEN}✅ All intents embed *BaseIntent${NC}"
fi

echo ""

# ============================================
# 7. GODOC COMPLETENESS CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6. GODOC COMPLETENESS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Find exported functions (starts with capital letter)
    EXPORTED_FUNCS=$(grep -n "^func.*) [A-Z]" "$file" 2>/dev/null || true)
    
    if [ -n "$EXPORTED_FUNCS" ]; then
        # Check each exported function for godoc
        while IFS= read -r func_line; do
            LINE_NUM=$(echo "$func_line" | cut -d: -f1)
            FUNC_NAME=$(echo "$func_line" | grep -oP 'func \([^)]+\) \K[A-Za-z]+' || true)
            
            # Check if previous line has comment
            PREV_LINE=$((LINE_NUM - 1))
            HAS_COMMENT=$(sed -n "${PREV_LINE}p" "$file" | grep -q "^//" && echo "yes" || echo "no")
            
            if [ "$HAS_COMMENT" = "no" ] && [ -n "$FUNC_NAME" ]; then
                echo -e "${YELLOW}⚠️  WARNING: Missing godoc${NC}"
                echo "   File: $file:$LINE_NUM"
                echo "   Function: $FUNC_NAME"
                echo ""
                WARNINGS=$((WARNINGS+1))
            fi
        done <<< "$EXPORTED_FUNCS"
    fi
done

if [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ Godoc coverage looks good${NC}"
fi

echo ""

# ============================================
# 8. MODAL OVERLAY PATTERN CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "7. MODAL OVERLAY PATTERN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK7_ISSUES=0
for file in $INTENT_FILES; do
    # Check if intent has modals
    HAS_MODALS=$(grep -q "Modal.*\*.*Modal" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_MODALS" = "yes" ]; then
        # Check if using RenderModalOverlay
        USES_OVERLAY=$(grep -q "RenderModalOverlay" "$file" && echo "yes" || echo "no")
        
        if [ "$USES_OVERLAY" = "no" ]; then
            report_issue "$file" "Modal Overlay" "Modal without RenderModalOverlay" \
                "Use behaviors.RenderModalOverlay() for modal rendering" \
                "Required pattern:
if i.modal != nil && i.modal.IsVisible() {
    return behaviors.RenderModalOverlay(i.modal, baseView)
}"
            CHECK7_ISSUES=$((CHECK7_ISSUES+1))
        fi
    fi
done

if [ $CHECK7_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ Modal patterns correct${NC}"
fi

echo ""

# ============================================
# 9. SCREENRESULTHANDLER IMPLEMENTATION CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "8. SCREENRESULTHANDLER IMPLEMENTATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK8_ISSUES=0
for file in $INTENT_FILES; do
    # Check if intent uses screens
    HAS_SCREENS=$(grep -q "screens\.Screen" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_SCREENS" = "yes" ]; then
        # Check for ScreenResultHandler implementation
        HAS_HANDLER=$(grep -q "var _ ScreenResultHandler" "$file" && echo "yes" || echo "no")
        
        if [ "$HAS_HANDLER" = "no" ]; then
            report_issue "$file" "ScreenResultHandler" "Missing ScreenResultHandler" \
                "Intents using screens must implement ScreenResultHandler" \
                "Required:
var _ ScreenResultHandler = (*YourIntent)(nil)

func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd { ... }
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd { ... }
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd { ... }
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd { ... }"
            CHECK8_ISSUES=$((CHECK8_ISSUES+1))
        fi
    fi
done

if [ $CHECK8_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ ScreenResultHandler implementations correct${NC}"
fi

echo ""

# ============================================
# 9. LEAN INTENT PATTERN CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "9. LEAN INTENT PATTERN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK9_ISSUES=0
for file in $INTENT_FILES; do
    # Check for direct database/service calls in Update() method
    # Intent should delegate to screens/services, not contain business logic
    UPDATE_METHOD=$(sed -n '/^func.*Update.*tea\.Msg/,/^func /p' "$file" 2>/dev/null || true)
    
    if [ -n "$UPDATE_METHOD" ]; then
        # Check for SQL queries in Update
        HAS_SQL=$(echo "$UPDATE_METHOD" | grep -E "\.Query\(|\.Exec\(|\.QueryRow\(|INSERT INTO|SELECT.*FROM|UPDATE.*SET|DELETE FROM" || true)
        
        if [ -n "$HAS_SQL" ]; then
            report_issue "$file" "Lean Intent" "Business logic in Update()" \
                "Intents should orchestrate, not contain business logic" \
                "Move business logic to:
- Service layer for data operations
- Screen layer for UI logic"
            CHECK9_ISSUES=$((CHECK9_ISSUES+1))
        fi
        
        # Check for complex computations (heuristic: multiple nested loops or complex math)
        COMPLEX_LOGIC=$(echo "$UPDATE_METHOD" | grep -E "for.*for.*{|\.Calculate|\.Process.*{" | wc -l)
        
        if [ "$COMPLEX_LOGIC" -gt 2 ]; then
            echo -e "${YELLOW}⚠️  WARNING: Complex logic in Update()${NC}"
            echo "   File: $file"
            echo "   Recommendation: Move complex computations to service layer"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    fi
done

if [ $CHECK9_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ Intents follow lean orchestration pattern${NC}"
fi

echo ""

# ============================================
# 10. REQUIRED STATE FIELD CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "10. REQUIRED STATE FIELD"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK10_ISSUES=0
for file in $INTENT_FILES; do
    # Check if intent has a state field
    HAS_STATE=$(grep -q "^\s*state\s\+\w" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_STATE" = "no" ]; then
        report_issue "$file" "Required State" "Missing state field" \
            "All intents must have a state field (typed enum)" \
            "Required:
state MyIntentState  // Typed state enum"
        CHECK10_ISSUES=$((CHECK10_ISSUES+1))
    fi
done

if [ $CHECK10_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ All intents have state fields${NC}"
fi

echo ""

# ============================================
# 11. EXPLICIT SCREEN FIELDS CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "11. EXPLICIT SCREEN FIELDS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK11_ISSUES=0
for file in $INTENT_FILES; do
    # Check if intent uses screens
    HAS_ACTIVE_SCREEN=$(grep -q "activeScreen.*screens\.Screen" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_ACTIVE_SCREEN" = "yes" ]; then
        # Check for explicitly typed screen fields (listScreen, detailScreen, etc.)
        TYPED_SCREENS=$(grep "Screen\s\+\*.*\..*Screen" "$file" 2>/dev/null | wc -l)
        
        if [ "$TYPED_SCREENS" -eq "0" ]; then
            report_issue "$file" "Explicit Screens" "Only generic activeScreen field" \
                "Intents must declare explicit typed screen fields" \
                "Required pattern:
listScreen   *myfeature.ListScreen
detailScreen *myfeature.DetailScreen
activeScreen screens.Screen  // Points to one of the above"
            CHECK11_ISSUES=$((CHECK11_ISSUES+1))
        fi
    fi
done

if [ $CHECK11_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ Screen management patterns correct${NC}"
fi

echo ""

# ============================================
# 12. EXPLICIT MODAL FIELDS CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "12. EXPLICIT MODAL FIELDS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check if intent uses RenderModalOverlay (implies modals)
    USES_MODALS=$(grep -q "RenderModalOverlay\|\.IsVisible()" "$file" && echo "yes" || echo "no")
    
    if [ "$USES_MODALS" = "yes" ]; then
        # Check for explicitly typed modal fields
        TYPED_MODALS=$(grep -E "Modal\s+\*.*Modal" "$file" 2>/dev/null | wc -l)
        
        if [ "$TYPED_MODALS" -eq "0" ]; then
            echo -e "${YELLOW}⚠️  WARNING: Using modals without explicit fields${NC}"
            echo "   File: $file"
            echo "   Recommendation: Declare explicit modal fields in struct"
            echo ""
            echo "   Example:"
            echo "   deleteModal  *feedback.ConfirmModal"
            echo "   successModal *feedback.SuccessModal"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    fi
done

if [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ Modal fields properly declared${NC}"
fi

echo ""

# ============================================
# 13. KEY HANDLING PATTERN CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "13. KEY HANDLING PATTERN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check for string-based key comparisons (should use HandleGlobalKeys)
    STRING_KEY_MATCH=$(grep -n 'keyMsg\.String()\s*==' "$file" 2>/dev/null | grep -v "HandleGlobalKeys\|tea.KeyMsg" || true)
    
    if [ -n "$STRING_KEY_MATCH" ]; then
        echo -e "${YELLOW}⚠️  WARNING: String-based key matching${NC}"
        echo "   File: $file"
        echo "   Recommendation: Use HandleGlobalKeys() for common keys"
        echo ""
        echo "   Found:"
        echo "$STRING_KEY_MATCH" | head -3 | sed 's/^/   /'
        echo ""
        echo "   Instead of: if keyMsg.String() == \"q\""
        echo "   Use: switch HandleGlobalKeys(keyMsg) { case KeyQuit: ... }"
        echo ""
        WARNINGS=$((WARNINGS+1))
    fi
done

echo ""

# ============================================
# 14. CONTEXT FIELD CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "14. CONTEXT FIELD REQUIREMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK14_ISSUES=0
for file in $INTENT_FILES; do
    INTENT_NAME=$(basename "$file" _intent.go | sed 's/_/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1' | sed 's/ //g')
    
    # Check if intent has a context field
    HAS_CONTEXT=$(grep -q "context\s\+\*${INTENT_NAME}Context" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_CONTEXT" = "no" ]; then
        # Check if it's using generic context
        HAS_ANY_CONTEXT=$(grep -q "context\s\+\*.*Context" "$file" && echo "yes" || echo "no")
        
        if [ "$HAS_ANY_CONTEXT" = "no" ]; then
            report_issue "$file" "Context Field" "Missing context field" \
                "Intents should have a context field for input parameters" \
                "Required:
context *${INTENT_NAME}Context"
            CHECK14_ISSUES=$((CHECK14_ISSUES+1))
        fi
    fi
done

if [ $CHECK14_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ All intents have context fields${NC}"
fi

echo ""

# ============================================
# 15. HANDLEGLOBALKEYS USAGE CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "15. HANDLEGLOBALKEYS USAGE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check if intent handles keyboard input
    HAS_KEYMSG=$(grep -q "tea\.KeyMsg" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_KEYMSG" = "yes" ]; then
        # Check if it uses HandleGlobalKeys
        USES_HANDLEGLOBALKEYS=$(grep -q "HandleGlobalKeys" "$file" && echo "yes" || echo "no")
        
        if [ "$USES_HANDLEGLOBALKEYS" = "no" ]; then
            echo -e "${YELLOW}⚠️  WARNING: Not using HandleGlobalKeys${NC}"
            echo "   File: $file"
            echo "   Recommendation: Use HandleGlobalKeys() for q, ?, m keys"
            echo ""
            echo "   Example:"
            echo "   switch HandleGlobalKeys(keyMsg) {"
            echo "       case KeyQuit: return tea.Quit"
            echo "       case KeyHelp: i.helpModal.Toggle()"
            echo "   }"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    fi
done

if [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ All intents use HandleGlobalKeys${NC}"
fi

echo ""

# ============================================
# 16. FILE SEPARATION CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "16. FILE SEPARATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK16_ISSUES=0
for file in $INTENT_FILES; do
    BASENAME=$(basename "$file" .go)
    DIRNAME=$(dirname "$file")
    
    # Check if ANY Context struct is defined in intent file (ending with Context)
    CONTEXT_STRUCTS=$(grep "^type.*Context struct" "$file" | grep -v "// " || true)
    
    if [ -n "$CONTEXT_STRUCTS" ]; then
        report_issue "$file" "File Separation" "Context struct(s) defined in intent file" \
            "Context structs must be in separate file" \
            "Found: $CONTEXT_STRUCTS

Required separation:
- Intent: ${BASENAME}.go (only intent struct and methods)
- Context: ${BASENAME%_intent}.go OR context.go (in subdirectory)"
        CHECK16_ISSUES=$((CHECK16_ISSUES+1))
    fi
    
    # Check if ANY Model struct is defined in intent file (ending with Model)
    MODEL_STRUCTS=$(grep "^type.*Model struct" "$file" | grep -v "// " || true)
    
    if [ -n "$MODEL_STRUCTS" ]; then
        report_issue "$file" "File Separation" "Model struct(s) defined in intent file" \
            "Model structs must be flattened into intent OR in separate file" \
            "Found: $MODEL_STRUCTS

RECOMMENDED: Flatten model fields directly into intent struct
ALTERNATIVE: Create separate file: ${BASENAME%_intent}_model.go"
        CHECK16_ISSUES=$((CHECK16_ISSUES+1))
    fi
    
    # Check if Screen structs are defined in intent file
    SCREEN_IN_INTENT=$(grep -q "^type.*Screen struct" "$file" && echo "yes" || echo "no")
    
    if [ "$SCREEN_IN_INTENT" = "yes" ]; then
        report_issue "$file" "File Separation" "Screen defined in intent file" \
            "Screen structs must be in screens/ package" \
            "Required:
- Screens: internal/cli/screens/myfeature/*.go
- Intent: internal/cli/intents/${BASENAME}.go"
        CHECK16_ISSUES=$((CHECK16_ISSUES+1))
    fi
    
    # Check if Modal structs are defined in intent file
    MODAL_IN_INTENT=$(grep -q "^type.*Modal struct" "$file" && echo "yes" || echo "no")
    
    if [ "$MODAL_IN_INTENT" = "yes" ]; then
        report_issue "$file" "File Separation" "Modal defined in intent file" \
            "Modal structs must be in screens/{feature}/modals/ or uikit/feedback/" \
            "Required:
- Feature modals: internal/cli/screens/myfeature/modals/*.go
- Common modals: internal/cli/uikit/feedback/*.go"
        CHECK16_ISSUES=$((CHECK16_ISSUES+1))
    fi
done

if [ $CHECK16_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ File separation correct${NC}"
fi

echo ""

# ============================================
# 17. SUBDIRECTORY STRUCTURE (Reference: PR #117 burst_management)
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "17. SUBDIRECTORY STRUCTURE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check for subdirectory-based intents
SUBDIRS=$(find internal/cli/intents -mindepth 1 -maxdepth 1 -type d 2>/dev/null | grep -v types || true)

CHECK17_ISSUES=0
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # ===========================================
        # REQUIRED CORE SOURCE FILES (5 files)
        # ===========================================
        REQUIRED_SOURCE=("context.go" "result.go" "constants.go" "messages.go" "intent.go")
        MISSING_SOURCE=""
        
        for req_file in "${REQUIRED_SOURCE[@]}"; do
            if [ ! -f "$intent_dir/$req_file" ]; then
                MISSING_SOURCE="$MISSING_SOURCE $req_file"
            fi
        done
        
        if [ -n "$MISSING_SOURCE" ]; then
            report_issue "$intent_dir" "Subdirectory Structure" "Missing required source files" \
                "Missing:$MISSING_SOURCE" \
                "REQUIRED SOURCE FILES (5 core):
intents/$INTENT_NAME/
├── constants.go    # State enum type and constants
├── context.go      # IntentContext struct + Validate()
├── result.go       # Result struct
├── messages.go     # ALL *Msg types
└── intent.go       # NewIntent, Init, Update, View, Result

Reference: burst_management (PR #117)"
            CHECK17_ISSUES=$((CHECK17_ISSUES+1))
        fi
        
        # ===========================================
        # REQUIRED TEST FILES (must have tests for core files)
        # ===========================================
        REQUIRED_TESTS=("constants_test.go" "context_test.go" "result_test.go" "messages_test.go" "intent_test.go")
        MISSING_TESTS=""
        
        for test_file in "${REQUIRED_TESTS[@]}"; do
            if [ ! -f "$intent_dir/$test_file" ]; then
                MISSING_TESTS="$MISSING_TESTS $test_file"
            fi
        done
        
        if [ -n "$MISSING_TESTS" ]; then
            report_issue "$intent_dir" "Test Coverage" "Missing required test files" \
                "Missing:$MISSING_TESTS" \
                "REQUIRED TEST FILES:
intents/$INTENT_NAME/
├── constants_test.go   # Tests for state constants
├── context_test.go     # Tests for context validation
├── result_test.go      # Tests for result struct
├── messages_test.go    # Tests for message types
└── intent_test.go      # Tests for intent logic

Reference: burst_management (PR #117)"
            CHECK17_ISSUES=$((CHECK17_ISSUES+1))
        fi
        
        # ===========================================
        # REQUIRED GINKGO TEST SUITE
        # ===========================================
        SUITE_FILE=$(find "$intent_dir" -maxdepth 1 -name "*_suite_test.go" 2>/dev/null | head -1)
        if [ -z "$SUITE_FILE" ]; then
            report_issue "$intent_dir" "Test Suite" "Missing Ginkgo test suite" \
                "No *_suite_test.go file found" \
                "REQUIRED: Ginkgo test suite entry point

Create: intents/$INTENT_NAME/${INTENT_NAME}_suite_test.go

Example content:
package ${INTENT_NAME}_test

import (
    \"testing\"
    . \"github.com/onsi/ginkgo/v2\"
    . \"github.com/onsi/gomega\"
)

func TestBrowseTimeline(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, \"${INTENT_NAME} Suite\")
}"
            CHECK17_ISSUES=$((CHECK17_ISSUES+1))
        fi
        
        # ===========================================
        # RECOMMENDED FILES (warn if large intent.go)
        # ===========================================
        RECOMMENDED_FILES=("types.go" "handlers.go" "helpers.go" "interfaces.go")
        MISSING_RECOMMENDED=""
        
        for rec_file in "${RECOMMENDED_FILES[@]}"; do
            if [ ! -f "$intent_dir/$rec_file" ]; then
                MISSING_RECOMMENDED="$MISSING_RECOMMENDED $rec_file"
            fi
        done
        
        # Only warn if missing recommended files AND intent.go is large
        if [ -n "$MISSING_RECOMMENDED" ]; then
            INTENT_FILE="$intent_dir/intent.go"
            if [ -f "$INTENT_FILE" ]; then
                LINE_COUNT=$(wc -l < "$INTENT_FILE")
                if [ $LINE_COUNT -gt 300 ]; then
                    echo -e "${YELLOW}⚠️  WARNING: Consider splitting large intent.go${NC}"
                    echo "   Intent: $INTENT_NAME ($LINE_COUNT lines)"
                    echo "   Missing recommended files:$MISSING_RECOMMENDED"
                    echo ""
                    echo "   RECOMMENDED FILES (for larger intents):"
                    echo "   ├── types.go      # Intent struct definition"
                    echo "   ├── handlers.go   # ScreenResultHandler methods"
                    echo "   ├── helpers.go    # Helper methods"
                    echo "   └── interfaces.go # Service interfaces for DI"
                    echo ""
                    echo "   Target: intent.go under 300 lines (broker only)"
                    echo ""
                    WARNINGS=$((WARNINGS+1))
                fi
            fi
        fi
        
        # ===========================================
        # E2E TEST FILE (recommended)
        # ===========================================
        E2E_FILE=$(find "$intent_dir" -maxdepth 1 -name "*_e2e_test.go" 2>/dev/null | head -1)
        if [ -z "$E2E_FILE" ]; then
            echo -e "${YELLOW}⚠️  WARNING: No E2E test file${NC}"
            echo "   Intent: $INTENT_NAME"
            echo "   Recommendation: Add ${INTENT_NAME}_e2e_test.go for integration tests"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    done
fi

# Check for OLD flat structure - these are LEGACY warnings (not blocking)
OLD_INTENTS=$(find internal/cli/intents -maxdepth 1 -name "*_intent.go" -type f 2>/dev/null || true)
OLD_COUNT=$(echo "$OLD_INTENTS" | grep -c "_intent.go" 2>/dev/null || echo 0)

if [ $OLD_COUNT -gt 0 ]; then
    echo -e "${CYAN}📋 LEGACY: $OLD_COUNT intents using old flat structure${NC}"
    echo "   These intents should be migrated to subdirectory structure:"
    echo ""
    echo "$OLD_INTENTS" | while read intent_file; do
        if [ -n "$intent_file" ]; then
            echo "   - $(basename "$intent_file" _intent.go)"
        fi
    done
    echo ""
    echo "   Migration guide: docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md"
    echo ""
    LEGACY_WARNINGS=$((LEGACY_WARNINGS+OLD_COUNT))
fi

if [ $CHECK17_ISSUES -eq 0 ] && [ $OLD_COUNT -eq 0 ]; then
    echo -e "${GREEN}✅ All intents use subdirectory structure${NC}"
elif [ $CHECK17_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ Subdirectory intents are complete${NC}"
fi

echo ""

# ============================================
# 18. INTENT FILE SIZE LIMIT
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "18. INTENT FILE SIZE LIMIT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK18_ISSUES=0

# Check subdirectory-based intents
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_FILE="$intent_dir/intent.go"
        if [ -f "$INTENT_FILE" ]; then
            LINE_COUNT=$(wc -l < "$INTENT_FILE")
            INTENT_NAME=$(basename "$intent_dir")
            
            if [ $LINE_COUNT -gt 600 ]; then
                report_issue "$INTENT_FILE" "File Size" "intent.go exceeds 600 lines ($LINE_COUNT lines)" \
                    "Intent should be broker only (orchestration)" \
                    "Required actions:
- Extract views to screens/$INTENT_NAME/
- Move business logic to context.go
- Extract helpers to appropriate packages

Target: 200-400 lines
Maximum: 600 lines (hard limit)"
                CHECK18_ISSUES=$((CHECK18_ISSUES+1))
            elif [ $LINE_COUNT -gt 400 ]; then
                echo -e "${YELLOW}⚠️  WARNING: intent.go approaching limit${NC}"
                echo "   File: $INTENT_FILE ($LINE_COUNT lines)"
                echo "   Target: 200-400 lines (broker only)"
                echo "   Consider: Extract views to screens/$INTENT_NAME/"
                echo ""
                WARNINGS=$((WARNINGS+1))
            fi
        fi
    done
fi

# Also check old flat structure files - LEGACY warnings
if [ -n "$OLD_INTENTS" ]; then
    for intent_file in $OLD_INTENTS; do
        if [ -f "$intent_file" ] && [ -n "$intent_file" ]; then
            LINE_COUNT=$(wc -l < "$intent_file")
            INTENT_NAME=$(basename "$intent_file" _intent.go)
            
            if [ $LINE_COUNT -gt 1000 ]; then
                echo -e "${CYAN}📋 LEGACY: Intent file exceeds 1,000 lines${NC}"
                echo "   File: $intent_file ($LINE_COUNT lines)"
                echo "   Action: Migrate to subdirectory structure"
                echo ""
                LEGACY_WARNINGS=$((LEGACY_WARNINGS+1))
            fi
        fi
    done
fi

if [ $CHECK18_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ All intent files within size limits${NC}"
fi

echo ""

# ============================================
# 19. NO RENDERING IN INTENT
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "19. NO RENDERING IN INTENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK19_ISSUES=0
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_FILE="$intent_dir/intent.go"
        if [ -f "$INTENT_FILE" ]; then
            INTENT_NAME=$(basename "$intent_dir")
            
            # Count render methods (allow View() + 1 helper max)
            RENDER_COUNT=$(grep -c "^func.*render\|^func.*Render" "$INTENT_FILE" 2>/dev/null || echo 0)
            VIEW_COUNT=$(grep -c "^func.*View() string" "$INTENT_FILE" 2>/dev/null || echo 0)
            TOTAL_RENDER=$((RENDER_COUNT + VIEW_COUNT))
            
            # Allow View() + 1 helper = 2 methods max
            if [ $TOTAL_RENDER -gt 2 ]; then
                FOUND_METHODS=$(grep -n "^func.*render\|^func.*Render\|^func.*View" "$INTENT_FILE" 2>/dev/null | head -5 || true)
                report_issue "$INTENT_FILE" "No Rendering" "Rendering methods in intent.go ($TOTAL_RENDER render methods)" \
                    "Intent should only have View() that delegates to screens" \
                    "Found methods:
$FOUND_METHODS

Required action:
- Extract rendering to screens/$INTENT_NAME/
- Keep only View() in intent.go
- View() should delegate: return i.activeScreen.View()

See: docs/guides/SCREEN_EXTRACTION_GUIDE.md"
                CHECK19_ISSUES=$((CHECK19_ISSUES+1))
            fi
        fi
    done
fi

if [ $CHECK19_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ No rendering logic in intent files${NC}"
fi

echo ""

# ============================================
# 20. TYPE LOCATION ENFORCEMENT
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "20. TYPE LOCATION ENFORCEMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK20_ISSUES=0
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # Check Context location - must be in context.go, NOT in intent.go or types.go
        if [ -f "$intent_dir/intent.go" ]; then
            CONTEXT_IN_INTENT=$(grep "^type.*Context struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$CONTEXT_IN_INTENT" ]; then
                report_issue "$intent_dir/intent.go" "Type Location" "Context defined in intent.go" \
                    "Context must be in context.go" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            CONTEXT_IN_TYPES=$(grep "^type.*Context struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$CONTEXT_IN_TYPES" ]; then
                report_issue "$intent_dir/types.go" "Type Location" "Context defined in types.go" \
                    "Context must be in context.go" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        # Check Result location - must be in result.go
        if [ -f "$intent_dir/intent.go" ]; then
            RESULT_IN_INTENT=$(grep "^type.*Result struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$RESULT_IN_INTENT" ]; then
                report_issue "$intent_dir/intent.go" "Type Location" "Result defined in intent.go" \
                    "Result must be in result.go" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            RESULT_IN_TYPES=$(grep "^type.*Result struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$RESULT_IN_TYPES" ]; then
                report_issue "$intent_dir/types.go" "Type Location" "Result defined in types.go" \
                    "Result must be in result.go" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        # Check State enum location - must be in constants.go
        if [ -f "$intent_dir/intent.go" ]; then
            STATE_IN_INTENT=$(grep "^type.*State string" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$STATE_IN_INTENT" ]; then
                report_issue "$intent_dir/intent.go" "Type Location" "State enum defined in intent.go" \
                    "State enum must be in constants.go" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            STATE_IN_TYPES=$(grep "^type.*State string" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$STATE_IN_TYPES" ]; then
                report_issue "$intent_dir/types.go" "Type Location" "State enum defined in types.go" \
                    "State enum must be in constants.go" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        # Check Msg types location (CRITICAL - must be in messages.go)
        if [ -f "$intent_dir/intent.go" ]; then
            MSG_IN_INTENT=$(grep "^type.*Msg struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$MSG_IN_INTENT" ]; then
                report_issue "$intent_dir/intent.go" "Type Location" "Msg types defined in intent.go" \
                    "ALL *Msg types must be in messages.go" \
                    "Found: $MSG_IN_INTENT"
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            MSG_IN_TYPES=$(grep "^type.*Msg struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$MSG_IN_TYPES" ]; then
                report_issue "$intent_dir/types.go" "Type Location" "Msg types defined in types.go" \
                    "ALL *Msg types must be in messages.go" \
                    "Found: $MSG_IN_TYPES"
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        # Check if Msg types are in constants.go (WRONG - should be messages.go)
        if [ -f "$intent_dir/constants.go" ]; then
            MSG_IN_CONSTANTS=$(grep "^type.*Msg struct" "$intent_dir/constants.go" 2>/dev/null || true)
            if [ -n "$MSG_IN_CONSTANTS" ]; then
                report_issue "$intent_dir/constants.go" "Type Location" "Msg types in constants.go" \
                    "Msg types must be in messages.go (not constants.go)" ""
                CHECK20_ISSUES=$((CHECK20_ISSUES+1))
            fi
        fi
        
        # Check Intent struct location - allowed in intent.go OR types.go
        HAS_INTENT_STRUCT=false
        if [ -f "$intent_dir/intent.go" ]; then
            INTENT_STRUCT_IN_INTENT=$(grep "^type.*Intent struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$INTENT_STRUCT_IN_INTENT" ]; then
                HAS_INTENT_STRUCT=true
            fi
        fi
        if [ -f "$intent_dir/types.go" ]; then
            INTENT_STRUCT_IN_TYPES=$(grep "^type.*Intent struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$INTENT_STRUCT_IN_TYPES" ]; then
                HAS_INTENT_STRUCT=true
            fi
        fi
        
        if [ "$HAS_INTENT_STRUCT" = false ]; then
            report_issue "$intent_dir" "Type Location" "No Intent struct found" \
                "Intent struct must be in intent.go OR types.go" ""
            CHECK20_ISSUES=$((CHECK20_ISSUES+1))
        fi
    done
fi

if [ $CHECK20_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ All types in correct locations${NC}"
fi

echo ""

# ============================================
# 21. UIKIT COMPONENT USAGE
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "21. UIKIT COMPONENT USAGE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK21_ISSUES=0

# Check for deprecated components usage
for file in $INTENT_FILES; do
    # Check for components.KeyBadge (should use primitives.HelpKeyBadge)
    KEYBADGE_USAGE=$(grep "components\.KeyBadge" "$file" 2>/dev/null || true)
    if [ -n "$KEYBADGE_USAGE" ]; then
        report_issue "$file" "UIKit Usage" "Deprecated components.KeyBadge" \
            "Use primitives.HelpKeyBadge() instead" \
            "Migration:
OLD: components.KeyBadge(\"key\", \"label\")
NEW: primitives.HelpKeyBadge(\"key\", \"label\", theme)"
        CHECK21_ISSUES=$((CHECK21_ISSUES+1))
    fi
    
    # Check for components.StandardView (should use layout.NewScreenLayout)
    STANDARDVIEW_USAGE=$(grep "components\.StandardView" "$file" 2>/dev/null || true)
    if [ -n "$STANDARDVIEW_USAGE" ]; then
        report_issue "$file" "UIKit Usage" "Deprecated components.StandardView" \
            "Use layout.NewScreenLayout() instead" \
            "Migration:
OLD: components.StandardView{...}
NEW: layout.NewScreenLayout(theme).WithContent(...).Render()"
        CHECK21_ISSUES=$((CHECK21_ISSUES+1))
    fi
    
    # Check for raw lipgloss.NewStyle() usage (should use theme or UIKit)
    RAW_LIPGLOSS=$(grep -E "lipgloss\.NewStyle\(\)|lipgloss\.Color\(" "$file" 2>/dev/null || true)
    if [ -n "$RAW_LIPGLOSS" ]; then
        LIPGLOSS_COUNT=$(echo "$RAW_LIPGLOSS" | wc -l)
        if [ $LIPGLOSS_COUNT -gt 5 ]; then
            echo -e "${YELLOW}⚠️  WARNING: Excessive raw lipgloss usage${NC}"
            echo "   File: $file ($LIPGLOSS_COUNT instances)"
            echo "   Recommendation: Use theme system or UIKit components"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    fi
    
    # Check for correct modal usage (feedback package for common modals)
    CUSTOM_MODAL=$(grep "type.*Modal struct" "$file" 2>/dev/null || true)
    if [ -n "$CUSTOM_MODAL" ]; then
        # Check if it's a common pattern (confirm, error, loading, etc.)
        if echo "$CUSTOM_MODAL" | grep -qE "Confirm|Delete|Error|Loading|Success|Warning"; then
            echo -e "${YELLOW}⚠️  WARNING: Custom modal for common pattern${NC}"
            echo "   File: $file"
            echo "   Found: $CUSTOM_MODAL"
            echo "   Recommendation: Use centralized modals from uikit/feedback/"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    fi
done

# Check subdirectory-based intents for UIKit usage
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_FILE="$intent_dir/intent.go"
        if [ -f "$INTENT_FILE" ]; then
            # Same checks for subdirectory intents
            KEYBADGE_USAGE=$(grep "components\.KeyBadge" "$INTENT_FILE" 2>/dev/null || true)
            if [ -n "$KEYBADGE_USAGE" ]; then
                report_issue "$INTENT_FILE" "UIKit Usage" "Deprecated components.KeyBadge" \
                    "Use primitives.HelpKeyBadge() instead" ""
                CHECK21_ISSUES=$((CHECK21_ISSUES+1))
            fi
            
            STANDARDVIEW_USAGE=$(grep "components\.StandardView" "$INTENT_FILE" 2>/dev/null || true)
            if [ -n "$STANDARDVIEW_USAGE" ]; then
                report_issue "$INTENT_FILE" "UIKit Usage" "Deprecated components.StandardView" \
                    "Use layout.NewScreenLayout() instead" ""
                CHECK21_ISSUES=$((CHECK21_ISSUES+1))
            fi
        fi
    done
fi

if [ $CHECK21_ISSUES -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ UIKit components used correctly${NC}"
elif [ $CHECK21_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ No UIKit violations (warnings present)${NC}"
fi

echo ""

# ============================================
# 22. DEPRECATED MODELS PACKAGE FOR FORMS
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "22. DEPRECATED MODELS PACKAGE FOR FORMS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK22_ISSUES=0

# Check for models.*Form usage in intents
for file in $INTENT_FILES; do
    # Check for models import
    MODELS_IMPORT=$(grep "\"github.com/baphled/kariya/internal/cli/models\"" "$file" 2>/dev/null || true)
    
    if [ -n "$MODELS_IMPORT" ]; then
        # Check if it's being used for forms
        MODELS_FORM_USAGE=$(grep "models\.\w*Form\|models\.New\w*Form" "$file" 2>/dev/null || true)
        
        if [ -n "$MODELS_FORM_USAGE" ]; then
            report_issue "$file" "Deprecated Models" "Deprecated models/ package for forms" \
                "Use screens/*FormScreen instead of models.*Form" \
                "Migration:
OLD: import \"github.com/baphled/kariya/internal/cli/models\"
     type MyIntent struct {
         form *models.CaptureForm
     }

NEW: import \"github.com/baphled/kariya/internal/cli/screens/myfeature\"
     type MyIntent struct {
         formScreen *myfeature.FormScreen
     }

See: docs/FORMS_GUIDE.md, docs/rules/FORMS_WORKFLOW_GUIDE.md"
            CHECK22_ISSUES=$((CHECK22_ISSUES+1))
        fi
    fi
done

if [ $CHECK22_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ No deprecated models/ usage for forms${NC}"
fi

echo ""

# ============================================
# 23. SCREENS DIRECTORY STRUCTURE (Reference: PR #117 burst_management)
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "23. SCREENS DIRECTORY STRUCTURE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHECK23_ISSUES=0

# Check for screens that correspond to subdirectory intents
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # Map intent name to expected screen directory
        # browse_timeline -> timeline, burst_management -> burst_management, etc.
        SCREEN_DIR_NAME=$(echo "$INTENT_NAME" | sed 's/browse_//')
        
        # Check if intent uses screens (has activeScreen or *Screen fields)
        USES_SCREENS=false
        for go_file in "$intent_dir"/*.go; do
            if [ -f "$go_file" ]; then
                if grep -q "screens\.Screen\|activeScreen\|Screen\s\+\*" "$go_file" 2>/dev/null; then
                    USES_SCREENS=true
                    break
                fi
            fi
        done
        
        if [ "$USES_SCREENS" = true ]; then
            # Check for corresponding screens directory
            SCREEN_DIR="internal/cli/screens/$SCREEN_DIR_NAME"
            ALT_SCREEN_DIR="internal/cli/screens/$INTENT_NAME"
            
            FOUND_SCREEN_DIR=""
            if [ -d "$SCREEN_DIR" ]; then
                FOUND_SCREEN_DIR="$SCREEN_DIR"
            elif [ -d "$ALT_SCREEN_DIR" ]; then
                FOUND_SCREEN_DIR="$ALT_SCREEN_DIR"
            fi
            
            if [ -z "$FOUND_SCREEN_DIR" ]; then
                report_issue "$intent_dir" "Screens Directory" "Intent uses screens but no screens directory found" \
                    "Expected: $SCREEN_DIR or $ALT_SCREEN_DIR" \
                    "REQUIRED SCREENS STRUCTURE:
screens/$SCREEN_DIR_NAME/
├── {name}_list_screen.go         # List view screen
├── {name}_list_screen_test.go    # Tests
├── {name}_detail_screen.go       # Detail view (if needed)
├── {name}_detail_screen_test.go  # Tests
├── {feature}_suite_test.go       # Ginkgo test suite
└── modals/                       # Modal subdirectory
    ├── modals_suite_test.go      # Ginkgo test suite
    ├── {modal}_modal.go          # Modal files
    └── {modal}_modal_test.go     # Modal tests

Reference: screens/burst_management/ (PR #117)"
                CHECK23_ISSUES=$((CHECK23_ISSUES+1))
            else
                # ===========================================
                # CHECK SCREENS DIRECTORY CONTENTS
                # ===========================================
                
                # Check for Ginkgo test suite
                SCREEN_SUITE=$(find "$FOUND_SCREEN_DIR" -maxdepth 1 -name "*_suite_test.go" 2>/dev/null | head -1)
                if [ -z "$SCREEN_SUITE" ]; then
                    echo -e "${YELLOW}⚠️  WARNING: No Ginkgo test suite in screens directory${NC}"
                    echo "   Directory: $FOUND_SCREEN_DIR"
                    echo "   Expected: *_suite_test.go"
                    echo ""
                    WARNINGS=$((WARNINGS+1))
                fi
                
                # Check for screen files
                SCREEN_FILES=$(find "$FOUND_SCREEN_DIR" -maxdepth 1 -name "*_screen.go" -type f 2>/dev/null | wc -l)
                if [ "$SCREEN_FILES" -eq 0 ]; then
                    echo -e "${YELLOW}⚠️  WARNING: No screen files found${NC}"
                    echo "   Directory: $FOUND_SCREEN_DIR"
                    echo "   Expected: *_list_screen.go, *_detail_screen.go, etc."
                    echo ""
                    WARNINGS=$((WARNINGS+1))
                fi
                
                # Check for test files for each screen
                for screen_file in "$FOUND_SCREEN_DIR"/*_screen.go; do
                    if [ -f "$screen_file" ]; then
                        TEST_FILE="${screen_file%.go}_test.go"
                        if [ ! -f "$TEST_FILE" ]; then
                            echo -e "${YELLOW}⚠️  WARNING: Missing test file for screen${NC}"
                            echo "   Screen: $screen_file"
                            echo "   Expected: $TEST_FILE"
                            echo ""
                            WARNINGS=$((WARNINGS+1))
                        fi
                    fi
                done
                
                # ===========================================
                # CHECK MODALS SUBDIRECTORY
                # ===========================================
                USES_MODALS=false
                for go_file in "$intent_dir"/*.go; do
                    if [ -f "$go_file" ] && grep -q "Modal\s\+\*\|modalRegistry" "$go_file" 2>/dev/null; then
                        USES_MODALS=true
                        break
                    fi
                done
                
                if [ "$USES_MODALS" = true ]; then
                    MODALS_DIR="$FOUND_SCREEN_DIR/modals"
                    
                    if [ ! -d "$MODALS_DIR" ]; then
                        report_issue "$FOUND_SCREEN_DIR" "Modals Directory" "Intent uses modals but no modals subdirectory" \
                            "Expected: $MODALS_DIR" \
                            "REQUIRED MODALS STRUCTURE:
$FOUND_SCREEN_DIR/modals/
├── modals_suite_test.go      # REQUIRED: Ginkgo test suite
├── {modal}_modal.go          # Each modal
├── {modal}_modal_test.go     # Tests for each modal
├── helpers.go                # OPTIONAL: Shared modal helpers
└── helpers_test.go           # Tests for helpers

Reference: screens/burst_management/modals/ (PR #117)"
                        CHECK23_ISSUES=$((CHECK23_ISSUES+1))
                    else
                        # Check modals directory contents
                        MODALS_SUITE=$(find "$MODALS_DIR" -maxdepth 1 -name "*_suite_test.go" 2>/dev/null | head -1)
                        if [ -z "$MODALS_SUITE" ]; then
                            echo -e "${YELLOW}⚠️  WARNING: No Ginkgo test suite in modals directory${NC}"
                            echo "   Directory: $MODALS_DIR"
                            echo "   Expected: modals_suite_test.go"
                            echo ""
                            WARNINGS=$((WARNINGS+1))
                        fi
                        
                        # Check for test files for each modal
                        for modal_file in "$MODALS_DIR"/*_modal.go; do
                            if [ -f "$modal_file" ]; then
                                TEST_FILE="${modal_file%.go}_test.go"
                                if [ ! -f "$TEST_FILE" ]; then
                                    echo -e "${YELLOW}⚠️  WARNING: Missing test file for modal${NC}"
                                    echo "   Modal: $modal_file"
                                    echo "   Expected: $TEST_FILE"
                                    echo ""
                                    WARNINGS=$((WARNINGS+1))
                                fi
                            fi
                        done
                    fi
                fi
            fi
        fi
    done
fi

# Check for screens with modals in root (should be in modals subdirectory)
SCREEN_DIRS=$(find internal/cli/screens -mindepth 1 -maxdepth 1 -type d 2>/dev/null | grep -v base || true)

if [ -n "$SCREEN_DIRS" ]; then
    for screen_dir in $SCREEN_DIRS; do
        SCREEN_NAME=$(basename "$screen_dir")
        
        # Check if screen has modal files in root but no modals subdirectory
        MODAL_FILES=$(find "$screen_dir" -maxdepth 1 -name "*modal*.go" -type f 2>/dev/null | wc -l)
        MODALS_SUBDIR="$screen_dir/modals"
        
        if [ "$MODAL_FILES" -gt 2 ] && [ ! -d "$MODALS_SUBDIR" ]; then
            echo -e "${YELLOW}⚠️  WARNING: Modal files in screen root (should be in modals/)${NC}"
            echo "   Screen: $SCREEN_NAME ($MODAL_FILES modal files)"
            echo "   Recommendation: Create $MODALS_SUBDIR and move modal files"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    done
fi

if [ $CHECK23_ISSUES -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ Screens structure looks good${NC}"
elif [ $CHECK23_ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ No screen violations (warnings present)${NC}"
fi

echo ""

# ============================================
# 24. INTENT STRUCT REQUIRED FIELDS (Reference: browse_timeline, burst_management)
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "24. INTENT STRUCT REQUIRED FIELDS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        TYPES_FILE="$intent_dir/types.go"
        INTENT_FILE="$intent_dir/intent.go"
        
        # Find the file containing the Intent struct
        STRUCT_FILE=""
        if [ -f "$TYPES_FILE" ] && grep -q "^type.*Intent struct" "$TYPES_FILE" 2>/dev/null; then
            STRUCT_FILE="$TYPES_FILE"
        elif [ -f "$INTENT_FILE" ] && grep -q "^type.*Intent struct" "$INTENT_FILE" 2>/dev/null; then
            STRUCT_FILE="$INTENT_FILE"
        fi
        
        if [ -n "$STRUCT_FILE" ]; then
            # Check for required fields in Intent struct
            
            # Check for 'active bool' field
            HAS_ACTIVE=$(grep -q "active\s\+bool" "$STRUCT_FILE" && echo "yes" || echo "no")
            if [ "$HAS_ACTIVE" = "no" ]; then
                report_issue "$STRUCT_FILE" "Required Fields" "Missing 'active bool' field" \
                    "Intent struct must have 'active bool' field for lifecycle management" \
                    "Required pattern:
type Intent struct {
    *intents.BaseIntent
    context *IntentContext
    state   State
    active  bool  // <- REQUIRED
    result  *intents.IntentResult[*Result]
    ...
}"
            fi
            
            # Check for 'result *intents.IntentResult' field
            HAS_RESULT=$(grep -q "result\s\+\*intents\.IntentResult" "$STRUCT_FILE" && echo "yes" || echo "no")
            if [ "$HAS_RESULT" = "no" ]; then
                report_issue "$STRUCT_FILE" "Required Fields" "Missing 'result *intents.IntentResult' field" \
                    "Intent struct must have typed result field for completion" \
                    "Required pattern:
result *intents.IntentResult[*Result]  // <- REQUIRED"
            fi
            
            # Check for modalRegistry if intent uses modals
            USES_MODALS=$(grep -q "Modal\s\+\*" "$STRUCT_FILE" && echo "yes" || echo "no")
            HAS_REGISTRY=$(grep -q "modalRegistry\s\+\*intents\.ModalRegistry" "$STRUCT_FILE" && echo "yes" || echo "no")
            
            if [ "$USES_MODALS" = "yes" ] && [ "$HAS_REGISTRY" = "no" ]; then
                echo -e "${YELLOW}⚠️  WARNING: Using modals without ModalRegistry${NC}"
                echo "   File: $STRUCT_FILE"
                echo "   Recommendation: Use ModalRegistry for unified modal management"
                echo ""
                echo "   Example:"
                echo "   modalRegistry *intents.ModalRegistry"
                echo ""
                WARNINGS=$((WARNINGS+1))
            fi
        fi
    done
fi

echo ""

# ============================================
# 25. CONTEXT VALIDATE METHOD (Reference: browse_timeline, burst_management)
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "25. CONTEXT VALIDATE METHOD"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        CONTEXT_FILE="$intent_dir/context.go"
        
        if [ -f "$CONTEXT_FILE" ]; then
            # Check for Validate method
            HAS_VALIDATE=$(grep -q "func.*IntentContext.*Validate" "$CONTEXT_FILE" && echo "yes" || echo "no")
            
            if [ "$HAS_VALIDATE" = "no" ]; then
                report_issue "$CONTEXT_FILE" "Context Validate" "Missing Validate() method on IntentContext" \
                    "IntentContext must have Validate() method for input validation" \
                    "Required pattern:
func (c *IntentContext) Validate() error {
    if c.Events == nil {
        c.Events = make([]*career.CareerEvent, 0)
    }
    // ... initialize nil fields
    return nil
}"
            fi
        fi
    done
fi

echo ""

# ============================================
# 26. SCREENRESULT DISPATCHER USAGE (Reference: browse_timeline, burst_management)
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "26. SCREENRESULT DISPATCHER USAGE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        HANDLERS_FILE="$intent_dir/handlers.go"
        INTENT_FILE="$intent_dir/intent.go"
        
        # Check if intent has handler methods
        HAS_HANDLERS=false
        for go_file in "$intent_dir"/*.go; do
            if [ -f "$go_file" ] && grep -q "func.*HandleCancel\|func.*HandleNavigate" "$go_file" 2>/dev/null; then
                HAS_HANDLERS=true
                break
            fi
        done
        
        if [ "$HAS_HANDLERS" = true ]; then
            # Check for NewScreenResultDispatcher usage
            USES_DISPATCHER=false
            for go_file in "$intent_dir"/*.go; do
                if [ -f "$go_file" ] && grep -q "NewScreenResultDispatcher" "$go_file" 2>/dev/null; then
                    USES_DISPATCHER=true
                    break
                fi
            done
            
            if [ "$USES_DISPATCHER" = false ]; then
                echo -e "${YELLOW}⚠️  WARNING: ScreenResultHandler without dispatcher${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Recommendation: Use intents.NewScreenResultDispatcher for type-safe dispatch"
                echo ""
                echo "   Example:"
                echo "   return intents.NewScreenResultDispatcher(i).Dispatch(screenResult)"
                echo ""
                WARNINGS=$((WARNINGS+1))
            fi
        fi
    done
fi

echo ""

# ============================================
# 27. INTERFACES FILE FOR DEPENDENCY INJECTION
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "27. INTERFACES FILE FOR DEPENDENCY INJECTION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        INTERFACES_FILE="$intent_dir/interfaces.go"
        CONTEXT_FILE="$intent_dir/context.go"
        
        # Check if intent uses services (has Service field in context or struct)
        USES_SERVICES=false
        for go_file in "$intent_dir"/*.go; do
            if [ -f "$go_file" ] && grep -q "Service\s" "$go_file" 2>/dev/null; then
                USES_SERVICES=true
                break
            fi
        done
        
        if [ "$USES_SERVICES" = true ]; then
            if [ ! -f "$INTERFACES_FILE" ]; then
                echo -e "${YELLOW}⚠️  WARNING: Uses services but no interfaces.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Recommendation: Define service interfaces for dependency injection"
                echo ""
                echo "   Create: $INTERFACES_FILE"
                echo ""
                echo "   Example:"
                echo "   type EventService interface {"
                echo "       DeleteEvent(ctx context.Context, eventID string) error"
                echo "       ListEvents(ctx context.Context, filters *Filters) ([]*Event, error)"
                echo "   }"
                echo ""
                WARNINGS=$((WARNINGS+1))
            fi
        fi
    done
fi

echo ""

# ============================================
# 28. TEST COVERAGE FOR CORE FILES
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "28. TEST COVERAGE FOR CORE FILES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        MISSING_TESTS=""
        
        # Check for test files for core files
        CORE_FILES=("context" "result" "constants" "messages" "intent")
        
        for core in "${CORE_FILES[@]}"; do
            SOURCE_FILE="$intent_dir/${core}.go"
            TEST_FILE="$intent_dir/${core}_test.go"
            
            if [ -f "$SOURCE_FILE" ] && [ ! -f "$TEST_FILE" ]; then
                MISSING_TESTS="$MISSING_TESTS ${core}_test.go"
            fi
        done
        
        if [ -n "$MISSING_TESTS" ]; then
            echo -e "${YELLOW}⚠️  WARNING: Missing test files${NC}"
            echo "   Intent: $INTENT_NAME"
            echo "   Missing:$MISSING_TESTS"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
        
        # Check for Ginkgo test suite file
        SUITE_FILE=$(find "$intent_dir" -maxdepth 1 -name "*_suite_test.go" 2>/dev/null | head -1)
        if [ -z "$SUITE_FILE" ]; then
            echo -e "${YELLOW}⚠️  WARNING: Missing Ginkgo test suite${NC}"
            echo "   Intent: $INTENT_NAME"
            echo "   Expected: ${INTENT_NAME}_suite_test.go"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    done
fi

echo ""

# ============================================
# 29. SCREENS MODALS SUBDIRECTORY
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "29. SCREENS MODALS SUBDIRECTORY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # Find corresponding screens directory
        SCREEN_DIR_NAME=$(echo "$INTENT_NAME" | sed 's/browse_//' | sed 's/_management$//')
        SCREEN_DIR="internal/cli/screens/$SCREEN_DIR_NAME"
        ALT_SCREEN_DIR="internal/cli/screens/$INTENT_NAME"
        
        FOUND_SCREEN_DIR=""
        if [ -d "$SCREEN_DIR" ]; then
            FOUND_SCREEN_DIR="$SCREEN_DIR"
        elif [ -d "$ALT_SCREEN_DIR" ]; then
            FOUND_SCREEN_DIR="$ALT_SCREEN_DIR"
        fi
        
        if [ -n "$FOUND_SCREEN_DIR" ]; then
            # Check if intent uses modals
            USES_MODALS=false
            for go_file in "$intent_dir"/*.go; do
                if [ -f "$go_file" ] && grep -q "Modal\s\+\*" "$go_file" 2>/dev/null; then
                    USES_MODALS=true
                    break
                fi
            done
            
            if [ "$USES_MODALS" = true ]; then
                MODALS_DIR="$FOUND_SCREEN_DIR/modals"
                if [ ! -d "$MODALS_DIR" ]; then
                    echo -e "${YELLOW}⚠️  WARNING: Intent uses modals but no modals/ directory${NC}"
                    echo "   Intent: $INTENT_NAME"
                    echo "   Expected: $MODALS_DIR"
                    echo "   Recommendation: Extract feature-specific modals to screens package"
                    echo ""
                    WARNINGS=$((WARNINGS+1))
                fi
            fi
        fi
    done
fi

echo ""

# ============================================
# SUMMARY
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Total checks run: 29"
echo ""
echo -e "Violations (NEW intents):     ${RED}$VIOLATIONS${NC}"
echo -e "Legacy warnings (EXISTING):   ${CYAN}$LEGACY_WARNINGS${NC}"
echo -e "Warnings (recommendations):   ${YELLOW}$WARNINGS${NC}"
echo ""

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ] && [ $LEGACY_WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL ARCHITECTURE CHECKS PASSED${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ PASSED - No blocking violations${NC}"
    echo ""
    if [ $LEGACY_WARNINGS -gt 0 ]; then
        echo -e "${CYAN}📋 Legacy intents have $LEGACY_WARNINGS issues (non-blocking)${NC}"
        echo "   These should be fixed during migration but don't block commits."
    fi
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  $WARNINGS recommendations for improvement${NC}"
    fi
    echo ""
    echo "See: docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
else
    echo -e "${RED}❌ FAILED: $VIOLATIONS violations in NEW intents must be fixed${NC}"
    echo ""
    echo "New intents must follow the architecture rules."
    echo "Reference implementations: browse_timeline, burst_management"
    echo ""
    if [ $LEGACY_WARNINGS -gt 0 ]; then
        echo -e "${CYAN}📋 Legacy intents: $LEGACY_WARNINGS issues (allowed)${NC}"
    fi
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  $WARNINGS recommendations${NC}"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Fix violations before committing."
    echo "See: docs/checklists/INTENT_DEVELOPMENT_CHECKLIST.md"
    echo ""
    exit 1
fi
