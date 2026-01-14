---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Product Requirements Document: KaRiya Career Journal & CV Generator

## 1. Introduction/Overview

KaRiya is a web-based platform designed to help senior engineers, consultants, and technical professionals create a low-friction, trustworthy method of capturing and generating career narratives. The platform addresses the challenge of maintaining a comprehensive, accurate, and easily manageable career history that can be transformed into tailored CVs for different roles and audiences.

## 2. Goals

1. Enable users to capture career events with minimal friction
2. Provide an intelligent system for grouping and contextualizing career events
3. Generate role and audience-specific CV views
4. Maintain a single source of truth for career history
5. Ensure transparency and traceability in CV generation

## 3. User Stories

1. As a senior engineer, I want to quickly log my professional achievements so that I can maintain an up-to-date career journal without spending excessive time on documentation.

2. As a consultant, I want to generate tailored CVs for different clients and roles, highlighting the most relevant experiences and achievements.

3. As a technical professional, I want to see how my career events are grouped and interpreted, so I can understand the broader context of my professional journey.

4. As a job seeker, I want to easily export my CV in multiple formats while maintaining the integrity of my original career events.

## 4. Functional Requirements

1. Career Event Capture
   1.1. Support multiple input modes:
     - Timeline journaling
     - CV backfill (importing existing CVs)
     - Manual event entry
     - Metadata-driven event classification
   1.2. Allow users to add optional metadata (date, company, project, tags)
   1.3. Validate input to ensure data quality (e.g., date constraints, text length)

2. Burst Generation
   2.1. Automatically suggest groupings of related career events
   2.2. Allow user confirmation and manual editing of event groupings
   2.3. Generate inferred facts from event groups

3. CV Generation
   3.1. Support CV generation for multiple roles:
     - Principal
     - Staff Engineer
     - Engineering Manager
     - Senior Individual Contributor
   3.2. Generate CV views for different audiences:
     - Hiring Managers
     - Recruiters
     - Peers
   3.3. Implement bullet caps and compression rules per role
   3.4. Ensure all CV bullets are traceable to source events

4. Metadata and Fact Extraction
   4.1. Infer competencies from career events
   4.2. Classify role fit and audience relevance
   4.3. Provide strength signals for achievements

5. Export Capabilities
   5.1. Support multiple export formats:
     - PDF
     - Markdown
     - YAML
   5.2. Maintain traceability of exported content to original events

## 5. Non-Goals (Out of Scope)

1. Automatic job application submission
2. Social networking features
3. Real-time collaboration
4. Advanced graphic design for CVs
5. Integration with external job boards (potential future enhancement)

## 6. Design Considerations

1. Web application with responsive design
2. Clean, minimalist interface focusing on content capture and generation
3. Mobile-friendly layout
4. Intuitive event input and CV generation workflow

## 7. Technical Considerations

1. Performance requirements:
   - CV generation ≤ 2 seconds for ≤500 events
   - Support ≥10,000 events per user

2. Security and Reliability:
   - Encrypt data at rest and in transit
   - ≥99.5% availability
   - Audit trail for all edits
   - Option to rollback changes

3. Potential future integrations:
   - LinkedIn CV import
   - Optional HR system connections

## 8. Success Metrics

1. Primary Metric: Number of career events captured
2. Secondary Metrics:
   - CV generation frequency
   - User engagement rate
   - Accuracy of fact extraction (qualitative assessment)

## 9. Open Questions

1. What additional export formats might be valuable?
2. How granular should the metadata classification be?
3. What level of AI-assisted fact extraction is desired?

## 10. Acceptance Criteria

1. Users can capture career events with minimal friction
2. CV generation respects role and audience-specific rules
3. All CV bullets must trace to at least one source event
4. Bursts cannot be saved with fewer than two events
5. Invalid inputs (tags, dates) are rejected with clear feedback
6. Metadata edits propagate correctly across the system
7. No aspirational language or unlinked metrics in generated content
8. CV generation must match the following structure:
    ```yaml
    ---
    first_name: Yomi
    last_name: Colledge
    email: yomi@boodah.net
    location: Remote (UK)
    phone: +44 7894 987 855
    links:
      - label: LinkedIn
        url: https://www.linkedin.com/in/yomicolledge
      - label: GitHub
        url: http://github.com/baphled
      - label: Portfolio
        url: http://boodah.net
    summary: |
      **Senior Ruby Engineer / Lead Developer / Principal Technical Consultant**

      Senior engineer with 15+ years’ experience delivering mission-critical Ruby systems across startups, scale-ups, and government platforms.

    highlights: |
      Acting as Principal Engineer and technical lead for an early-stage startup, owning architecture, delivery, and technical strategy.

    jobs:
      - company: "FullSpektrum®"
        position: "Principal Technical Consultant"
        start_date: "Jul 2024"
        end_date: "Nov 2025"
        description: |
          - Defined and executed end-to-end technical strategy and delivery roadmap for a startup, ensuring alignment with business objectives
          - Advised on technology selection and architecture to showcase the company’s core USP effectively
          - Designed MVP development and demonstration strategy to accelerate time-to-market
          - Prepared comprehensive technical due-diligence materials, supporting successful investor engagement
          - Developed and implemented a strategy for beta version development, optimising product feedback and iteration

      - company: "Mindful Chef"
        position: "Senior Ruby Developer"
        start_date: "Mar 2022"
        end_date: "Apr 2024"
        description: |
          - Engineered platform migration to Ruby 3.1, reducing server load and improving system performance
          - Designed and implemented innovative delivery flow, driving **30% increase in weekly delivery volume**
          - Re-architected complex Kafka-based pipelines, resulting in:
            - **45% reduction in delivery-box debt collection**
            - Eliminated payment failures
            - Minimised data loss risks
          - Automated critical workflows, generating **£132,000 annual operational cost savings**
          - Delivered cross-functional features spanning ordering, subscriptions, logistics, and CMS

      - company: "Spice Rack"
        position: "Ruby DevOps Consultant"
        start_date: "Jan 2022"
        end_date: "Mar 2022"
        description: |
          - Rapidly developed AWS-backed provisioning and content delivery services using Ruby + Docker
          - Compressed project timeline from 6 to 3 months while maintaining high-quality delivery
          - Implemented robust backup + rollback system, reducing infrastructure setup time by **80%**
          - Delivered production-ready infrastructure with minimal client overhead
          - Enabled seamless handover to internal teams for ongoing maintenance
          - Ensured high availability and scalability for anticipated traffic spikes

      - company: "Freelance Technical Consultant"
        position: "Senior Ruby Developer & Architect"
        start_date: "Oct 2020"
        end_date: "Dec 2021"
        description: |
          - Architected scale-able white-label SaaS platforms and high-performance APIs
          - Implemented ELK stack for real-time monitoring of multi-million event streams
          - Developed specialised tools including pension calculation systems and advanced publishing workflows
          - Designed REST APIs that enhanced user on-boarding and reduced client churn

      - company: "NTTData / BEIS / Defra / MAS"
        position: "Senior Ruby Development Consultant"
        start_date: "Jun 2019"
        end_date: "Oct 2020"
        description: |
          - Led migration of critical government platforms, ensuring reliability, compliance, and accessibility
          - Standardised Ruby development practices across multiple teams
          - Integrated complex REST and SOAP APIs for inter-agency data exchange
          - Stabilised large-scale public-sector services under high-traffic conditions
          - Rapidly integrated with cross-functional team to deliver strategic digital tools
          - Introduced defra-ruby-style guide (https://github.com/DEFRA/defra-ruby-style), enhancing code quality and consistency across teams
          - Implemented Dough open-source framework (https://github.com/moneyadviceservice/dough)
          - Developed and deployed Pensions Calculator Tool
          - Maintained internal CMS with high efficiency, consistently delivering ahead of schedule

      - company: "Diverse Technical Consulting Portfolio"
        position: "Senior Ruby Developer & Technical Architect"
        start_date: "Jun 2012"
        end_date: "Aug 2018"
        description: |
          - Established ELK stack solutions for advanced data analysis and operational insights
          - Architected flexible white-label platforms and high-performance APIs
          - Managed complex migrations, integrations, and infrastructure scaling across varied technology stacks
          - Demonstrated adaptability across multiple industries and technological challenges
          - Contributed to HSBC Global Connections changes and enhanced the Atlas project.
    projects:
      - name: "n-vyro.io"
        description: |
          Architected a comprehensive IoT ecosystem demonstrating end-to-end technical leadership.
          Developed custom hardware and software solutions, implementing CI/CD across 14+ repositories
          to ensure reliable, scalable product delivery. Showcased ability to manage complex,
          multi-domain technical projects from prototype to production.
        url: "https://n-vyro.io"
        key_achievements:
          - Managed 14+ interconnected repositories
          - Implemented robust CI/CD pipelines
          - Bridged hardware and software development domains

      - name: "QuikCV"
        description: |
          Open-source CV generation tool highlighting technical versatility and commitment to
          best practices. Migrated from Ruby to Vue.js, demonstrating technology adaptation
          and architectural evolution. Serves as a living portfolio of clean architecture
          and modern development practices.
        url: "https://github.com/boodah-consulting/quik-cv"
        key_achievements:
          - Successfully migrated between technology stacks
          - Maintained open-source project as a technical showcase
          - Demonstrated CLI tooling and architectural best practices
    skills:
      Languages:
        - Ruby
        - Go
        - C/C++
        - JavaScript / ECMAScript

      Frameworks:
        - Ruby on Rails
        - Sinatra
        - Grape
        - React
        - Vue.js
        - Quasar Framework

      "Data & Messaging":
        - Kafka
        - PostgreSQL
        - MySQL
        - Redis
        - ELK Stack
        - GraphQL
        - MongoDB
        - Elasticsearch

      "Cloud & DevOps":
        - AWS (EC2, S3, RDS, CloudFront)
        - Docker
        - CI/CD (Jenkins, Travis CI, GitHub Actions)
        - Semantic Versioning
        - Capistrano
        - Chef
        - Ansible
        - Logstash / Kibana
        - StatsD / Graphite

      "Engineering Practices":
        - Agile / Scrum / Kanban
        - TDD / BDD
        - API Design
        - Object-Oriented Programming
        - Component Driven Development
        - Pair Programming
        - Technical Mentoring
        - Strategic Consulting
        - Observability & Performance Optimisation
        - Platform Modernisation
        - Rescue Missions
        - Cost Reduction Strategies
        - Release Management
        - Infrastructure as Code
        - Cross-functional Collaboration
    ```

## 11. Phased Rollout Approach

### Phase 1: Core Functionality
- Basic event capture
- Manual event grouping
- Simple CV generation

### Phase 2: Advanced Features
- AI-assisted burst generation
- More sophisticated fact extraction
- Enhanced export capabilities

### Phase 3: Integrations and Advanced Analytics
- External system integrations
- Comprehensive user insights and metrics

