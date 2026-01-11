//go:build ignore

package main

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("=== FORM SUBMISSION FLOW TEST ===")

	// Create necessary services
	repo, err := career.NewSQLiteRepository(":memory:")
	if err != nil {
		fmt.Printf("ERROR creating repository: %v\n", err)
		return
	}

	svc := careerservice.NewService(repo)
	cliService := service.NewCLIEventService(svc)

	// Create intent context
	ctx := &intents.CaptureEventContext{
		CaptureStrategy: "manual",
		CLIEventService: cliService,
		CareerService:   svc,
	}

	// Create intent
	intent, err := intents.NewCaptureEventIntent(ctx)
	if err != nil {
		fmt.Printf("ERROR creating intent: %v\n", err)
		return
	}

	fmt.Println("✅ Intent created")
	fmt.Printf("   Initial state: %v\n", intent.GetState())

	// Initialize intent
	fmt.Println("\n--- STEP 1: Initialize Intent ---")
	cmd := intent.Init()
	if cmd != nil {
		result := cmd()
		fmt.Printf("Init() returned command that emitted: %T\n", result)
	}
	fmt.Printf("State after Init: %v\n", intent.GetState())

	// Transition to form state
	fmt.Println("\n--- STEP 2: Transition to Form State ---")
	intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	fmt.Printf("State after selecting strategy: %v\n", intent.GetState())

	// Now we're in form state. Let's try to submit with Ctrl+S
	fmt.Println("\n--- STEP 3: Press Ctrl+S (without entering data) ---")
	cmd = intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	fmt.Printf("Update(Ctrl+S) returned command: %v\n", cmd != nil)

	if cmd != nil {
		fmt.Println("   Executing command...")
		result := cmd()
		fmt.Printf("   Command emitted: %T\n", result)

		if submitMsg, ok := result.(models.SubmitMsg); ok {
			fmt.Printf("   SubmitMsg received:\n")
			fmt.Printf("     Event: %v\n", submitMsg.Event)
			fmt.Printf("     Error: %v\n", submitMsg.Err)

			// Now pass the SubmitMsg back to the intent
			fmt.Println("\n--- STEP 4: Pass SubmitMsg back to intent ---")
			cmd2 := intent.Update(result)
			fmt.Printf("Update(SubmitMsg) returned command: %v\n", cmd2 != nil)
			fmt.Printf("State after SubmitMsg: %v\n", intent.GetState())
		} else {
			fmt.Printf("   ERROR: Expected SubmitMsg but got %T: %v\n", result, result)
		}
	} else {
		fmt.Println("   ERROR: Ctrl+S did not return a command!")
	}

	// Now let's try with actual form data
	fmt.Println("=== ATTEMPT 2: WITH FORM DATA ===")

	intent2, _ := intents.NewCaptureEventIntent(ctx)
	intent2.Init()
	intent2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

	fmt.Println("--- STEP 1: Fill in form data ---")
	formModel := intent2.GetForm()

	if formModel == nil {
		fmt.Println("ERROR: Form model is nil!")
		return
	}

	fmt.Printf("Form created\n")

	// Try to set form data by sending text input messages
	fmt.Println("Sending text input to form...")
	for _, ch := range "Test event" {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
		intent2.Update(msg)
	}

	// Note: GetInputValue is not available via interface, form values are accessed differently
	fmt.Printf("Text input sent to form\n")

	// Send Tab to move to next field
	fmt.Println("Sending Tab to move to date field...")
	intent2.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Send date
	for _, ch := range "2025-01-03" {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
		intent2.Update(msg)
	}

	fmt.Printf("Date input sent to form\n")

	fmt.Println("\n--- STEP 2: Press Ctrl+S with data ---")
	cmd = intent2.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	fmt.Printf("Update(Ctrl+S) returned command: %v\n", cmd != nil)

	if cmd != nil {
		fmt.Println("   Executing command...")
		result := cmd()
		fmt.Printf("   Command emitted: %T\n", result)

		if submitMsg, ok := result.(models.SubmitMsg); ok {
			fmt.Printf("   ✅ SubmitMsg received:\n")
			fmt.Printf("      Event: %v\n", submitMsg.Event)
			fmt.Printf("      Error: %v\n", submitMsg.Err)

			if submitMsg.Err == nil && submitMsg.Event != nil {
				fmt.Printf("      Event Text: '%s'\n", submitMsg.Event.Text)
				fmt.Printf("      Event Date: %v\n", submitMsg.Event.Date)
			}

			fmt.Println("\n--- STEP 3: Pass SubmitMsg back to intent ---")
			cmd2 := intent2.Update(result)
			fmt.Printf("Update(SubmitMsg) returned command: %v\n", cmd2 != nil)
			fmt.Printf("State after SubmitMsg: %v\n", intent2.GetState())

			if intent2.GetState() == "review" {
				fmt.Println("✅ SUCCESS: Form submission worked! State is now 'review'")
			} else {
				fmt.Printf("❌ PROBLEM: State is '%v' but should be 'review'\n", intent2.GetState())
			}
		} else {
			fmt.Printf("❌ ERROR: Expected SubmitMsg but got %T\n", result)
		}
	} else {
		fmt.Println("❌ ERROR: Ctrl+S did not return a command!")
	}
}
