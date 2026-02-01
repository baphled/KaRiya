package clean

type localStruct struct {
	Name string
}

func helperWithLocalStruct() {
	_ = localStruct{Name: "test"}
}
