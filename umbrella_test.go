package umbrella_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Coverage of external binaries", func() {
	const (
		pkgPrefix        = "github.com/rosenhouse/umbrella/example"
		inProcessProfile = "some-example.coverprofile"
		externalProfile  = "some-example.coverprofile.external.coverprofile"
	)

	var (
		pkgsForCoverage string
		testPkg         string
		workDir         string
		cmd             *exec.Cmd
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

		testPkg = pkgPrefix + "/external_tests"

		cmd = exec.Command("go", "test",
			"-v",
			"-covermode", "set",
			"-coverpkg", pkgsForCoverage,
			"-coverprofile", inProcessProfile,
			testPkg)
		cmd.Dir = workDir
	})

	AfterEach(func() {
		Expect(os.RemoveAll(workDir)).To(Succeed())
	})

	AssertCoverageFileGetsGenerated := func(expectedDir string) string {
		Expect(runAndWait(cmd)).To(ContainSubstring("ok"))

		coverStats := runAndWait(
			exec.Command("go", "tool", "cover", "-func", filepath.Join(expectedDir, externalProfile)),
		)
		Expect(coverStats).To(MatchRegexp(`total:\s+\(statements\)\s+.*%`))
		return coverStats
	}

	It("generates the coverage file", func() {
		AssertCoverageFileGetsGenerated(workDir)
	})

	It("aggregates coverage data from multiple test executions", func() {
		coverStats := AssertCoverageFileGetsGenerated(workDir)
		// main has a conditional.  there are two tests, each executing the program to cover a single branch
		// this expectation verifies that the coverage results from both test runs
		// get aggregated together to show full coverage
		Expect(coverStats).To(MatchRegexp(`example/program/main.go:11:	main			100.0%`))
	})

	Context("when the tests live in the same package as the binary", func() {
		BeforeEach(func() {
			testPkg = pkgPrefix + "/program"
		})

		It("generates the coverage file", func() {
			AssertCoverageFileGetsGenerated(workDir)
		})
	})

	Context("when the in-process covermode is atomic", func() {
		BeforeEach(func() {
			cmd.Args = []string{"go", "test",
				"-covermode", "atomic",
				"-coverprofile", inProcessProfile,
				testPkg}
		})

		It("generates the external profile in the same mode", func() {
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

			cmd.Args = []string{"go", "test",
				"-covermode", "set",
				"-coverpkg", pkgsForCoverage,
				"-coverprofile", inProcessProfile,
				"-outputdir", outputDir,
				testPkg}
		})

		AfterEach(func() {
			Expect(os.RemoveAll(outputDir)).To(Succeed())
		})

		It("generates the coverage file", func() {
			AssertCoverageFileGetsGenerated(outputDir)
		})
	})
})

func runAndWait(cmd *exec.Cmd) string {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "10s").Should(gexec.Exit(0))

	return string(session.Out.Contents())
}
