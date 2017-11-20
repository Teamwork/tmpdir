// Package tmpdir manages temporary directories for a server.
package tmpdir // import "github.com/teamwork/tmpdir"

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

// MkTemp creates a unique filename in the temporary directory and returns a
// writable file descriptor. This function is safe for any arbitrary filename.
//
// Subdirectories are not supported. The Setup() function must be called first.
func MkTemp(filename string) (*os.File, error) {
	if Dir == "" {
		panic("Dir is empty. Setup() needs to be called first.")
	}

	// Trim spaces for sanity.
	filename = strings.TrimSpace(filename)

	// Dir always has a / appended.
	filename = strings.TrimLeft(filename, "/")

	// Slashes and NULL bytes are not allowed in file paths.
	filename = strings.Replace(
		strings.Replace(filename, "\x00", "", -1), "/", "-", -1)

	filename, ext := splitExt(filename)

	// Max path length is 255, with some padding for adding random data if
	// needed.
	if len(filename)+len(ext)+len(Dir)+1 > 250 {
		filename = filename[:250-len(Dir)-len(ext)-1]
	}

	path := Dir + filename
	if ext != "" {
		path += "." + ext
	}

	i := 0
	for {
		i++

		if i > 50 {
			return nil, errors.New("could not create temporary file: too many attempts at creating a unique filename")
		}

		fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			if isEEXISTS(err) {
				path = fmt.Sprintf("%v%v-1", Dir, filename)
				if ext != "" {
					path += "." + ext
				}
				continue
			}
			return nil, err
		}

		return fp, nil
	}
}

// splitExt splits a path in the pathname without extension and the extension.
func splitExt(path string) (string, string) {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[:i], path[i+1:]
		}
	}
	return path, ""
}
