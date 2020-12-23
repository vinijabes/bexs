package utils

import "testing"

func Assert(t *testing.T, expected, got interface{}) {
	if expected != got {
		t.Fatalf("expected %v. got %v", expected, got)
	}
}
