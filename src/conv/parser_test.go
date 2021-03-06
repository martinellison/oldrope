// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
// parser_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser1(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "fred", lineNumber: 1}}
	//	assert.False(theParser.expectIdent("fred"))
	assert.Equal("fred", theParser.tokText())
	assert.Equal(textTokenType, theParser.tokTyp())
}

func TestParserStop1(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "fred", lineNumber: 1}}
	assert.False(theParser.stopped([]string{}))
}

func TestParserStop2(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: identTokenType, text: "fred", lineNumber: 1}}
	assert.False(theParser.stopped([]string{}))
}
func TestParserStop3(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: identTokenType, text: "fred", lineNumber: 1}}
	assert.True(theParser.stopped([]string{"fred"}))
}
func TestParserStop4(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: eofTokenType, text: "", lineNumber: 1}}
	assert.True(theParser.stopped([]string{"fred"}))
}
func TestParserIdent1(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: identTokenType, text: "fred", lineNumber: 1}}
	assert.True(theParser.tokIsIdent("fred"))
}
func TestParserIdent2(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "fred", lineNumber: 1}}
	assert.False(theParser.tokIsIdent("fred"))
}
func TestParserIdent3(t *testing.T) {
	assert := assert.New(t)
	theParser := &parser{theCurrentToken: token{theType: identTokenType, text: "fred", lineNumber: 1}}
	assert.False(theParser.tokIsIdent("bill"))
}

func TestParserString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("span", spanFragType.String())
	assert.Equal("div", divFragType.String())
	assert.Equal("para", paraFragType.String())
	assert.Equal("jsCode", jsCodeFragType.String())
	assert.Equal("jsExpr", jsExprFragType.String())
	assert.Equal("text", textFragType.String())
	assert.Equal("link", linkFragType.String())
	assert.Equal("html", htmlFragType.String())
	assert.Equal("include", includeFragType.String())
}
func makeTokens(theTokens []token) {
	for _, theToken := range theTokens {
		logfIfLogging("pushing token '%s'", theToken.text)
		tokenChan <- theToken
	}
}

func TestParserGetToken1(t *testing.T) {
	//logging = true
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	assert.Equal(textTokenType, theParser.tokTyp())
	assert.NotNil(tokenChan)
	theParser.getToken()
	assert.Equal(eofTokenType, theParser.tokTyp())
}
func TestParserGetToken2(t *testing.T) {
	//logging = true
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: identTokenType, text: "bill"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	assert.Equal(textTokenType, theParser.tokTyp())
	assert.NotNil(tokenChan)
	theParser.getToken()
	assert.Equal(identTokenType, theParser.tokTyp())
	assert.Equal("bill", theParser.tokText())
	theParser.getToken()
	assert.Equal(eofTokenType, theParser.tokTyp())
}

func TestParserParseBody1(t *testing.T) {
	testParseBody("div", divFragType, t)
}
func TestParserParseBody2(t *testing.T) {
	testParseBody("span", spanFragType, t)
}
func testParseBody(tag string, theFragType fragType, t *testing.T) {
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: identTokenType, text: tag}, {theType: textTokenType, text: "fred"}, {theType: identTokenType, text: "end"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(theFragType, theFragments[1].theFragType)
}

func TestParserParseInclude1(t *testing.T) {
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: identTokenType, text: "include"}, {theType: textTokenType, text: "fred"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(includeFragType, theFragments[1].theFragType)
	assert.Equal("fred", theFragments[1].auxName)
}

func TestParserParseCode1(t *testing.T) {
	assert := assert.New(t)
	theFragment := testParseFrag(jsCodeTokenType, jsCodeFragType, assert)
	assert.Equal("fred", theFragment.text)
}
func TestParserParseCode2(t *testing.T) {
	assert := assert.New(t)
	theFragment := testParseFrag(jsExprTokenType, jsExprFragType, assert)
	assert.Equal("fred", theFragment.text)
}
func TestParserParseCode3(t *testing.T) {
	assert := assert.New(t)
	theFragment := testParseFrag(htmlTokenType, htmlFragType, assert)
	assert.Equal("fred", theFragment.text)
}
func testParseFrag(theTokenType tokenType, theFragType fragType, assert *assert.Assertions) (theFragmant *fragment) {
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: theTokenType, text: "fred"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(theFragType, theFragments[1].theFragType)
	return theFragments[1]
}
func TestParserParseLink1(t *testing.T) {
	//logging = true
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: identTokenType, text: "link"}, {theType: identTokenType, text: "fred"}, {theType: identTokenType, text: "goto"}, {theType: identTokenType, text: "bill"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(linkFragType, theFragments[1].theFragType)
	assert.Equal("fred", theFragments[1].name)
	assert.Equal("bill", theFragments[1].auxName)
}
func TestParserParseLink2(t *testing.T) {
	//logging = true
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: identTokenType, text: "link"}, {theType: identTokenType, text: "fred"}, {theType: identTokenType, text: "act"}, {theType: identTokenType, text: "end"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(linkFragType, theFragments[1].theFragType)
	assert.Equal("fred", theFragments[1].name)
	assert.Equal("", theFragments[1].auxName)
}
func TestParserParseLink3(t *testing.T) {
	//logging = true
	assert := assert.New(t)
	tokenChan = make(chan token)
	go makeTokens([]token{{theType: identTokenType, text: "link"}, {theType: identTokenType, text: "fred"}, {theType: identTokenType, text: "end"}, {theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(linkFragType, theFragments[1].theFragType)
	assert.Equal("fred", theFragments[1].name)
	assert.Equal("", theFragments[1].auxName)
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
