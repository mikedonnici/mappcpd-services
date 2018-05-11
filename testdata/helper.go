package testdata

import "testing"

const success = "\u2713"
const failure = "\u2717"

type Helper struct {}

func NewHelper() *Helper {
	return &Helper{}
}

func (h Helper) Result(t *testing.T, expect, result interface{}) {

	if result != expect {
		t.Fatalf("\n\t%s expected: %v, result: %v", failure, expect, result)
	}
	t.Logf("\n\t%s expected result: %v", success, result)
}

func (h Helper) MessageResult(t *testing.T, message string, expect, result interface{}) {

	if result != expect {
		t.Fatalf("%s\n\t%s expected: %v, result: %v", message, failure, expect, result)
	}
	t.Logf("%s\n\t%s expected result: %v", message, success, result)
}
