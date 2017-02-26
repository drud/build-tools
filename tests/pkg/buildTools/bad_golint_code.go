package buildTools

import "fmt"

func DummyExported_function(s string) {
	// Note that this exported function deliberately does *not* have a comment; that will trigger golint
	fmt.Println(s)
}
