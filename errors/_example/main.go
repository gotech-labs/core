package main

import (
	"io/ioutil"
	"os"

	"github.com/gotech-labs/core/errors"
	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/core/log/formats"
)

var (
	FileAccessError = errors.Define("file i/o error")
	ValidationError = errors.Define("file i/o error")
	NetworkError    = errors.Define("network i/o error")
)

func main() {
	log.SetGlobalLogger(os.Stdout, log.WithFormat(formats.JSON))
	path := "/abc"
	err := call_1(path)
	log.Error(err.Error(), "error", err)
}

func call_1(dir string) error {
	if _, err := ioutil.ReadDir(dir); err != nil {
		return FileAccessError.Wrap(err, "failed to read directory", "userID", 123)
	}
	return nil
}
