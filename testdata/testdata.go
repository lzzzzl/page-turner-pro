package testdata

import (
	"path/filepath"
	"runtime"
)

var basepath string

const (
	TestDataUser = "users.yaml"
)

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

func Path(rel string) string {
	return filepath.Join(basepath, rel)
}
