package dirtyComplex

import "fmt"

// SomeExportedFunction doesn't do anything but introduce an unreachable error for govet
func SomeExportedFunction(s string) {
	fmt.Println(s)
	// Note that this exported function deliberately does *not* have a comment; that will trigger golint
	return
	panic("unreachable") // This is deliberately unreachable to test go vet functionality
}
