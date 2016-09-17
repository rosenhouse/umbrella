package umbrella

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Collector interface {
	Build(pkgPath string) (string, error)
	CleanupBuildArtifacts()
}

var DefaultCollector = New(os.Getenv("GOPATH"))

func Build(pkgPath string) (string, error) { return DefaultCollector.Build(pkgPath) }
func CleanupBuildArtifacts()               { DefaultCollector.CleanupBuildArtifacts() }

func New(goPath string) Collector {
	return &collector{goPath: goPath}
}

type collector struct {
	mutex   sync.Mutex
	binDirs []string
	goPath  string
}

func getInProcessFilePath() string {
	filename := flag.Lookup("test.coverprofile").Value.String()
	if filename == "" {
		return ""
	}
	dir := flag.Lookup("test.outputdir").Value.String()
	if dir == "" {
		return filename
	}
	return filepath.Join(dir, filename)
}

func getProfilePath() (string, error) {
	orig := getInProcessFilePath()
	if orig == "" {
		return "", nil
	}
	dir := filepath.Dir(orig)

	name := filepath.Base(orig)
	if strings.HasSuffix(name, ".coverprofile") {
		name = strings.Replace(name, ".coverprofile", ".external.coverprofile", -1)
	} else {
		name += ".external.coverprofile"
	}

	return filepath.Abs(filepath.Join(dir, name))
}

func (c *collector) Build(pkgPath string) (string, error) {
	profilePath, err := getProfilePath()
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir("", "build")
	if err != nil {
		return "", err
	}

	outPath := filepath.Join(dir, filepath.Base(pkgPath))

	buildCmd := exec.Command("go", "test",
		"-covermode", "set",
		"-c",
		"-o", outPath,
		"-tags", "umbrella_testrunmain",
		"-ldflags", fmt.Sprintf("-X %s.coverProfilePath=%s", pkgPath, profilePath),
		pkgPath,
	)
	if c.goPath != "" {
		buildCmd.Env = []string{"GOPATH=" + c.goPath}
	}
	msg, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go test -c failed: \n%s", msg)
	}

	if err := c.verifyBinary(outPath); err != nil {
		return "", err
	}

	c.mutex.Lock()
	c.binDirs = append(c.binDirs, dir)
	c.mutex.Unlock()

	return outPath, nil
}

var ErrMissingHook = errors.New("program source missing umbrella hook")

func (c *collector) verifyBinary(binPath string) error {
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

func (c *collector) CleanupBuildArtifacts() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, dir := range c.binDirs {
		os.RemoveAll(dir)
	}
}
