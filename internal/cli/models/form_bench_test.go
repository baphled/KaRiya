package models

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// BenchmarkFormView measures the performance of form rendering
func BenchmarkFormView(b *testing.B) {
	repo := careerrepo.NewMemoryRepository()
	svc := careerservice.NewService(repo)
	cliSvc := service.NewCLIEventService(svc)
	form := NewFormModel(cliSvc)

	// Fill form with data
	form.inputs[0].SetValue("Test event for performance benchmarking")
	form.inputs[1].SetValue("today")
	form.inputs[2].SetValue("TestCorp")
	form.inputs[3].SetValue("Performance Testing")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = form.View()
	}
}

// BenchmarkFormUpdate measures the performance of form updates
func BenchmarkFormUpdate(b *testing.B) {
	repo := careerrepo.NewMemoryRepository()
	svc := careerservice.NewService(repo)
	cliSvc := service.NewCLIEventService(svc)
	form := NewFormModel(cliSvc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.inputs[0].SetValue("Test")
	}
}

// BenchmarkCharacterCountTracking measures character count update performance
func BenchmarkCharacterCountTracking(b *testing.B) {
	repo := careerrepo.NewMemoryRepository()
	svc := careerservice.NewService(repo)
	cliSvc := service.NewCLIEventService(svc)
	form := NewFormModel(cliSvc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.charCount = (form.charCount + 1) % 2001
	}
}

// BenchmarkFieldErrorManagement measures error map operations
func BenchmarkFieldErrorManagement(b *testing.B) {
	repo := careerrepo.NewMemoryRepository()
	svc := careerservice.NewService(repo)
	cliSvc := service.NewCLIEventService(svc)
	form := NewFormModel(cliSvc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		form.fieldErrors[TextField] = "Test error"
		delete(form.fieldErrors, TextField)
	}
}
