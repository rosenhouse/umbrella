package umbrella_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestUmbrella(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Umbrella Suite")
}

var pathToGenBinary string

var _ = BeforeSuite(func() {
	var err error
	pathToGenBinary, err = gexec.Build("github.com/rosenhouse/umbrella/cmd/umbrella")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
