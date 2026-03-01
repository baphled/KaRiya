package configure_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/configure"
)

func TestConfigureScreens(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configure Screens Suite")
}

var _ = Describe("DomainSelectScreen", func() {
	var (
		screen  *configure.DomainSelectScreen
		domains []configtypes.ConfigurationDomain
	)

	BeforeEach(func() {
		domains = []configtypes.ConfigurationDomain{
			configtypes.DomainSystem,
			configtypes.DomainProfile,
			configtypes.DomainExport,
			configtypes.DomainUI,
		}
		screen = configure.NewDomainSelectScreen(domains)
		screen.SetTerminalInfo(120, 40)
	})

	Describe("Initialization", func() {
		It("should create screen with all domains", func() {
			Expect(screen).NotTo(BeNil())
		})

		It("should start with first domain selected", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("System"))
		})
	})

	Describe("View Rendering", func() {
		It("should show title", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Select Configuration Domain"))
		})

		It("should show all four domains", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("System"))
			Expect(view).To(ContainSubstring("Profile"))
			Expect(view).To(ContainSubstring("Export"))
			Expect(view).To(ContainSubstring("UI"))
		})
	})

	Describe("Navigation", func() {
		It("should handle down arrow", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(result).To(BeNil()) // Navigation doesn't return a result
		})

		It("should handle up arrow", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(result).To(BeNil())
		})
	})

	Describe("Selection", func() {
		It("should return NavigateResult with selected domain on Enter", func() {
			// Select first domain (system) with Enter
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal(configtypes.DomainSystem))
		})
	})

	Describe("Cancellation", func() {
		It("should return CancelResult on Esc", func() {
			_, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})
	})
})
