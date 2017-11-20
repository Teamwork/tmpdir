// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package tmpdir

import (
	"os"
	"syscall"
)

// Check if an error is the EEXISTS errno.
func isEEXISTS(err error) bool {
	if err == nil {
		return false
	}

	perr, ok := err.(*os.PathError)
	if !ok {
		return false
	}

	oserr, ok := perr.Err.(syscall.Errno)
	if !ok {
		return false
	}

	return oserr == syscall.EEXIST
}
