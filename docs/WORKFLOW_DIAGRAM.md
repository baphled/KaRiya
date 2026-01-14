---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
```mermaid
%%{init: {'theme': 'base', 'themeVariables': {
    'primaryColor': '#2C3E50',
    'primaryTextColor': '#ECF0F1',
    'primaryBorderColor': '#34495E',
    'lineColor': '#7F8C8D',
    'secondaryColor': '#3498DB',
    'tertiaryColor': '#2980B9'
}}}%%

flowchart TD
    %% Main Menu
    MainMenu["Main Menu"] --> IntentStartCapture["Intent: Start Capture"]
    MainMenu --> IntentViewTimeline["Intent: View Timeline"]
    MainMenu --> IntentGenerateCV["Intent: Generate CV"]

    %% Event Capture Flow
    IntentStartCapture --> EventCaptureMode["Choose Mode"]
    EventCaptureMode --> InputData["Input Data"]
    InputData --> ReviewFacts["Review Facts"]
    ReviewFacts --> TagAndCategorise["Tag & Categorise"]
    TagAndCategorise --> ReviewSummary["Review Summary"]
    ReviewSummary --> ConfirmSave["Confirm Save"]

    %% Timeline View
    IntentViewTimeline --> TimelineView["Timeline View"]
    TimelineView --> IntentViewEvent["Intent: View Event"]
    IntentViewEvent --> EventDetail["Event Detail"]

    %% CV Composition Flow
    IntentGenerateCV --> SelectConfig["Select Config"]
    SelectConfig --> SelectAudience["Select Audience"]
    SelectAudience --> SelectFormat["Select Format"]
    SelectFormat --> Preview["Preview"]
    Preview --> ReviewEdit["Review/Edit"]
    ReviewEdit --> ConfirmExport["Confirm Export"]

    %% Styling
    classDef mainMenu fill:#2C3E50,color:#ECF0F1;
    classDef flow fill:#3498DB,color:#FFFFFF;
    classDef leaf fill:#2980B9,color:#FFFFFF;

    class MainMenu mainMenu;
    class EventCaptureMode,SelectConfig flow;
    class TimelineView,EventDetail,ConfirmExport leaf;
```

# KaRiya Workflow Diagram

## Workflow Overview

The diagram illustrates the primary user workflows in the KaRiya Career Journal CLI application, highlighting the key user intents and their corresponding flows.

### Main Menu Intents
1. **Start Capture**: Initiates the event capture workflow
2. **View Timeline**: Allows browsing of existing career events
3. **Generate CV**: Starts the CV composition process

### Event Capture Flow
- Choose capture strategy (Quick or Manual)
- Input event data
- Review extracted facts
- Tag and categorize events
- Review summary
- Confirm and save event

### Timeline View
- Browse timeline of events
- Select and view individual event details

### CV Composition Flow
- Select configuration
- Choose target audience
- Select export format
- Preview generated CV
- Review and edit
- Confirm export

## Design Principles
- Modular workflow design
- Clear user navigation
- Extensible architecture
- User-centric feature progression

## Workflow Notes
- Each workflow is designed to be independent yet interconnected
- Provides clear user guidance
- Supports multiple capture and generation modes

