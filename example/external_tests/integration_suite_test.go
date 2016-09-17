package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rosenhouse/umbrella"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Some Integration Suite for the Fake Project")
}

var pathToProgram string

var _ = BeforeSuite(func() {
	var err error
	pathToProgram, err = umbrella.Build("github.com/rosenhouse/umbrella/example/program")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	umbrella.CleanupBuildArtifacts()
})
