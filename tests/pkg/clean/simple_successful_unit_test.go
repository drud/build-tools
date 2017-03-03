package clean

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessfulMath(t *testing.T) {
	assert := assert.New(t)

	assert.EqualValues(2+2, 4)
}
