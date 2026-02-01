package violations

import "fake/career"

func helperWithInlineEvent() {
	_ = career.Event{ID: "1", Text: "test"} // want `inline career\.Event\{\} in test file; use fixtures package instead`
}

func helperWithInlineFact() {
	_ = career.Fact{ID: "1"} // want `inline career\.Fact\{\} in test file; use fixtures package instead`
}

func helperWithInlineBurst() {
	_ = career.Burst{ID: "1", Name: "burst"} // want `inline career\.Burst\{\} in test file; use fixtures package instead`
}

func helperWithInlineSkill() {
	_ = career.Skill{ID: "1", Name: "Go"} // want `inline career\.Skill\{\} in test file; use fixtures package instead`
}

func helperWithInlineCVView() {
	_ = career.CVView{ID: "1"} // want `inline career\.CVView\{\} in test file; use fixtures package instead`
}

func helperWithInlineCVSection() {
	_ = career.CVSection{ID: "1"} // want `inline career\.CVSection\{\} in test file; use fixtures package instead`
}

func helperWithInlineCVBullet() {
	_ = career.CVBullet{ID: "1"} // want `inline career\.CVBullet\{\} in test file; use fixtures package instead`
}

func helperWithInlineCVConfig() {
	_ = career.CVConfig{Name: "test"} // want `inline career\.CVConfig\{\} in test file; use fixtures package instead`
}

func helperWithInlineSectionContentGroup() {
	_ = career.SectionContentGroup{Header: "test"} // want `inline career\.SectionContentGroup\{\} in test file; use fixtures package instead`
}

func helperWithPointerLiteral() {
	_ = &career.Event{ID: "1"} // want `inline career\.Event\{\} in test file; use fixtures package instead`
}

func helperWithSliceLiteral() {
	_ = []*career.Event{
		{ID: "1"}, // want `inline career\.Event\{\} in test file; use fixtures package instead`
		{ID: "2"}, // want `inline career\.Event\{\} in test file; use fixtures package instead`
	}
}
