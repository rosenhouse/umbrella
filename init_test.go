package umbrella_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/rosenhouse/umbrella"

	"testing"
)

func TestUmbrella(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Umbrella Suite")
}

var pathToBumbershoot string

var _ = BeforeSuite(func() {
	var err error
	pathToBumbershoot, err = gexec.Build("github.com/rosenhouse/umbrella/cmd/bumbershoot")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	umbrella.CleanupBuildArtifacts()
})
