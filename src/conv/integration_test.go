// integration_test.go
package main

// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration1(t *testing.T) {
	testIntegration([]string{"$[page fred]$Hello"}, "pages = { fred: { init: function() { }, display: function(parts) { parts.push(\"Hello\"); }, fix: function(parts) { }, }, };", t)
}
func testIntegration(inLines []string, expected string, t *testing.T) {
	assert := assert.New(t)
	lineChan = make(chan scanLine)
	tokenChan = make(chan token)
	var theTokeniser tokeniser
	theTokeniser.init()
	go theTokeniser.getTokens()
	theMockFileReader := makeMockFileReader(inLines)
	go sendInlines(theMockFileReader, assert)
	theParseChan := make(chan *pageSet)
	var theParser parser
	go theParser.parse(theParseChan)
	thePageSet := <-theParseChan
	var theGenerator generator
	theGenerator.makeTemplate()
	var theOutData outData
	theOutData.makeGenData(thePageSet)
	testBuffer := new(bytes.Buffer)
	theGenerator.expandTemplate(testBuffer, theOutData)
	assert.Equal(expected, compress(testBuffer.String()))
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
