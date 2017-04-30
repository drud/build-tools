package main

import (
	"fmt"
	"github.com/drud/build-tools/tests/pkg/version"
)

func main() {
	fmt.Println("This is build_tools_dummy.go version=")
	fmt.Println("Version:", version.VERSION)
	fmt.Println("Commit:", version.COMMIT)
	fmt.Println("Buildinfo:", version.BUILDINFO)
}
