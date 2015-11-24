package main

import "testing"

func TestParseFlag(t *testing.T) {
	args := []string{"letter", "-g", "**/*.go", "-c", "ls {{.File}}"}
	globs, commands, err := ParseFlag(args)
	if err != nil {
		t.Error(err)
	}
	if len(globs) != 1 {
		t.Errorf("len(globs) should be 1, but got %d", len(globs))
	}
	if len(commands) != 1 {
		t.Errorf("len(commands) should be 1, but got %d", len(commands))
	}
}
