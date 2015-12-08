// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
// tokens_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFileReader struct {
	inlines   []string
	nextIndex int
}

func makeMockFileReader(theInlines []string) *mockFileReader {
	return &mockFileReader{inlines: theInlines}
}
func (theMockFileReader *mockFileReader) getLine(assert *assert.Assertions) (theLine string, eof bool) {
	numLines := len(theMockFileReader.inlines)
	assert.True(theMockFileReader.nextIndex < numLines+1)
	if theMockFileReader.nextIndex == numLines {
		theMockFileReader.nextIndex++
		theLine = ""
		eof = true
		//	log.Print("end of mock input")
		return
	}
	//log.Printf("using mock line %d of %d", theMockFileReader.nextIndex, numLines)
	theLine = theMockFileReader.inlines[theMockFileReader.nextIndex]
	theMockFileReader.nextIndex++
	return
}
func (theMockFileReader *mockFileReader) verify(assert *assert.Assertions) {
	assert.True(theMockFileReader.nextIndex == len(theMockFileReader.inlines)+1)
}
func sendInlines(theMockFileReader *mockFileReader, assert *assert.Assertions) {
	//logging = true
	theLine, eof := theMockFileReader.getLine(assert)
	lineNumber := 1
	for !eof {
		//log.Printf("scan line %d", lineNumber)
		lineChan <- scanLine{text: theLine, number: lineNumber, eof: false}
		theLine, eof = theMockFileReader.getLine(assert)
		lineNumber++
	}
	//log.Print("end of scan lines")
	lineChan <- scanLine{text: "", number: lineNumber, eof: true}
}
func TestConv1(t *testing.T) {
	assert := assert.New(t)
	lineChan = make(chan scanLine)
	tokenChan = make(chan token)
	go getTokens()
	theMockFileReader := makeMockFileReader([]string{"test"})
	go sendInlines(theMockFileReader, assert)
	testToken := <-tokenChan
	checkToken(&testToken, textTokenType, "test", 2, assert) //??
	select {
	case testToken = <-tokenChan:
		assert.Equal(eofTokenType, testToken.theType, "should be at eof")
	default:
		assert.Fail("block on token channel")
	}
	theMockFileReader.verify(assert)
}
func TestConv2(t *testing.T) {
	assert := assert.New(t)
	lineChan = make(chan scanLine)
	tokenChan = make(chan token)
	go getTokens()
	inlines := []string{
		"/* comment */ $[test]$",
		"$(expr)$ some text",
		"$/doSomeJS();/$ $<html>$"}
	theMockFileReader := makeMockFileReader(inlines)
	go sendInlines(theMockFileReader, assert)
	testToken := <-tokenChan
	//	checkToken(&testToken, textTokenType, "", 1, assert)
	//	testToken = <-tokenChan
	checkToken(&testToken, textTokenType, " ", 1, assert)
	testToken = <-tokenChan
	checkToken(&testToken, identTokenType, "test", 1, assert)
	testToken = <-tokenChan
	//	checkToken(&testToken, textTokenType, "", 2, assert)
	//	testToken = <-tokenChan
	checkToken(&testToken, jsExprTokenType, "expr", 2, assert)
	testToken = <-tokenChan
	checkToken(&testToken, textTokenType, " some text", 3, assert) //??
	testToken = <-tokenChan
	checkToken(&testToken, jsCodeTokenType, "doSomeJS();", 3, assert)
	testToken = <-tokenChan
	checkToken(&testToken, textTokenType, " ", 3, assert)
	testToken = <-tokenChan
	checkToken(&testToken, htmlTokenType, "html", 3, assert)
	select {
	case testToken = <-tokenChan:
		assert.Equal(eofTokenType, testToken.theType, "should be at eof")
	default:
		assert.Fail("block on token channel")
	}
	theMockFileReader.verify(assert)
}
func checkToken(testToken *token, theTokenType tokenType, theText string, theLineNum int, assert *assert.Assertions) {
	//log.Printf("checking token: '%s'", testToken.text)
	assert.Equal(theTokenType, testToken.theType, "wrong token type for '%s'", testToken.text)
	assert.Equal(theText, testToken.text, "wrong token")
	assert.Equal(theLineNum, testToken.lineNumber, "wrong line number for '%s'", testToken.text)
}
// This file is part of Foobar. Foobar is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. Foobar is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with Foobar. If not, see <http://www.gnu.org/licenses/>.