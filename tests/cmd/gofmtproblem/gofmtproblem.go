package main

import (
	"fmt"
	"github.com/drud/build-tools/tests/pkg/version"
)

func main() {
	fmt.Println("This is gofmtproblem.go version=" + version.VERSION + ". It has a gofmt problem.")          // Comment is way out to the right


}
