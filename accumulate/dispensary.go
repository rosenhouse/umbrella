package accumulate

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type dispensary struct {
	lock     sync.Mutex
	dir      string
	acquired []string
}

func (a *dispensary) AcquireOne() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	name := fmt.Sprintf("%08d", len(a.acquired))
	path := filepath.Join(a.dir, name)
	a.acquired = append(a.acquired, path)

	return path, nil
}

func (a *dispensary) ListAll() []string {
	a.lock.Lock()
	defer a.lock.Unlock()

	report := make([]string, len(a.acquired))
	copy(report, a.acquired)

	return report
}

func (a *dispensary) Cleanup() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	defer func() { a.dir = "" }()

	return os.RemoveAll(a.dir)
}
