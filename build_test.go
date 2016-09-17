package umbrella_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rosenhouse/umbrella"
)

var _ = Describe("Build", func() {
	var (
		goPath    string
		pkgPath   string
		collector umbrella.Collector
	)

	BeforeEach(func() {
		var err error
		goPath, err = ioutil.TempDir("", "gopath")
		Expect(err).NotTo(HaveOccurred())

		pkgPath = filepath.Join(goPath, "src", "something")
		Expect(os.MkdirAll(pkgPath, 0777)).To(Succeed())

		collector = umbrella.New(goPath)
		Expect(err).NotTo(HaveOccurred())

		const program = `package main; func main() {}`
		Expect(ioutil.WriteFile(filepath.Join(pkgPath, "main.go"),
			[]byte(program), 0600)).To(Succeed())

		const test = `package main_test

import "testing"

func Test1(t *testing.T) {}
`
		Expect(ioutil.WriteFile(filepath.Join(pkgPath, "main_test.go"),
			[]byte(test), 0600)).To(Succeed())

		bumbershootCmd := exec.Command(pathToBumbershoot, "-o", filepath.Join(pkgPath, "hook_test.go"))
		Expect(bumbershootCmd.Run()).To(Succeed())
	})

	AfterEach(func() {
		collector.CleanupBuildArtifacts()
	})

	Context("when the main_test.go file is missing", func() {
		BeforeEach(func() {
			Expect(os.RemoveAll(filepath.Join(pkgPath, "main_test.go"))).To(Succeed())
		})

		It("builds correctly", func() {
			binPath, err := collector.Build("something")
			Expect(err).NotTo(HaveOccurred())
			Expect(binPath).To(BeAnExistingFile())
		})
	})

	Context("when the hook file is missing", func() {
		BeforeEach(func() {
			Expect(os.RemoveAll(filepath.Join(pkgPath, "hook_test.go"))).To(Succeed())
		})

		It("fails to build the binary, printing a useful error message", func() {
			_, err := collector.Build("something")
			Expect(err).To(MatchError("program source missing umbrella hook"))
		})
	})

	Context("when no _test.go files are present in the package", func() {
		BeforeEach(func() {
			Expect(os.RemoveAll(filepath.Join(pkgPath, "hook_test.go"))).To(Succeed())
			Expect(os.RemoveAll(filepath.Join(pkgPath, "main_test.go"))).To(Succeed())
		})

		It("fails to build the binary, printing a useful error message", func() {
			_, err := collector.Build("something")
			Expect(err).To(MatchError("program source missing umbrella hook"))
		})
	})
})
