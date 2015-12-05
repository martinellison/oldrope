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
