package clean

import (
	"fmt"
	"github.com/golang/example/stringutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessfulReverse(t *testing.T) {
	assert := assert.New(t)

	s := fmt.Sprintf(stringutil.Reverse("hello"))
	assert.EqualValues(s, "olleh")
}
