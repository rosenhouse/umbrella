package main_test

import (
	"os/exec"

	"testing"

	"github.com/rosenhouse/umbrella"
)

func TestFromInsideThePackage(t *testing.T) {
	coverageCollector, err := umbrella.New()
	if err != nil {
		t.Errorf("init umbrella: %s", err)
		return
	}

	pathToProgram, err := coverageCollector.Build("github.com/rosenhouse/umbrella/example/program")
	if err != nil {
		t.Errorf("build binary: %s", err)
		return
	}
	defer coverageCollector.CleanupBuildArtifacts()

	cmd := exec.Command(pathToProgram, "one", "two", "three")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("running binary: %s", err)
		return
	}

	if string(output) != "onetwothree\n" {
		t.Errorf("unexpected output from test binary: %s", string(output))
		return
	}
}
