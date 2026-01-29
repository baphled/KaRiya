# KaRiya Architecture Diagrams

**Version**: 1.0
**Last Updated**: 2026-01-28
**Status**: Current

This document provides comprehensive visual diagrams of the KaRiya architecture. For detailed implementation guidance, see:
- [TUI Architecture](./TUI_ARCHITECTURE.md)
- [Intent Diagram](../TUI_INTENT_DIAGRAM.md)
- [UIKit Guide](../UIKIT_GUIDE.md)

---

## Table of Contents

1. [System Architecture Overview](#1-system-architecture-overview)
2. [Layer Dependency Flow](#2-layer-dependency-flow)
3. [Data Flow Diagram](#3-data-flow-diagram)
4. [Intent Lifecycle](#4-intent-lifecycle)
5. [Screen/Modal Hierarchy](#5-screenmodal-hierarchy)
6. [Service Layer Architecture](#6-service-layer-architecture)
7. [Domain Model](#7-domain-model)
8. [UIKit Component Tree](#8-uikit-component-tree)
9. [Package Structure](#9-package-structure)

---

## 1. System Architecture Overview

High-level view of all layers and their relationships:

```mermaid
flowchart TB
    subgraph Entry["Entry Point"]
        CMD[cmd/cli/main.go]
    end

    subgraph AppLayer["Application Layer"]
        APP[app.Model]
        ROUTER[IntentRouter]
        REG[IntentRegistrar]
    end

    subgraph Intents["Intents Layer (State Machines)"]
        direction LR
        BT[BrowseTimeline]
        SM[SkillsManagement]
        FM[FactManagement]
        BM[BurstManagement]
        CE[CaptureEvent]
        GCV[GenerateCV]
        CS[ConfigureSystem]
    end

    subgraph Screens["Screens Layer (Stateless Views)"]
        direction TB
        subgraph ScreenBase["Base Screens"]
            BS[BaseScreen]
            FS[FormScreen]
            DS[DetailScreen]
            PS[ProgressScreen]
            SS[SelectScreen]
            CFS[ConfirmScreen]
        end
        subgraph ScreenFeatures["Feature Screens"]
            TL[timeline/]
            SK[skills/]
            CV[cv/]
            CAP[capture/]
            CFG[configure/]
            EXP[export/]
        end
        subgraph Modals["Feature Modals"]
            MOD[modals/]
        end
    end

    subgraph UIKit["UIKit Component Library"]
        direction TB
        subgraph Primitives
            BTN[Button]
            BDG[Badge]
            TXT[Text]
            INP[Input]
        end
        subgraph Containers
            BOX[Box]
            OVL[Overlay]
        end
        subgraph Feedback
            MDL[Modal]
            HLP[HelpModal]
            CFM[ConfirmModal]
        end
        subgraph Layout
            SCL[ScreenLayout]
            HDR[Header]
            FTR[Footer]
        end
    end

    subgraph Behaviors["Behaviors Layer"]
        TB[TableBehavior]
        SR[ScreenResult]
        MH[ModalHelpers]
    end

    subgraph Forms["Forms Layer"]
        FRM[forms.NewInput/Select/Form]
        FSPEC[CaptureEventForm<br/>FactForm<br/>SkillForm<br/>BurstForm]
    end

    subgraph Themes["Theme System"]
        THM[themes.Theme]
        STY[Styles]
    end

    subgraph CLIService["CLI Services"]
        EVTSVC[EventService]
    end

    subgraph CareerService["Career Services"]
        direction TB
        CSVC[Service]
        subgraph SubServices
            CVSVC[cv/]
            BFSVC[burstfact/]
            TECHSVC[technology/]
            CLSSVC[classification/]
        end
    end

    subgraph Repository["Repository Layer"]
        direction LR
        REPO[Repository Interface]
        SQLITE[SQLite Implementation]
        MIG[Migrations]
    end

    subgraph Domain["Domain Layer"]
        direction LR
        EVENT[Event]
        FACT[Fact]
        SKILL[Skill]
        BURST[Burst]
        CVD[CV]
    end

    subgraph Support["Supporting Services"]
        BOOT[Bootstrap]
        NAV[Navigation]
        CFG2[Config]
        LOG[Logger]
    end

    %% Flow connections
    CMD --> BOOT
    BOOT --> APP
    APP --> ROUTER
    ROUTER --> REG
    ROUTER --> Intents

    %% Intent to Screen flow
    Intents --> Screens
    Screens --> UIKit
    Screens --> Behaviors
    Screens --> Forms

    %% Behaviors uses UIKit
    Behaviors --> UIKit

    %% Theme integration
    UIKit --> Themes
    Screens --> Themes
    Forms --> Themes

    %% Service layer
    Intents --> CLIService
    Intents --> CareerService
    CLIService --> CareerService
    CSVC --> SubServices

    %% Persistence
    CareerService --> Repository
    Repository --> SQLITE
    SQLITE --> MIG

    %% Domain
    Repository --> Domain
    CareerService --> Domain
    Intents --> Domain

    %% Support services
    APP --> NAV
    APP --> LOG
    BOOT --> CFG2

    %% Styling
    classDef entryStyle fill:#e1f5fe,stroke:#01579b
    classDef appStyle fill:#fff3e0,stroke:#e65100
    classDef intentStyle fill:#f3e5f5,stroke:#7b1fa2
    classDef screenStyle fill:#e8f5e9,stroke:#2e7d32
    classDef uikitStyle fill:#fce4ec,stroke:#c2185b
    classDef serviceStyle fill:#e3f2fd,stroke:#1565c0
    classDef repoStyle fill:#fff8e1,stroke:#f9a825
    classDef domainStyle fill:#efebe9,stroke:#5d4037

    class CMD entryStyle
    class APP,ROUTER,REG appStyle
    class BT,SM,FM,BM,CE,GCV,CS intentStyle
    class BS,FS,DS,PS,SS,CFS,TL,SK,CV,CAP,CFG,EXP,MOD screenStyle
    class BTN,BDG,TXT,INP,BOX,OVL,MDL,HLP,CFM,SCL,HDR,FTR uikitStyle
    class CSVC,CVSVC,BFSVC,TECHSVC,CLSSVC,EVTSVC serviceStyle
    class REPO,SQLITE,MIG repoStyle
    class EVENT,FACT,SKILL,BURST,CVD domainStyle
```

---

## 2. Layer Dependency Flow

Simplified view showing allowed import directions:

```mermaid
graph TD
    A[App/Router] --> B[Intents]
    B --> C[Screens]
    C --> D[UIKit]
    C --> E[Behaviors]
    C --> F[Forms]
    B --> G[Services]
    G --> H[Repository]
    H --> I[Domain]
    
    D --> J[Themes]
    E --> D
    F --> J
    F -.->|wraps| K[huh library]

    style A fill:#fff3e0
    style B fill:#f3e5f5
    style C fill:#e8f5e9
    style D fill:#fce4ec
    style G fill:#e3f2fd
    style H fill:#fff8e1
    style I fill:#efebe9
```

### Dependency Rules

| Layer | Can Import | NEVER Import |
|-------|------------|--------------|
| `intents/` | screens, uikit, behaviors, forms, services | - |
| `screens/` | uikit, behaviors, forms | **intents** |
| `uikit/` | themes only | **screens, intents** |
| `behaviors/` | uikit, themes | **screens, intents** |
| `forms/` | huh (only here), themes | screens, intents |

---

## 3. Data Flow Diagram

How data flows through the application during a typical user interaction:

```mermaid
flowchart LR
    subgraph User["User Interaction"]
        KEY[Keyboard Input]
    end

    subgraph TUI["TUI Layer"]
        direction TB
        APP[App Model]
        ROUTER[Intent Router]
        INTENT[Active Intent]
        SCREEN[Active Screen]
    end

    subgraph Business["Business Layer"]
        direction TB
        SVC[Career Service]
        BFSVC[BurstFact Service]
        CVSVC[CV Service]
    end

    subgraph Persistence["Persistence Layer"]
        REPO[Repository]
        DB[(SQLite DB)]
    end

    subgraph Output["Output"]
        VIEW[Rendered View]
        EXPORT[Exported Files]
    end

    KEY --> APP
    APP --> ROUTER
    ROUTER --> INTENT
    INTENT --> SCREEN
    
    SCREEN -->|ScreenResult| INTENT
    INTENT -->|IntentResult| ROUTER
    ROUTER -->|Navigate/Complete| APP
    
    INTENT -->|CRUD Operations| SVC
    SVC --> BFSVC
    SVC --> CVSVC
    SVC --> REPO
    REPO --> DB
    
    DB -->|Domain Entities| REPO
    REPO -->|Domain Entities| SVC
    SVC -->|Domain Entities| INTENT
    INTENT -->|Display Data| SCREEN
    
    SCREEN --> VIEW
    CVSVC --> EXPORT

    style KEY fill:#e3f2fd
    style VIEW fill:#e8f5e9
    style DB fill:#fff8e1
```

### Message Flow Detail

```mermaid
sequenceDiagram
    participant U as User
    participant A as App
    participant R as Router
    participant I as Intent
    participant S as Screen
    participant SVC as Service
    participant DB as Database

    U->>A: KeyMsg (e.g., 'enter')
    A->>R: HandleMessage(msg)
    R->>I: Update(msg)
    I->>S: Update(msg)
    S-->>I: ScreenResult (Navigate/Submit)
    
    alt Submit Result
        I->>SVC: SaveEvent(event)
        SVC->>DB: Insert/Update
        DB-->>SVC: Success
        SVC-->>I: Domain Entity
    end
    
    I-->>R: tea.Cmd / IntentResult
    R-->>A: Update state
    A->>A: View()
    A-->>U: Rendered output
```

---

## 4. Intent Lifecycle

How intents manage state and transitions:

```mermaid
stateDiagram-v2
    direction TB
    
    [*] --> Inactive: Intent Created
    Inactive --> Initializing: ActivateIntent()
    Initializing --> Active: Init() returns
    
    state Active {
        direction LR
        [*] --> ScreenActive
        ScreenActive --> ScreenActive: Update(msg)
        ScreenActive --> ModalActive: Open Modal
        ModalActive --> ScreenActive: Close Modal
        ScreenActive --> Transitioning: ScreenResult
        Transitioning --> ScreenActive: New Screen
    }
    
    Active --> Completing: IntentResult ready
    Completing --> [*]: Result returned
    
    Active --> Cancelling: Escape/Cancel
    Cancelling --> [*]: Cancelled result
```

### Intent State Pattern (Reference: BrowseTimeline)

```mermaid
stateDiagram-v2
    direction TB
    
    [*] --> StateList: Init()
    
    StateList --> StateDetail: Select Event
    StateList --> StateList: Filter/Sort
    StateList --> [*]: Escape
    
    StateDetail --> StateList: Back
    StateDetail --> StateDelete: Delete Action
    StateDetail --> StateEdit: Edit Action
    
    StateEdit --> StateDetail: Save/Cancel
    StateDelete --> StateList: Confirm/Cancel
    
    note right of StateList
        ListScreen active
        TableBehavior manages items
    end note
    
    note right of StateDetail
        DetailScreen active
        Shows selected event
    end note
```

---

## 5. Screen/Modal Hierarchy

Organization of screens and modals by feature:

```mermaid
flowchart TB
    subgraph Base["Base Screens (screens/base/)"]
        direction LR
        BS[BaseScreen]
        FS[FormScreen]
        DS[DetailScreen]
        PS[ProgressScreen]
        SS[SelectScreen]
        CS[ConfirmScreen]
    end

    subgraph Timeline["Timeline (screens/timeline/)"]
        direction TB
        EL[EventListScreen]
        ED[EventDetailScreen]
        EDC[EventDeleteConfirmScreen]
        
        subgraph TLModals["modals/"]
            FM[FilterModal]
            SM[SortModal]
            SEM[SearchModal]
            EM[EditModal]
            DM[DetailModal]
            QAM[QuickAddModal]
        end
    end

    subgraph Skills["Skills (screens/skills/)"]
        direction TB
        SL[SkillListScreen]
        SDS[SkillDetailScreen]
        SFS[SkillFormScreen]
        SDelS[SkillDeleteScreen]
        
        subgraph SKModals["modals/"]
            SKF[FilterModal]
            SKS[SortModal]
            SKSE[SearchModal]
            SKAE[AddEditModal]
            SKD[DetailModal]
            SKE[EventsModal]
        end
    end

    subgraph CV["CV Generation (screens/cv/)"]
        direction TB
        AS[AudienceSelectScreen]
        PRS[ProfileSelectScreen]
        GS[GeneratingScreen]
        PVS[PreviewScreen]
        RS[ReviewScreen]
    end

    subgraph Capture["Capture (screens/capture/)"]
        direction TB
        STRS[StrategySelectScreen]
        EFS[EventFormScreen]
        ERS[EventReviewScreen]
        ESS[EventSubmitScreen]
    end

    subgraph Configure["Configure (screens/configure/)"]
        direction TB
        DSS[DomainSelectScreen]
        ESS2[EditSettingsScreen]
        RCS[ReviewChangesScreen]
        COS[ConfirmScreen]
        CMS[CompleteScreen]
        FLS[FailedScreen]
    end

    %% Inheritance
    BS --> EL
    BS --> ED
    BS --> SL
    DS --> SDS
    FS --> SFS
    SS --> AS
    SS --> PRS
    PS --> GS
    SS --> STRS
    FS --> EFS

    style Base fill:#e3f2fd
    style Timeline fill:#e8f5e9
    style Skills fill:#fff3e0
    style CV fill:#f3e5f5
    style Capture fill:#fce4ec
    style Configure fill:#efebe9
```

---

## 6. Service Layer Architecture

Detailed view of the service layer:

```mermaid
flowchart TB
    subgraph CLI["CLI Service Layer"]
        CLISVC[CLIEventService]
    end

    subgraph Career["Career Service (internal/service/career/)"]
        SVC[Service]
        
        subgraph CV["cv/"]
            CVGEN[CVGenerationService]
            BULLET[BulletGenerator]
            EXPORT[ExportService]
            PROFILE[ProfileInference]
            SECTION[SectionBuilder]
            CONFIG[ConfigManager]
            TRACE[TraceabilityService]
        end
        
        subgraph BurstFact["burstfact/"]
            DETECT[Detector]
            EXTRACT[Extractor]
            CLASSIFY[Classifier]
            TEMPORAL[TemporalGrouper]
            SIMILAR[SimilarityScorer]
        end
        
        subgraph Tech["technology/"]
            TECHEXT[Extractor]
            FOCUS[FocusArea]
        end
        
        subgraph Class["classification/"]
            EVTCLASS[Classifier]
        end
    end

    subgraph Repo["Repository Layer"]
        EVENTREPO[EventRepository]
        FACTREPO[FactRepository]
        BURSTREPO[BurstRepository]
        SKILLREPO[SkillRepository]
    end

    CLISVC --> SVC
    SVC --> CV
    SVC --> BurstFact
    SVC --> Tech
    SVC --> Class
    
    SVC --> EVENTREPO
    SVC --> FACTREPO
    SVC --> BURSTREPO
    SVC --> SKILLREPO
    
    CVGEN --> BULLET
    CVGEN --> SECTION
    CVGEN --> PROFILE
    CVGEN --> EXPORT
    
    DETECT --> TEMPORAL
    DETECT --> SIMILAR
    EXTRACT --> CLASSIFY

    style CLI fill:#e1f5fe
    style Career fill:#e8f5e9
    style Repo fill:#fff8e1
```

### Service Operations

```mermaid
flowchart LR
    subgraph EventOps["Event Operations"]
        CE[CreateEvent]
        UE[UpdateEvent]
        DE[DeleteEvent]
        LE[ListEvents]
        FE[FilterEvents]
    end

    subgraph FactOps["Fact Operations"]
        CF[CreateFact]
        UF[UpdateFact]
        DF[DeleteFact]
        LF[ListFacts]
        EF[ExtractFacts]
    end

    subgraph BurstOps["Burst Operations"]
        CB[CreateBurst]
        UB[UpdateBurst]
        DB[DeleteBurst]
        LB[ListBursts]
        DETB[DetectBursts]
    end

    subgraph SkillOps["Skill Operations"]
        CS[CreateSkill]
        US[UpdateSkill]
        DS[DeleteSkill]
        LS[ListSkills]
        EXTS[ExtractSkills]
    end

    subgraph CVOps["CV Operations"]
        GEN[GenerateCV]
        PREV[PreviewCV]
        EXPRT[ExportCV]
        VARS[GetVariants]
    end

    style EventOps fill:#e3f2fd
    style FactOps fill:#e8f5e9
    style BurstOps fill:#fff3e0
    style SkillOps fill:#f3e5f5
    style CVOps fill:#fce4ec
```

---

## 7. Domain Model

Core domain entities and their relationships:

```mermaid
erDiagram
    CareerEvent ||--o{ Fact : "generates"
    CareerEvent ||--o{ Burst : "belongs to"
    CareerEvent }o--o{ Skill : "demonstrates"
    Burst ||--o{ Fact : "generates"
    Fact }o--o{ CVView : "included in"
    CareerEvent }o--o{ CVView : "source for"
    Skill }o--o{ CVView : "featured in"

    CareerEvent {
        string ID PK
        string Text
        datetime Date
        string Company
        string Project
        array Tags
        array Categories
        array SkillIDs FK
        datetime CreatedAt
        datetime UpdatedAt
    }

    Fact {
        string ID PK
        string Text
        array CompetencyCategories
        string RoleFit
        array AudienceRelevance
        string StrengthSignal
        string SourceEventID FK
        string SourceBurstID FK
        datetime CreatedAt
        datetime UpdatedAt
    }

    Burst {
        string ID PK
        string Name
        string Description
        array EventIDs FK
        boolean Confirmed
        datetime ConfirmedAt
        datetime CreatedAt
        datetime UpdatedAt
    }

    Skill {
        string ID PK
        string Name
        string Category
        string Level
        int YearsUsed
        datetime LastUsed
        datetime CreatedAt
        datetime UpdatedAt
    }

    CVView {
        string ID PK
        string Name
        string TargetRole
        string TargetAudience
        map EventFilters
        datetime GeneratedAt
        int SourceEventCount
        int SourceFactCount
        array Sections
    }
```

### Domain Concepts

```mermaid
flowchart TB
    subgraph Core["Core Entities"]
        EVENT[CareerEvent<br/>Professional milestone]
        SKILL[Skill<br/>Technology/competency]
    end

    subgraph Derived["Derived Entities"]
        BURST[Burst<br/>Related event grouping]
        FACT[Fact<br/>Extracted competency]
    end

    subgraph Generated["Generated Artifacts"]
        CV[CVView<br/>Generated CV]
        SECTION[CVSection<br/>CV section]
    end

    EVENT -->|"grouped into"| BURST
    EVENT -->|"facts extracted"| FACT
    BURST -->|"facts extracted"| FACT
    EVENT -->|"demonstrates"| SKILL
    
    FACT -->|"assembled into"| CV
    EVENT -->|"source for"| CV
    SKILL -->|"featured in"| CV
    CV -->|"contains"| SECTION

    style Core fill:#e3f2fd
    style Derived fill:#e8f5e9
    style Generated fill:#f3e5f5
```

---

## 8. UIKit Component Tree

Organization of UIKit components:

```mermaid
flowchart TB
    subgraph UIKit["uikit/"]
        direction TB
        
        subgraph Primitives["primitives/"]
            TEXT[Text<br/>Title, Body, Label]
            BTN[Button]
            BTNG[ButtonGroup]
            BADGE[Badge<br/>HelpKeyBadge]
            INPUT[Input]
            KV[KeyValue]
        end
        
        subgraph Containers["containers/"]
            BOX[Box]
            OVERLAY[Overlay]
        end
        
        subgraph Feedback["feedback/"]
            MODAL[Modal]
            MODALC[ModalContainer]
            HELP[HelpModal]
            CONFIRM[ConfirmModal]
            DETAIL[DetailModal]
            INFO[InfoModal]
        end
        
        subgraph Layout["layout/"]
            SCREEN[ScreenLayout]
            HEADER[Header]
            FOOTER[Footer]
        end
        
        subgraph Navigation["navigation/"]
            BREAD[Breadcrumb]
        end
        
        subgraph Selectors["selectors/"]
            CAT[CategorySelector]
            TAG[TagSelector]
            SKILLSEL[SkillSelector]
            AUD[AudienceRelevanceSelector]
        end
        
        subgraph Display["display/"]
            LOGO[Logo]
        end
        
        subgraph Widgets["widgets/"]
            DETAILV[DetailView]
        end
        
        subgraph Theme["theme/"]
            THM[Theme]
            AWARE[Aware]
        end
    end

    %% Composition relationships
    SCREEN --> HEADER
    SCREEN --> FOOTER
    SCREEN --> BREAD
    
    MODAL --> BOX
    MODAL --> OVERLAY
    
    HELP --> MODAL
    CONFIRM --> MODAL
    DETAIL --> MODAL
    INFO --> MODAL
    
    MODALC --> MODAL
    
    DETAILV --> KV
    DETAILV --> TEXT

    style Primitives fill:#fce4ec
    style Containers fill:#e3f2fd
    style Feedback fill:#fff3e0
    style Layout fill:#e8f5e9
    style Selectors fill:#f3e5f5
    style Theme fill:#efebe9
```

---

## 9. Package Structure

Directory layout of the codebase:

```mermaid
flowchart TB
    subgraph Root["KaRiya/"]
        CMD["cmd/<br/>Entry points"]
        INTERNAL["internal/<br/>Core application"]
        DOCS["docs/<br/>Documentation"]
        SCRIPTS["scripts/<br/>Build tools"]
    end

    subgraph Internal["internal/"]
        CLI["cli/<br/>TUI application"]
        DOMAIN["domain/<br/>Domain entities"]
        SERVICE["service/<br/>Business logic"]
        REPO["repository/<br/>Data access"]
        CONFIG["config/<br/>Configuration"]
        LOGGER["logger/<br/>Logging"]
        CONSTANTS["constants/<br/>Constants"]
        TESTUTIL["testutil/<br/>Test utilities"]
    end

    subgraph CLIDetail["cli/"]
        APP["app/<br/>Application model"]
        INTENTS["intents/<br/>State machines"]
        SCREENS["screens/<br/>View components"]
        UIKIT["uikit/<br/>UI primitives"]
        BEHAVIORS["behaviors/<br/>Reusable behaviors"]
        FORMS["forms/<br/>Form builders"]
        MODELS["models/<br/>(deprecated)"]
        COMPONENTS["components/<br/>(legacy)"]
        THEMES["themes/<br/>Theme system"]
        NAV["navigation/<br/>Key mappings"]
        BOOTSTRAP["bootstrap/<br/>Initialization"]
        TERMINAL["terminal/<br/>Terminal utils"]
        IMPORTER["importer/<br/>Data import"]
    end

    subgraph IntentsDetail["intents/"]
        BTDIR["browsetimeline/"]
        SMDIR["skillsmanagement/"]
        FMDIR["factmanagement/"]
        LEGACY["*_intent.go<br/>(legacy single-file)"]
    end

    subgraph ScreensDetail["screens/"]
        BASE["base/<br/>Base screens"]
        TIMELINE["timeline/"]
        SKILLS["skills/"]
        CVS["cv/"]
        CAPTURE["capture/"]
        CONFIGURE["configure/"]
        EXPORT["export/"]
    end

    CMD --> INTERNAL
    INTERNAL --> CLI
    CLI --> CLIDetail
    INTENTS --> IntentsDetail
    SCREENS --> ScreensDetail

    style Root fill:#e1f5fe
    style Internal fill:#e8f5e9
    style CLIDetail fill:#fff3e0
    style IntentsDetail fill:#f3e5f5
    style ScreensDetail fill:#fce4ec
```

---

## Quick Reference

### Current Intent Implementations

| Intent | Location | Pattern | Status |
|--------|----------|---------|--------|
| BrowseTimeline | `intents/browsetimeline/` | Subdirectory (reference) | Modern |
| SkillsManagement | `intents/skillsmanagement/` | Subdirectory | Modern |
| FactManagement | `intents/factmanagement/` | Subdirectory | Modern |
| BurstManagement | `intents/burst_management*.go` | Single file | Legacy |
| CaptureEvent | `intents/capture_event*.go` | Single file | Legacy |
| GenerateCV | `intents/generate_cv*.go` | Single file | Legacy |
| ConfigureSystem | `intents/configure_system*.go` | Single file | Legacy |

### Key File Patterns

| Type | Location Pattern | Example |
|------|-----------------|---------|
| Intent (modern) | `intents/{feature}/intent.go` | `intents/browsetimeline/intent.go` |
| Screen | `screens/{feature}/{type}.go` | `screens/timeline/event_list.go` |
| Modal | `screens/{feature}/modals/{action}_modal.go` | `screens/timeline/modals/filter_modal.go` |
| Form | `forms/{entity}_form.go` | `forms/skill_form.go` |
| Service | `service/career/{domain}/` | `service/career/cv/` |

---

## Related Documentation

- [TUI Architecture](./TUI_ARCHITECTURE.md) - Detailed architecture guide
- [Intent Diagram](../TUI_INTENT_DIAGRAM.md) - Intent state flows
- [UIKit Guide](../UIKIT_GUIDE.md) - Component usage
- [Screen Specification](./SCREEN_SPECIFICATION.md) - Screen contracts
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md) - Intent patterns
