package career

func helperWithExportedType() {
	_ = Event{ID: "1"} // want `inline career\.Event\{\} in test file; use fixtures package instead`
}

func helperWithUnexportedType() {
	_ = unexportedTestDouble{name: "test"}
}
