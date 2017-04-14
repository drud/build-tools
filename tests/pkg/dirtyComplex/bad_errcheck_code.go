package dirtyComplex

import "os"

// AFuncWithMissingErrCheck doesn't do anything useful - it just fails, and the err is lost
func AFuncWithMissingErrCheck(s string) {
	os.Chown("/never/will/this/file/exist", 99999, 99999)
}
