package dirtyUnit

import (
	"fmt"

	"github.com/golang/example/stringutil"
)

func ExampleFailingReverse() {
	fmt.Println(stringutil.Reverse("hello"))
	// Output: ollehnopeNotTheCorrectOutput
}
