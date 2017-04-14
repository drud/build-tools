package dirtyComplex

// SomeExportedFunction doesn't do anything but introduce an unused return for staticcheck
func yetAnotherExportedFunction() int {
	x := 1
	return x
}
