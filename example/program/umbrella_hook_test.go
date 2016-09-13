// +build testrunmain

package main

import (
	"os"
	"testing"

	"github.com/rosenhouse/umbrella"
)

func TestRunWithCoverage(t *testing.T) {
	main()
	os.Stdout, _ = os.Create(os.DevNull)
}

func TestMain(m *testing.M) {
	umbrella.PrepCoverage()
	os.Exit(m.Run())
}
