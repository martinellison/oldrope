// conv_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestITE1(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("a", ifThenElse(true, "a", "b"))
	assert.Equal("b", ifThenElse(false, "a", "b"))
}

func TestLOL1(t *testing.T) {
	assert := assert.New(t)
	logging = true
	logIfLogging("this should appear in the test output")
	logging = false
	logIfLogging("this should NOT appear in the test output")
	assert.True(true) //dummy to keep Go happy
}
