package umbrella_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

func TestUmbrella(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Umbrella Acceptance Tests")
}

var _ = Describe("Coverage of external binaries", func() {
	const (
		pkgPrefix        = "github.com/rosenhouse/umbrella/example"
		inProcessProfile = "some-example.coverprofile"
		externalProfile  = "some-example.external.coverprofile"
	)

	var (
		pkgsForCoverage string
		testPkg         string
		workDir         string
	)

	BeforeEach(func() {
		var err error
		workDir, err = ioutil.TempDir("", "work-dir")
		Expect(err).NotTo(HaveOccurred())

		pkgsForCoverage = strings.Join(
			[]string{
				pkgPrefix + "/lib",
				pkgPrefix + "/program",
			}, ",",
		)
		testPkg = pkgPrefix + "/tests"
	})

	AfterEach(func() {
		Expect(os.RemoveAll(workDir)).To(Succeed())
	})

	It("collects coverage profile data from an external binary run", func() {
		cmd := exec.Command("go", "test",
			"-v",
			"-covermode", "set",
			"-coverpkg", pkgsForCoverage,
			"-coverprofile", inProcessProfile,
			testPkg)
		cmd.Dir = workDir

		Expect(runAndWait(cmd)).To(ContainSubstring("ok"))

		Expect(runAndWait(
			exec.Command("go", "tool", "cover", "-func", filepath.Join(workDir, externalProfile))),
		).To(MatchRegexp(`total:\s+\(statements\)\s+66\.7\%`))
	})

	Context("when the in-process covermode is atomic", func() {
		It("generates the external profile in the same mode", func() {
			cmd := exec.Command("go", "test",
				"-covermode", "atomic",
				"-coverprofile", inProcessProfile,
				testPkg)
			cmd.Dir = workDir

			Expect(runAndWait(cmd)).To(ContainSubstring("ok"))
			Expect(ioutil.ReadFile(filepath.Join(workDir, inProcessProfile))).To(HavePrefix("mode: atomic"))
		})
	})

	Context("when an outputdir is provided", func() {
		var outputDir string
		BeforeEach(func() {
			var err error
			outputDir, err = ioutil.TempDir("", "outputs")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(outputDir)).To(Succeed())
		})

		It("stores the external coverage profile in the outputdir", func() {
			Expect(runAndWait(
				exec.Command("go", "test",
					"-covermode", "set",
					"-coverpkg", pkgsForCoverage,
					"-coverprofile", inProcessProfile,
					"-outputdir", outputDir,
					testPkg),
			)).To(ContainSubstring("ok"))

			Expect(runAndWait(
				exec.Command("go", "tool", "cover", "-func", filepath.Join(outputDir, externalProfile))),
			).To(MatchRegexp(`total:\s+\(statements\)\s+66\.7\%`))
		})
	})
})

func runAndWait(cmd *exec.Cmd) string {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "10s").Should(gexec.Exit(0))

	return string(session.Out.Contents())
}
