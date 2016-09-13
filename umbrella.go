package umbrella

import (
	"flag"
	"fmt"
	"go/build"
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

func New() (Collector, error) {
	accPath, err := getAccumulatorPath()
	if err != nil {
		return nil, err
	}
	return &collector{
		accumulator: accPath,
	}, nil
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

func getAccumulatorPath() (string, error) {
	orig := getInProcessFilePath()
	if orig == "" {
		return "", nil
	}
	dir := filepath.Dir(orig)

	name := filepath.Base(orig)
	if strings.HasSuffix(name, ".coverprofile") {
		name = strings.Replace(name, ".coverprofile", ".external.coverprofile", -1)
	} else {
		name += "-external"
	}

	return filepath.Abs(filepath.Join(dir, name))
}

type collector struct {
	mutex       sync.Mutex
	binDirs     []string
	accumulator string
}

func (c *collector) Build(pkgPath string) (string, error) {
	dir, err := ioutil.TempDir("", "build")
	if err != nil {
		return "", err
	}

	err = c.addTestHook(pkgPath)
	if err != nil {
		return "", fmt.Errorf("adding test hook: %s", err)
	}

	outPath := filepath.Join(dir, filepath.Base(pkgPath))

	buildCmd := exec.Command("go", "test",
		"-covermode", "set",
		"-c",
		"-o", outPath,
		"-tags", "testrunmain",
		"-ldflags", "-X github.com/rosenhouse/umbrella.coverProfilePath="+c.accumulator,
		pkgPath,
	)
	msg, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go test -c failed: \n%s", msg)
	}

	c.mutex.Lock()
	c.binDirs = append(c.binDirs, dir)
	c.mutex.Unlock()

	return outPath, nil
}

func (c *collector) hookContents() string {
	template :=
		`// +build testrunmain

package main

import (
	"os"
	"testing"

	"github.com/rosenhouse/umbrella"
)

func TestRunWithCoverage(t *testing.T) {
	main()
	os.Stdout, _ = os.Create(os.DevNull)
}

func TestMain(m *testing.M) {
	umbrella.PrepCoverage()
	os.Exit(m.Run())
}
`
	return template
}

const hookFileName = "umbrella_hook_test.go"

func (c *collector) getHookFilePath(pkgPath string) (string, error) {
	pkg, err := build.Default.Import(pkgPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}

	return filepath.Join(pkg.Dir, hookFileName), nil
}

func (c *collector) addTestHook(pkgPath string) error {
	hookFile, err := c.getHookFilePath(pkgPath)
	if err != nil {
		return err
	}

	hookContents := c.hookContents()
	return ioutil.WriteFile(hookFile, []byte(hookContents), 0666)
}

func (c *collector) removeTestHook(pkgPath string) error {
	hookFile, err := c.getHookFilePath(pkgPath)
	if err != nil {
		return err
	}
	return os.Remove(hookFile)
}

func (c *collector) CleanupBuildArtifacts() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, dir := range c.binDirs {
		os.RemoveAll(dir)
	}
}

// path to save coverage data
// override this at link-time, e.g.
//    go build -ldflags '-X github.com/rosenhouse/umbrella.coverProfilePath=/tmp/some/path'
var coverProfilePath = "REPLACEME.coverprofile"

func PrepCoverage() {
	if coverProfilePath != "" {
		flag.Set("test.coverprofile", coverProfilePath)
	}
	flag.Set("test.run", "TestRunWithCoverage")
	origArgs := os.Args[1:]
	os.Args = append([]string{os.Args[0], "spacer"}, origArgs...)
	flag.Parse()
	os.Args = append([]string{os.Args[0]}, origArgs...)
}
