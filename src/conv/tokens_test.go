// tokens_test.go
package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConv(t *testing.T) {
	assert := assert.New(t)
	//log.Print("start\n")
	lineChan = make(chan scanLine, 1)
	tokenChan = make(chan token)
	go getTokens()
	theScanLine := scanLine{text: "test", number: 1, eof: false}
	//log.Print("to line chan\n")
	lineChan <- theScanLine
	var testToken token	
	theScanLine = scanLine{text: "", number: 2, eof: true}
	//log.Print("to line chan\n")
	lineChan <- theScanLine
	//log.Print("to token chan\n")
	testToken = <-tokenChan
	assert.Equal(textTokenType, testToken.theType)
	assert.Equal("test", testToken.text)
	assert.Equal(2, testToken.lineNumber)
	//log.Print("done\n")
}
