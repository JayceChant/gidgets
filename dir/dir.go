// Wrapper for the directory functions in os package. (mainly for os.Stat so far)
// Since the error returned by os.Stat indicate both REAL exception and NORMAL state,
// it's error-prone in use, encapsulation of the best practice is required.
package dir

import (
	"os"
)

// PathExists returns if the given path exists.
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir returns if the given path is a folder or not.
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return fi.IsDir()
}

// IsFile returns if the given path is a file or not.
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return !fi.IsDir()
}
