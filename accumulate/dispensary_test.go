package accumulate

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dispensary", func() {
	var disp dispensary

	BeforeEach(func() {
		dataDir, err := ioutil.TempDir("", "umbrella.datadir")
		Expect(err).NotTo(HaveOccurred())

		disp = dispensary{dir: dataDir}
	})

	AssertContainingDirectoryExists := func(path string) {
		Expect(filepath.Dir(path)).To(BeADirectory())
	}

	It("dispenses temporary files, keeping track of each one", func() {
		f1, err := disp.AcquireOne()
		Expect(err).NotTo(HaveOccurred())
		AssertContainingDirectoryExists(f1)

		f2, err := disp.AcquireOne()
		Expect(err).NotTo(HaveOccurred())
		AssertContainingDirectoryExists(f2)

		Expect(f1).NotTo(Equal(f2))

		all := disp.ListAll()
		Expect(all).To(Equal([]string{f1, f2}))

		Expect(ioutil.WriteFile(f1, []byte("foo"), 0600)).To(Succeed())
		Expect(ioutil.WriteFile(f2, []byte("bar"), 0600)).To(Succeed())

		Expect(disp.Cleanup()).To(Succeed())

		Expect(f1).NotTo(BeAnExistingFile())
		Expect(f2).NotTo(BeAnExistingFile())
	})
})
