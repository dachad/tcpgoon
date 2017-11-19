package cmdutil

import (
	"testing"
)

func TestStringInSlice(t *testing.T) {
	testString := "test"
	testList := []string{"test", "tset"}
	if !stringInSlice(testString, testList) {
		t.Error("stringInSlice")
	}
}
