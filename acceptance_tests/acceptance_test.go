package acceptance_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Coverage of external binaries", func() {
	It("collects coverage profile data from an external binary run", func() {
		Expect(true).To(BeTrue())

		const pkgPrefix = "github.com/rosenhouse/umbrella/example"
		pkgsForCoverage := strings.Join(
			[]string{
				pkgPrefix + "/lib",
				pkgPrefix + "/program",
			}, ",",
		)
		testPkg := pkgPrefix + "/tests"

		profileDir, err := ioutil.TempDir("", "profiles")
		Expect(err).NotTo(HaveOccurred())

		coverProfilePath := filepath.Join(profileDir, "testrun.coverprofile")

		By("running the example test suite")
		Expect(runAndWait(
			exec.Command("go", "test",
				"-covermode", "set",
				"-coverpkg", pkgsForCoverage,
				"-coverprofile", coverProfilePath,
				testPkg),
		)).To(ContainSubstring("ok"))

		By("analyzing the code coverage from that test run")
		Expect(runAndWait(
			exec.Command("go", "tool", "cover", "-func", coverProfilePath)),
		).To(ContainSubstring("total: (statements) 100.0%"))
	})
})

func runAndWait(cmd *exec.Cmd) string {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "5s").Should(gexec.Exit(0))

	return string(session.Out.Contents())
}
