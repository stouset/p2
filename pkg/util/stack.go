package util

import (
	"path"
	"path/filepath"
)

type Caller struct {
	Filename string
}

// Get a caller object. This is intented to be passed runtime.Caller()
// from the target
func From(pc uintptr, file string, line int, ok bool) *Caller {
	if !ok {
		return nil
	}
	return &Caller{file}
}

func (c *Caller) Dirname() string {
	return path.Dir(c.Filename)
}

func (c *Caller) ExpandPath(pathstr string) string {
	return path.Join(c.Dirname(), pathstr)
}

func (c *Caller) Glob(pattern string) ([]string, error) {
	return filepath.Glob(path.Join(c.Dirname(), pattern))
}
