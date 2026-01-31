package nontest

import "career"

func productionCodeWithCareerLiteral() career.Event {
	return career.Event{ID: "1", Text: "allowed in non-test files"}
}
