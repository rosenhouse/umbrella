package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Some Integration Suite for the Fake Project")
}

var pathToProgram string

var _ = BeforeSuite(func() {
	var err error
	pathToProgram, err = gexec.Build("github.com/rosenhouse/umbrella/example/program")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
