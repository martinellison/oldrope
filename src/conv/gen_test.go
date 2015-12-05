// gen_test.go
package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen1(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	genStart(testBuffer)
	testGen("<!DOCTYPE html>", testBuffer, t)
}

func TestGen2(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	genHeader(testBuffer)
	testGen("/* created by program on ", testBuffer, t)
}
func TestGen3(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	genJsStart(testBuffer)
	testGen("var gd = {};", testBuffer, t)
}
func TestGen4(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	genJsEnd(testBuffer)
	testGen("setPage('start');", testBuffer, t)
}
func TestGen5(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	genEnd(testBuffer)
	testGen("</body>", testBuffer, t)
}
func TestGen6(t *testing.T) {
	theOutData.Pages = make([]*outPage, 0, 0)
	testBuffer := new(bytes.Buffer)
	expandTemplate(testBuffer)
	testGen("pages = {", testBuffer, t)
}
func testGen(expected string, loadedBuffer *bytes.Buffer, t *testing.T) {
	assert := assert.New(t)
	testBuffer := new(bytes.Buffer)
	makeTemplate()
	genHeader(testBuffer)
	expectedStart := "/* created by program on "
	actualStart := testBuffer.String()[0:len(expectedStart)]
	assert.Equal(expectedStart, actualStart)
}
