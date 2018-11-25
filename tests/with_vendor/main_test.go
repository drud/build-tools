package main

import (
	"testing"
	asrt "github.com/stretchr/testify/assert"

)

func TestSomething(t *testing.T) {
	t.Log("Yup, this is working")
	assert := asrt.New(t)

	assert.True(true)
}
