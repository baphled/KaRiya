#!/bin/bash

# Intent Architecture Enforcement
# This script enforces strict architectural rules for intent implementations.
# It catches violations that would otherwise slip through code review.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

VIOLATIONS=0
WARNINGS=0

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏛️  INTENT ARCHITECTURE ENFORCEMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
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
        echo -e "${RED}❌ VIOLATION: Untyped state field${NC}"
        echo "   File: $file"
        echo "   Line: $UNTYPED_STATE"
        echo "   Rule: State fields must use typed enum, not raw string"
        echo ""
        echo "   Required pattern:"
        echo "   type ${INTENT_NAME^}State string"
        echo ""
        echo "   const ("
        echo "       State... ${INTENT_NAME^}State = \"...\""
        echo "   )"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
    
    # Check if intent has state field but missing typed enum definition
    HAS_STATE=$(grep -q "state.*State" "$file" && echo "yes" || echo "no")
    HAS_TYPED_ENUM=$(grep -q "type.*State string" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_STATE" = "yes" ] && [ "$HAS_TYPED_ENUM" = "no" ]; then
        echo -e "${RED}❌ VIOLATION: Missing typed state enum definition${NC}"
        echo "   File: $file"
        echo "   Has state field but no 'type <Name>State string' declaration"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

echo ""

# ============================================
# 2. FLATTENED STATE MODEL CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2. FLATTENED STATE MODEL"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

WRAPPED_STATE=$(grep -n "^\s*state\s\+\*.*Model" internal/cli/intents/*_intent.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$WRAPPED_STATE" ]; then
    echo -e "${RED}❌ VIOLATION: Wrapped state model found${NC}"
    echo "   Rule: Intent state must be flattened directly into the intent struct"
    echo "   Found: $WRAPPED_STATE"
    echo ""
    echo "   BAD:  state *BrowseTimelineModel"
    echo "   GOOD: Flatten all model fields directly into intent"
    echo ""
    VIOLATIONS=$((VIOLATIONS+1))
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

BG_CONTEXT=$(grep -n "context\.Background()" internal/cli/intents/*_intent.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$BG_CONTEXT" ]; then
    echo -e "${RED}❌ VIOLATION: context.Background() in intent${NC}"
    echo "   Rule: Intents must use i.getContext() for cancellation support"
    echo "   Found:"
    echo "$BG_CONTEXT" | sed 's/^/   /'
    echo ""
    echo "   BAD:  ctx := context.Background()"
    echo "   GOOD: ctx := i.getContext()"
    echo ""
    VIOLATIONS=$((VIOLATIONS+1))
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

DEAD_CODE=$(grep -n "should not be reached\|LEGACY.*not.*reached\|This case is kept for backward compatibility but should not be reached" internal/cli/intents/*.go 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$DEAD_CODE" ]; then
    echo -e "${RED}❌ VIOLATION: Dead code markers found${NC}"
    echo "   Rule: Remove unreachable code or create cleanup task"
    echo "   Found:"
    echo "$DEAD_CODE" | sed 's/^/   /'
    echo ""
    VIOLATIONS=$((VIOLATIONS+1))
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
    echo -e "${RED}❌ VIOLATION: Intent missing *BaseIntent${NC}"
    echo "   Files:"
    echo "$MISSING_BASE" | sed 's/^/   /'
    echo ""
    echo "   Required: All intents must embed *BaseIntent"
    echo ""
    VIOLATIONS=$((VIOLATIONS+1))
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

for file in $INTENT_FILES; do
    # Check if intent has modals
    HAS_MODALS=$(grep -q "Modal.*\*.*Modal" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_MODALS" = "yes" ]; then
        # Check if using RenderModalOverlay
        USES_OVERLAY=$(grep -q "RenderModalOverlay" "$file" && echo "yes" || echo "no")
        
        if [ "$USES_OVERLAY" = "no" ]; then
            echo -e "${RED}❌ VIOLATION: Modal without RenderModalOverlay${NC}"
            echo "   File: $file"
            echo "   Rule: Use behaviors.RenderModalOverlay() for modal rendering"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ Modal patterns correct${NC}"
fi

echo ""

# ============================================
# 9. SCREENRESULTHANDLER IMPLEMENTATION CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "8. SCREENRESULTHANDLER IMPLEMENTATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check if intent uses screens
    HAS_SCREENS=$(grep -q "screens\.Screen" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_SCREENS" = "yes" ]; then
        # Check for ScreenResultHandler implementation
        HAS_HANDLER=$(grep -q "var _ ScreenResultHandler" "$file" && echo "yes" || echo "no")
        
        if [ "$HAS_HANDLER" = "no" ]; then
            echo -e "${RED}❌ VIOLATION: Missing ScreenResultHandler${NC}"
            echo "   File: $file"
            echo "   Rule: Intents using screens must implement ScreenResultHandler"
            echo ""
            echo "   Add: var _ ScreenResultHandler = (*YourIntent)(nil)"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ ScreenResultHandler implementations correct${NC}"
fi

echo ""

# ============================================
# 9. LEAN INTENT PATTERN CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "9. LEAN INTENT PATTERN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check for direct database/service calls in Update() method
    # Intent should delegate to screens/services, not contain business logic
    UPDATE_METHOD=$(sed -n '/^func.*Update.*tea\.Msg/,/^func /p' "$file" 2>/dev/null || true)
    
    if [ -n "$UPDATE_METHOD" ]; then
        # Check for SQL queries in Update
        HAS_SQL=$(echo "$UPDATE_METHOD" | grep -E "\.Query\(|\.Exec\(|\.QueryRow\(|INSERT INTO|SELECT.*FROM|UPDATE.*SET|DELETE FROM" || true)
        
        if [ -n "$HAS_SQL" ]; then
            echo -e "${RED}❌ VIOLATION: Business logic in Update()${NC}"
            echo "   File: $file"
            echo "   Rule: Intents should orchestrate, not contain business logic"
            echo "   Found: Direct SQL/database calls in Update() method"
            echo ""
            echo "   Move business logic to:"
            echo "   - Service layer for data operations"
            echo "   - Screen layer for UI logic"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
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

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ Intents follow lean orchestration pattern${NC}"
fi

echo ""

# ============================================
# 10. REQUIRED STATE FIELD CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "10. REQUIRED STATE FIELD"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check if intent has a state field
    HAS_STATE=$(grep -q "^\s*state\s\+\w" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_STATE" = "no" ]; then
        echo -e "${RED}❌ VIOLATION: Missing state field${NC}"
        echo "   File: $file"
        echo "   Rule: All intents must have a state field (typed enum)"
        echo ""
        echo "   Required:"
        echo "   state MyIntentState  // Typed state enum"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ All intents have state fields${NC}"
fi

echo ""

# ============================================
# 11. EXPLICIT SCREEN FIELDS CHECK
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "11. EXPLICIT SCREEN FIELDS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for file in $INTENT_FILES; do
    # Check if intent uses screens
    HAS_ACTIVE_SCREEN=$(grep -q "activeScreen.*screens\.Screen" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_ACTIVE_SCREEN" = "yes" ]; then
        # Check for explicitly typed screen fields (listScreen, detailScreen, etc.)
        TYPED_SCREENS=$(grep "Screen\s\+\*.*\..*Screen" "$file" 2>/dev/null | wc -l)
        
        if [ "$TYPED_SCREENS" -eq "0" ]; then
            echo -e "${RED}❌ VIOLATION: Only generic activeScreen field${NC}"
            echo "   File: $file"
            echo "   Rule: Intents must declare explicit typed screen fields"
            echo ""
            echo "   Required pattern:"
            echo "   listScreen   *myfeature.ListScreen"
            echo "   detailScreen *myfeature.DetailScreen"
            echo "   activeScreen screens.Screen  // Points to one of the above"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
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

for file in $INTENT_FILES; do
    INTENT_NAME=$(basename "$file" _intent.go | sed 's/_/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1' | sed 's/ //g')
    
    # Check if intent has a context field
    HAS_CONTEXT=$(grep -q "context\s\+\*${INTENT_NAME}Context" "$file" && echo "yes" || echo "no")
    
    if [ "$HAS_CONTEXT" = "no" ]; then
        # Check if it's using generic context
        HAS_ANY_CONTEXT=$(grep -q "context\s\+\*.*Context" "$file" && echo "yes" || echo "no")
        
        if [ "$HAS_ANY_CONTEXT" = "no" ]; then
            echo -e "${RED}❌ VIOLATION: Missing context field${NC}"
            echo "   File: $file"
            echo "   Rule: Intents should have a context field for input parameters"
            echo ""
            echo "   Add: context *${INTENT_NAME}Context"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
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

for file in $INTENT_FILES; do
    BASENAME=$(basename "$file" .go)
    DIRNAME=$(dirname "$file")
    
    # Check if ANY Context struct is defined in intent file (ending with Context)
    CONTEXT_STRUCTS=$(grep "^type.*Context struct" "$file" | grep -v "// " || true)
    
    if [ -n "$CONTEXT_STRUCTS" ]; then
        echo -e "${RED}❌ VIOLATION: Context struct(s) defined in intent file${NC}"
        echo "   File: $file"
        echo "   Rule: Context structs must be in separate file"
        echo ""
        echo "   Found:"
        echo "$CONTEXT_STRUCTS" | sed 's/^/   /'
        echo ""
        echo "   Required separation:"
        echo "   - Intent: ${BASENAME}.go (only intent struct and methods)"
        echo "   - Context: ${BASENAME%_intent}.go OR ${BASENAME%_intent}_context.go"
        echo ""
        echo "   Context = Input parameters (Events, Services, Config, etc.)"
        echo "   Context should be in its own file, separate from intent implementation"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
    
    # Check if ANY Model struct is defined in intent file (ending with Model)
    MODEL_STRUCTS=$(grep "^type.*Model struct" "$file" | grep -v "// " || true)
    
    if [ -n "$MODEL_STRUCTS" ]; then
        echo -e "${RED}❌ VIOLATION: Model struct(s) defined in intent file${NC}"
        echo "   File: $file"
        echo "   Rule: Model structs must be flattened into intent OR in separate file"
        echo ""
        echo "   Found:"
        echo "$MODEL_STRUCTS" | sed 's/^/   /'
        echo ""
        echo "   RECOMMENDED: Flatten model fields directly into intent struct"
        echo "   ALTERNATIVE: Create separate file: ${BASENAME%_intent}_model.go"
        echo ""
        echo "   Model = State wrapper (should be flattened per Check #2)"
        echo "   If you must keep it, put it in a separate file"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
    
    # Check if Screen structs are defined in intent file
    SCREEN_IN_INTENT=$(grep -q "^type.*Screen struct" "$file" && echo "yes" || echo "no")
    
    if [ "$SCREEN_IN_INTENT" = "yes" ]; then
        echo -e "${RED}❌ VIOLATION: Screen defined in intent file${NC}"
        echo "   File: $file"
        echo "   Rule: Screen structs must be in screens/ package"
        echo ""
        echo "   Required:"
        echo "   - Screens: internal/cli/screens/myfeature/*.go"
        echo "   - Intent: internal/cli/intents/${BASENAME}.go"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
    
    # Check if Modal structs are defined in intent file
    MODAL_IN_INTENT=$(grep -q "^type.*Modal struct" "$file" && echo "yes" || echo "no")
    
    if [ "$MODAL_IN_INTENT" = "yes" ]; then
        echo -e "${RED}❌ VIOLATION: Modal defined in intent file${NC}"
        echo "   File: $file"
        echo "   Rule: Modal structs must be in components/ or uikit/feedback/"
        echo ""
        echo "   Required:"
        echo "   - Modals: internal/cli/components/*.go or internal/cli/uikit/feedback/*.go"
        echo "   - Intent: internal/cli/intents/${BASENAME}.go"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ File separation correct${NC}"
fi

echo ""

# ============================================
# 17. SUBDIRECTORY STRUCTURE
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "17. SUBDIRECTORY STRUCTURE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check for subdirectory-based intents
SUBDIRS=$(find internal/cli/intents -mindepth 1 -maxdepth 1 -type d 2>/dev/null | grep -v types || true)

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # Check for required core files (5 required)
        REQUIRED_FILES=("context.go" "result.go" "constants.go" "messages.go" "intent.go")
        MISSING_FILES=""
        
        for req_file in "${REQUIRED_FILES[@]}"; do
            if [ ! -f "$intent_dir/$req_file" ]; then
                MISSING_FILES="$MISSING_FILES $req_file"
            fi
        done
        
        if [ -n "$MISSING_FILES" ]; then
            echo -e "${RED}❌ VIOLATION: Incomplete subdirectory structure${NC}"
            echo "   Intent: $INTENT_NAME"
            echo "   Missing required files:$MISSING_FILES"
            echo ""
            echo "   Required structure (5 CORE FILES):"
            echo "   intents/$INTENT_NAME/"
            echo "   ├── context.go    (IntentContext struct + Validate())"
            echo "   ├── result.go     (Result struct)"
            echo "   ├── constants.go  (State enum)"
            echo "   ├── messages.go   (ALL *Msg types)"
            echo "   └── intent.go     (NewIntent, Init, Update, View, Result)"
            echo ""
            echo "   Optional recommended files:"
            echo "   ├── types.go      (Intent struct definition)"
            echo "   ├── handlers.go   (ScreenResultHandler methods)"
            echo "   ├── helpers.go    (Helper methods)"
            echo "   ├── filters.go    (Domain-specific filter logic)"
            echo "   └── interfaces.go (Service interfaces)"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
        
        # Check for recommended files (WARNING only)
        RECOMMENDED_FILES=("types.go" "handlers.go" "helpers.go")
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
                    echo "   Recommendation: Extract to keep intent.go under 300 lines"
                    echo ""
                    WARNINGS=$((WARNINGS+1))
                fi
            fi
        fi
    done
fi

# Check for OLD flat structure (grace period - WARNING only)
OLD_INTENTS=$(find internal/cli/intents -maxdepth 1 -name "*_intent.go" -type f 2>/dev/null || true)
OLD_COUNT=$(echo "$OLD_INTENTS" | grep -c "_intent.go" 2>/dev/null || echo 0)

if [ $OLD_COUNT -gt 0 ]; then
    echo -e "${YELLOW}⚠️  WARNING: $OLD_COUNT intents using old flat structure${NC}"
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
    WARNINGS=$((WARNINGS+1))
fi

if [ $VIOLATIONS -eq 0 ] && [ $OLD_COUNT -eq 0 ]; then
    echo -e "${GREEN}✅ All intents use subdirectory structure${NC}"
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ Subdirectory intents are complete${NC}"
fi

echo ""

# ============================================
# 18. INTENT FILE SIZE LIMIT
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "18. INTENT FILE SIZE LIMIT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check subdirectory-based intents
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_FILE="$intent_dir/intent.go"
        if [ -f "$INTENT_FILE" ]; then
            LINE_COUNT=$(wc -l < "$INTENT_FILE")
            INTENT_NAME=$(basename "$intent_dir")
            
            if [ $LINE_COUNT -gt 600 ]; then
                echo -e "${RED}❌ VIOLATION: intent.go exceeds 600 lines${NC}"
                echo "   File: $INTENT_FILE ($LINE_COUNT lines)"
                echo "   Rule: Intent should be broker only (orchestration)"
                echo ""
                echo "   Required actions:"
                echo "   - Extract views to screens/$INTENT_NAME/"
                echo "   - Move business logic to context.go"
                echo "   - Extract helpers to appropriate packages"
                echo ""
                echo "   Target: 200-400 lines"
                echo "   Maximum: 600 lines (hard limit)"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
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

# Also check old flat structure files
if [ -n "$OLD_INTENTS" ]; then
    for intent_file in $OLD_INTENTS; do
        if [ -f "$intent_file" ] && [ -n "$intent_file" ]; then
            LINE_COUNT=$(wc -l < "$intent_file")
            INTENT_NAME=$(basename "$intent_file" _intent.go)
            
            if [ $LINE_COUNT -gt 1000 ]; then
                echo -e "${YELLOW}⚠️  WARNING: Legacy intent file exceeds 1,000 lines${NC}"
                echo "   File: $intent_file ($LINE_COUNT lines)"
                echo "   Action: Migrate to subdirectory structure"
                echo ""
                WARNINGS=$((WARNINGS+1))
            fi
        fi
    done
fi

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ All intent files within size limits${NC}"
fi

echo ""

# ============================================
# 19. NO RENDERING IN INTENT
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "19. NO RENDERING IN INTENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

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
                echo -e "${RED}❌ VIOLATION: Rendering methods in intent.go${NC}"
                echo "   File: $INTENT_FILE ($TOTAL_RENDER render methods)"
                echo "   Rule: Intent should only have View() that delegates to screens"
                echo ""
                echo "   Found methods:"
                grep -n "^func.*render\|^func.*Render\|^func.*View" "$INTENT_FILE" 2>/dev/null | head -10 | sed 's/^/   /' || true
                echo ""
                echo "   Required action:"
                echo "   - Extract rendering to screens/$INTENT_NAME/"
                echo "   - Keep only View() in intent.go"
                echo "   - View() should delegate: return i.activeScreen.View()"
                echo ""
                echo "   See: docs/guides/SCREEN_EXTRACTION_GUIDE.md"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
    done
fi

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ No rendering logic in intent files${NC}"
fi

echo ""

# ============================================
# 20. TYPE LOCATION ENFORCEMENT
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "20. TYPE LOCATION ENFORCEMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # Check Context location - must be in context.go, NOT in intent.go or types.go
        if [ -f "$intent_dir/intent.go" ]; then
            CONTEXT_IN_INTENT=$(grep "^type.*Context struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$CONTEXT_IN_INTENT" ]; then
                echo -e "${RED}❌ VIOLATION: Context defined in intent.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: Context must be in context.go"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            CONTEXT_IN_TYPES=$(grep "^type.*Context struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$CONTEXT_IN_TYPES" ]; then
                echo -e "${RED}❌ VIOLATION: Context defined in types.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: Context must be in context.go"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        # Check Result location - must be in result.go
        if [ -f "$intent_dir/intent.go" ]; then
            RESULT_IN_INTENT=$(grep "^type.*Result struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$RESULT_IN_INTENT" ]; then
                echo -e "${RED}❌ VIOLATION: Result defined in intent.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: Result must be in result.go"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            RESULT_IN_TYPES=$(grep "^type.*Result struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$RESULT_IN_TYPES" ]; then
                echo -e "${RED}❌ VIOLATION: Result defined in types.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: Result must be in result.go"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        # Check State enum location - must be in constants.go
        if [ -f "$intent_dir/intent.go" ]; then
            STATE_IN_INTENT=$(grep "^type.*State string" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$STATE_IN_INTENT" ]; then
                echo -e "${RED}❌ VIOLATION: State enum defined in intent.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: State enum must be in constants.go"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            STATE_IN_TYPES=$(grep "^type.*State string" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$STATE_IN_TYPES" ]; then
                echo -e "${RED}❌ VIOLATION: State enum defined in types.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: State enum must be in constants.go"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        # Check Msg types location (CRITICAL - must be in messages.go)
        if [ -f "$intent_dir/intent.go" ]; then
            MSG_IN_INTENT=$(grep "^type.*Msg struct" "$intent_dir/intent.go" 2>/dev/null || true)
            if [ -n "$MSG_IN_INTENT" ]; then
                echo -e "${RED}❌ VIOLATION: Msg types defined in intent.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: ALL *Msg types must be in messages.go"
                echo ""
                echo "   Found:"
                echo "$MSG_IN_INTENT" | sed 's/^/   /'
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        if [ -f "$intent_dir/types.go" ]; then
            MSG_IN_TYPES=$(grep "^type.*Msg struct" "$intent_dir/types.go" 2>/dev/null || true)
            if [ -n "$MSG_IN_TYPES" ]; then
                echo -e "${RED}❌ VIOLATION: Msg types defined in types.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: ALL *Msg types must be in messages.go"
                echo ""
                echo "   Found:"
                echo "$MSG_IN_TYPES" | sed 's/^/   /'
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
        
        # Check if Msg types are in constants.go (WRONG - should be messages.go)
        if [ -f "$intent_dir/constants.go" ]; then
            MSG_IN_CONSTANTS=$(grep "^type.*Msg struct" "$intent_dir/constants.go" 2>/dev/null || true)
            if [ -n "$MSG_IN_CONSTANTS" ]; then
                echo -e "${RED}❌ VIOLATION: Msg types in constants.go${NC}"
                echo "   Intent: $INTENT_NAME"
                echo "   Rule: Msg types must be in messages.go (not constants.go)"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
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
            echo -e "${RED}❌ VIOLATION: No Intent struct found${NC}"
            echo "   Intent: $INTENT_NAME"
            echo "   Rule: Intent struct must be in intent.go OR types.go"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    done
fi

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ All types in correct locations${NC}"
fi

echo ""

# ============================================
# 21. UIKIT COMPONENT USAGE
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "21. UIKIT COMPONENT USAGE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check for deprecated components usage
for file in $INTENT_FILES; do
    # Check for components.KeyBadge (should use primitives.HelpKeyBadge)
    KEYBADGE_USAGE=$(grep "components\.KeyBadge" "$file" 2>/dev/null || true)
    if [ -n "$KEYBADGE_USAGE" ]; then
        echo -e "${RED}❌ VIOLATION: Deprecated components.KeyBadge${NC}"
        echo "   File: $file"
        echo "   Rule: Use primitives.HelpKeyBadge() instead"
        echo ""
        echo "   Found:"
        echo "$KEYBADGE_USAGE" | head -3 | sed 's/^/   /'
        echo ""
        echo "   Migration:"
        echo "   OLD: components.KeyBadge(\"key\", \"label\")"
        echo "   NEW: primitives.HelpKeyBadge(\"key\", \"label\", theme)"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
    fi
    
    # Check for components.StandardView (should use layout.NewScreenLayout)
    STANDARDVIEW_USAGE=$(grep "components\.StandardView" "$file" 2>/dev/null || true)
    if [ -n "$STANDARDVIEW_USAGE" ]; then
        echo -e "${RED}❌ VIOLATION: Deprecated components.StandardView${NC}"
        echo "   File: $file"
        echo "   Rule: Use layout.NewScreenLayout() instead"
        echo ""
        echo "   Found:"
        echo "$STANDARDVIEW_USAGE" | head -3 | sed 's/^/   /'
        echo ""
        echo "   Migration:"
        echo "   OLD: components.StandardView{...}"
        echo "   NEW: layout.NewScreenLayout(theme).WithContent(...).Render()"
        echo ""
        VIOLATIONS=$((VIOLATIONS+1))
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
            echo "   Instead of: lipgloss.NewStyle().Foreground(lipgloss.Color(\"#ff0000\"))"
            echo "   Use theme: style := lipgloss.NewStyle().Foreground(theme.Error())"
            echo "   Or UIKit: primitives.ErrorText(\"message\", theme)"
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
            echo "   Available centralized modals:"
            echo "   - feedback.NewConfirmModal() - Confirmation dialogs"
            echo "   - feedback.NewErrorModal() - Error messages"
            echo "   - feedback.NewLoadingModal() - Loading states"
            echo "   - feedback.NewSuccessModal() - Success messages"
            echo "   - feedback.NewWarningModal() - Warnings"
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
                echo -e "${RED}❌ VIOLATION: Deprecated components.KeyBadge${NC}"
                echo "   File: $INTENT_FILE"
                echo "   Rule: Use primitives.HelpKeyBadge() instead"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
            
            STANDARDVIEW_USAGE=$(grep "components\.StandardView" "$INTENT_FILE" 2>/dev/null || true)
            if [ -n "$STANDARDVIEW_USAGE" ]; then
                echo -e "${RED}❌ VIOLATION: Deprecated components.StandardView${NC}"
                echo "   File: $INTENT_FILE"
                echo "   Rule: Use layout.NewScreenLayout() instead"
                echo ""
                VIOLATIONS=$((VIOLATIONS+1))
            fi
        fi
    done
fi

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ UIKit components used correctly${NC}"
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ No UIKit violations (warnings present)${NC}"
fi

echo ""

# ============================================
# 22. DEPRECATED MODELS PACKAGE FOR FORMS
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "22. DEPRECATED MODELS PACKAGE FOR FORMS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check for models.*Form usage in intents
for file in $INTENT_FILES; do
    # Check for models import
    MODELS_IMPORT=$(grep "\"github.com/baphled/kariya/internal/cli/models\"" "$file" 2>/dev/null || true)
    
    if [ -n "$MODELS_IMPORT" ]; then
        # Check if it's being used for forms
        MODELS_FORM_USAGE=$(grep "models\.\w*Form\|models\.New\w*Form" "$file" 2>/dev/null || true)
        
        if [ -n "$MODELS_FORM_USAGE" ]; then
            echo -e "${RED}❌ VIOLATION: Deprecated models/ package for forms${NC}"
            echo "   File: $file"
            echo "   Rule: Use screens/*FormScreen instead of models.*Form"
            echo ""
            echo "   Found (DEPRECATED):"
            echo "$MODELS_FORM_USAGE" | head -3 | sed 's/^/   /'
            echo ""
            echo "   Migration:"
            echo "   OLD: import \"github.com/baphled/kariya/internal/cli/models\""
            echo "        type MyIntent struct {"
            echo "            form *models.CaptureForm"
            echo "        }"
            echo ""
            echo "   NEW: import \"github.com/baphled/kariya/internal/cli/screens/myfeature\""
            echo "        type MyIntent struct {"
            echo "            formScreen *myfeature.FormScreen"
            echo "        }"
            echo ""
            echo "   Rationale:"
            echo "   - models/ package is DEPRECATED for forms"
            echo "   - Use screens/ package with embedded forms.Form"
            echo "   - Cleaner: intents → screens → forms (not intents → models → forms)"
            echo ""
            echo "   See: docs/FORMS_GUIDE.md, docs/rules/FORMS_WORKFLOW_GUIDE.md"
            echo ""
            VIOLATIONS=$((VIOLATIONS+1))
        fi
    fi
done

if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ No deprecated models/ usage for forms${NC}"
fi

echo ""

# ============================================
# 23. SCREENS DIRECTORY STRUCTURE
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "23. SCREENS DIRECTORY STRUCTURE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check for screens that correspond to subdirectory intents
if [ -n "$SUBDIRS" ]; then
    for intent_dir in $SUBDIRS; do
        INTENT_NAME=$(basename "$intent_dir")
        
        # Map intent name to expected screen directory
        # browse_timeline -> timeline, burst_management -> burst, etc.
        SCREEN_DIR_NAME=$(echo "$INTENT_NAME" | sed 's/_intent$//' | sed 's/browse_//' | sed 's/_management$//')
        
        # Check if intent uses screens (has activeScreen or *Screen fields)
        USES_SCREENS=false
        for go_file in "$intent_dir"/*.go; do
            if [ -f "$go_file" ]; then
                if grep -q "screens\.Screen\|\..*Screen\s" "$go_file" 2>/dev/null; then
                    USES_SCREENS=true
                    break
                fi
            fi
        done
        
        if [ "$USES_SCREENS" = true ]; then
            # Check for corresponding screens directory
            SCREEN_DIR="internal/cli/screens/$SCREEN_DIR_NAME"
            
            if [ ! -d "$SCREEN_DIR" ]; then
                # Try alternative naming conventions
                ALT_SCREEN_DIR="internal/cli/screens/$INTENT_NAME"
                if [ ! -d "$ALT_SCREEN_DIR" ]; then
                    echo -e "${YELLOW}⚠️  WARNING: Intent uses screens but no screens directory found${NC}"
                    echo "   Intent: $INTENT_NAME"
                    echo "   Expected: $SCREEN_DIR or $ALT_SCREEN_DIR"
                    echo "   Recommendation: Extract screen components to screens package"
                    echo ""
                    WARNINGS=$((WARNINGS+1))
                fi
            fi
        fi
    done
fi

# Check for screens with modals subdirectory (recommended pattern)
SCREEN_DIRS=$(find internal/cli/screens -mindepth 1 -maxdepth 1 -type d 2>/dev/null | grep -v base || true)

if [ -n "$SCREEN_DIRS" ]; then
    for screen_dir in $SCREEN_DIRS; do
        SCREEN_NAME=$(basename "$screen_dir")
        
        # Check if screen has modal files but no modals subdirectory
        MODAL_FILES=$(find "$screen_dir" -maxdepth 1 -name "*modal*.go" -type f 2>/dev/null | wc -l)
        MODALS_SUBDIR="$screen_dir/modals"
        
        if [ "$MODAL_FILES" -gt 2 ] && [ ! -d "$MODALS_SUBDIR" ]; then
            echo -e "${YELLOW}⚠️  WARNING: Multiple modal files without modals subdirectory${NC}"
            echo "   Screen: $SCREEN_NAME ($MODAL_FILES modal files)"
            echo "   Recommendation: Create $MODALS_SUBDIR and move modal files"
            echo ""
            WARNINGS=$((WARNINGS+1))
        fi
    done
fi

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ Screens structure looks good${NC}"
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ No screen violations (warnings present)${NC}"
fi

echo ""

# ============================================
# SUMMARY
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Total checks run: 23"
echo "Violations: $VIOLATIONS"
echo "Warnings: $WARNINGS"
echo ""

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL ARCHITECTURE CHECKS PASSED${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  PASSED with $WARNINGS warnings${NC}"
    echo ""
    echo "Warnings indicate technical debt or pending migrations."
    echo "Consider addressing to improve code quality."
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
else
    echo -e "${RED}❌ FAILED: $VIOLATIONS violations must be fixed${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  $WARNINGS warnings${NC}"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Fix these violations before committing."
    echo "See: docs/checklists/INTENT_DEVELOPMENT_CHECKLIST.md"
    echo ""
    exit 1
fi
