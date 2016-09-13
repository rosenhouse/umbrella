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

var coverageCollector umbrella.Collector
var pathToProgram string

var _ = BeforeSuite(func() {
	var err error
	coverageCollector, err = umbrella.New()
	Expect(err).NotTo(HaveOccurred())

	pathToProgram, err = coverageCollector.Build("github.com/rosenhouse/umbrella/example/program")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	coverageCollector.CleanupBuildArtifacts()
})
