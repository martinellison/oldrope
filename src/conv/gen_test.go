// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
// gen_test.go
package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen1(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	var theGenerator generator
	theGenerator.genStart(testBuffer)
	testGen("<!DOCTYPE html>", testBuffer, t)
}

func TestGen2(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	var theGenerator generator
	theGenerator.genHeader(testBuffer)
	testGen("/* created by program on ", testBuffer, t)
}
func TestGen3(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	var theGenerator generator
	theGenerator.genJsStart(testBuffer)
	testGen("var gd = {};", testBuffer, t)
}
func TestGen4(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	var theGenerator generator
	theGenerator.genJsEnd(testBuffer)
	testGen("setPage('start');", testBuffer, t)
}
func TestGen5(t *testing.T) {
	testBuffer := new(bytes.Buffer)
	var theGenerator generator
	theGenerator.genEnd(testBuffer)
	testGen("</body>", testBuffer, t)
}
func TestGen6(t *testing.T) {
	var theOutData outData
	theOutData.Pages = make([]*outPage, 0, 0)
	var theGenerator generator
	testBuffer := new(bytes.Buffer)
	theGenerator.makeTemplate()
	theGenerator.expandTemplate(testBuffer, theOutData)
	testGen("pages = {", testBuffer, t)
}
func testGen(expected string, loadedBuffer *bytes.Buffer, t *testing.T) {
	assert := assert.New(t)
	testBuffer := new(bytes.Buffer)
	var theGenerator generator
	theGenerator.makeTemplate()
	theGenerator.genHeader(testBuffer)
	expectedStart := "/* created by program on "
	actualStart := testBuffer.String()[0:len(expectedStart)]
	assert.Equal(expectedStart, actualStart)
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
