// Package tmpdir manages temporary directories for a server.
package tmpdir // import "github.com/teamwork/tmpdir"

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Dir is the temporary directory that Setup() created.
var Dir string

// Setup creates an unique temporary directory.
func Setup(path string) (err error) {
	Dir, err = ioutil.TempDir("", path)
	if err != nil {
		_ = Cleanup()
		return errors.WithStack(err)
	}
	Dir += "/"
	return nil
}

// Cleanup removes the tmp directory.
func Cleanup() error {
	return os.RemoveAll(Dir)
}
