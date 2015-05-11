package main

import (
	"testing"
)

func TestTrivial(t *testing.T) {
	got := 1 + 1
	want := 2
	if got != want {
		t.Errorf("1 + 1 is %v, want %v", got, want)
	}
}
