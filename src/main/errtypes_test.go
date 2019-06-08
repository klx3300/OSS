package main

import (
	"errtypes"
	"testing"
)

func TestErrEqual(t *testing.T) {
	if errtypes.GenTestError() != errtypes.ETEST {
		t.Fatal("Inequal:", errtypes.GenTestError(), errtypes.ETEST)
	}
}
