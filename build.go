package umbrella

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var ErrMissingHook = errors.New("program source missing umbrella hook")

type builder struct {
	goPath string
	goRoot string
}

func (b *builder) Build(outPath, testServerAddr, pkgPath string, coverPkg []string) error {
	buildCmd := exec.Command("go", "test",
		"-covermode", "set",
		"-coverpkg", getCoverPkg(coverPkg),
		"-c",
		"-o", outPath,
		"-tags", "umbrella_testrunmain",
		"-ldflags", fmt.Sprintf("-X %s.testServerAddr=%s", pkgPath, testServerAddr),
		pkgPath,
	)
	buildCmd.Env = []string{"PATH=" + os.Getenv("PATH")}
	if b.goPath != "" {
		buildCmd.Env = append(buildCmd.Env, "GOPATH="+b.goPath)
	}
	if b.goRoot != "" {
		buildCmd.Env = append(buildCmd.Env, "GOROOT="+b.goRoot)
	}
	msg, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go test -c failed: \n%s", msg)
	}

	return b.verifyBinary(outPath)
}

func (b *builder) verifyBinary(binPath string) error {
	binData, err := ioutil.ReadFile(binPath)
	if err != nil {
		return ErrMissingHook
	}

	evidenceOfHook := []byte("TestRunWithUmbrellaCoverage")

	if !bytes.Contains(binData, evidenceOfHook) {
		return ErrMissingHook
	}

	return nil
}

func getCoverPkg(coverPkg []string) string {
	if len(coverPkg) == 0 {
		return ""
	}
	return strings.Join(coverPkg, ",")
}
