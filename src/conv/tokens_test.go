// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
// tokens_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConv1(t *testing.T) {
	assert := assert.New(t)
	lineChan = make(chan scanLine)
	tokenChan = make(chan token)
	var theTokeniser tokeniser
	theTokeniser.init()
	go theTokeniser.getTokens()
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
	var theTokeniser tokeniser
	theTokeniser.init()
	go theTokeniser.getTokens()
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

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
