package helpers

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

func PrintMyPath() {
	fmt.Println(RootDir())
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
