package dirtyUnit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailingMath(t *testing.T) {
	assert := assert.New(t)

	assert.EqualValues(2+2, 5)
}
