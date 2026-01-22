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
# SUMMARY
# ============================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $VIOLATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL ARCHITECTURE CHECKS PASSED${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
elif [ $VIOLATIONS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  $WARNINGS WARNING(S) - Consider fixing${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
else
    echo -e "${RED}❌ $VIOLATIONS VIOLATION(S) FOUND${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  $WARNINGS WARNING(S)${NC}"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Fix these violations before committing."
    echo "See: docs/INTENT_ARCHITECTURE_GUIDE.md"
    exit 1
fi
