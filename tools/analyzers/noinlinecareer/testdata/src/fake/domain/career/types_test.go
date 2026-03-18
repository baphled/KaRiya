package career

func helperWithExportedType() {
	_ = Event{ID: "1"}
}

func helperWithUnexportedType() {
	_ = unexportedTestDouble{name: "test"}
}
