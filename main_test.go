package main

import "testing"

func TestValidStates(t *testing.T) {
	if isValidState("WRONG") {
		t.Error("WRONG is not a valid state")
	}

	if !isValidState("up") {
		t.Error("up should be a valid state")
	}
}
