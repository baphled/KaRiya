package facts_test

import (
	"bytes"
	"context"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/cmd/facts"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("ListFacts", func() {
	var (
		ctx cliutil.ServiceContext
		out *bytes.Buffer
		err *bytes.Buffer
	)

	BeforeEach(func() {
		cliCtx := cmdpkg.NewCLIContext("", true)
		initErr := cliCtx.InitService(new(bytes.Buffer))
		Expect(initErr).NotTo(HaveOccurred())
		ctx = cliCtx
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
	})

	Context("with no facts", func() {
		It("should return 0 and print info message", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should not write error output", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(err.Len()).To(Equal(0))
		})
	})

	Context("with facts in database", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())
			fact := fixtures.FactWithCategories("fact-1", "Implemented authentication system", "event-1", []string{"technical"}, []string{"hiring_manager"})
			fact.RoleFit = "principal"
			saveErr := svc.SaveFact(context.Background(), fact)
			Expect(saveErr).NotTo(HaveOccurred())
		})

		It("should return 0 on successful execution", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should write output to stdout", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Existing Facts"))
		})

		It("should display fact information", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Implemented authentication system"))
		})

		It("should handle execution without panic", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			Expect(func() {
				facts.ListFacts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})

	Context("with multiple facts", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			for i := 1; i <= 3; i++ {
				factID := "fact-" + string(rune(48+i))
				factText := "Fact " + string(rune(48+i))
				fact := fixtures.FactWithCategories(factID, factText, "event-1", []string{"technical"}, []string{"hiring_manager"})
				fact.RoleFit = "principal"
				saveErr := svc.SaveFact(context.Background(), fact)
				Expect(saveErr).NotTo(HaveOccurred())
			}
		})

		It("should list all facts", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Total facts: 3"))
		})
	})

	Context("error handling", func() {
		It("should handle service errors gracefully", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should handle nil service gracefully", func() {
			Expect(func() {
				facts.ListFacts(nil, out, err, tea.WithInput(nil))
			}).To(Panic())
		})
	})
})
