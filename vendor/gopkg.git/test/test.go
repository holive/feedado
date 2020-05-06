package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// Load is used by tests to load mocks. It used the runtime.Caller to get the request file directory
// and load all the files from a testdata folder at the same level.
func Load(name string) []byte {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("could not get the caller that invoked load")
	}

	path := fmt.Sprintf("%s/testdata/%s", filepath.Dir(file), name)
	f, err := os.Open(path)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("error during open '%s'", path)))
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("error during read '%s'", path)))
	}
	return content
}
