package dirtyComplex

import "fmt"

// AnotherExportedFunction doesn't do anything but introduce an unused return for staticcheck
func AnotherExportedFunction(s string) {
	num, err := fmt.Println(s)
	num, err = fmt.Println(s)
	fmt.Println(num, err)
}
