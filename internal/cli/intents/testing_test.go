package intents

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestingPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testing Package Suite")
}

var _ = Describe("TestIntentFactory", func() {
	Describe("NewTestIntentFactory", func() {
		It("creates a valid factory", func() {
			factory := NewTestIntentFactory()
			Expect(factory).NotTo(BeNil())
		})
	})

	Describe("Register and Create", func() {
		It("work together correctly", func() {
			factory := NewTestIntentFactory()
			mock := NewMockIntent()
			factory.Register("test", func() Intent {
				return mock
			})

			intent := factory.Create("test")
			Expect(intent).NotTo(BeNil())
			Expect(intent).To(Equal(mock))
		})
	})

	Describe("Create", func() {
		It("returns nil for unregistered intent", func() {
			factory := NewTestIntentFactory()
			intent := factory.Create("nonexistent")
			Expect(intent).To(BeNil())
		})

		It("supports multiple registered intents", func() {
			factory := NewTestIntentFactory()
			mock1 := NewMockIntent()
			mock2 := NewMockIntent()

			factory.Register("intent1", func() Intent { return mock1 })
			factory.Register("intent2", func() Intent { return mock2 })

			intent1 := factory.Create("intent1")
			intent2 := factory.Create("intent2")

			Expect(intent1).To(Equal(mock1))
			Expect(intent2).To(Equal(mock2))
		})
	})
})

var _ = Describe("IntentWithState", func() {
	Describe("NewIntentWithState", func() {
		It("creates a valid wrapper", func() {
			mock := NewMockIntent()
			state := "test state"
			wrapper := NewIntentWithState(mock, state)

			Expect(wrapper).NotTo(BeNil())
			Expect(wrapper.GetState()).To(Equal(state))
		})
	})

	Describe("Delegation", func() {
		It("delegates to wrapped intent", func() {
			mock := NewMockIntent()
			wrapper := NewIntentWithState(mock, "state")

			wrapper.Init()
			Expect(mock.initCalled).To(BeTrue())

			wrapper.Update(tea.KeyMsg{})
			Expect(mock.updateCalled).To(BeGreaterThan(0))

			wrapper.View()
			Expect(mock.viewCalled).To(BeTrue())

			wrapper.Result()
			// Result just returns nil from mock, which is fine
		})
	})

	Describe("SetState", func() {
		It("changes the state", func() {
			mock := NewMockIntent()
			wrapper := NewIntentWithState(mock, "initial")

			wrapper.SetState("updated")
			Expect(wrapper.GetState()).To(Equal("updated"))
		})
	})
})

var _ = Describe("IntentTestHarness", func() {
	var (
		harness *IntentTestHarness
		mock    *MockIntent
	)

	BeforeEach(func() {
		mock = NewMockIntent()
		harness = NewIntentTestHarness(GinkgoT(), mock)
	})

	Describe("NewIntentTestHarness", func() {
		It("creates a valid harness", func() {
			Expect(harness).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("calls the intent's Init method", func() {
			harness.Init()
			Expect(mock.initCalled).To(BeTrue())
		})
	})

	Describe("SendMessage", func() {
		It("calls the intent's Update method", func() {
			harness.SendMessage(tea.KeyMsg{})
			Expect(mock.updateCalled).To(BeGreaterThan(0))
		})
	})

	Describe("GetView", func() {
		It("returns the intent's view", func() {
			view := harness.GetView()
			Expect(view).To(Equal("Mock Intent View"))
		})
	})

	Describe("GetResult", func() {
		It("returns the intent's result", func() {
			expectedResult := &IntentResult[interface{}]{Status: Completed}
			mock.SetResult(expectedResult)
			harness = NewIntentTestHarness(GinkgoT(), mock)

			result := harness.GetResult()
			Expect(result).To(Equal(expectedResult))
		})
	})

	Describe("AssertViewContains", func() {
		It("passes with matching substring", func() {
			harness.AssertViewContains("Mock")
			harness.AssertViewContains("Intent")
		})
	})

	Describe("AssertViewNotContains", func() {
		It("passes with non-matching substring", func() {
			harness.AssertViewNotContains("goodbye")
		})
	})
})

var _ = Describe("IntentRouterTestHelper", func() {
	var (
		router *DefaultIntentRouter
		helper *IntentRouterTestHelper
	)

	BeforeEach(func() {
		router = NewDefaultIntentRouter()
		helper = NewIntentRouterTestHelper(GinkgoT(), router)
	})

	Describe("NewIntentRouterTestHelper", func() {
		It("creates a valid helper", func() {
			Expect(helper).NotTo(BeNil())
		})
	})

	Describe("ActivateIntent", func() {
		It("activates a registered intent", func() {
			mock := NewMockIntent()
			router.RegisterIntent("test", func() Intent { return mock })
			helper = NewIntentRouterTestHelper(GinkgoT(), router)

			cmd := helper.ActivateIntent("test", nil)
			Expect(cmd).To(BeNil()) // Init returns nil

			activeIntent := helper.GetActiveIntent()
			Expect(activeIntent).To(Equal(mock))
		})
	})

	Describe("AssertIntentActive", func() {
		It("passes when intent is active", func() {
			mock := NewMockIntent()
			router.RegisterIntent("test", func() Intent { return mock })
			helper = NewIntentRouterTestHelper(GinkgoT(), router)

			helper.ActivateIntent("test", nil)
			helper.AssertIntentActive()
		})
	})

	Describe("GetHistory", func() {
		It("returns the navigation history", func() {
			mock1 := NewMockIntent()
			mock2 := NewMockIntent()
			router.RegisterIntent("test1", func() Intent { return mock1 })
			router.RegisterIntent("test2", func() Intent { return mock2 })
			helper = NewIntentRouterTestHelper(GinkgoT(), router)

			helper.ActivateIntent("test1", nil)
			helper.ActivateIntent("test2", nil)

			history := helper.GetHistory()
			Expect(history).To(HaveLen(1))
			Expect(history[0]).To(Equal(mock1))
		})
	})

	Describe("AssertHistoryLength", func() {
		It("validates history length", func() {
			mock1 := NewMockIntent()
			mock2 := NewMockIntent()
			router.RegisterIntent("test1", func() Intent { return mock1 })
			router.RegisterIntent("test2", func() Intent { return mock2 })
			helper = NewIntentRouterTestHelper(GinkgoT(), router)

			helper.ActivateIntent("test1", nil)
			helper.AssertHistoryLength(0)

			helper.ActivateIntent("test2", nil)
			helper.AssertHistoryLength(1)
		})
	})

	Describe("GoBack", func() {
		It("navigates to previous intent", func() {
			mock1 := NewMockIntent()
			mock2 := NewMockIntent()
			router.RegisterIntent("test1", func() Intent { return mock1 })
			router.RegisterIntent("test2", func() Intent { return mock2 })
			helper = NewIntentRouterTestHelper(GinkgoT(), router)

			helper.ActivateIntent("test1", nil)
			helper.ActivateIntent("test2", nil)
			helper.GoBack()

			activeIntent := helper.GetActiveIntent()
			Expect(activeIntent).To(Equal(mock1))
		})
	})

	Describe("AssertGoBackFails", func() {
		It("fails when at root", func() {
			mock := NewMockIntent()
			router.RegisterIntent("test", func() Intent { return mock })
			helper = NewIntentRouterTestHelper(GinkgoT(), router)

			helper.ActivateIntent("test", nil)
			helper.AssertGoBackFails()
		})
	})
})

var _ = Describe("contains helper function", func() {
	It("returns true for matching substring", func() {
		Expect(contains("hello world", "hello")).To(BeTrue())
		Expect(contains("hello world", "world")).To(BeTrue())
		Expect(contains("hello world", "lo wo")).To(BeTrue())
	})

	It("returns false for non-matching substring", func() {
		Expect(contains("hello world", "goodbye")).To(BeFalse())
		Expect(contains("hello world", "HELLO")).To(BeFalse())
	})

	It("handles empty strings", func() {
		Expect(contains("hello", "")).To(BeTrue())
		Expect(contains("", "hello")).To(BeFalse())
		Expect(contains("", "")).To(BeTrue())
	})
})
