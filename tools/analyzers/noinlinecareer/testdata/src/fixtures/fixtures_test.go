package fixtures

import "fake/domain/career"

func Event(id string) *career.Event {
	return &career.Event{ID: id, Text: "fixture event " + id}
}

func Fact(id string) *career.Fact {
	return &career.Fact{ID: id, Text: "fixture fact " + id}
}
