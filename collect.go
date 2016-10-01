package umbrella

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/rosenhouse/umbrella/accumulate"
	"github.com/rosenhouse/umbrella/merge"
)

func Build(pkgPath string, coverPkg ...string) (string, error) {
	return DefaultCollector.Build(pkgPath, coverPkg...)
}

func CleanupBuildArtifacts() {
	DefaultCollector.CleanupBuildArtifacts()
}

var DefaultCollector = New(os.Getenv("GOPATH"), os.Getenv("GOROOT"))

func New(goPath, goRoot string) Collector {
	return &collector{
		builder: builder{
			goPath: goPath,
			goRoot: goRoot,
		},
	}
}

type Collector interface {
	Build(pkgPath string, coverPkg ...string) (string, error)
	CleanupBuildArtifacts()
}

type collector struct {
	mutex             sync.Mutex
	binDirs           []string
	builder           builder
	server            accumulate.Server
	outputProfilePath string
}

func (c *collector) ensureServer() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.server != nil {
		return nil
	}

	var err error
	c.server, err = accumulate.NewServer()
	if err != nil {
		return err
	}

	return c.server.GoListenAndServe()
}

func (c *collector) Build(pkgPath string, coverPkg ...string) (string, error) {
	err := c.ensureServer()
	if err != nil {
		return "", err
	}

	c.outputProfilePath, err = getProfilePath()
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir("", "build")
	if err != nil {
		return "", err
	}

	outPath := filepath.Join(dir, filepath.Base(pkgPath))

	err = c.builder.Build(outPath, c.server.Address(), pkgPath, coverPkg)
	if err != nil {
		return "", err
	}

	c.mutex.Lock()
	c.binDirs = append(c.binDirs, dir)
	c.mutex.Unlock()

	return outPath, nil
}

func (c *collector) CleanupBuildArtifacts() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, dir := range c.binDirs {
		os.RemoveAll(dir)
	}

	if c.server != nil {
		coverFiles := c.server.ListAll()
		if c.outputProfilePath != "" && len(coverFiles) > 0 {
			fmt.Printf("collecting %d cover profiles into %s\n", len(coverFiles), c.outputProfilePath)
			err := merge.Files(c.outputProfilePath, coverFiles)
			if err != nil {
				panic(fmt.Errorf("merging coverage files: %s", err))
			}
		}

		c.server.Close()
	}
}
