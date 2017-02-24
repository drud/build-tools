package main

import (
	"fmt"
	"github.com/drud/build-tools/tests/pkg/version"
)

func main() {
	fmt.Println("This is gofmtproblem.go version=" + version.VERSION + ". It has a gofmt problem.")          // Comment is way out to the right to make gofmt complain
}

func SomeExportedFunction(s string) {
	fmt.Println(s)
	// Note that this exported function deliberately does *not* have a comment; that will trigger golint
	return
	panic("unreachable")  // This is deliberately unreachable to test go vet functionality
}