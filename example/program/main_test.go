package main_test

import (
	"fmt"
	"os/exec"

	"testing"

	"github.com/rosenhouse/umbrella"
)

// This is an example of a in-package test of the external binary
// Before running this test, make sure that you first 'go generate'
// to create the required hook file

//go:generate bumbershoot

func TestFromInsideThePackage(t *testing.T) {
	fmt.Println("building the binary with coverage")
	pathToProgram, err := umbrella.Build("github.com/rosenhouse/umbrella/example/program")
	if err != nil {
		t.Errorf("build binary: %s", err)
		return
	}
	defer umbrella.CleanupBuildArtifacts()

	fmt.Println("executing the binary")
	cmd := exec.Command(pathToProgram, "one", "two", "three")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("running binary: %s", err)
		return
	}

	fmt.Println("verifying coverage")
	if string(output) != "onetwothree\n" {
		t.Errorf("unexpected output from test binary: %s", string(output))
		return
	}
	fmt.Println("complete")
}
