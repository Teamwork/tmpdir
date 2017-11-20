// +build windows

package tmpdir

// Check if an error is the EEXISTS errno.
//
// This feature is not implemented on Windows, and provided only for
// compatibility.
func isEEXISTS(err error) bool {
	return false
}
