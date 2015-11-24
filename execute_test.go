package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestExitStatus(t *testing.T) {
	c, err := ExitStatus(nil)
	if err != nil {
		t.Error(err)
	}
	if c != 0 {
		t.Errorf("Status code should be 0, but got %d", c)
	}

	c, err = ExitStatus(fmt.Errorf("hoge"))
	t.Log(err)
	if err == nil {
		t.Errorf("Error should not be nil, but got nil")
	}

	err = exec.Command("true").Run()
	c, err = ExitStatus(err)
	if err != nil {
		t.Error(err)
	}
	if c != 0 {
		t.Errorf("Status code should be 0, but got %d", c)
	}

	err = exec.Command("false").Run()
	c, err = ExitStatus(err)
	if err != nil {
		t.Error(err)
	}
	if c != 1 {
		t.Errorf("Status code should be 1, but got %d", c)
	}
}
