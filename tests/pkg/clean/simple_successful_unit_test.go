package clean

import (
	"fmt"

	"github.com/golang/example/stringutil"
)

func ExampleSuccessfulReverse() {
	fmt.Println(stringutil.Reverse("hello"))
	// Output: olleh
}
