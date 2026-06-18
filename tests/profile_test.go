package tests

import (
	"testing"

	generate "hera/core"
)

func TestProfileCreation(t *testing.T) {
	p := &generate.Profile{Name: "HER"}
	if p.Name != "HER" {
		t.Errorf("expected HER, got %s", p.Name)
	}
}

func TestProfileLettersRange(t *testing.T) {
	letters := [3]rune{'A', 'B', 'C'}
	name := string(letters[:])
	if name != "ABC" {
		t.Errorf("expected ABC, got %s", name)
	}
}
