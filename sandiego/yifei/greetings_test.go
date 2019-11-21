package mdchecker

import (
	"testing"
)

func TestGreetingFile(t *testing.T) {
	if err := GreetingCheck(); err != nil {
		t.Errorf(err.Error())
	}
}
