package clean

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessfulMath(t *testing.T) {
	a := assert.New(t)

	a.EqualValues(2+2, 4)
}
