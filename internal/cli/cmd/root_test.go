package cmd_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careerservice "github.com/baphled/kariya/internal/service/career"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Root Command", func() {
	Describe("NewRootCmd", func() {
		It("creates a root command with correct properties", func() {
			cmd := cmdpkg.NewRootCmd("1.0.0")
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("kariya"))
			Expect(cmd.Short).To(ContainSubstring("Career Event Capture"))
			Expect(cmd.Version).To(Equal("1.0.0"))
		})

		It("has all required subcommands", func() {
			cmd := cmdpkg.NewRootCmd("1.0.0")
			subcommands := cmd.Commands()
			subcommandNames := make([]string, len(subcommands))
			for i, sub := range subcommands {
				subcommandNames[i] = sub.Name()
			}
			Expect(subcommandNames).To(ContainElement("import"))
			Expect(subcommandNames).To(ContainElement("bursts"))
			Expect(subcommandNames).To(ContainElement("facts"))
			Expect(subcommandNames).To(ContainElement("skills"))
		})

		It("has persistent flags for db and mode", func() {
			cmd := cmdpkg.NewRootCmd("1.0.0")
			dbFlag := cmd.PersistentFlags().Lookup("db")
			modeFlag := cmd.PersistentFlags().Lookup("mode")
			inMemoryFlag := cmd.PersistentFlags().Lookup("in-memory")

			Expect(dbFlag).NotTo(BeNil())
			Expect(modeFlag).NotTo(BeNil())
			Expect(inMemoryFlag).NotTo(BeNil())
		})

		It("sets SilenceUsage to true", func() {
			cmd := cmdpkg.NewRootCmd("1.0.0")
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})

	Describe("NewImportCmd", func() {
		It("creates an import command with correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewImportCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("import <file>"))
			Expect(cmd.Short).To(ContainSubstring("Import career events"))
		})

		It("requires exactly one argument", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewImportCmd(ctx)
			Expect(cmd.Args).NotTo(BeNil())
			// Verify it's ExactArgs by checking the function behavior
			err := cmd.Args(cmd, []string{})
			Expect(err).To(HaveOccurred())
		})

		It("has RunE function", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewImportCmd(ctx)
			Expect(cmd.RunE).NotTo(BeNil())
		})

		It("returns error when import fails", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())

			cmd := cmdpkg.NewImportCmd(ctx)
			// Try to import a non-existent file
			err := cmd.RunE(cmd, []string{"/nonexistent/file.csv"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("import failed"))
		})
	})

	Describe("NewBurstsCmd", func() {
		It("creates a bursts command with correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewBurstsCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("bursts"))
			Expect(cmd.Short).To(ContainSubstring("Manage bursts"))
		})

		It("has detect and list subcommands", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewBurstsCmd(ctx)
			subcommands := cmd.Commands()
			subcommandNames := make([]string, len(subcommands))
			for i, sub := range subcommands {
				subcommandNames[i] = sub.Name()
			}
			Expect(subcommandNames).To(ContainElement("detect"))
			Expect(subcommandNames).To(ContainElement("list"))
		})

		It("detect subcommand has correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewBurstsCmd(ctx)
			detectCmd, _, _ := cmd.Find([]string{"detect"})
			Expect(detectCmd).NotTo(BeNil())
			Expect(detectCmd.Use).To(Equal("detect"))
			Expect(detectCmd.Short).To(ContainSubstring("Detect bursts"))
		})

		It("list subcommand has correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewBurstsCmd(ctx)
			listCmd, _, _ := cmd.Find([]string{"list"})
			Expect(listCmd).NotTo(BeNil())
			Expect(listCmd.Use).To(Equal("list"))
			Expect(listCmd.Short).To(ContainSubstring("List all bursts"))
		})
	})

	Describe("NewFactsCmd", func() {
		It("creates a facts command with correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewFactsCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("facts"))
			Expect(cmd.Short).To(ContainSubstring("Manage facts"))
		})

		It("has extract and list subcommands", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewFactsCmd(ctx)
			subcommands := cmd.Commands()
			subcommandNames := make([]string, len(subcommands))
			for i, sub := range subcommands {
				subcommandNames[i] = sub.Name()
			}
			Expect(subcommandNames).To(ContainElement("extract"))
			Expect(subcommandNames).To(ContainElement("list"))
		})

		It("extract subcommand has correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewFactsCmd(ctx)
			extractCmd, _, _ := cmd.Find([]string{"extract"})
			Expect(extractCmd).NotTo(BeNil())
			Expect(extractCmd.Use).To(Equal("extract"))
			Expect(extractCmd.Short).To(ContainSubstring("Extract facts"))
		})

		It("list subcommand has correct properties", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			cmd := cmdpkg.NewFactsCmd(ctx)
			listCmd, _, _ := cmd.Find([]string{"list"})
			Expect(listCmd).NotTo(BeNil())
			Expect(listCmd.Use).To(Equal("list"))
			Expect(listCmd.Short).To(ContainSubstring("List all facts"))
		})
	})

	Describe("Command RunE execution", func() {
		var (
			ctx      *cmdpkg.CLIContext
			progress cliutil.MockProgressRunner
		)

		BeforeEach(func() {
			ctx = &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())
			progress = cliutil.MockProgressRunner{}
		})

		It("bursts detect command executes without error", func() {
			// Mock successful progress runner
			progress.RunWithSpinnerFn = func(message string, fn func() error, opts ...tea.ProgramOption) error {
				return fn()
			}

			cmd := cmdpkg.NewBurstsCmd(ctx)
			detectCmd, _, _ := cmd.Find([]string{"detect"})

			// Re-inject a handler with mock runner
			detectCmd.RunE = func(cmd *cobra.Command, args []string) error {
				code := cmdpkg.DetectBursts(ctx.Service, cmd.OutOrStdout(), cmd.ErrOrStderr(), progress)
				if code != 0 {
					return errors.New("burst detection failed")
				}
				return nil
			}

			err := detectCmd.RunE(detectCmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("bursts list command executes without error", func() {
			cmd := cmdpkg.NewBurstsCmd(ctx)
			listCmd, _, _ := cmd.Find([]string{"list"})
			err := listCmd.RunE(listCmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("bursts detect command returns error when detection fails", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())

			// Mock the service to return an error
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error")).
				AnyTimes()

			mockService := careerservice.NewService(mockEventRepo)
			ctx.Service = mockService

			cmd := cmdpkg.NewBurstsCmd(ctx)
			detectCmd, _, _ := cmd.Find([]string{"detect"})

			// Use actual handler which now uses DefaultProgressRunner
			// In test environment, DefaultProgressRunner should be replaced or we use our mock
			detectCmd.RunE = func(cmd *cobra.Command, args []string) error {
				code := cmdpkg.DetectBursts(ctx.Service, cmd.OutOrStdout(), cmd.ErrOrStderr(), progress)
				if code != 0 {
					return errors.New("burst detection failed")
				}
				return nil
			}

			err := detectCmd.RunE(detectCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst detection failed"))
		})

		It("bursts list command returns error when listing fails", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())

			cmd := cmdpkg.NewBurstsCmd(ctx)
			listCmd, _, _ := cmd.Find([]string{"list"})

			// Mock the service to return an error
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockBurstRepo := mockrepo.NewMockBurstRepository(ctrl)
			mockBurstRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error")).
				Times(1)

			mockService := careerservice.NewService(mockEventRepo)
			mockService.SetBurstRepository(mockBurstRepo)
			ctx.Service = mockService

			err := listCmd.RunE(listCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("listing bursts failed"))
		})

		It("facts extract command executes without error", func() {
			cmd := cmdpkg.NewFactsCmd(ctx)
			extractCmd, _, _ := cmd.Find([]string{"extract"})

			// Re-inject mock runner
			extractCmd.RunE = func(cmd *cobra.Command, args []string) error {
				code := cmdpkg.ExtractFacts(ctx.Service, cmd.OutOrStdout(), cmd.ErrOrStderr(), progress)
				if code != 0 {
					return errors.New("fact extraction failed")
				}
				return nil
			}

			err := extractCmd.RunE(extractCmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("facts list command executes without error", func() {
			cmd := cmdpkg.NewFactsCmd(ctx)
			listCmd, _, _ := cmd.Find([]string{"list"})
			err := listCmd.RunE(listCmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("facts extract command returns error when extraction fails", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())

			// Mock the service to return an error
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error")).
				AnyTimes()

			mockService := careerservice.NewService(mockEventRepo)
			ctx.Service = mockService

			cmd := cmdpkg.NewFactsCmd(ctx)
			extractCmd, _, _ := cmd.Find([]string{"extract"})

			// Re-inject mock runner
			extractCmd.RunE = func(cmd *cobra.Command, args []string) error {
				code := cmdpkg.ExtractFacts(ctx.Service, cmd.OutOrStdout(), cmd.ErrOrStderr(), progress)
				if code != 0 {
					return errors.New("fact extraction failed")
				}
				return nil
			}

			err := extractCmd.RunE(extractCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact extraction failed"))
		})

		It("facts list command returns error when listing fails", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())

			cmd := cmdpkg.NewFactsCmd(ctx)
			listCmd, _, _ := cmd.Find([]string{"list"})

			// Mock the service to return an error
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockFactRepo := mockrepo.NewMockFactRepository(ctrl)
			mockFactRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error")).
				Times(1)

			mockService := careerservice.NewService(mockEventRepo)
			mockService.SetFactRepository(mockFactRepo)
			ctx.Service = mockService

			err := listCmd.RunE(listCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("listing facts failed"))
		})
	})
})
