package dirtyComplex

import "fmt"

// ADummyFunctionWithBadCommand is just a function with a gofmt error - misformatted spaces before line comment
func ADummyFunctionWithBadCommand(s string) {
	fmt.Println(s)                    // This command is way out to the right to make gofmt complain. Leave it here.
}
