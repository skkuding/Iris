package testcase

import (
	"encoding/json"
	"fmt"
)

type TestcaseElement struct {
	In  string
	Out string
}

type Testcase struct {
	Data []TestcaseElement
}

func (t *Testcase) Count() int {
	fmt.Println(t.Data)
	return len(t.Data)
	// return 3
}

func (t Testcase) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Testcase) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
