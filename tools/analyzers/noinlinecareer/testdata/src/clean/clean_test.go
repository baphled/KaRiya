package clean

import (
	"time"
)

type localStruct struct {
	Name string
}

func helperWithLocalStruct() {
	_ = localStruct{Name: "test"}
}

func helperWithExternalNonCareerType() {
	_ = time.Time{}
}
