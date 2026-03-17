package harness

import (
	"time"

	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/tui/app"
	"github.com/baphled/kariya/internal/tui/intents/burst_management"
	"github.com/baphled/kariya/internal/tui/intents/captureevent"
	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
	"github.com/baphled/kariya/internal/tui/intents/skillsmanagement"
)

// SubmitEvent sends a SubmitMsg directly to the model with the given event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitEvent(event *career.Event) *TestEnv {
	e.T.Helper()

	return e.SendMessage(captureevent.SubmitMsg{Event: event, Error: nil})
}

// SubmitEventWithError sends a SubmitMsg with an error to trigger validation error handling.
//
// Expected:
//   - err should describe the validation failure.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitEventWithError(event *career.Event, err error) *TestEnv {
	e.T.Helper()

	return e.SendMessage(captureevent.SubmitMsg{Event: event, Error: err})
}

// SubmitSkill creates a skill in the repository and sends a SkillCreatedMsg.
// This bypasses the UI form submission path, similar to SubmitEvent.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Creates the skill in the repository.
func (e *TestEnv) SubmitSkill(skill *career.Skill) *TestEnv {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	err := skillRepo.Create(e.Ctx, skill)

	return e.SendMessage(skillsmanagement.SkillCreatedMsg{Skill: skill, Error: err})
}

// SubmitSkillWithError sends a SkillCreatedMsg with an error to trigger validation error handling.
//
// Expected:
//   - err should describe the validation failure.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitSkillWithError(skill *career.Skill, err error) *TestEnv {
	e.T.Helper()

	return e.SendMessage(skillsmanagement.SkillCreatedMsg{Skill: skill, Error: err})
}

// SubmitSkillUpdate updates a skill in the repository and sends a SkillUpdatedMsg.
// This bypasses the UI form submission path for skill editing.
//
// Expected:
//   - skill must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the skill in the repository.
func (e *TestEnv) SubmitSkillUpdate(skill *career.Skill) *TestEnv {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	err := skillRepo.Update(e.Ctx, skill)

	return e.SendMessage(skillsmanagement.SkillUpdatedMsg{Skill: skill, Error: err})
}

// SubmitSkillUpdateWithError sends a SkillUpdatedMsg with an error.
//
// Expected:
//   - err should describe the update failure.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SubmitSkillUpdateWithError(skill *career.Skill, err error) *TestEnv {
	e.T.Helper()

	return e.SendMessage(skillsmanagement.SkillUpdatedMsg{Skill: skill, Error: err})
}

// SubmitFact creates a fact in the repository and sends a FactSavedMsg.
// This bypasses the UI form submission path for fact creation.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Creates the fact in the repository.
//
//nolint:dupl // Acceptable duplication in test helpers - similar to SubmitSkill pattern
func (e *TestEnv) SubmitFact(fact *career.Fact) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	err := factRepo.Create(e.Ctx, fact)
	if err != nil {
		return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: true, Message: err.Error()})
	}

	return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: true, Message: "Fact saved successfully"})
}

// SubmitFactUpdate updates a fact in the repository and sends a FactSavedMsg.
// This bypasses the UI form submission path for fact editing.
//
// Expected:
//   - fact must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the fact in the repository.
//
//nolint:dupl // Acceptable duplication in test helpers - similar to SubmitSkill pattern
func (e *TestEnv) SubmitFactUpdate(fact *career.Fact) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	err := factRepo.Update(e.Ctx, fact)
	if err != nil {
		return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: false, Message: err.Error()})
	}

	return e.SendMessage(factmanagement.FactSavedMsg{Fact: fact, IsNew: false, Message: "Fact updated successfully"})
}

// SubmitBurstUpdate updates a burst in the repository and sends a BurstEditCompleteMsg.
// This bypasses the UI form submission path for burst editing.
//
// Expected:
//   - burst must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the burst in the repository.
func (e *TestEnv) SubmitBurstUpdate(burst *career.Burst) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	err := burstRepo.Update(e.Ctx, burst)

	return e.SendMessage(burst_management.BurstEditCompleteMsg{Burst: burst, Cancelled: false, Error: err})
}

// ConfirmBurst confirms a burst and shows the loading modal for fact extraction.
// This bypasses the UI confirmation modal but triggers the loading state that tests expect.
//
// Expected:
//   - burst must be valid with existing ID.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - Updates the burst as confirmed in the repository.
//   - Shows loading modal (fact extraction runs but modal stays visible for test assertions).
func (e *TestEnv) ConfirmBurst(burst *career.Burst) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	// Bypass service ConfirmBurst (which triggers async fact extraction)
	// Just update the burst directly in the repository
	burst.Confirmed = true
	now := time.Now()
	burst.ConfirmedAt = &now
	burst.UpdatedAt = now

	err := burstRepo.Update(e.Ctx, burst)
	if err != nil {
		e.T.Fatalf("failed to update burst: %v", err)
	}

	// Don't send any message - tests will check repository state
	return e
}

// DismissSuccessModal bypasses the auto-dismiss countdown and immediately
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) DismissSuccessModal() *TestEnv {
	e.T.Helper()

	return e.SendMessage(captureevent.DismissModalMsg{})
}

// AssertEventCount verifies the number of events in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertEventCount(expected int) *TestEnv {
	e.T.Helper()

	events, err := e.Service.ListEvents(e.Ctx, careerrepo.EventListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get events: %v", err)
	}

	if len(events) != expected {
		e.T.Errorf("expected %d events, got %d", expected, len(events))
	}

	return e
}

// AssertBurstCount verifies the number of bursts in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertBurstCount(expected int) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	bursts, err := burstRepo.List(e.Ctx, careerrepo.BurstListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get bursts: %v", err)
	}

	if len(bursts) != expected {
		e.T.Errorf("expected %d bursts, got %d", expected, len(bursts))
	}

	return e
}

// AssertFactCount verifies the number of facts in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertFactCount(expected int) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	facts, err := factRepo.List(e.Ctx, careerrepo.FactListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get facts: %v", err)
	}

	if len(facts) != expected {
		e.T.Errorf("expected %d facts, got %d", expected, len(facts))
	}

	return e
}

// GetEvents returns all events from the database.
//
// Returns:
//   - A []*career.Event value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetEvents() []*career.Event {
	e.T.Helper()

	events, err := e.Service.ListEvents(e.Ctx, careerrepo.EventListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get events: %v", err)
	}

	return events
}

// GetBursts returns all bursts from the database.
//
// Returns:
//   - A []*career.Burst value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetBursts() []*career.Burst {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	bursts, err := burstRepo.List(e.Ctx, careerrepo.BurstListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get bursts: %v", err)
	}

	return bursts
}

// GetFacts returns all facts from the database.
//
// Returns:
//   - A []*career.Fact value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetFacts() []*career.Fact {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	facts, err := factRepo.List(e.Ctx, careerrepo.FactListFilters{})
	if err != nil {
		e.T.Fatalf("failed to get facts: %v", err)
	}

	return facts
}

// GetSkills returns all skills from the database.
//
// Returns:
//   - A []*career.Skill value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetSkills() []*career.Skill {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	skills, err := skillRepo.List(e.Ctx, nil)
	if err != nil {
		e.T.Fatalf("failed to get skills: %v", err)
	}

	return skills
}

// AssertSkillCount verifies the number of skills in the database.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertSkillCount(expected int) *TestEnv {
	e.T.Helper()

	skills := e.GetSkills()

	if len(skills) != expected {
		e.T.Errorf("expected %d skills, got %d", expected, len(skills))
	}

	return e
}

// AddSkill creates a skill in the database.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddSkill(skill *career.Skill) *TestEnv {
	e.T.Helper()

	skillRepo := e.Service.GetSkillRepository()
	if skillRepo == nil {
		e.T.Fatal("skill repository not set")
	}

	err := skillRepo.Create(e.Ctx, skill)
	if err != nil {
		e.T.Fatalf("failed to create skill: %v", err)
	}

	return e
}

// SimulateRestart recreates the application model while preserving the database.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SimulateRestart() *TestEnv {
	e.T.Helper()

	if e.DB == nil {
		e.T.Fatal("SimulateRestart requires SQLite persistence (use Setup, not SetupWithMemory)")
	}

	// Create new repositories pointing to the same database using ORM
	repos, err := careersql.NewRepositories(e.DB)
	if err != nil {
		e.T.Fatalf("failed to create repositories: %v", err)
	}
	eventRepo := repos.Event
	burstRepo := repos.Burst
	factRepo := repos.Fact
	skillRepo := repos.Skill

	// Create new service
	svc := careerservice.NewService(eventRepo)
	svc.SetBurstRepository(burstRepo)
	svc.SetFactRepository(factRepo)
	svc.SetSkillRepository(skillRepo)

	// Create new CLI service
	cliService := service.NewCLIEventService(svc)

	// Create bootstrap result (skipping onboarding for tests)
	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	// Create new application model
	model := app.NewModel(cliService, svc, bootstrapResult)

	// Update environment
	e.EventRepo = eventRepo
	e.BurstRepo = burstRepo
	e.FactRepo = factRepo
	e.SkillRepo = skillRepo
	e.Service = svc
	e.CLIService = cliService
	e.Model = model

	return e
}

// AddEvent creates an event in the database.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddEvent(event *career.Event) *TestEnv {
	e.T.Helper()

	eventRepo := e.Service.GetEventRepository()
	if eventRepo == nil {
		e.T.Fatal("event repository not set")
	}

	err := eventRepo.Create(e.Ctx, event)
	if err != nil {
		e.T.Fatalf("failed to create event: %v", err)
	}

	return e
}

// AddBurst creates a burst in the database.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddBurst(burst *career.Burst) *TestEnv {
	e.T.Helper()

	burstRepo := e.Service.GetBurstRepository()
	if burstRepo == nil {
		e.T.Fatal("burst repository not set")
	}

	err := burstRepo.Create(e.Ctx, burst)
	if err != nil {
		e.T.Fatalf("failed to create burst: %v", err)
	}

	return e
}

// AddFact creates a fact in the database.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AddFact(fact *career.Fact) *TestEnv {
	e.T.Helper()

	factRepo := e.Service.GetFactRepository()
	if factRepo == nil {
		e.T.Fatal("fact repository not set")
	}

	err := factRepo.Create(e.Ctx, fact)
	if err != nil {
		e.T.Fatalf("failed to create fact: %v", err)
	}

	return e
}
