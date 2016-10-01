package accumulate_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rosenhouse/umbrella/accumulate"
)

var _ = Describe("RPC", func() {
	It("allows a remote client to acquire a file", func(done Done) {
		server, err := accumulate.NewServer()
		Expect(err).NotTo(HaveOccurred())

		Expect(server.GoListenAndServe()).To(Succeed())

		path, err := accumulate.AcquireRemote(server.Address())
		Expect(err).NotTo(HaveOccurred())

		Expect(path).NotTo(BeEmpty())

		Expect(server.ListAll()).To(Equal([]string{path}))

		Expect(server.Close()).To(Succeed())
		close(done)
	}, 2 /* <-- overall spec timeout in seconds */)
})
