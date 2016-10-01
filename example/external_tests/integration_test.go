package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("The Program", func() {
	var cmd *exec.Cmd

	BeforeEach(func() {
		cmd = exec.Command(pathToProgram, "one", "two", "three")
		cmd.Env = []string{}
	})

	It("concatenates the arguments", func() {
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(ContainSubstring("onetwothree"))
	})

	Context("when the UPPERCASE env var is set", func() {
		BeforeEach(func() {
			cmd.Env = []string{"UPPERCASE=true"}
		})

		It("upper-cases the result", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out.Contents()).To(ContainSubstring("ONETWOTHREE"))
		})
	})
})
