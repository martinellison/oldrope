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
	go makeTokens([]token{token{theType: eofTokenType}})
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
	go makeTokens([]token{token{theType: identTokenType, text: "bill"}, token{theType: eofTokenType}})
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
	go makeTokens([]token{token{theType: identTokenType, text: tag}, token{theType: textTokenType, text: "fred"}, token{theType: identTokenType, text: "end"}, token{theType: eofTokenType}})
	theParser := &parser{theCurrentToken: token{theType: textTokenType, text: "", lineNumber: 1}}
	theFragments := theParser.parseBody([]string{"end"})
	assert.Equal(2, len(theFragments))
	assert.Equal(textFragType, theFragments[0].theFragType)
	assert.Equal(theFragType, theFragments[1].theFragType)
}
