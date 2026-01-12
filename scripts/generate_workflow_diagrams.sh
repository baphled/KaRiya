#!/bin/bash
# Generate workflow diagrams for KaRiya TUI documentation
# 
# This script generates Mermaid diagrams for the two most complex workflows:
# 1. CV Generation (10 states)
# 2. Event Capture (4 states + 3 modals)
#
# Diagrams are extracted from the actual implementation code to ensure accuracy.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DIAGRAMS_DIR="$PROJECT_ROOT/docs/workflows/diagrams"

# Create diagrams directory if it doesn't exist
mkdir -p "$DIAGRAMS_DIR"

echo "================================================"
echo "🎨 GENERATING WORKFLOW DIAGRAMS"
echo "================================================"
echo ""
echo "Output directory: $DIAGRAMS_DIR"
echo ""

# Function to generate CV Generation diagram
generate_cv_diagram() {
    echo "📊 Generating CV Generation workflow diagram..."
    
    cat > "$DIAGRAMS_DIR/cv_generation_flow.mermaid" << 'EOF'
%%{ init: { 'theme': 'base', 'themeVariables': { 'fontSize': '16px' } } }%%
graph TD
    Start([Start: Generate CV])
    
    SelectProfile[1. Select Profile]
    SelectAudience[2. Select Audience]
    Generating[3. Generating CV...]
    Preview[4. Preview CV]
    Review[5. Review/Edit CV]
    Confirm[6. Confirm]
    
    ExportFormat[7. Select Export Format]
    ExportLocation[8. Select Save Location]
    Exporting[9. Exporting...]
    ExportComplete[10. Export Complete]
    
    Complete([Complete])
    Cancel([Cancel])
    
    Start --> SelectProfile
    
    SelectProfile -->|Enter: Select| SelectAudience
    SelectProfile -->|Esc/m: Cancel| Cancel
    
    SelectAudience -->|Enter: Generate| Generating
    SelectAudience -->|Esc: Back| SelectProfile
    SelectAudience -->|m: Main Menu| Cancel
    
    Generating -->|Success| Preview
    Generating -->|Error| SelectAudience
    Generating -->|Esc: Let complete| SelectAudience
    
    Preview -->|e/Enter: Edit| Review
    Preview -->|c: Continue| Confirm
    Preview -->|Esc: Back| SelectAudience
    
    Review -->|Enter: Continue| Confirm
    Review -->|Esc: Back| Preview
    
    Confirm -->|y/Enter: Complete| Complete
    Confirm -->|e/x: Export| ExportFormat
    Confirm -->|n/Esc: Back| Review
    
    ExportFormat -->|Enter: Select| ExportLocation
    ExportFormat -->|Esc: Back| Confirm
    
    ExportLocation -->|Enter: Export| Exporting
    ExportLocation -->|Cancel: Back| Confirm
    ExportLocation -->|Esc: Back| ExportFormat
    
    Exporting -->|Success| ExportComplete
    Exporting -->|Error| ExportComplete
    Exporting -->|Esc: Let complete| ExportLocation
    
    ExportComplete -->|Enter: Done| Complete
    ExportComplete -->|Esc: Retry| ExportLocation
    
    style Generating fill:#ffd700,stroke:#333,stroke-width:2px
    style Exporting fill:#ffd700,stroke:#333,stroke-width:2px
    style Complete fill:#90EE90,stroke:#333,stroke-width:2px
    style Cancel fill:#FFB6C1,stroke:#333,stroke-width:2px
EOF
    
    echo "✅ CV Generation diagram created: cv_generation_flow.mermaid"
}

# Function to generate Event Capture diagram
generate_capture_diagram() {
    echo "📊 Generating Event Capture workflow diagram..."
    
    cat > "$DIAGRAMS_DIR/event_capture_flow.mermaid" << 'EOF'
%%{ init: { 'theme': 'base', 'themeVariables': { 'fontSize': '16px' } } }%%
graph TD
    Start([Start: Capture Event])
    
    ChooseStrategy[1. Choose Strategy<br/>Quick or Manual]
    Form[2. Event Form<br/>Huh-based input]
    Review[3. Review Inferred Event]
    Submit[4. Submit Event]
    
    EditMetadata[Edit Metadata Modal]
    EditBursts[Edit Bursts Modal]
    EditFacts[Edit Facts Modal]
    
    Success([Event Saved])
    Cancel([Cancel])
    
    Start --> ChooseStrategy
    
    ChooseStrategy -->|Enter: Select Quick| Form
    ChooseStrategy -->|Enter: Select Manual| Form
    ChooseStrategy -->|Esc/m: Cancel| Cancel
    
    Form -->|Ctrl+S/Submit: Quick path| Submit
    Form -->|Enter: Review first| Review
    Form -->|Esc: Back| ChooseStrategy
    Form -->|Ctrl+O: Toggle fields<br/>Manual mode only| Form
    
    Review -->|Ctrl+S/Enter: Confirm| Submit
    Review -->|e: Edit metadata| EditMetadata
    Review -->|b: Edit bursts| EditBursts
    Review -->|f: Edit facts| EditFacts
    Review -->|Esc: Back| Form
    
    EditMetadata -->|Enter: Save| Review
    EditMetadata -->|Esc: Cancel| Review
    
    EditBursts -->|Enter: Save| Review
    EditBursts -->|Esc: Cancel| Review
    
    EditFacts -->|Enter: Save| Review
    EditFacts -->|Esc: Cancel| Review
    
    Submit -->|Success| Success
    Submit -->|Error| Submit
    Submit -->|r: Retry| Submit
    Submit -->|Esc: Back with error| Review
    
    style Form fill:#87CEEB,stroke:#333,stroke-width:2px
    style Submit fill:#ffd700,stroke:#333,stroke-width:2px
    style Success fill:#90EE90,stroke:#333,stroke-width:2px
    style Cancel fill:#FFB6C1,stroke:#333,stroke-width:2px
    style EditMetadata fill:#DDA0DD,stroke:#333,stroke-width:2px
    style EditBursts fill:#DDA0DD,stroke:#333,stroke-width:2px
    style EditFacts fill:#DDA0DD,stroke:#333,stroke-width:2px
EOF
    
    echo "✅ Event Capture diagram created: event_capture_flow.mermaid"
}

# Generate both diagrams
generate_cv_diagram
echo ""
generate_capture_diagram
echo ""

echo "================================================"
echo "✅ DIAGRAM GENERATION COMPLETE"
echo "================================================"
echo ""
echo "Generated diagrams:"
echo "  1. $DIAGRAMS_DIR/cv_generation_flow.mermaid"
echo "  2. $DIAGRAMS_DIR/event_capture_flow.mermaid"
echo ""
echo "View diagrams:"
echo "  • On GitHub (automatic Mermaid rendering)"
echo "  • In VS Code with Mermaid preview extension"
echo "  • Online: https://mermaid.live"
echo ""
echo "Integration:"
echo "  • CV Generation: docs/workflows/CV_GENERATION_WORKFLOW.md"
echo "  • Event Capture: docs/workflows/EVENT_CAPTURE_WORKFLOW.md"
echo ""
